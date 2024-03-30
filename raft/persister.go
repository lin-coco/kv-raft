package raft

import (
	"encoding/json"

	"kv-raft/raft/common"
)

func (r *Raft) Persist() {
	// 持久化到commitIndex
	var data []string
	for _, entry := range r.logs[:r.commitIndex+1] {
		data = append(data, entry.Command)
	}
	// 调用persister
	if err := r.persister.Persist(data); err != nil {
		panic(err)
	}
	r.logs = r.logs[r.commitIndex+1:]
}

func (r *Raft) ReplaceSnapshot(data []byte) {
	err := r.persister.ReplaceSnapshot(data)
	if err != nil {
		panic(err)
	}
}

func (r *Raft) Snapshot() []byte {
	bytes, err := r.persister.Snapshot()
	if err != nil {
		panic(err)
	}
	return bytes
}

func (r *Raft) SaveState() {
	state := &common.State{
		CurrentTerm:      r.currentTerm,
		VotedFor:         r.votedFor,
		LastIncludeIndex: r.lastIncludeIndex,
		LastIncludeTerm:  r.lastIncludeTerm,
	}
	bytes, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}
	if err = r.persister.SaveState(bytes); err != nil {
		panic(err)
	}
}

func (r *Raft) ReadState() *common.State {
	bytes, err := r.persister.ReadState()
	if err != nil {
		panic(err)
	}
	var state common.State
	if err = json.Unmarshal(bytes, &state); err != nil {
		panic(err)
	}
	return &state
}
