package raft

import (
	"fmt"
	"net"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"kv-raft/raft/common"
	"kv-raft/raft/rpc"
	"kv-raft/raft/state_machine_interface"
)

func NewRaft(me int, addrs []string,
	persister state_machine_interface.Persister,
	reset state_machine_interface.Reset,
	rwJudge state_machine_interface.RWJudge,
	apply state_machine_interface.Apply,
	isInitReset bool,
	clientCommands <-chan string) (*Raft, error) {
	var r Raft
	listen, err := net.Listen("tcp", addrs[me])
	if err != nil {
		return nil, fmt.Errorf("net.Listen err: %v", err)
	}
	s := grpc.NewServer()
	rpc.RegisterRaftRpcServer(s, &r)
	go func() {
		if err = s.Serve(listen); err != nil {
			panic(err)
		}
	}()
	var peers []rpc.RaftRpcClient
	for i, addr := range addrs {
		if i == me {
			peers = append(peers, nil)
		} else {
			peers = append(peers, rpc.NewRpcClient(addr))
		}
	}
	log.Debug("has collected all rpc client.")
	r.mutex = sync.Mutex{}
	r.me = me // 自己索引
	log.Infof("me: %d", r.me)
	r.leaderId = -1
	r.peers = peers                     // 所有节点的通讯
	r.status = common.Follower          // 初始时都为follower
	r.logs = make([]common.LogEntry, 0) // logEntry
	r.persister = persister             // 上层状态机的持久话实现
	r.reset = reset                     // 上层状态机的重置实现
	r.rwJudge = rwJudge                 // 上层状态机的命令读写判断实现
	r.clientCommands = clientCommands   // 上层状态机传递的用户层命令
	r.apply = apply                     // 上层状态机的应用实现
	state, _ := r.readState()           // 阅读状态
	if state == nil {
		r.currentTerm = 0
		r.votedFor = -1
		r.lastIncludeIndex = 0
		r.lastIncludeTerm = 0
		log.Infof("没有读取到快照 初始化 currentTerm: %d,votedFor: %d,lastIncludeIndex: %d,lastIncludeTerm: %d. ", r.currentTerm, r.votedFor, r.lastIncludeIndex, r.lastIncludeTerm)
		r.commitIndex = 0
		r.lastApplied = 0
	} else {
		r.currentTerm = state.CurrentTerm
		r.votedFor = state.VotedFor
		r.lastIncludeIndex = state.LastIncludeIndex
		r.lastIncludeTerm = state.LastIncludeTerm
		log.Infof("从快照中初始化 currentTerm: %d,votedFor: %d,lastIncludeIndex: %d,lastIncludeTerm: %d. ", r.currentTerm, r.votedFor, r.lastIncludeIndex, r.lastIncludeTerm)
		if !isInitReset {
			r.Reset()
		}
		r.commitIndex = r.lastIncludeIndex // 上层状态机已经应用了快照
		r.lastApplied = r.lastIncludeTerm  // 上层状态机已经应用了快照
	}
	r.nextIndex = make([]int, len(peers))
	for i := 0; i < len(r.nextIndex); i++ {
		r.nextIndex[i] = 1 // 初始化设置为 1，因为不确定其他节点的状态
	}
	r.matchIndex = make([]int, len(peers))       // 初始化matchIndex都为0
	r.lastElectionTime = time.Now().UnixMilli()  // 可以初始化为当前时间，也可以初始化为0
	r.lastHeartbeatTime = time.Now().UnixMilli() // 可以初始化为当前时间，也可以初始化为0
	return &r, nil
}

func (r *Raft) Start() {
	go r.startElection()     // 开始leader选举
	go r.startHeartbeat()    //开始心跳
	go r.startFetchCommand() // 开始抓取应用层命令
}

// IsLeader 为上层状态机判断是否是leader
// 不需要加锁，即使获得了旧的数据，也不影响
func (r *Raft) IsLeader() (bool, int) {
	if r.status == common.Leader {
		return true, r.leaderId
	}
	return false, r.leaderId
}
