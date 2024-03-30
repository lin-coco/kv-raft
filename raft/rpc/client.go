package rpc

import (
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// addr = flag.String("addr", "127.0.0.1:8972", "the address to connect to")

func NewRpcClient(addr string) RaftRpcClient {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("collect addr %s failed: %v, retry...", addr, err)
		panic(err)
	}
	return NewRaftRpcClient(conn)
}
