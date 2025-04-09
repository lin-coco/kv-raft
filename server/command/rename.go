package command

import (
	"encoding/json"
	"errors"
	"kv-raft/server/global"
)

type RenameCMD struct {
	Key string
	NewName string
}

type RenameResponse struct {
	Exist bool `json:"exist"`
	DelKey string `json:"del_key"`
	NewKey string `json:"new_key"`
}

func renameCheck(split []string) error {
	if len(split) != 2 {
		return errors.New("command is incorrent")
	}
	return nil
}

func renameUnmarshal(split []string) RenameCMD {
	return RenameCMD{
		Key: split[0],
		NewName: split[1],
	}
}

func (g RenameCMD) Marshal() error {
	return nil
}
func (g RenameCMD) GetKeys() []string {
	return []string{g.Key,g.NewName}
}
func (g RenameCMD) ExecCMD() string {
	var resp RenameResponse
	v,_ := global.StorageEngine.Get(g.Key)
	if v == "" {
		resp.Exist = false
		r,_ := json.Marshal(resp)
		return string(r)
	}
	global.StorageEngine.Del(g.Key)
	global.StorageEngine.Put(g.NewName, v)
	resp.Exist = true
	resp.DelKey = g.Key
	resp.NewKey = g.NewName
	r,_ := json.Marshal(resp)
	return string(r)
}
func (g RenameCMD) ReadOnly() bool {
	return false
}

func (g RenameCMD) Name() string {
	return "rename"
}