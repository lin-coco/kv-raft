package server

import (
	"errors"
	"fmt"
	"strings"
)

var KV map[string]string

func init() {
	KV = make(map[string]string)
}

func Put(key, value string) {
	KV[key] = value
}

func Get(key string) string {
	return KV[key]
}

func Del(key string) {
	delete(KV, key)
}

// ExecCommand 执行之前必须先check
func ExecCommand(command string) string {
	split := strings.Split(command, " ")
	switch split[0] {
	case "Put":
		Put(split[1], split[2])
		return ""
	case "Get":
		return Get(split[1])
	case "Del":
		Del(split[1])
		return ""
	}
	return ""
}

// CheckCommand 检查
func CheckCommand(command string) error {
	split := strings.Split(command, " ")
	if len(split) < 1 {
		return errors.New("command is incorrect")
	}
	// 只支持put Get Del
	switch split[0] {
	case "Put":
		if len(split) != 3 {
			return errors.New("number of command parameters is incorrect")
		}
		return nil
	case "Get":
		if len(split) != 2 {
			return errors.New("number of command parameters is incorrect")
		}
		return nil
	case "Del":
		if len(split) != 2 {
			return errors.New("number of command parameters is incorrect")
		}
		return nil
	default:
		return fmt.Errorf("unsupported command: %s", split[0])
	}
}
