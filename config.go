package kv_raft

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Me    int `json:"me,omitempty"`
	Addrs []struct {
		Server string `json:"server,omitempty"`
		Raft   string `json:"raft,omitempty"`
	} `json:"addrs,omitempty"`
	LogStorageFile   string `json:"log_storage_file,omitempty"`
	StateStorageFile string `json:"state_storage_file,omitempty"`
	ClientRetry      int    `json:"client_retry,omitempty"`
}

func ReadConfigFile(path string) (*Config, error) {
	var config Config
	configContent, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile err: %v", err)
	}
	if err = json.Unmarshal(configContent, &config); err != nil {
		return nil, fmt.Errorf("json.Unmarshal err: %v", err)
	}
	return &config, nil
}

func (c *Config) GetRaftAddrs() []string {
	var raftAddrs []string
	for _, addr := range c.Addrs {
		raftAddrs = append(raftAddrs, addr.Raft)
	}
	return raftAddrs
}
