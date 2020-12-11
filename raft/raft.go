/*+-----------------------------+
 *| Author: Zoueature           |
 *+-----------------------------+
 *| Email: zoueature@gmail.com  |
 *+-----------------------------+
 */
package raft

import (
	"github.com/google/uuid"
	"log"
	"net/rpc"
	"sync"
	"time"
)

type logEntity struct {
	term     int64
	command  *Command
	commited bool
}

// persistenceStatus 持久性状态
type persistenceStatus struct {
	currentTerm int64        //当前任期
	votedFor    string       //选举投给的候选者id
	log         []*logEntity //日志条目
}

// volatileStatus 易失性状态
type volatileStatus struct {
	commitIndex int //已提交的最高的日志条目的索引
	lastApplied int //已经被应用到状态机的最高的日志条目的索引
}

// leaderStatus leader服务器状态， 只有leader服务器节点需要维护这个状态， 也是易失性的
type leaderStatus struct {
	nextIndex  map[int64]int64 //发送到每台服务器的下一个日志条目的索引（初始值为leader最后日志索引+1）
	matchIndex map[int64]int64 //已经复制到每台服务器的最高日志索引（初始值为0，单调递增）
}

type Role int

const (
	Follower  Role = iota + 1 //跟随者
	Leader                    //领导者
	Candidate                 //候选者

	voteTimeoutMinTime = 150 * time.Millisecond //选举超时最小时间
	voteTimeoutMaxTime = 300 * time.Millisecond //选举超时最大时间
)

type Cluster struct {
	host    string
	cli     *rpc.Client
	healthy bool
}

type raft struct {
	persistenceStatus
	volatileStatus
	id              string              //节点id
	leadStatus      *leaderStatus       //领导人状态
	role            Role                //角色
	roleLock        sync.RWMutex        //角色读写锁
	voteTimeout     time.Duration       //选举超时时间
	rpc             *rpc.Server         //rpc服务
	cluster         map[string]*Cluster //与其他节点的连接
	statMachine     *statMachine        //状态机
	dialChan        chan<- string       //重试连接节点通道
	latestHeartBeat int64               //上次心跳包的时间
}

type Raft interface {
	// 开始新一轮选举， 未选举成功会一直阻塞
	StartVote()

	Commit(startIndex int, length ...int) error
}

// initialPersistenceStatus 获取存储在文件系统中的节点状态
func (rf *raft) initialPersistenceStatus() {

}

type Conf struct {
	Host    string   `yaml:"host"`
	Cluster []string `yaml:"cluster"`
}

func NewRaft(cfg *Conf) Raft {
	id := uuid.New().String()
	rf := &raft{
		id:       id,
		role:     Follower,
		roleLock: sync.RWMutex{},
		volatileStatus: volatileStatus{
			commitIndex: 0,
			lastApplied: 0,
		},
		statMachine: NewStatMachine(),
		cluster:     make(map[string]*Cluster),
	}
	err := InitialRpcServer(cfg.Host, &Rpc{rf: rf})
	if err != nil {
		log.Fatalf("start %s server error: %s", cfg.Host, err.Error())
		return nil
	}
	rf.startDialNodes(cfg.Cluster)
	rf.initialPersistenceStatus()
	//rf.handlerHearBeat()
	return rf
}

func (rf *raft) changeRole(r Role) {
	rf.roleLock.Lock()
	rf.role = r
	rf.roleLock.Unlock()
}

// Commit 提交日志条目
func (rf *raft) Commit(startIndex int, length ...int) error {
	l := 0
	if len(length) > 0 && length[0] > 0 {
		l = length[0]
	}
	//todo 提交日志条目
	for _, logEntry := range rf.log[startIndex-1 : startIndex+l] {
		_ = rf.statMachine.Commit(logEntry.command)
		logEntry.commited = true
	}
	return nil
}
