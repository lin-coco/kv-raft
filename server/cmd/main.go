package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"kv-raft"
	"kv-raft/raft"
	"kv-raft/server"
	"kv-raft/server/impl"
)

func main() {
	app := cli.NewApp()
	app.Name = "mykv-server"
	app.Usage = "kv storage system based on raft"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "f",
			Usage: "-f config file",
		},
	}
	app.Before = func(context *cli.Context) error {
		// Log as JSON instead of the default ASCII formatter.
		log.SetFormatter(&log.JSONFormatter{})
		log.SetOutput(os.Stdout)
		log.SetLevel(log.DebugLevel)
		return nil
	}
	app.Action = func(ctx *cli.Context) {
		configFilePath := ctx.String("f")
		// 阅读配置
		config, err := kv_raft.ReadConfigFile(configFilePath)
		if err != nil {
			log.Errorf("kv_raft.ReadConfigFile err: %v", err)
			return
		}
		log.Debugf("config loaded: %v.", config)
		if err = Server(config); err != nil {
			log.Errorf("server run failed: %v", err)
		}
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func Server(config *kv_raft.Config) error {
	// 创建持久化的文件描述符
	logFile, err := os.OpenFile(config.LogStorageFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("os.OpenFile err: %v", err)
	}
	defer func() {
		_ = logFile.Close()
	}()
	stateFile, err := os.OpenFile(config.StateStorageFile, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("os.OpenFile err: %v", err)
	}
	defer func() {
		_ = logFile.Close()
	}()
	// 初始化持久化组件
	persister := impl.NewFilePersister(logFile, stateFile)
	// 初始化重置组件
	reset := impl.NewReset(logFile)
	// 初始化读写判断命令组件
	rwJudge := impl.NewRWJudge()
	// 创建客户端命令通道
	clientCommands := make(chan string)
	// 创建raft节点
	r, err := raft.NewRaft(config.Me, config.GetRaftAddrs(), persister, reset, rwJudge, false, clientCommands)
	if err != nil {
		return fmt.Errorf("raft.NewRaft err: %v", err)
	}
	// 启动raft
	r.Start()
	// 创建http Server
	if err = RunKVServer(config, r, clientCommands); err != nil {
		return fmt.Errorf("RunKVServer err: %v", err)
	}
	return nil
}

const (
	Success = 1
	Failed  = 2
	Forward = 3
)

func RunKVServer(config *kv_raft.Config, r *raft.Raft, clientCommands chan string) error {
	http.HandleFunc("/receive", func(writer http.ResponseWriter, request *http.Request) {
		// 判断当前是否是leader节点
		b, leaderId := r.IsLeader()
		if !b {
			_, err := writer.Write([]byte{Forward, byte(leaderId)})
			if err != nil {
				_, err = writer.Write([]byte(err.Error()))
				if err != nil {
					log.Error("writer.Write err: %v", err)
					return
				}
			}
		}
		// 获取请求体
		bytes, err := io.ReadAll(request.Body)
		if err != nil {
			reason := []byte(err.Error())
			reason = append([]byte{Failed}, reason...)
			_, err = writer.Write(reason)
			if err != nil {
				log.Error("writer.Write err: %v", err)
				return
			}
		}
		// 校验requestId
		requestID := request.Header.Get("KV-Raft-Request-ID")
		if len(requestID) != 36 {
			reason := []byte("length of the request ID is incorrect")
			reason = append([]byte{Failed}, reason...)
			_, err = writer.Write(reason)
			if err != nil {
				log.Error("writer.Write err: %v", err)
				return
			}
		}
		// 校验命令的正确性
		command := string(bytes)
		if err = server.CheckCommand(command); err != nil {
			reason := []byte(err.Error())
			reason = append([]byte{Failed}, reason...)
			_, err = writer.Write(reason)
			if err != nil {
				log.Error("writer.Write err: %v", err)
				return
			}
		}
		// 发送命令
		command += " " + requestID
		clientCommands <- command
		// 等待提交

	})
	if err := http.ListenAndServe(config.Addrs[config.Me].Server, nil); err != nil {
		return fmt.Errorf("http.ListenAndServe err: %v", err)
	}
	return nil
}
