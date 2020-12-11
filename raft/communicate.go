/*+-----------------------------+
 *| Author: Zoueature           |
 *+-----------------------------+
 *| Email: zoueature@gmail.com  |
 *+-----------------------------+
 */
package raft

import (
	"github.com/zoueature/ekv/raft/utils"
	"log"
	"net"
	"net/rpc"
	"time"
)

const (
	dialClusterTimeout = 300 * time.Millisecond
)

type Server interface {
	Vote(msg *VoteMsg, reply *VoteReplyMsg) error
	Log(log *LogMsg, reply *LogReply) error
}

type Rpc struct {
	rf *raft
}

// startDialNodes 尝试与集群中的其他节点建立连接
func (rf *raft) startDialNodes(hosts []string) {
	ch := make(chan string, len(hosts))
	for _, host := range hosts {
		ch <- host
	}
	rf.dialChan = ch
	go dialWorker(ch, rf)
}

func dialWorker(queue <-chan string, rf *raft) {
	for {
		host := <-queue
		conn, err := net.DialTimeout("tcp", host, dialClusterTimeout)
		if err == nil {
			cli := rpc.NewClient(conn)
			old, ok := rf.cluster[host]
			if ok {
				old.cli = cli
				old.healthy = true
			} else {
				rf.cluster[host] = &Cluster{
					host:    host,
					cli:     cli,
					healthy: true,
				}
			}
			return
		}
		log.Printf("dial node: %s error: %s", host, err.Error())
		go func() {
			//1秒后重试
			timer := time.NewTicker(1 * time.Second)
			<-timer.C
			rf.dialChan <- host
		}()
	}
}

// InitialRpcServer 初始化rpc服务
func InitialRpcServer(host string, rf Server) error {
	listener, err := net.Listen("tcp", host)
	if err != nil {
		return err
	}
	server := rpc.NewServer()
	err = server.Register(rf)
	if err != nil {
		return err
	}
	go server.Accept(listener)
	return nil
}

// Vote 候选者调用，用于投票选举领导人
func (r *Rpc) Vote(msg *VoteMsg, reply *VoteReplyMsg) (err error) {
	log.Printf("receive vote request: %+v", msg)
	reply.Term = r.rf.currentTerm
	if msg.Term < r.rf.currentTerm {
		reply.VoteGranted = false
		return
	}
	if (r.rf.votedFor == "" || r.rf.votedFor == msg.CandidateId) && (msg.LastLogIndex >= r.rf.commitIndex) {
		reply.VoteGranted = true
		goto follow
	}
	if r.rf.commitIndex > r.rf.lastApplied {
		_ = r.rf.Commit(r.rf.commitIndex, r.rf.lastApplied-r.rf.commitIndex)
	}
follow:
	if msg.Term > r.rf.currentTerm {
		r.rf.currentTerm = msg.Term
		r.rf.changeRole(Follower)
	}
	return
}

// Log leader调用， 用于复制日志条目
func (r *Rpc) Log(l *LogMsg, reply *LogReply) error {
	log.Printf("recive heart beat from %s, %v", l.LeaderID, l)
	r.rf.latestHeartBeat = time.Now().UnixNano()
	if l.Term >= r.rf.currentTerm {
		r.rf.currentTerm = l.Term
		r.rf.changeRole(Follower)
	}
	return nil
}

func (rf *raft) handlerHearBeat() {
	for {
		lastTime := rf.latestHeartBeat
		utils.BlockStaticTime(300*time.Millisecond)
		if rf.latestHeartBeat <= lastTime {
			rf.StartVote()
			break
		}
	}
}
