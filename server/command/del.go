package command

import (
	"encoding/json"
	"errors"
	"kv-raft/server/global"
)

type DelCMD struct {
	Keys []string
}

type DelResponse struct {
	DelKeys []string `json:"del_keys"`
	Num int `json:"num"`
}

func delCheck(split []string) error {
	if len(split) == 0 {
		return errors.New("command is incorrent")
	}
	return nil
}

func delUnmarshal(split []string) DelCMD {
	return DelCMD{
		Keys: split,
	}
}

func (g DelCMD) Marshal() error {
	return nil
}
func (g DelCMD) GetKeys() []string {
	return g.Keys
}
func (g DelCMD) ExecCMD() string {
	var resp DelResponse
	resp.DelKeys = make([]string, 0)
	for i := 0;i < len(g.Keys);i++ {
		v,_ := global.StorageEngine.Get(g.Keys[i])
		if v == "" {
			// 不存在
			continue
		}
		global.StorageEngine.Del(g.Keys[i])
		resp.DelKeys = append(resp.DelKeys, g.Keys[i])
	}
	resp.Num = len(resp.DelKeys)
	r,_ := json.Marshal(resp)
	return string(r)
}
func (g DelCMD) ReadOnly() bool {
	return false
}
func (g DelCMD) Name() string {
	return "del"
}