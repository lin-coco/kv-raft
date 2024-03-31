package raft

// 应用状态机，从r.lastApplied 到 commitIndex
func (r *Raft) applyLog() {
	for i := r.lastApplied + 1; i < r.commitIndex; i++ {
		r.apply.ApplyCommand(r.logs[r.getSliceIndexByLogIndex(i)].Command)
	}
}
