package impl

import (
	log "github.com/sirupsen/logrus"

	"kv-raft/server"
)

type ApplyCommand struct {
}

func NewApplyCommand() ApplyCommand {
	return NewApplyCommand()
}

func (a ApplyCommand) ApplyCommand(command string) string {
	reqID := command[len(command)-36:]
	command = command[:len(command)-37]
	result := server.ExecCommand(command)
	server.Reqs[reqID] <- result // 允许返回
	log.Infof("server listen apply command: %s, req id: %s, result: %s", command, reqID, result)
	return result
}
