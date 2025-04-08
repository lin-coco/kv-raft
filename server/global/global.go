package global

import (
	kv_raft "kv-raft"
	"kv-raft/raft"
	"kv-raft/server/storage"
)

var (
	Config *kv_raft.Config
	// raft实例
	R *raft.Raft
	// server发送给raft
	ClientCommands chan string
	// 存储引擎
	StorageEngine storage.StorageEngineInterface
	// Reqs key: requestID value: result
	Reqs = make(map[string]chan string)
)

const (
	Success = "success"
	Failed  = "failed"
	Forward = "forward"
)
