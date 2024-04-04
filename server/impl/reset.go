package impl

import (
	"bufio"
	"os"

	log "github.com/sirupsen/logrus"

	"kv-raft/server"
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
		command := string(bytes)
		if len(command) < 37 {
			log.Warnf("the log parsing in the snapshot failed: %s, and will be skipped", "length < 37")
			continue
		}
		command = command[:len(command)-37]
		if err := server.CheckCommand(command); err != nil {
			log.Warnf("the log parsing in the snapshot failed: %v, and will be skipped", err)
			continue
		}
		// 执行
		_ = server.ExecCommand(command)
		log.Debugf("执行完成命令: %v", command)
	}
	return nil
}
