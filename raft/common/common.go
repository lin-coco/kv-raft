package common

type Status int

const (
	Leader    = 1 // 领导者
	Follower  = 2 // 跟随者
	Candidate = 3 // 候选者

	ElectionBaseTimeout     = 150 // 选举基准超时时间（毫秒）
	ElectionMaxExtraTimeout = 150 // 选举额外最大超时时间（毫秒）
	HeartbeatTimeout        = 25  // 心跳时间超时时间（毫秒）TODO 待确定

	RpcTimeout = 10 // rpc超时时间，设置100毫秒
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
