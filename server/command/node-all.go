package command

import (
	"encoding/json"
	"errors"
	"kv-raft/server/global"
)

type NodeAllCMD struct {
}

func nodeAllCheck(split []string) error {
	if len(split) != 0 {
		return errors.New("command is incorrent")
	}
	return nil
}

func nodeAllUnmarshal(split []string) NodeAllCMD {
	_ = split
	return NodeAllCMD{}
}

func (g NodeAllCMD) Marshal() error {
	return nil
}
func (g NodeAllCMD) GetKeys() []string {
	return nil
}
func (g NodeAllCMD) ExecCMD() string {
	nodeInfos := global.R.GetNodeInfos()
	v,_ := json.Marshal(nodeInfos)
	return string(v)
}
func (g NodeAllCMD) ReadOnly() bool {
	return true
}

func (g NodeAllCMD) Name() string {
	return "node-all"
}