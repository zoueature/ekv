/*+-----------------------------+
 *| Author: Zoueature           |
 *+-----------------------------+
 *| Email: zoueature@gmail.com  |
 *+-----------------------------+
 */
package raft

import "net"

type raftServer struct {
	listener net.Listener
}

type RaftServer interface {
	Start()
	Stop()
	CreateNode()
	DeleteNode()
	SetNode()
	GetNode()
}

func newRaftService(listener net.Listener) RaftServer {
	return &raftServer{listener: listener}
}

// Start 启动服务
func (rfSvc *raftServer) Start() {
	go func() {
		for {
			_, err := rfSvc.listener.Accept()
			if err != nil {
				//todo return error
				continue
			}
			//todo 解析命令， 调用对应的方法

		}
	}()
}

func (rfSvc *raftServer) Stop() {
	defer rfSvc.listener.Close()
}

// CreateNode 创建kv节点
func (rfSvc *raftServer) CreateNode() {

}

// DeleteNode 删除kv节点
func (rfSvc *raftServer) DeleteNode() {

}

// SetNode 设置节点value
func (rfSvc *raftServer) SetNode() {

}

// GetNode 获取节点值
func (rfSvc *raftServer) GetNode() {

}
