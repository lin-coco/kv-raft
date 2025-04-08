package command

import (
	"encoding/json"
	"errors"
	"kv-raft/server/global"
)

type GetCMD struct {
	Keys []string
}

type GetResponse struct {
	Kvs []KV `json:"kvs"`
	Num int `json:"num"`
}

type KV struct {
	Key string `json:"key"`
	Value string `json:"value"`
}

func getCheck(split []string) error {
	if len(split) == 0 {
		return errors.New("command is incorrent")
	}
	return nil
}

func getUmarshal(split []string) GetCMD {
	return GetCMD{
		Keys: split,
	}
}

func (g GetCMD) Marshal() error {
	return nil
}
func (g GetCMD) GetKeys() []string {
	return nil
}
func (g GetCMD) ExecCMD() string {
	var resp GetResponse
	resp.Num = len(g.Keys)
	resp.Kvs = make([]KV, 0, len(g.Keys))
	for i := 0;i < len(g.Keys);i++ {
		v,_ := global.StorageEngine.Get(g.Keys[i])
		resp.Kvs = append(resp.Kvs, KV{Key: g.Keys[i],Value: v})
	}
	r,_ := json.Marshal(resp)
	return string(r)
}
func (g GetCMD) ReadOnly() bool {
	return true
}

func (g GetCMD) Name() string {
	return "get"
}