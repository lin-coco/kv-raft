.PHONY: generate

generate: raft/rpc/raft.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative $<

raft/rpc/raft.proto:
	@echo "Error: raft/rpc/raft.proto 文件不存在"
	@exit 1
