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
		if err := server.CheckCommand(command); err != nil {
			log.Warnf("the log parsing in the snapshot failed: %v, and will be skipped", err)
			continue
		}
		// 执行
		_ = server.ExecCommand(command)
	}
	return nil
}
