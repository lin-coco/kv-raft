package handler

import (
	"encoding/json"
	"kv-raft/server/command"
	"kv-raft/server/components/acltool"
	"kv-raft/server/components/apikey"
	"kv-raft/server/components/httptool"
	"kv-raft/server/global"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func ReceiveHandler(w http.ResponseWriter,r *http.Request) {
	// 日志记录 请求耗时
	start := time.Now()
	defer func() {
		duration := time.Since(start)
        log.Debugf("%s %s %s completed in %v", start.GoString(),r.URL.Path,r.Method, duration)
	}()
	// RequestId 校验
	requestId := r.Header.Get("kv-raft-request-id")
	if len(requestId) != 36 {
		httptool.BadResponse(w)
		return
	}
	// IsLeader 领导者检查
	b,leaderId := global.R.IsLeader()
	if !b {
		httptool.SeeOtherResponse(w,global.Config.Addrs[leaderId].Server)
		return
	}
	// Cmd校验命令正确性
	cmd, err := httptool.StringBody(r)
	if err != nil {
		httptool.ErrorResponse(w, "校验命令正确性失败 " + err.Error())
		return
	}
	normalCmd := command.Normalize(cmd)
	// if strings.HasPrefix("acl-add")
	if strings.HasPrefix(string(normalCmd),"acl-add") || 
		strings.HasPrefix(string(normalCmd),"acl-update") {
			normalCmd = command.Normalize(cmd + " " + apikey.GenerateApiKey(32))
	}
	err = command.Check(normalCmd)
	if err != nil {
		httptool.BadResponse(w)
		return
	}
	exactCmdStr := command.ExactCmdStr(normalCmd)
	// 获取身份信息
	apiKey := r.Header.Get("Authorization")
	if len(apiKey) != 32 || strings.Contains(apiKey, " ") {
		httptool.UnauthorizedResponse(w)
		return
	}
	getAccountExactCmdStr := command.ExactCmdStr(command.Normalize("get system:acl:apikey:" + apiKey))
	getAccountCmd := command.Unmarshal(getAccountExactCmdStr)
	jsonData := getAccountCmd.ExecCMD()
	getResp := command.GetResponse{}
	json.Unmarshal([]byte(jsonData), &getResp)
	account := getResp.Kvs[0].Value
	if account == "" {
		httptool.UnauthorizedResponse(w)
		return
	}
	if exactCmdStr == "acl-get" {
		exactCmdStr = command.ExactCmdStr(command.Normalize("acl-get " + account))
	}
	aclGetExactCmdStr := command.ExactCmdStr(command.Normalize("acl-get " + account))
	aclGetExactCmd := command.Unmarshal(aclGetExactCmdStr)
	jsonData = aclGetExactCmd.ExecCMD()
	aclInfo := command.AclGetResponse{}
	json.Unmarshal([]byte(jsonData), &aclInfo)
	if !aclInfo.Exist {
		httptool.ErrorResponse(w, "存在该账号" + account + ", 但不存在账号acl？" + aclGetExactCmd.ExecCMD())
	}
	// 检查命令权限
	//（acl-del、acl-update、acl-add禁止被非root用户）
	//（acl-add、acl-update需要补充apikey）
	//（api-get非root账号只能访问自己，root账号访问可以有参数）
	exactCmd := command.Unmarshal(exactCmdStr)
	if exactCmd.Name() == "acl-del" ||  exactCmd.Name() == "acl-update" || 
		exactCmd.Name() == "acl-add" || exactCmd.Name() == "acl-all" {
		if account != "root" {
			httptool.ForbiddenResponse(w)
			return
		}
	}
	if account != "root" && exactCmd.Name() == "acl-get" {
		exactCmd = command.Unmarshal(command.ExactCmdStr(command.Normalize("acl-get " + account)))
	}
	// node系统命令直接返回
	if exactCmd.Name() == "node-get" ||  exactCmd.Name() == "node-leader" || exactCmd.Name() == "node-all" {
		httptool.OKResponse(w, []byte(exactCmd.ExecCMD()))
		return
	}
	// 检查身份操作键权限
	if account != "root" {
		for _,key := range exactCmd.GetKeys() {
			permission := acltool.KeyPermission(key,aclInfo.Rules)
			if permission == acltool.Deny {
				httptool.ForbiddenResponse(w)
				return
			}
			if permission == acltool.Read && !exactCmd.ReadOnly() {
				httptool.ForbiddenResponse(w)
				return
			}
		}
	}
	// 只读快速返回
	if exactCmd.ReadOnly() {
		b, leaderId := global.R.IsLeader()
		if !b {
			httptool.SeeOtherResponse(w, global.Config.Addrs[leaderId].Server)
			return
		}
		result := exactCmd.ExecCMD()
		httptool.OKResponse(w, []byte(result))
		return
	}
	// 将命令发送给raft
	if global.Reqs[requestId] != nil {
		httptool.BadResponse(w)
	}
	global.Reqs[requestId] = make(chan string)
	raftCmd := string(exactCmdStr) + " " + requestId
	global.ClientCommands <- raftCmd
	result := <- global.Reqs[requestId]
	httptool.OKResponse(w, []byte(result))
}