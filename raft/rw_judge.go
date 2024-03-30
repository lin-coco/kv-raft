package raft

func (r *Raft) IsCommandReadOnly(command string) bool {
	return r.rwJudge.ReadOnly(command)
}
