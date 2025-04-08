package command

import (
	"encoding/json"
	"errors"
	"kv-raft/server/components/acltool"
	"kv-raft/server/global"
	"strings"
)

type AclUpdateCMD struct {
	Account string
	Rules []acltool.ACLRule
	ApiKey  string
}

type AclUpdateResponse struct {
	Exist bool `json:"exist"`
	ApiKey string `json:"apikey"`
}

func aclUpdateCheck(split []string) error {
	if len(split) <= 2 {
		return errors.New("command is incorrent")
	}
	rules := split[1 : len(split)-1]
	for i := 0; i < len(rules); i++ {
		// asd*,rw
		split := strings.Split(rules[i], ",")
		if len(split) != 2 {
			return errors.New("command is incorrent")
		}
		permission := acltool.Permission(strings.TrimSpace(split[1]))
		if !acltool.HasPermission(permission) {
			return errors.New("command is incorrent")
		}
	}
	return nil
}

func aclUpdateUnmarshal(split []string) AclUpdateCMD {
	account := split[0]
	apiKey := split[len(split)-1]
	rules := split[1:len(split)-1]
	aclRules := make([]acltool.ACLRule, 0, len(rules))
	for i := 0; i < len(rules); i++ {
		// asd*,rw
		split := strings.Split(rules[i], ",")
		pattern := strings.TrimSpace(split[0])
		permission := acltool.Permission(strings.TrimSpace(split[1]))
		aclRules = append(aclRules, acltool.ACLRule{
			Pattern:    pattern,
			Permission: permission,
		})
	}
	return AclUpdateCMD{
		Account: account,
		Rules:   aclRules,
		ApiKey:  apiKey,
	}
}

func (g AclUpdateCMD) Marshal() error {
	return nil
}
func (g AclUpdateCMD) GetKeys() []string {
	return []string{"system:acl:" + g.Account}
}
func (g AclUpdateCMD) ExecCMD() string {
	var resp AclUpdateResponse
	// 删除旧apikey
	a,_ := global.StorageEngine.Get("system:acl:" + g.Account)
	if a == "" {
		resp.Exist = false
		r,_ := json.Marshal(resp)
		return string(r)
	}
	var oldAcl AclAddCMD
	json.Unmarshal([]byte(a), &oldAcl)
	global.StorageEngine.Del("system:acl:apikey:"+oldAcl.ApiKey)
	// 覆盖acl，新增apikey
	rules, _ := json.Marshal(g)
	global.StorageEngine.Put("system:acl:"+g.Account, string(rules))
	global.StorageEngine.Put("system:acl:apikey:" + g.ApiKey, g.Account)
	resp.Exist = true
	resp.ApiKey = g.ApiKey
	r,_ := json.Marshal(resp)
	return string(r)
}
func (g AclUpdateCMD) ReadOnly() bool {
	return false
}
func (g AclUpdateCMD) Name() string {
	return "acl-update"
}