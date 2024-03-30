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

func put(key, value string) {
	KV[key] = value
}

func get(key string) string {
	return KV[key]
}

func del(key string) {
	delete(KV, key)
}

// ExecCommand 执行之前必须先check
func ExecCommand(command string) string {
	split := strings.Split(command, " ")
	switch split[0] {
	case "put":
		put(split[1], split[2])
		return ""
	case "get":
		return get(split[1])
	case "del":
		del(split[1])
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
	// 只支持put get del
	switch split[0] {
	case "put":
		if len(split) != 3 {
			return errors.New("number of command parameters is incorrect")
		}
		return nil
	case "get":
		if len(split) != 2 {
			return errors.New("number of command parameters is incorrect")
		}
		return nil
	case "del":
		if len(split) != 2 {
			return errors.New("number of command parameters is incorrect")
		}
		return nil
	default:
		return fmt.Errorf("unsupported command: %s", split[0])
	}
}
