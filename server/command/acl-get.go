package command

import (
	"encoding/json"
	"errors"
	"kv-raft/server/components/acltool"
	"kv-raft/server/global"
)

type AclGetCMD struct {
	Account string
}

type AclGetResponse struct {
	Exist bool `json:"exist"`
	Account string            `json:"account"`
	Rules   []acltool.ACLRule `json:"rules"`
	ApiKey  string            `json:"apikey"`
}

func aclGetCheck(split []string) error {
	if len(split) > 1 {
		return errors.New("command is incorrent")
	}
	return nil
}

func aclGetUnmarshal(split []string) AclGetCMD {
	if split[0] != "" {
		return AclGetCMD{
			Account: split[0],
		}
	}
	return AclGetCMD{}
}

func (g AclGetCMD) Marshal() error {
	return nil
}
func (g AclGetCMD) GetKeys() []string {
	return []string{"system:acl:" + g.Account}
}
func (g AclGetCMD) ExecCMD() string {
	var resp AclGetResponse
	a, _ := global.StorageEngine.Get("system:acl:" + g.Account)
	if a == "" {
		resp.Exist = false
		resp.Account = g.Account
		r,_ := json.Marshal(resp)
		return string(r)
	}
	var acl AclAddCMD
	json.Unmarshal([]byte(a), &acl)
	resp.Exist = true
	resp.Account = acl.Account
	resp.ApiKey = acl.ApiKey
	resp.Rules = acl.Rules
	r,_ := json.Marshal(resp)
	return string(r)
}
func (g AclGetCMD) ReadOnly() bool {
	return true
}

func (g AclGetCMD) Name() string {
	return "acl-get"
}