package impl

import (
	"bufio"
	"os"

	log "github.com/sirupsen/logrus"

	"kv-raft/server/command"
)

type Reset struct {
	file *os.File
}

func NewReset(file *os.File) Reset {
	return Reset{file: file}
}

func (r Reset) Reset() error {
	scanner := bufio.NewScanner(r.file)
	for scanner.Scan() {
		bytes := scanner.Bytes()
		cmd := string(bytes)
		if len(cmd) < 37 {
			log.Warnf("the log parsing in the snapshot failed: %s, and will be skipped", "length < 37")
			continue
		}
		cmd = cmd[:len(cmd)-37]
		kvcmd := command.Unmarshal(command.ExactCmdStr(cmd))
		kvcmd.ExecCMD()
		log.Debugf("执行完成命令: %v", cmd)
	}
	return nil
}
