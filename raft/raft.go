package raft

import (
	"sync"

	"kv-raft/raft/common"
	"kv-raft/raft/rpc"
	"kv-raft/raft/state_machine_interface"
)

type Raft struct {
	mutex    sync.Mutex
	me       int
	peers    []rpc.RaftRpcClient // 与其他节点通信的rpc入口
	status   common.Status       // 身份
	leaderId int
	// Persistent state on all servers:
	currentTerm int               // 当前term
	votedFor    int               // 记录当前term给谁投票了
	logs        []common.LogEntry // 日志数组 第一个log.Index为1
	// Volatile state on all servers:
	commitIndex int // 提交日志索引，表示单个节点可以被提交的最高日志索引
	lastApplied int // 最后一次应用给状态机log的index，它会将从 lastApplied+1 到 commitIndex 之间的所有日志条目应用到状态机中去，以确保状态机的一致性
	// Volatile state on leaders:
	// leader维护 这两个状态的下标(是元素值不是物理下标)1开始，因为通常commitIndex和lastApplied从0开始，应该是一个无效的index，因此下标从1开始
	nextIndex  []int // 记录下一次将要发送给Follower的日志条目的索引
	matchIndex []int // 记录已经复制给该Follower的最高日志条目的索引，之前的已经全部被接收

	// 记录选举和心跳的时间
	lastElectionTime  int64
	lastHeartbeatTime int64
	// 日志快照 仅存储已经提交的log
	persister        state_machine_interface.Persister
	lastIncludeIndex int // 在日志快照中最后一个logIndex
	lastIncludeTerm  int // 在日志快照中最后一个logIndex的term
	// 重置状态机
	reset state_machine_interface.Reset
	// 上层客户端的指令 leader取
	clientCommands <-chan string
	// 判断是否是只读命令
	rwJudge state_machine_interface.RWJudge
	rpc.UnimplementedRaftRpcServer
}

func (r *Raft) getLastLogIndexAndTerm() (lastLogIndex, lastLogTerm int) {
	if len(r.logs) == 0 {
		return 0, 0
	}
	return r.logs[len(r.logs)-1].Index, r.logs[len(r.logs)-1].Term
}
