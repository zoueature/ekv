/*+-----------------------------+
 *| Author: Zoueature           |
 *+-----------------------------+
 *| Email: zoueature@gmail.com  |
 *+-----------------------------+
 */
package raft

type LogMsg struct {
	Term         int64
	LeaderID     string
	PrevLogIndex int64
	PrevLogTerm  int64
	Entries      []string
	LeaderCommit int64
}

type LogReply struct {
	Term    int64
	Success bool
}

type VoteMsg struct {
	Term         int64  //候选人任期
	CandidateId  string //候选人id
	LastLogIndex int64    //候选人最后日志条目索引值
	LastLogTerm  int64  //候选人最后日志条目的任期
}

type VoteReplyMsg struct {
	Term        int64 //当前任期票
	VoteGranted bool  //是否获得选票
}
