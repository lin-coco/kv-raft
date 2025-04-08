package command

import (
	"encoding/json"
	"errors"
	"kv-raft/server/global"
)

type AclAllCMD struct {
}

type AclAllResponse struct {
	Acls []AclAddCMD `json:"acls"`
	Num int `json:"nums"`
}

func aclAllCheck(split []string) error {
	if len(split) > 0 {
		return errors.New("command is incorrent")
	}
	return nil
}

func aclAllUnmarshal(split []string) AclAllCMD {
	_ = split
	return AclAllCMD{}
}

func (g AclAllCMD) Marshal() error {
	return nil
}
func (g AclAllCMD) GetKeys() []string {
	return nil
}
func (g AclAllCMD) ExecCMD() string {
	var resp AclAllResponse
	keys := global.StorageEngine.Prefix("system:acl:")
	resp.Acls = make([]AclAddCMD, 0, len(keys))
	for i := 0;i < len(keys);i++ {
		a,_ := global.StorageEngine.Get(keys[i])
		var acl AclAddCMD
		json.Unmarshal([]byte(a), &acl)
		if acl.Account == "" {
			continue
		}
		resp.Acls = append(resp.Acls, acl)
	}
	resp.Num = len(resp.Acls)
	r,_ := json.Marshal(resp)
	return string(r)
}
func (g AclAllCMD) ReadOnly() bool {
	return true
}

func (g AclAllCMD) Name() string {
	return "acl-all"
}