package impl

import (
	log "github.com/sirupsen/logrus"

	"kv-raft/server"
)

type ApplyCommand struct {
}

func NewApplyCommand() ApplyCommand {
	return ApplyCommand{}
}

func (a ApplyCommand) ApplyCommand(command string) string {
	reqID := command[len(command)-36:]
	command = command[:len(command)-37]
	result := server.ExecCommand(command)
	// 检查是不是leader节点接收了客户端的请求，是的话通知状态机返回响应
	if server.Reqs[reqID] != nil {
		server.Reqs[reqID] <- result
	}
	log.Infof("server listen apply command: %s, req id: %s, result: %s", command, reqID, result)
	return result
}
