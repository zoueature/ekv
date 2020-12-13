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
			//"127.0.0.1:8085",
			//"127.0.0.1:8086",
			//"127.0.0.1:8087",
			//"127.0.0.1:8088",
		},
		MinNodeNum: 3,
	}
	_ = NewRaft(cfg)
	cfg.Host = "127.0.0.1:8082"
	_ = NewRaft(cfg)
	cfg.Host = "127.0.0.1:8083"
	_ = NewRaft(cfg)
	cfg.Host = "127.0.0.1:8084"
	_ = NewRaft(cfg)
	//cfg.Host = "127.0.0.1:8085"
	//_ = NewRaft(cfg)
	//cfg.Host = "127.0.0.1:8086"
	//_ = NewRaft(cfg)
	//cfg.Host = "127.0.0.1:8087"
	//_ = NewRaft(cfg)
	//cfg.Host = "127.0.0.1:8088"
	//_ = NewRaft(cfg)
	time.Sleep(1000 * time.Second)
}
