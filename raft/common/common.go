package common

type Status int

const (
	Leader    = 1 // 领导者
	Follower  = 2 // 跟随者
	Candidate = 3 // 候选者

	ElectionBaseTimeout     = 1500 // 选举基准超时时间（毫秒）
	ElectionMaxExtraTimeout = 1500 // 选举额外最大超时时间（毫秒）
	HeartbeatTimeout        = 250  // 心跳时间超时时间（毫秒）

	//RpcTimeout = 10 // rpc超时时间，设置100毫秒
)

type LogEntry struct {
	Command string
	Term    int
	Index   int
}

// State 需要保存的状态
type State struct {
	CurrentTerm      int
	VotedFor         int
	LastIncludeIndex int
	LastIncludeTerm  int
}
