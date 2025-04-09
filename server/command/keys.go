package command

import (
	"encoding/json"
	"errors"
	"kv-raft/server/global"
)

type KeysCMD struct {
}

type KeysResponse struct {
	Kvs []KV `json:"kvs"`
	Num int `json:"num"`
}

func keysCheck(split []string) error {
	if len(split) != 0 {
		return errors.New("command is incorrent")
	}
	return nil
}

func keysUnmarshal(split []string) KeysCMD {
	_ = split
	return KeysCMD{}
}

func (g KeysCMD) Marshal() error {
	return nil
}
func (g KeysCMD) GetKeys() []string {
	return nil
}
func (g KeysCMD) ExecCMD() string {
	//  目前只支持root账号
	var resp KeysResponse
	keys := global.StorageEngine.Prefix("")
	resp.Num = len(keys)
	resp.Kvs = make([]KV, 0, len(keys))
	for i := 0;i < len(keys);i++ {
		v,_ := global.StorageEngine.Get(keys[i])
		resp.Kvs = append(resp.Kvs, KV{Key: keys[i],Value: v})
	}
	r,_ := json.Marshal(resp)
	return string(r)
}
func (g KeysCMD) ReadOnly() bool {
	return true
}

func (g KeysCMD) Name() string {
	return "get"
}