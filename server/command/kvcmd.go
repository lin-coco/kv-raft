package command

import (
	"errors"
	"strings"
)

/*
get key1 key2
put key1 value1
del key1 key2
exist key1
rename key1 key2
export fuzzy1 fuzzy2 key1 
import file
acl-add account rule1 rule2 xxxxx...
acl-get account
acl-update account rule1 rule2
acl-del account
acl-all
node-get
node-leader
node-all
*/

type KVCMD interface {
	// 序列化命令
	Marshal() error
	// 获取该命令所有键
	GetKeys() []string
	// 执行命令
	ExecCMD() string
	// 命令是否只读
	ReadOnly() bool
	// 命令名称
	Name() string
}

// 规范化的命令
type NormalCmdStr string
// 正确的命令
type ExactCmdStr string

func Normalize(command string) NormalCmdStr {
	// 去除前后空格，引号
	command = strings.TrimSpace(command)
	command = strings.ReplaceAll(command, "\"", "")
	command = strings.ReplaceAll(command, "'", "")
	// command转换为切片
	split := strings.Split(command, " ")
	// 切片去除空元素
	result := make([]string, 0)
	for _, str := range split {
		if str != "" {
			result = append(result, str)
		}
	}
	// 重新组合string
	normalCmdStr := NormalCmdStr(strings.Join(result, " "))
	return normalCmdStr
}

func Check(command NormalCmdStr) error {
	split := strings.Split(string(command), " ")
	cmd := split[0]
	split = split[1:]
	switch cmd {
	case "get":
		return getCheck(split)
	case "put":
		return putCheck(split)
	case "del":
		return delCheck(split)
	case "exist":
		return existCheck(split)
	case "rename":
		return renameCheck(split)
	case "keys":
		return keysCheck(split)
	case "acl-add":
		return aclAddCheck(split)
	case "acl-get":
		return aclGetCheck(split)
	case "acl-update":
		return aclUpdateCheck(split)
	case "acl-del":
		return aclDelCheck(split)
	case "acl-all":
		return aclAllCheck(split)
	case "node-get":
		return nodeGetCheck(split)
	case "node-leader":
		return nodeLeaderCheck(split)
	case "node-all":
		return nodeAllCheck(split)
	}
	return errors.New("non-existent command")
}

func Unmarshal(command ExactCmdStr) KVCMD {
	split := strings.Split(string(command), " ")
	cmd := split[0]
	split = split[1:]
	// 解析命令
	switch cmd {
	case "get":
		return getUmarshal(split)
	case "put":
		return putUnmarshal(split)
	case "del":
		return delUnmarshal(split)
	case "exist":
		return existUnmarshal(split)
	case "rename":
		return renameUnmarshal(split)
	case "keys":
		return keysUnmarshal(split)
	case "acl-add":
		return aclAddUnmarshal(split)
	case "acl-get":
		return aclGetUnmarshal(split)
	case "acl-update":
		return aclUpdateUnmarshal(split)
	case "acl-del":
		return aclDelUnmarshal(split)
	case "acl-all":
		return aclAllUnmarshal(split)
	case "node-get":
		return nodeGetUnmarshal(split)
	case "node-leader":
		return nodeLeaderUnmarshal(split)
	case "node-all":
		return nodeAllUnmarshal(split)
	}
	return nil
}