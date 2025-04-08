package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"kv-raft/client"
	"os"
	"path/filepath"
	"strings"

	"github.com/peterh/liner"
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
		apiKey := ctx.String("k")
		client, err := client.NewClient(servers, apiKey)
		if err != nil {
			panic(err)
		}
		fmt.Println("servers配置: ", strings.Join(client.Servers," "))
		fmt.Println("api key配置: ", client.ApiKey)
		RunCmd(client)
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

var (
	historyFile = filepath.Join(os.TempDir(), ".liner_example_history")
	commands = []string{"get","del","put","exist","rename",
						"acl-get","acl-add","acl-update","acl-del","acl-all",
						"node-get","node-leader","node-all",
						"debug","debug-allkey"}
)

func RunCmd(client client.Client) {
	line := liner.NewLiner()
	defer line.Close()
	line.SetCtrlCAborts(true)
	// 设置自动补全
	line.SetCompleter(func(line string) (c []string) {
		for _, cmd := range commands {
			if cmd == line {
				return []string{cmd}
			}
		}
		return
	})
	// 加载历史记录
	if f, err := os.Open(historyFile); err == nil {
		line.ReadHistory(f)
		f.Close()
	}
	// 保存历史记录
	defer func() {
		if f, err := os.Create(historyFile); err != nil {
			fmt.Println("Error writing history file: ", err)
		} else {
			line.WriteHistory(f)
			f.Close()
		}
	}()
	// 执行命令行
	for {
		if input, err := line.Prompt("kv-cli> "); err == nil {
			input = strings.TrimSpace(input)
			// 处理输入命令
			if input == "" {
				continue
			} else if input == "exit" || input == "quit" || input == "q" {
				fmt.Println("bye")
				break
			} else if input == "help" {
				fmt.Println("Available commands:")
				for _, cmd := range commands {
					fmt.Println(" -", cmd)
				}
			} else if input == "history" {
				line.WriteHistory(os.Stdout)
			} else {
				userCommand := input
				result, err := client.SendCommand(userCommand)
				if err != nil {
					fmt.Println("Failed!\n" + err.Error())
					continue
				}
				var buf bytes.Buffer
				json.Indent(&buf,[]byte(result),"", "  ")
				fmt.Println("OK!\n" + buf.String())
			}

			// 添加到历史记录
			line.AppendHistory(input)
		} else if err == liner.ErrPromptAborted {
			fmt.Println("Aborted")
			break
		} else {
			fmt.Println("Error reading line: ", err)
			break
		}
	}
}