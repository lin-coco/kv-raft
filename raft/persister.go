package raft

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"kv-raft/raft/common"
)

func (r *Raft) persist() {
	// 持久化到commitIndex
	var data []string
	commitSliceIndex := r.getSliceIndexByLogIndex(r.commitIndex)
	for _, entry := range r.logs[:commitSliceIndex+1] {
		data = append(data, entry.Command)
	}
	// 调用persister
	if err := r.persister.Persist(data); err != nil {
		panic(err)
	}
	r.lastIncludeIndex = r.logs[commitSliceIndex].Index
	r.lastIncludeTerm = r.logs[commitSliceIndex].Term
	r.saveState()
	r.logs = r.logs[commitSliceIndex+1:]
}

func (r *Raft) replaceSnapshot(data []byte) {
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

func (r *Raft) saveState() {
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

func (r *Raft) readState() (*common.State, error) {
	bytes, err := r.persister.ReadState()
	if err != nil {
		panic(err)
	}
	log.Infof("last state: %v", string(bytes))
	var state common.State
	if err = json.Unmarshal(bytes, &state); err != nil {
		return nil, fmt.Errorf("json.Unmarshal err: %v", err)
	}
	return &state, nil
}
