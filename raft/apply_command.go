package raft

import log "github.com/sirupsen/logrus"

// 应用状态机，从r.lastApplied 到 commitIndex
func (r *Raft) applyLog() {
	log.Infof("执行应用日志... 从lastApplied+1:%v到commitIndex:%v", r.lastApplied+1, r.commitIndex)
	for i := r.lastApplied + 1; i <= r.commitIndex; i++ {
		log.Infof("日志索引:%d", i)
		applyLogSliceIndex := r.getSliceIndexByLogIndex(i)
		command := r.logs[applyLogSliceIndex].Command
		term := r.logs[applyLogSliceIndex].Term
		r.apply.ApplyCommand(command)
		log.Infof("应用日志 command: %v,logTerm: %v,logIndex: %v", command, term, i)
	}
	log.Info("执行应用日志结束...")
	r.lastApplied = r.commitIndex
}
