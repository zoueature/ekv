/*+-----------------------------+
 *| Author: Zoueature           |
 *+-----------------------------+
 *| Email: zoueature@gmail.com  |
 *+-----------------------------+
 */
package raft

import (
	"testing"
	"time"
)

func TestNewRaft(t *testing.T) {
	cfg := &Conf{
		Host: "127.0.0.1:8081",
		Cluster: []string{
			"127.0.0.1:8082",
			"127.0.0.1:8083",
			"127.0.0.1:8084",
			"127.0.0.1:8085",
		},
	}
	raft := NewRaft(cfg)
	raft.StartVote()
	cfg = &Conf{
		Host: "127.0.0.1:8082",
		Cluster: []string{
			"127.0.0.1:8081",
			"127.0.0.1:8083",
			"127.0.0.1:8084",
			"127.0.0.1:8085",
		},
	}
	raft1 := NewRaft(cfg)
	raft1.StartVote()
	cfg = &Conf{
		Host: "127.0.0.1:8083",
		Cluster: []string{
			"127.0.0.1:8081",
			"127.0.0.1:8082",
			"127.0.0.1:8084",
			"127.0.0.1:8085",
		},
	}
	raft2 := NewRaft(cfg)
	raft2.StartVote()
	time.Sleep(1000 * time.Second)
}
