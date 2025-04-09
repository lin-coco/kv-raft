package command

import (
	"encoding/json"
	"errors"
	"kv-raft/server/global"
)

type ExistCMD struct {
	Key string
}

type ExistResponse struct {
	Exist bool `json:"exist"`
}

func existCheck(split []string) error {
	if len(split) != 1 {
		return errors.New("command is incorrent")
	}
	return nil
}

func existUnmarshal(split []string) ExistCMD {
	return ExistCMD{
		Key: split[0],
	}
}

func (g ExistCMD) Marshal() error {
	return nil
}
func (g ExistCMD) GetKeys() []string {
	return []string{g.Key}
}
func (g ExistCMD) ExecCMD() string {
	var resp ExistResponse
	v,_ := global.StorageEngine.Get(g.Key)
	if v == "" {
		resp.Exist = false
		r,_ := json.Marshal(resp)
		return string(r)
	}
	resp.Exist = true
	r,_ := json.Marshal(resp)
	return string(r)
}
func (g ExistCMD) ReadOnly() bool {
	return true
}

func (g ExistCMD) Name() string {
	return "exist"
}