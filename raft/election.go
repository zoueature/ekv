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
	"sync"
	"time"
)

const (
	heartBeatRate = 3 * time.Second
)

// StartVote 开始新一轮的领导人选举
func (rf *raft) StartVote() {
	rf.setCurrentTerm(rf.currentTerm + 1)
	rf.setRole(Candidate)
	rf.setVote(rf.id)
	voteResult := make(map[string]bool)
	nodeNum := len(rf.cluster)
	wg := sync.WaitGroup{}
	for _, v := range rf.cluster {
		//并行发送投票选举请求
		wg.Add(1)
		request := &VoteMsg{
			Term:         rf.currentTerm,
			CandidateId:  rf.id,
			LastLogIndex: rf.logLastIndex,
			LastLogTerm:  rf.logLastItem,
		}
		go func(node *Cluster) {
			defer wg.Done()
			reply := new(VoteReplyMsg)
			err := node.cli.Call("Rpc.Vote", request, reply)
			if err != nil {
				log.Printf("request host: %s to vote error: %s", node.host, err.Error())
				return
			}
			log.Printf("%s => %+v", rf.host, reply)
			if reply.Term > rf.getTerm() {
				rf.setCurrentTerm(reply.Term)
			}
			if reply.VoteGranted {
				voteResult[node.host] = true
			}
		}(v)
	}
	wg.Wait()
	voteNum := len(voteResult)
	if voteNum >= nodeNum/2 {
		//获得大多数选票, 选举成功， 成为领导人
		rf.ChangeToLeader()
		return
	}
	//起定时器， 超时未获得其他领导人的心跳包， 则重新起一轮领导人选举
	utils.BlockRandTime(1500*time.Millisecond, 3000*time.Millisecond)
	if rf.getRole() == Follower {
		return
	}
	rf.StartVote()
}

func (rf *raft) ChangeToLeader() {
	rf.setRole(Leader)
	go rf.HeartBeat()
}

// HeartBeat 发送心跳包
func (rf *raft) HeartBeat() {
	for rf.getRole() == Leader {
		wg := sync.WaitGroup{}
		for _, v := range rf.cluster {
			wg.Add(1)
			//给集群中的其他节点发送心跳包， 展示领导人权威
			request := &LogMsg{
				Term:         rf.getTerm(),
				LeaderID:     rf.id,
				PrevLogIndex: 0,
				PrevLogTerm:  0,
				Entries:      nil,
				LeaderCommit: 0,
			}
			go func(node *Cluster) {
				defer wg.Done()
				reply := new(LogReply)
				err := node.cli.Call("Rpc.Log", request, reply)
				if err != nil {
					log.Printf("request host %s heart beat error: %s", node.host, err.Error())
					return
				}
				//todo 处理返回值
			}(v)
		}
		wg.Wait()
		utils.BlockStaticTime(heartBeatRate)
	}
}
