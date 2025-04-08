package command

import (
	"encoding/json"
	"errors"
	"kv-raft/server/global"
)

type NodeGetCMD struct {
}

func nodeGetCheck(split []string) error {
	if len(split) != 0 {
		return errors.New("command is incorrent")
	}
	return nil
}

func nodeGetUnmarshal(split []string) NodeGetCMD {
	_ = split
	return NodeGetCMD{}
}

func (g NodeGetCMD) Marshal() error {
	return nil
}
func (g NodeGetCMD) GetKeys() []string {
	return nil
}
func (g NodeGetCMD) ExecCMD() string {
	nodeInfos := global.R.GetNodeInfo()
	v,_ := json.Marshal(nodeInfos)
	return string(v)
}
func (g NodeGetCMD) ReadOnly() bool {
	return false
}

func (g NodeGetCMD) Name() string {
	return "node-get"
}