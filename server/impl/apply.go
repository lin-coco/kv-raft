package impl

import (
	log "github.com/sirupsen/logrus"

	"kv-raft/server/command"
	"kv-raft/server/global"
)

type ApplyCommand struct {
}

func NewApplyCommand() ApplyCommand {
	return ApplyCommand{}
}

func (a ApplyCommand) ApplyCommand(strCommand string) string {
	reqID := strCommand[len(strCommand)-36:]
	strCommand = strCommand[:len(strCommand)-37]
	exactCmdStr := command.ExactCmdStr(strCommand)

	cmd := command.Unmarshal(exactCmdStr)
	result := cmd.ExecCMD()
	log.Infof("应用日志: %s", exactCmdStr)
	// result := command.ExecCommand(strCommand)
	// 检查是不是leader节点接收了客户端的请求，是的话通知状态机返回响应
	if global.Reqs[reqID] != nil {
		global.Reqs[reqID] <- result
	}
	log.Infof("server listen apply command: %s, req id: %s, result: %s", exactCmdStr, reqID, result)
	return result
}
