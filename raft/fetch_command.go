package raft

import (
	log "github.com/sirupsen/logrus"

	"kv-raft/raft/common"
)

func (r *Raft) startFetchCommand() {
	for command := range r.clientCommands {
		r.mutex.Lock()
		if r.status == common.Leader {
			log.Debugf("收到状态机命令: %v, 是leader 将加入日志", command)
			lastLogIndex, _ := r.getLastLogIndexAndTerm()
			r.logs = append(r.logs, common.LogEntry{
				Command: command,
				Term:    r.currentTerm,
				Index:   lastLogIndex + 1,
			})
		} else {
			log.Debugf("收到状态机命令: %v, 不是leader 将丢弃", command)
		}
		r.mutex.Unlock()
	}
}
