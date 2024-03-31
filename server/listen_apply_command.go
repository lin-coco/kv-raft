package server

import (
	log "github.com/sirupsen/logrus"

	"kv-raft/raft"
)

// Reqs key: requestID value: result
var Reqs map[string]chan string

func StartListenApplyCommand(r *raft.Raft) {
	for applyCommand := range r.ApplyCommands {
		reqID := applyCommand[len(applyCommand)-36:]
		applyCommand = applyCommand[:len(applyCommand)-37]
		result := ExecCommand(applyCommand)
		Reqs[reqID] <- result // 允许返回
		log.Infof("server listen apply command: %s, req id: %s, result: %s", applyCommand, reqID, result)
	}
}
