package impl

import "strings"

type RWJudge struct {
}

func NewRWJudge() RWJudge {
	return RWJudge{}
}

func (r RWJudge) ReadOnly(command string) bool {
	if strings.HasPrefix(command, "get") {
		return true
	}
	return false
}
