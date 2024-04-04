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
			log.Debugf("我是leader，开始新一轮心跳")
			r.doHeartbeat()
		} else {
			//log.Debugf("我非leader，不进行心跳")
			r.mutex.Unlock()
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
		//prevLogIndex, prevLogTerm := r.getServerPrevLogIndexAndTerm(i)
		prevLogIndex := r.getServerPrevLogIndex(i)
		if r.lastIncludeIndex > prevLogIndex { // 落后太久，使用快照同步
			log.Warnf("落后太久，使用快照同步, server: %v,r.lastIncludeIndex: %v,prevLogIndex: %v", i, r.lastIncludeIndex, prevLogIndex)
			snapshot := r.Snapshot()
			req := &rpc.InstallSnapshotReq{
				Term:             int64(r.currentTerm),
				LeaderId:         int64(r.me),
				LastIncludeIndex: int64(r.lastIncludeIndex),
				LastIncludeTerm:  int64(r.lastIncludeTerm),
				Data:             snapshot,
			}
			go func(server int) {
				resp, err := r.peers[server].InstallSnapshot(context.Background(), req)
				if err != nil {
					log.Errorf("发送快照心跳失败 节点id:%d err: %v", server, err)
					return
				}
				log.Debugf("发送快照心跳 节点id:%d req.Term:%d,LastIncludeIndex:%d,LastIncludeTerm:%d", server, req.Term, req.LastIncludeIndex, req.LastIncludeTerm)
				r.handleInstallSnapshotResp(r.lastIncludeIndex, server, resp)
			}(i)
		} else { // 使用日志同步
			prevLogTerm := r.getLogTermByIndex(prevLogIndex)
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
				resp, err := r.peers[server].AppendEntries(context.Background(), req)
				if err != nil {
					log.Errorf("发送日志心跳失败 节点id:%d err: %v", server, err)
					return
				}
				log.Debugf("发送日志心跳 节点id:%d req.Term:%d,PrevLogIndex:%d,PrevLogTerm:%d,logSize:%d,LeaderCommit:%d", server, req.Term, req.PrevLogIndex, req.PrevLogTerm, len(req.Entries), req.LeaderCommit)
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
	if r.status == common.Follower {
		log.Debugf("我已经成为Follower，忽略收到的日志心跳响应")
		return
	}
	if int(resp.Term) > r.currentTerm {
		r.currentTerm = int(resp.Term)
		r.status = common.Follower
		r.votedFor = -1
		r.saveState()
		log.Debugf("收到server:%d日志心跳回复，它的term:%d比我:%d的大，更新自己term，成为Follower", server, resp.Term, r.currentTerm)
		return
	}
	if resp.Success { // 成功
		// 为follower更新nextIndex和matchIndex
		if sendLastLogIndex == 0 { // 没有发送log，不更新nextIndex和matchIndex
			log.Debugf("收到server:%d日志心跳回复，成功", server)
			return
		} else { // 发送了log，更新nextIndex和matchIndex
			r.nextIndex[server] = sendLastLogIndex + 1
			r.matchIndex[server] = sendLastLogIndex
			log.Debugf("收到server:%d日志心跳回复，成功，nextIndex:%d, matchIndex:%d", server, r.nextIndex[server], r.matchIndex[server])
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
					r.persist()
					r.saveState()
					log.Debugf("日志索引:%d, 超过半数match，提交索引并应用", r.commitIndex)
					break
				}
			}
			return
		}
	} else { // 失败
		// 可能日志不一致失败，递减nextIndex
		r.nextIndex[server]--
		log.Debugf("收到server:%d日志心跳回复，发生了日志不一样的错误，递减nextIndex:%d", server, r.nextIndex[server])
		return
	}
}

func (r *Raft) handleInstallSnapshotResp(sendLastLogIndex int, server int, resp *rpc.InstallSnapshotResp) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.status == common.Follower {
		log.Debugf("我已经成为Follower，忽略收到的快照心跳响应")
		return
	}
	if int(resp.Term) > r.currentTerm {
		r.currentTerm = int(resp.Term)
		r.status = common.Follower
		r.votedFor = -1
		r.saveState()
		log.Debugf("收到server:%d快照心跳回复，它的term:%d比我:%d的大，更新自己term，成为Follower", server, resp.Term, r.currentTerm)
		return
	}
	// 为follower更新nextIndex和matchIndex
	r.nextIndex[server] = sendLastLogIndex + 1
	r.matchIndex[server] = sendLastLogIndex
	log.Debugf("收到server:%d快照心跳回复，成功，设置此server的nextIndex:%d, matchIndex:%d", server, r.nextIndex[server], r.matchIndex[server])
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
			r.persist()
			r.saveState()
			log.Debugf("日志索引:%d, 超过半数match，提交索引并应用", r.commitIndex)
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
		log.Debugf("收到server:%d日志心跳，它的Term:%d比我的旧:%d，让它更新自己", req.LeaderId, req.Term, r.currentTerm)
		return &resp, nil
	}
	if int(req.Term) > r.currentTerm {
		r.currentTerm = int(req.Term)
		r.status = common.Follower
		r.votedFor = -1
		r.saveState()
	}
	// req.Term == r.currentTerm
	r.leaderId = int(req.LeaderId)
	lastLogIndex, _ := r.getLastLogIndexAndTerm()
	if req.PrevLogIndex != 0 {
		if lastLogIndex < int(req.PrevLogIndex) { // 本节点不包含prevLogIndex，返回false
			resp.Term = int64(r.currentTerm)
			resp.Success = false
			log.Debugf("收到server:%d日志心跳，prevLogIndex:%d还不存在，拒绝同步日志, lastLogIndex:%d", req.LeaderId, req.PrevLogIndex, lastLogIndex)
			return &resp, nil
		}
		term := r.getLogTermByIndex(int(req.PrevLogIndex))
		if int(req.PrevLogTerm) != term { // 日志term与prevLogTerm不符合，返回false
			resp.Term = int64(r.currentTerm)
			resp.Success = false
			log.Debugf("收到server:%d日志心跳，PrevLogTerm:%d不等于本节点相同位置处的Term:%d，日志不统一，拒绝同步日志", req.LeaderId, req.PrevLogTerm, term)
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
		r.persist()
		r.saveState()
		log.Debugf("日志索引:%d, 提交索引并应用", r.commitIndex)
	}
	resp.Term = int64(r.currentTerm)
	resp.Success = true
	log.Debugf("收到server:%d日志心跳，同意同步日志", req.LeaderId)
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
		log.Debugf("收到server:%d快照心跳，它的Term:%d比我的旧:%d，让它更新自己", req.LeaderId, req.Term, r.currentTerm)
		return &resp, nil
	} else if int(req.Term) > r.currentTerm {
		r.currentTerm = int(req.Term)
		r.status = common.Follower
		r.votedFor = -1
		r.saveState()
	}
	// term == r.currentTerm
	if r.lastIncludeIndex != int(req.LastIncludeIndex) || r.lastIncludeTerm != int(req.LastIncludeTerm) { // 丢掉logs
		r.logs = nil
	}
	r.replaceSnapshot(req.Data) // 丢弃旧的快照，存储新的快照
	r.lastIncludeIndex = int(req.LastIncludeIndex)
	r.lastIncludeTerm = int(req.LastIncludeTerm)
	r.saveState()
	// 重置状态机
	r.Reset()
	r.commitIndex = r.lastIncludeIndex
	r.lastApplied = r.lastIncludeIndex
	// 响应
	resp.Term = int64(r.currentTerm)
	log.Debugf("收到server:%d快照心跳，成功同步快照，commitIndex:%d", req.LeaderId, r.commitIndex)
	return &resp, nil
}

//	func (r *Raft) getServerPrevLogIndexAndTerm(server int) (prevLogIndex, prevLogTerm int) {
//		prevLogIndex = r.nextIndex[server] - 1
//		if prevLogIndex == 0 {
//			return 0, 0
//		}
//		if prevLogIndex == r.lastIncludeIndex {
//			prevLogTerm = r.lastIncludeTerm
//			return
//		}
//		log.Errorf("server: %v, r.nextIndex[server]: %v,prevLogIndex: %v, r.lastIncludeIndex: %v, r.lastIncludeTerm: %v", server, r.nextIndex[server], prevLogIndex, r.lastIncludeIndex, r.lastIncludeTerm)
//		prevLogTerm = r.logs[r.getSliceIndexByLogIndex(prevLogIndex)].Term
//		return
//	}
func (r *Raft) getServerPrevLogIndex(server int) (prevLogIndex int) {
	prevLogIndex = r.nextIndex[server] - 1
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

// 确保索引存在
func (r *Raft) getLogTermByIndex(logIndex int) int {
	if logIndex == 0 {
		return 0
	}
	if logIndex == r.lastIncludeIndex {
		return r.lastIncludeTerm
	}
	return r.logs[r.getSliceIndexByLogIndex(logIndex)].Term
}
