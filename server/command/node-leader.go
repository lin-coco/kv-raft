package command

import (
	"encoding/json"
	"errors"
	"kv-raft/server/global"
)

type NodeLeaderCMD struct {
}

func nodeLeaderCheck(split []string) error {
	if len(split) != 0 {
		return errors.New("command is incorrent")
	}
	return nil
}

func nodeLeaderUnmarshal(split []string) NodeLeaderCMD {
	_ = split
	return NodeLeaderCMD{}
}

func (g NodeLeaderCMD) Marshal() error {
	return nil
}
func (g NodeLeaderCMD) GetKeys() []string {
	return nil
}
func (g NodeLeaderCMD) ExecCMD() string {
	nodeLeaderInfo := global.R.GetNodeLeaderInfo()
	v,_ := json.Marshal(nodeLeaderInfo)
	return string(v)
}
func (g NodeLeaderCMD) ReadOnly() bool {
	return false
}

func (g NodeLeaderCMD) Name() string {
	return "node-leader"
}