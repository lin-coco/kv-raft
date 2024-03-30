package raft

func (r *Raft) Reset() {
	if err := r.reset.Reset(); err != nil {
		panic(err)
	}
}
