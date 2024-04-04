package raft

import (
	"context"
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"

	"kv-raft/raft/common"
	"kv-raft/raft/rpc"
)

/*
startElection 一段时间内没有收到心跳，就开始一次选举
*/
func (r *Raft) startElection() {
	for {
		// 睡眠
		flagTime := time.Now().UnixMilli()                                     // 标记开始时间
		time.Sleep(time.Duration(getRandElectionTimeout()) * time.Millisecond) // 睡眠本次选举超时时间
		// 检查是否选举要发生
		r.mutex.Lock()
		if r.status == common.Leader {
			// 自己是leader，不进行选举
			r.mutex.Unlock()
		} else if r.lastElectionTime > flagTime {
			// 自己不是leader，超时选举时间被刷新了，说明在此期间收到了心跳，开始下一次睡眠
			log.Debugf("从%d到now 收到了心跳，不进行选举", flagTime)
			r.mutex.Unlock()
		} else {
			// 自己不是leader，超时选举时间没有被刷新，说明leader不可用
			log.Debugf("从%d到now 没有收到心跳，开始执行选举", flagTime)
			r.mutex.Unlock()
			r.doElection()
		}
	}
}

// 由于一段时间没有收到心跳，leader超时了，执行选举
func (r *Raft) doElection() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	// 1. 增大当前term
	r.currentTerm++
	// 2. 转换到候选者
	r.status = common.Candidate
	// 3. 为自己投票
	r.votedFor = r.me
	r.saveState()
	voteNum := 1 // 自己获得的投票数
	log.Debugf("开始执行选举 转换成Candidate，任期CurrentTerm自增到:%d，为自己投票VotedNum:%d", r.currentTerm, voteNum)
	//4. 重置选举计时器
	//r.lastElectionTime = time.Now().UnixMilli()
	// 4. 并行地给集群中每个其它的服务器发起RequestVote RPC
	//ctx, _ := context.WithTimeout(context.Background(), common.RpcTimeout)
	lastLogIndex, lastLogTerm := r.getLastLogIndexAndTerm()
	req := &rpc.RequestVoteReq{
		Term:         int64(r.currentTerm),
		CandidateId:  int64(r.me),
		LastLogIndex: int64(lastLogIndex),
		LastLogTerm:  int64(lastLogTerm),
	}
	for i := 0; i < len(r.peers); i++ {
		if i == r.me {
			continue
		}
		go func(server int) {
			resp, err := r.peers[server].RequestVote(context.Background(), req)
			if err != nil {
				log.Errorf("发送请求投票失败 节点id:%v: err: %v", server, err)
				return
			}
			log.Debugf("发送请求投票 节点id:%v, data:%v", server, req)
			r.handleRequestVoteResp(server, resp, &voteNum)
		}(i)
	}
}

// 处理投票响应
func (r *Raft) handleRequestVoteResp(server int, resp *rpc.RequestVoteResp, voteNum *int) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	//log.Debugf("收到server:%d 投票回复 resp: %v",server,resp)
	// rpc请求或响应，都要检查term
	if r.status == common.Follower {
		log.Debugf("我已经成为Follower，收到票也没用了")
		return
	}
	if int(resp.Term) > r.currentTerm {
		r.currentTerm = int(resp.Term)
		r.status = common.Follower
		r.votedFor = -1
		r.saveState()
		log.Debugf("收到server:%d回复选举 serverTerm:%d比my.currentTerm:%d大 从Candidate变成Follower", server, resp.Term, r.currentTerm)
		return
	} else if int(resp.Term) < r.currentTerm { // 网络延迟较大的情况会出现term < currentTerm。投票过后成为leader会出现这种情况
		log.Debugf("收到server:%d回复选举 收到的term比自己小，可能已经有下一个leader被选举出，或者网络延迟较大 server:%d resp.Term:%d r.currentTerm:%d", server, resp.Term, r.currentTerm, r.currentTerm)
		return // 不处理
	}
	// term == currentTerm
	if resp.VoteGranted {
		*voteNum++
		// 已经是leader了，不会进行下一步处理
		if r.status == common.Leader {
			log.Debugf("收到server:%d回复选举 获得投票，不过已经成为leader了, voteNum:%d", server, *voteNum)
		} else {
			if *voteNum >= len(r.peers)/2+1 { // 获得了大多数的投票，变成leader
				r.status = common.Leader
				r.currentTerm++
				r.votedFor = -1
				r.leaderId = r.me
				r.saveState()
				lastLogIndex, _ := r.getLastLogIndexAndTerm()
				for i := 0; i < len(r.peers); i++ {
					if i == r.me {
						continue
					}
					r.nextIndex[i] = lastLogIndex + 1 // 初始化设置为lastLogIndex + 1
					r.matchIndex[i] = 0               // 初始化设置成0，表示新leader和follower之间还没有复制日志
				}
				log.Debugf("收到server:%d回复选举 获得投票，过半投票，变成leader，voteNum:%d, currentTerm:%d", server, *voteNum, r.currentTerm)
				// TODO 成为leader立刻进行一次心跳
				//r.doHeartbeat()
			} else {
				log.Debugf("收到server:%d回复选举 获得投票，未过半，voteNum:%d", server, *voteNum)
			}
		}
	} else {
		log.Debugf("收到server:%d回复选举 未获得投票，可能我的日志没他新或者他已经投给别人了，voteNum:%d", server, *voteNum)
	}
}

/*
RequestVote 请求投票rpc
*/
func (r *Raft) RequestVote(_ context.Context, req *rpc.RequestVoteReq) (*rpc.RequestVoteResp, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	var resp rpc.RequestVoteResp
	if int(req.Term) < r.currentTerm { // 让对方更新自己的term
		resp.Term = int64(r.currentTerm)
		resp.VoteGranted = false
		log.Debugf("收到了server:%d选举 我的term:%d比它:%d大，拒绝投票，去更新一下term吧", req.CandidateId, r.currentTerm, req.Term)
		return &resp, nil
	}
	if int(req.Term) > r.currentTerm { // 更新自己的term
		r.currentTerm = int(req.Term)
		r.status = common.Follower
		r.votedFor = -1
		r.saveState()
	}
	// term == r.currentTerm
	// 检查lastLogIndex和lastLogTerm
	lastLogIndex, lastLogTerm := r.getLastLogIndexAndTerm()
	if lastLogTerm < int(req.LastLogTerm) || (lastLogTerm == int(req.LastLogTerm) && lastLogIndex <= int(req.LastLogIndex)) { // 我的日志比它的旧
		if r.votedFor == -1 || r.votedFor == int(req.CandidateId) { // 为什么有第二个条件？说明到新一个term没有重置voteFor || 已经投给candidate，但是candidate没有收到，又发生了选举
			resp.Term = int64(r.currentTerm)
			resp.VoteGranted = true
			r.votedFor = int(req.CandidateId)
			r.lastElectionTime = time.Now().UnixMilli() // 投完票，重置选举时间
			r.saveState()
			log.Debugf("收到了server:%d选举 我的日志比它旧或一样，同意投票", req.CandidateId)
			return &resp, nil
		} else { // 已经投给其他candidate了
			resp.Term = int64(r.currentTerm)
			resp.VoteGranted = false
			log.Debugf("收到了server:%d选举 你适合我投票，不过我已经投票给别人", req.CandidateId)
			return &resp, nil
		}
	} else { // 我的日志比它的新
		resp.Term = int64(r.currentTerm)
		resp.VoteGranted = false // 它的日志没有我的新，不进行投票
		log.Debugf("收到了server:%d选举 我的日志比它的新，不进行投票", req.CandidateId)
		return &resp, nil
	}
}

// 返回[150,300)
func getRandElectionTimeout() int64 {
	r := rand.New(rand.NewSource(time.Now().UnixMicro()))
	return common.ElectionBaseTimeout + r.Int63n(common.ElectionMaxExtraTimeout)
}
