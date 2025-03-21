package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	kv_raft "kv-raft"
)

func main() {
	app := cli.NewApp()
	app.Name = "mykv-cli"
	app.Usage = "kv storage system based on raft"
	app.Before = func(context *cli.Context) error {
		// Log as JSON instead of the default ASCII formatter.
		log.SetFormatter(&log.JSONFormatter{})
		log.SetOutput(os.Stdout)
		log.SetLevel(log.DebugLevel)
		return nil
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "f",
			Usage: "-f config file",
		},
	}
	app.Action = func(ctx *cli.Context) {
		configFilePath := ctx.String("f")
		// 阅读配置
		var err error
		config, err = kv_raft.ReadConfigFile(configFilePath)
		if err != nil {
			log.Errorf("kv_raft.ReadConfigFile err: %v", err)
			return
		}
		log.Debugf("config loaded: %v.", config)
		if err = Client(config); err != nil {
			log.Errorf("Client err: %v", err)
		}

	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

const (
	Success = 1
	Failed  = 2
	Forward = 3
)

var (
	config *kv_raft.Config
	w      *bufio.Writer
	r      *bufio.Reader
	retry  = 3
)

func Client(config *kv_raft.Config) error {
	if config.ClientRetry != 0 {
		retry = config.ClientRetry
	}
	var lastLeaderId int
	w = bufio.NewWriter(os.Stdout)
	r = bufio.NewReader(os.Stdin)
	for {
		writeAngle()
		userCommand, err := r.ReadBytes('\n') // 包含\n
		if err != nil {
			return fmt.Errorf("r.ReadString err: %v", err)
		}
		userCommand = userCommand[:len(userCommand)-1] // 丢掉最后的\n
		if len(userCommand) == 0 {
			// 空
			continue
		}
		if strCommand := string(userCommand); strCommand == "exit" || strCommand == "quit" || strCommand == "q" {
			fmt.Println("bye")
			break
		}
		requestID := uuid.New().String()
		if b := SendCommand(&lastLeaderId, userCommand, requestID); !b {
			for i := 0; i < retry; i++ {
				b = SendCommand(&lastLeaderId, userCommand, requestID)
				if b {
					break
				}
			}
		}
	}
	return nil
}

// SendCommand 返回结果是否成功发送并响应
func SendCommand(lastLeaderId *int, userCommand []byte, requestID string) bool {
	leaderAddr := config.Addrs[*lastLeaderId].Server
	req, _ := http.NewRequest("POST", "http://"+leaderAddr+"/receive", bytes.NewBuffer(userCommand))
	req.Header.Set("Content-Type", "application/stream+json")
	req.Header.Set("kv-raft-request-id", requestID)
	client := &http.Client{}
	resp, err := client.Do(req)
	//resp, err := http.Post("http://"+leaderAddr+"/receive", "application/stream+json", bytes.NewBuffer(userCommand))
	if err != nil {
		for i := 0; i < len(config.Addrs); i++ {
			var b bool
			resp, b = LookForCan(i, userCommand)
			if b {
				break
			}
		}
	}
	if resp == nil {
		panic("look for can failed")
	}
	var all []byte
	all, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("addr: %v, io.ReadAll err: %v", leaderAddr, err)
		return false
	}
	//log.Debugf("resp.Body: %v,%v", all[0], string(all[1:]))
	if all[0] == Success {
		log.Infof("请求id: %v, size: %v", requestID, len(requestID))
		writeSuccess(all[1:])
		return true
	} else if all[0] == Forward {
		*lastLeaderId = int(all[1])
		log.Debugf("leader节点变化: %v", *lastLeaderId)
		return false
	} else if all[0] == Failed {
		log.Infof("请求id: %v, size: %v", requestID, len(requestID))
		writeFailed(all[1:])
		return true
	} else {
		panic("unexpected response: " + string(all))
	}
}

// LookForCan 寻找能够通信的
func LookForCan(trialLeaderId int, userCommand []byte) (*http.Response, bool) {
	trialLeaderAddr := config.Addrs[trialLeaderId].Server
	resp, err := http.Post("http://"+trialLeaderAddr+"/receive", "application/stream+json", bytes.NewBuffer(userCommand))
	if err != nil {
		return nil, false
	}
	return resp, true
}

func writeAngle() {
	_, _ = w.Write([]byte("> "))
	_ = w.Flush()
}
func writeSuccess(result []byte) {
	result = append(result, '\n')
	_, _ = w.Write(append([]byte("success: \n"), result...))
	_ = w.Flush()
}

func writeFailed(result []byte) {
	result = append(result, '\n')
	_, _ = w.Write(append([]byte("failed: \n"), result...))
	_ = w.Flush()
}
