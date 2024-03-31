package raft

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"kv-raft/raft/common"
	"kv-raft/raft/rpc"
)

func (r *Raft) startHeartbeat() {
	for {
		time.Sleep(time.Duration(common.HeartbeatTimeout) * time.Millisecond) // 睡眠
		// 执行appendEntries，也是心跳
		r.mutex.Lock()
		if r.status == common.Leader {
			r.mutex.Unlock()
			r.doHeartbeat()
		}
	}
}

func (r *Raft) doHeartbeat() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	lastLogIndex, _ := r.getLastLogIndexAndTerm() // 即将要发送nextIndex[i]到lastLogIndex之间的log
	for i := 0; i < len(r.peers); i++ {
		if i == r.me {
			continue
		}
		prevLogIndex, prevLogTerm := r.getServerPrevLogIndexAndTerm(i)

		if r.lastIncludeIndex > prevLogIndex { // 落后太久，使用快照同步
			snapshot := r.Snapshot()
			req := &rpc.InstallSnapshotReq{
				Term:             int64(r.currentTerm),
				LeaderId:         int64(r.me),
				LastIncludeIndex: int64(r.lastIncludeIndex),
				LastIncludeTerm:  int64(r.lastIncludeTerm),
				Data:             snapshot,
			}
			go func(server int) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Duration(common.RpcTimeout)*time.Millisecond)
				defer cancel()
				resp, err := r.peers[i].InstallSnapshot(ctx, req)
				if err != nil {
					log.Error("r.peers[i].InstallSnapshot err: %v", err)
				}
				r.handleInstallSnapshotResp(r.lastIncludeIndex, server, resp)
			}(i)
		} else { // 使用日志同步
			req := &rpc.AppendEntriesReq{
				Term:         int64(r.currentTerm),
				LeaderId:     int64(r.me),
				PrevLogIndex: int64(prevLogIndex),
				PrevLogTerm:  int64(prevLogTerm),
				Entries:      nil,
				LeaderCommit: int64(r.commitIndex), // 集群中已经被提交的最高索引
			}
			if lastLogIndex >= r.nextIndex[i] {
				// 发送从[r.nextIndex[i],lastLogIndex] 之间的log
				for j := r.nextIndex[i]; j <= lastLogIndex; j++ {
					sliceIndex := r.getSliceIndexByLogIndex(j)
					req.Entries = append(req.Entries, &rpc.LogEntry{
						Command: r.logs[sliceIndex].Command,
						Term:    int64(r.logs[sliceIndex].Term),
						Index:   int64(r.logs[sliceIndex].Index),
					})
				}
			}
			go func(server int) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Duration(common.RpcTimeout)*time.Millisecond)
				defer cancel()
				resp, err := r.peers[server].AppendEntries(ctx, req)
				if err != nil {
					log.Error("r.peers[%v].AppendEntries err: %v", server, err)
					return
				}
				var sendLastLogIndex int
				if len(req.Entries) > 0 {
					sendLastLogIndex = int(req.Entries[len(req.Entries)-1].Index)
				}
				r.handleAppendEntriesResp(sendLastLogIndex, server, resp)
			}(i)
		}
	}
}

// leader处理appendEntries rpc的结果
func (r *Raft) handleAppendEntriesResp(sendLastLogIndex int, server int, resp *rpc.AppendEntriesResp) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if int(resp.Term) > r.currentTerm {
		r.currentTerm = int(resp.Term)
		r.status = common.Follower
		r.votedFor = -1
		r.SaveState()
		return
	}

	if resp.Success { // 成功
		// 为follower更新nextIndex和matchIndex
		if sendLastLogIndex == 0 { // 没有发送log，不更新nextIndex和matchIndex
			return
		} else { // 发送了log，更新nextIndex和matchIndex
			r.nextIndex[server] = sendLastLogIndex + 1
			r.matchIndex[server] = sendLastLogIndex
			for n := sendLastLogIndex; n > 0; n-- { // 更新commitIndex
				if n <= r.commitIndex {
					break
				}
				if r.getLogTermByIndex(n) != r.currentTerm { // Raft永远不会通过对副本数计数的方式提交之前term的条目。只有leader当前term的日志条目才能通过对副本数计数的方式被提交
					break
				}
				matchCount := 0
				for i := 0; i < len(r.peers); i++ {
					if r.matchIndex[i] >= n {
						matchCount++
					}
				}
				if matchCount >= len(r.peers)/2+1 {
					r.commitIndex = n
					// 应用指令
					r.applyLog()
					r.Persist()
					break
				}
			}
			return
		}
	} else { // 失败
		// 可能日志不一致失败，递减nextIndex
		r.nextIndex[server]--
		return
	}
}

func (r *Raft) handleInstallSnapshotResp(sendLastLogIndex int, server int, resp *rpc.InstallSnapshotResp) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if int(resp.Term) > r.currentTerm {
		r.currentTerm = int(resp.Term)
		r.status = common.Follower
		r.votedFor = -1
		return
	}
	// 为follower更新nextIndex和matchIndex
	r.nextIndex[server] = sendLastLogIndex + 1
	r.matchIndex[server] = sendLastLogIndex
	for n := sendLastLogIndex; n > 0; n-- { // 更新commitIndex
		if n <= r.commitIndex {
			break
		}
		if r.getLogTermByIndex(n) != r.currentTerm { // Raft永远不会通过对副本数计数的方式提交之前term的条目。只有leader当前term的日志条目才能通过对副本数计数的方式被提交
			break
		}
		matchCount := 0
		for i := 0; i < len(r.peers); i++ {
			if r.matchIndex[i] >= n {
				matchCount++
			}
		}
		if matchCount >= len(r.peers)/2+1 {
			r.commitIndex = n
			// 应用指令
			r.applyLog()
			break
		}
	}
	return
}

/*
AppendEntries 心跳rpc
*/
func (r *Raft) AppendEntries(_ context.Context, req *rpc.AppendEntriesReq) (*rpc.AppendEntriesResp, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.lastElectionTime = time.Now().UnixMilli() // 收到心跳，重置选举时间
	var resp rpc.AppendEntriesResp
	if int(req.Term) < r.currentTerm { // leader日志旧了
		resp.Term = int64(r.currentTerm)
		resp.Success = false
		return &resp, nil
	}
	if int(req.Term) > r.currentTerm {
		r.currentTerm = int(req.Term)
		r.status = common.Follower
	}
	// req.Term == r.currentTerm
	r.leaderId = int(req.LeaderId)
	lastLogIndex, _ := r.getLastLogIndexAndTerm()
	if req.PrevLogIndex != 0 {
		if lastLogIndex < int(req.PrevLogIndex) { // 本节点不包含prevLogIndex，返回false
			resp.Term = int64(r.currentTerm)
			resp.Success = false
			return &resp, nil
		}
		if int(req.Term) != r.getLogTermByIndex(int(req.PrevLogIndex)) { // 日志term与prevLogTerm不符合，返回false
			resp.Term = int64(r.currentTerm)
			resp.Success = false
			return &resp, nil
		}
	}
	if len(req.Entries) > 0 {
		var notExistSliceIndex = -1
		for i := 0; i < len(req.Entries); i++ { // 检查已经存在的日志条目和新的是否冲突
			if !r.isExistLogByIndex(int(req.Entries[i].Index)) { // 不存在
				notExistSliceIndex = i
				break
			}
			localTerm := r.getLogTermByIndex(int(req.Entries[i].Index))
			if int(req.Entries[i].Term) != localTerm { // 本地的entry和新的产生冲突，删除本地的entry及其只有所有的
				sliceIndex := r.getSliceIndexByLogIndex(int(req.Entries[i].Index))
				r.logs = r.logs[:sliceIndex]
				notExistSliceIndex = i
				break
			}
		}
		if notExistSliceIndex != -1 { // 不存在，或者产生冲突
			for i := notExistSliceIndex; i < len(req.Entries); i++ { // 附加尚未存在的新entry
				r.logs = append(r.logs, common.LogEntry{
					Command: req.Entries[i].Command,
					Term:    int(req.Entries[i].Term),
					Index:   int(req.Entries[i].Index),
				})
			}
		}
	}
	if int(req.LeaderCommit) > r.commitIndex {
		lastLogIndex, _ = r.getLastLogIndexAndTerm()
		// 提交索引
		r.commitIndex = min(int(req.LeaderCommit), lastLogIndex)
		// 应用指令
		r.applyLog()
	}
	resp.Term = int64(r.currentTerm)
	resp.Success = true
	return &resp, nil
}

/*
InstallSnapshot 同步快照rpc
*/
func (r *Raft) InstallSnapshot(_ context.Context, req *rpc.InstallSnapshotReq) (*rpc.InstallSnapshotResp, error) {
	var resp rpc.InstallSnapshotResp
	r.lastElectionTime = time.Now().UnixMilli() // 收到心跳，重置选举时间
	if int(req.Term) < r.currentTerm {
		resp.Term = int64(r.currentTerm)
		return &resp, nil
	} else if int(req.Term) > r.currentTerm {
		r.currentTerm = int(req.Term)
		r.status = common.Follower
		r.votedFor = -1
		r.SaveState()
	}
	// term == r.currentTerm
	if r.lastIncludeIndex != int(req.LastIncludeIndex) || r.lastIncludeTerm != int(req.LastIncludeTerm) { // 丢掉logs
		r.logs = nil
	}
	r.ReplaceSnapshot(req.Data) // 丢弃旧的快照，存储新的快照
	r.lastIncludeIndex = int(req.LastIncludeIndex)
	r.lastIncludeTerm = int(req.LastIncludeTerm)
	r.SaveState()
	// 重置状态机
	r.Reset()
	r.commitIndex = r.lastIncludeIndex
	r.lastApplied = r.lastIncludeIndex
	// 响应
	resp.Term = int64(r.currentTerm)
	return &resp, nil
}

func (r *Raft) getServerPrevLogIndexAndTerm(server int) (prevLogIndex, prevLogTerm int) {
	prevLogIndex = r.nextIndex[server] - 1
	if prevLogIndex == 0 {
		return 0, 0
	}
	prevLogTerm = r.logs[r.getSliceIndexByLogIndex(prevLogIndex)].Term
	return
}

// 确保logIndex >= 1
func (r *Raft) getSliceIndexByLogIndex(logIndex int) (sliceIndex int) {
	return logIndex - r.lastIncludeIndex - 1
}

func (r *Raft) isExistLogByIndex(logIndex int) bool {
	if logIndex <= 0 {
		return false
	}
	if i := r.getSliceIndexByLogIndex(logIndex); i >= len(r.logs) {
		return false
	}
	return true
}

// 确保先确保索引存在
func (r *Raft) getLogTermByIndex(logIndex int) (logTerm int) {
	return r.logs[r.getSliceIndexByLogIndex(logIndex)].Term
}
