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
func (rf *raft) startDialNodes() {
	ch := make(chan string, len(rf.conf.Cluster))
	for _, host := range rf.conf.Cluster {
		ch <- host
	}
	rf.dialChan = ch
	go dialWorker(ch, rf)
}

func dialWorker(queue <-chan string, rf *raft) {
	for {
		host := <-queue
		log.Printf("%s dial %s", rf.host, host)
		go func() {
			conn, err := net.DialTimeout("tcp", host, dialClusterTimeout)
			if err == nil {
				cli := rpc.NewClient(conn)
				rf.lock.Lock()
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
				rf.lock.Unlock()
				return
			}
			log.Printf("%s dial node: %s error: %s", rf.host, host, err.Error())
			go func() {
				//1秒后重试
				timer := time.NewTicker(1 * time.Second)
				<-timer.C
				rf.dialChan <- host
			}()
		}()
	}
}

func (rf *raft) InitialService() error {
	listener, err := net.Listen("tcp", rf.host)
	if err != nil {
		log.Fatalf("start listener error: %s", err.Error())
		return err
	}
	// 初始化集群rpc服务
	err = InitialRpcServer(listener, &Rpc{rf: rf})
	if err != nil {
		return err
	}
	// 初始化leader对外服务
	rfSvc := newRaftService(listener)
	rfSvc.Start()
	return nil
}

// InitialRpcServer 初始化rpc服务
func InitialRpcServer(listener net.Listener, rf Server) error {
	server := rpc.NewServer()
	err := server.Register(rf)
	if err != nil {
		return err
	}
	go server.Accept(listener)
	return nil
}

// Vote 候选者调用，用于投票选举领导人
// 投票原则
// 1、先到先服务原则， 在满足其他条件的前提下， 谁先请求选举就投票给谁
// 2、谁的任期大选谁， 当选举请求的任期大于节点当前任期时， 投票
// 3、候选人的日志需要不小于当前节点， 即lastLogTerm大于当前任期或者任期相等并且lastLogIndex大于等于当前节点的日志
func (r *Rpc) Vote(msg *VoteMsg, reply *VoteReplyMsg) (err error) {

	reply.Term = r.rf.getTerm()
	if msg.Term < r.rf.getTerm() {
		//候选人的任期小于等于节点当前任期， 拒绝投票
		log.Printf("候选人的任期小于等于节点当前任期， 拒绝投票")
		return
	} else if msg.Term > r.rf.getTerm() {
		goto follow
	}
	if r.rf.getVote() != "" && r.rf.getVote() != msg.CandidateId {
		//已经投票给了别的节点， 拒绝投票
		log.Printf("已经投票给了别的节点， 拒绝投票")
		return
	}
	if msg.LastLogTerm < r.rf.logLastItem {
		//候选者最新日志任期小于当前节点最新日志的任期， 拒绝投票
		log.Printf("候选者最新日志任期小于当前节点最新日志的任期， 拒绝投票")
		return
	}
	if msg.LastLogTerm == r.rf.logLastItem && msg.LastLogIndex < r.rf.logLastIndex {
		//候选者最新任期与当前节点相同， 但是最新索引更小， 拒绝投票
		log.Printf("候选者最新任期与当前节点相同， 但是最新索引更小， 拒绝投票")
		return
	}
follow:
	r.rf.setVote(msg.CandidateId)
	r.rf.setCurrentTerm(msg.Term)
	r.rf.setRole(Follower)
	reply.Term = r.rf.currentTerm
	reply.VoteGranted = true
	//拒绝投票
	log.Printf("Access Vote: \n %+v \n %+v \n %+v", msg, r.rf, reply)
	return
}

// Log leader调用， 用于复制日志条目
func (r *Rpc) Log(l *LogMsg, reply *LogReply) error {
	log.Printf("beat heat from leader: %s", l.LeaderID)
	r.rf.latestHeartBeat = time.Now().UnixNano() //更新心跳时间
	r.rf.setVote("")                             //已经确定了领导人，删除投票
	if l.Term >= r.rf.getTerm() {
		//领导人任期高， 更新任期， 并切换成跟随者
		r.rf.setCurrentTerm(l.Term)
		if r.rf.getRole() == Leader {
			log.Printf("%s change leader to follower", r.rf.host)
		}
		r.rf.setRole(Follower)
	}
	return nil
}

func (rf *raft) handlerHearBeat() {
	for {
		lastTime := rf.latestHeartBeat
		utils.BlockStaticTime(heartBeatRate)
		//一个选举周期内未捕捉到心跳包, 开始新的节点
		if rf.latestHeartBeat <= lastTime && len(rf.cluster) >= rf.conf.MinNodeNum-1 {
			log.Printf("%s start vote: %d", rf.host, len(rf.cluster))
			rf.StartVote()
			break
		}
	}
}
