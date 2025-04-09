package command

import (
	"encoding/json"
	"errors"
	"kv-raft/server/global"
)

type PutCMD struct {
	Keys []string
	Values []string
}

type PutResponse struct {
	Kvs []KV `json:"kvs"`
	Num int `json:"num"`
}

func putCheck(split []string) error {
	size := len(split)
	if size % 2 == 1 {
		return errors.New("command is incorrent")
	}
	return nil
}

func putUnmarshal(split []string) PutCMD {
	size := len(split)
	keys := make([]string, 0, size)
	values := make([]string, 0, size)
	for i := 0;i < len(split);i++ {
		if i%2 == 0 {
			keys = append(keys, split[i])
		} else {
			values = append(values, split[i])
		}
	}
	return PutCMD{
		Keys: keys,
		Values: values,
	}
}

func (g PutCMD) Marshal() error {
	return nil
}
func (g PutCMD) GetKeys() []string {
	return g.Keys
}
func (g PutCMD) ExecCMD() string {
	var resp PutResponse
	resp.Kvs = make([]KV, 0, len(g.Keys))
	for i := 0;i < len(g.Keys);i++ {
		global.StorageEngine.Put(g.Keys[i], g.Values[i])
		resp.Kvs = append(resp.Kvs, KV{g.Keys[i], g.Values[i]})
	}
	r,_ := json.Marshal(resp)
	return string(r)
}
func (g PutCMD) ReadOnly() bool {
	return false
}

func (g PutCMD) Name() string {
	return "put"
}