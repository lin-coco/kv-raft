package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
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
			Name: "s",
			Usage: "-s list of servers e.g. 127.0.0.1:8000,127.0.0.1:8001,127.0.0.1:8002",
			Required: true,
		},
		cli.StringFlag{
			Name: "k",
			Usage: "-k api key e.g. mh3UtppYmnKA7v0Xmnj1zxAfBZvfP6dw",
			Required: true,
		},
	}
	app.Action = func(ctx *cli.Context) {
		s := ctx.String("s")
		servers := strings.Split(s, ",")
		if len(servers) == 0 {
			log.Errorf("Client err: -s至少需要一个节点")
		}
		log.Info("servers配置: ", strings.Join(servers," "))
		apiKey := ctx.String("k")
		if apiKey == "" {
			log.Errorf("Client err: -k需要apikey")
		}
		log.Info("api key配置: ", apiKey)
		if err := Client(servers, apiKey); err != nil {
			log.Errorf("Client err: %v", err)
		}
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

var (
	w      *bufio.Writer
	r      *bufio.Reader
	retry  = 3
)

func Client(servers []string, apiKey string) error {
	lastLeaderAddr := servers[0]
	w = bufio.NewWriter(os.Stdout)
	r = bufio.NewReader(os.Stdin)
	for {
		writeAngle()
		
		userCommand, err := r.ReadString('\n') // 包含\n
		if err != nil {
			return fmt.Errorf("r.ReadString err: %v", err)
		}
		userCommand = strings.TrimSpace(userCommand)
		userCommand = strings.TrimSpace(userCommand) // 删除前后的空白字符
		if len(userCommand) == 0 {
			continue
		}
		if strCommand := string(userCommand); strCommand == "exit" || strCommand == "quit" || strCommand == "q" {
			fmt.Println("bye")
			break
		}
		requestID := uuid.New().String()
		if b := SendCommand(servers,&lastLeaderAddr, userCommand, requestID, apiKey); !b {
			for i := 0; i < retry; i++ {
				b = SendCommand(servers,&lastLeaderAddr, userCommand, requestID, apiKey)
				if b {
					break
				}
			}
		}
	}
	return nil
}

// SendCommand 返回结果是否成功发送并响应
func SendCommand(servers []string,leaderAddr *string, userCommand string, requestID string, apiKey string) bool {
	path := "receive"
	if strings.HasPrefix(userCommand, "debug-allkey") {
		path = "debug-allkey"
	} else if strings.HasPrefix(userCommand, "debug") {
		path = "debug"
		userCommand  = userCommand[6:]
	}
	req, _ := http.NewRequest("POST", "http://"+*leaderAddr+"/" + path, bytes.NewBuffer([]byte(userCommand)))
	req.Header.Set("Content-Type", "application/stream+json")
	req.Header.Set("Authorization", apiKey)
	req.Header.Set("kv-raft-request-id", requestID)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		for i := 0; i < len(servers); i++ {
			var b bool
			resp, b = LookForCan(servers[i], userCommand)
			if b {
				break
			}
		}
	}
	if resp == nil {
		panic("look for can failed")
	}
	if resp.StatusCode == http.StatusSeeOther {
		// 转发
		*leaderAddr = resp.Header.Get("leader-id")
		// log.Debugf("leader节点变化: %s", leaderAddr)
		return false
	} else if resp.StatusCode == http.StatusBadRequest {
		// 命令错误
		// log.Infof("请求id: %v, size: %v", requestID, len(requestID))
		writeFailed("命令错误", http.StatusBadRequest)
		return true
	} else if resp.StatusCode == http.StatusForbidden {
		// 禁止
		writeFailed("命令或键无权限访问", http.StatusForbidden)
		return true
	} else if resp.StatusCode == http.StatusUnauthorized {
		// 认证失败
		writeFailed("认证失败", http.StatusUnauthorized)
		return true
	} else if resp.StatusCode == http.StatusInternalServerError {
		// 服务端错误
		writeFailed("服务端错误，请调试服务端: " + resp.Header.Get("message"), http.StatusInternalServerError)
		return true
	} else if resp.StatusCode == http.StatusOK {
		// 成功响应
		body, _ := io.ReadAll(resp.Body)
		writeSuccess(strings.TrimSpace(string(body)))
		return true
	} else {
		// 未知错误
		writeFailed("未知错误", resp.StatusCode)
		return true
	}
}

// LookForCan 寻找能够通信的节点
func LookForCan(trialLeaderAddr string, userCommand string) (*http.Response, bool) {
	resp, err := http.Post("http://"+trialLeaderAddr+"/receive", "application/stream+json", bytes.NewBuffer([]byte(userCommand)))
	if err != nil {
		fmt.Println(trialLeaderAddr,err)
		return nil, false
	}
	return resp, true
}

func writeAngle() {
	w.WriteString(">>")
	w.Flush()
}
func writeSuccess(body string) {
	// 将紧凑JSON格式转换为带缩进的格式
	var buf bytes.Buffer
	json.Indent(&buf,[]byte(body),"", "  ")
	body = "OK!\n" + buf.String() + "\n"
	w.WriteString(body)
	w.Flush()
}

type Failed struct {
	ErrMsg string `json:"err_msg"`
	Code int `json:"code"`
}

func writeFailed(errMsg string,code int) {
	failed := Failed{
		ErrMsg: errMsg,
		Code: code,
	}
	jsonData, _ := json.MarshalIndent(failed, "", "    ") // 缩进为4个空格
	message := "Failed!\n" + string(jsonData) + "\n"
	w.WriteString(message)
	w.Flush()
}