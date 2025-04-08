package command

import (
	"encoding/json"
	"errors"
	"kv-raft/server/global"
)

type AclDelCMD struct {
	Account string
}

type AclDelResponse struct {
	Exist bool `json:"exist"`
	DelAcl string `json:"del_acl"`
}

func aclDelCheck(split []string) error {
	if len(split) != 1 {
		return errors.New("command is incorrent")
	}
	return nil
}

func aclDelUnmarshal(split []string) AclDelCMD {
	return AclDelCMD{
		Account: split[0],
	}
}

func (g AclDelCMD) Marshal() error {
	return nil
}
func (g AclDelCMD) GetKeys() []string {
	return []string{"system:acl:" + g.Account}
}
func (g AclDelCMD) ExecCMD() string {
	var resp AclDelResponse
	a,_ := global.StorageEngine.Get("system:acl:" + g.Account)
	if a == "" {
		resp.Exist = false
		r,_ := json.Marshal(resp)
		return string(r)
	}
	var acl AclAddCMD
	json.Unmarshal([]byte(a), &acl)
	global.StorageEngine.Del("system:acl:apikey:" + acl.ApiKey)
	global.StorageEngine.Del("system:acl:" + g.Account)
	resp.DelAcl = acl.Account
	resp.Exist = true
	r,_ := json.Marshal(resp)
	return string(r)
}
func (g AclDelCMD) ReadOnly() bool {
	return false
}

func (g AclDelCMD) Name() string {
	return "acl-del"
}