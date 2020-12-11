/*+-----------------------------+
 *| Author: Zoueature           |
 *+-----------------------------+
 *| Email: zoueature@gmail.com  |
 *+-----------------------------+
 */
package raft

import "errors"

type operate int

const (
	Set  operate = iota + 1 //覆盖
	Plus                    //加
	Sub                     //减
	Mul                     //乘
	Div                     //除
)

var (
	ErrUnknowOperate = errors.New("unknow operate")
)

type OpFunc func(k, v string) error

// Command 状态机执行指令
type Command struct {
	Key   string
	Value string
	Op    operate
}

type Storage interface {
	Set(k, v string) error
	Plus(k, v string) error
	Sub(k, v string) error
	Mul(k, v string) error
	Div(k, v string) error
}

// statMachine 状态机
type statMachine struct {
	memLimit int64   //内存限制
	storage  Storage //存储引擎
}

func NewStatMachine() *statMachine {
	return &statMachine{
		memLimit: 0,
		storage:  nil,
	}
}
func (sm *statMachine) getOpFun(com *Command) OpFunc {
	switch com.Op {
	case Set:
		return sm.storage.Set
	case Plus:
		return sm.storage.Plus
	case Sub:
		return sm.storage.Sub
	case Mul:
		return sm.storage.Mul
	case Div:
		return sm.storage.Div
	}
	return nil
}

// Commit 提交状态机指令
func (sm *statMachine) Commit(com *Command) error {
	op := sm.getOpFun(com)
	if op == nil {
		return ErrUnknowOperate
	}
	return op(com.Key, com.Value)
}
