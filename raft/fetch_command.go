package raft

import (
	log "github.com/sirupsen/logrus"

	"kv-raft/raft/common"
)

func (r *Raft) startFetchCommand() {
	for command := range r.clientCommands {
		r.mutex.Lock()
		if r.status == common.Leader {
			lastLogIndex, _ := r.getLastLogIndexAndTerm()
			r.logs = append(r.logs, common.LogEntry{
				Command: command,
				Term:    r.currentTerm,
				Index:   lastLogIndex + 1,
			})
		} else {
			log.Warnf("not a leader, will be discarded: %v", command)
		}
		r.mutex.Unlock()
	}
}
