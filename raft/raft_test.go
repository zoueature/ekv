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
		},
		MinNodeNum: 3,
	}
	_ = NewRaft(cfg)
	cfg = &Conf{
		Host: "127.0.0.1:8082",
		Cluster: []string{
			"127.0.0.1:8081",
			"127.0.0.1:8083",
		},
		MinNodeNum: 3,
	}
	_ = NewRaft(cfg)
	cfg = &Conf{
		Host: "127.0.0.1:8083",
		Cluster: []string{
			"127.0.0.1:8081",
			"127.0.0.1:8082",
		},
		MinNodeNum: 3,
	}
	_ = NewRaft(cfg)
	time.Sleep(1000 * time.Second)
}
