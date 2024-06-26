// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.19.4
// source: raft/rpc/raft.proto

package rpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	RaftRpc_RequestVote_FullMethodName     = "/rpc.RaftRpc/RequestVote"
	RaftRpc_AppendEntries_FullMethodName   = "/rpc.RaftRpc/AppendEntries"
	RaftRpc_InstallSnapshot_FullMethodName = "/rpc.RaftRpc/InstallSnapshot"
)

// RaftRpcClient is the client API for RaftRpc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RaftRpcClient interface {
	// candidate收集票
	// 接受者实现:
	// 1. 回复false如果term < currentTerm
	// 2. 如果 votedFor 为空或者为 candidateId，并且候选人的日志至少和自己一样新，那么就投票给他，也就是自己还没投给别人，且这个候选人合法的话，就投票给他
	RequestVote(ctx context.Context, in *RequestVoteReq, opts ...grpc.CallOption) (*RequestVoteResp, error)
	// leader同步日志，也被作为心跳
	// 接收者实现:
	// 1. 回复false如果term < currentTerm
	// 2. 如果本节点日志中不包含prevLogIndex或者日志任期和prevLogTerm不符合，返回false
	// 3. 如果已经存在的日志条目和新的产生冲突（索引值相同但是任期号不同），删除这一条和之后所有的
	// 4. 附加日志中尚未存在的任何新条目
	// 5. 如果 leaderCommit > commitIndex，令 commitIndex 等于 leaderCommit 和 last新日志条目索引值中较小的一个
	AppendEntries(ctx context.Context, in *AppendEntriesReq, opts ...grpc.CallOption) (*AppendEntriesResp, error)
	// leader将快照发送给落后的follower
	// 1. if term < r.currentTerm 回复当前term
	// 2. 将数据写入快照中，替代为新的快照，丢弃任何旧的快照
	// 3. 如果lastIncludeIndex、lastIncludeTerm和快照中最后的一样，则保留logs，否则丢弃之后的所有
	// 4. 重置状态机
	InstallSnapshot(ctx context.Context, in *InstallSnapshotReq, opts ...grpc.CallOption) (*InstallSnapshotResp, error)
}

type raftRpcClient struct {
	cc grpc.ClientConnInterface
}

func NewRaftRpcClient(cc grpc.ClientConnInterface) RaftRpcClient {
	return &raftRpcClient{cc}
}

func (c *raftRpcClient) RequestVote(ctx context.Context, in *RequestVoteReq, opts ...grpc.CallOption) (*RequestVoteResp, error) {
	out := new(RequestVoteResp)
	err := c.cc.Invoke(ctx, RaftRpc_RequestVote_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftRpcClient) AppendEntries(ctx context.Context, in *AppendEntriesReq, opts ...grpc.CallOption) (*AppendEntriesResp, error) {
	out := new(AppendEntriesResp)
	err := c.cc.Invoke(ctx, RaftRpc_AppendEntries_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftRpcClient) InstallSnapshot(ctx context.Context, in *InstallSnapshotReq, opts ...grpc.CallOption) (*InstallSnapshotResp, error) {
	out := new(InstallSnapshotResp)
	err := c.cc.Invoke(ctx, RaftRpc_InstallSnapshot_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RaftRpcServer is the server API for RaftRpc service.
// All implementations must embed UnimplementedRaftRpcServer
// for forward compatibility
type RaftRpcServer interface {
	// candidate收集票
	// 接受者实现:
	// 1. 回复false如果term < currentTerm
	// 2. 如果 votedFor 为空或者为 candidateId，并且候选人的日志至少和自己一样新，那么就投票给他，也就是自己还没投给别人，且这个候选人合法的话，就投票给他
	RequestVote(context.Context, *RequestVoteReq) (*RequestVoteResp, error)
	// leader同步日志，也被作为心跳
	// 接收者实现:
	// 1. 回复false如果term < currentTerm
	// 2. 如果本节点日志中不包含prevLogIndex或者日志任期和prevLogTerm不符合，返回false
	// 3. 如果已经存在的日志条目和新的产生冲突（索引值相同但是任期号不同），删除这一条和之后所有的
	// 4. 附加日志中尚未存在的任何新条目
	// 5. 如果 leaderCommit > commitIndex，令 commitIndex 等于 leaderCommit 和 last新日志条目索引值中较小的一个
	AppendEntries(context.Context, *AppendEntriesReq) (*AppendEntriesResp, error)
	// leader将快照发送给落后的follower
	// 1. if term < r.currentTerm 回复当前term
	// 2. 将数据写入快照中，替代为新的快照，丢弃任何旧的快照
	// 3. 如果lastIncludeIndex、lastIncludeTerm和快照中最后的一样，则保留logs，否则丢弃之后的所有
	// 4. 重置状态机
	InstallSnapshot(context.Context, *InstallSnapshotReq) (*InstallSnapshotResp, error)
	mustEmbedUnimplementedRaftRpcServer()
}

// UnimplementedRaftRpcServer must be embedded to have forward compatible implementations.
type UnimplementedRaftRpcServer struct {
}

func (UnimplementedRaftRpcServer) RequestVote(context.Context, *RequestVoteReq) (*RequestVoteResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestVote not implemented")
}
func (UnimplementedRaftRpcServer) AppendEntries(context.Context, *AppendEntriesReq) (*AppendEntriesResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AppendEntries not implemented")
}
func (UnimplementedRaftRpcServer) InstallSnapshot(context.Context, *InstallSnapshotReq) (*InstallSnapshotResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InstallSnapshot not implemented")
}
func (UnimplementedRaftRpcServer) mustEmbedUnimplementedRaftRpcServer() {}

// UnsafeRaftRpcServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RaftRpcServer will
// result in compilation errors.
type UnsafeRaftRpcServer interface {
	mustEmbedUnimplementedRaftRpcServer()
}

func RegisterRaftRpcServer(s grpc.ServiceRegistrar, srv RaftRpcServer) {
	s.RegisterService(&RaftRpc_ServiceDesc, srv)
}

func _RaftRpc_RequestVote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestVoteReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftRpcServer).RequestVote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RaftRpc_RequestVote_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftRpcServer).RequestVote(ctx, req.(*RequestVoteReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _RaftRpc_AppendEntries_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AppendEntriesReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftRpcServer).AppendEntries(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RaftRpc_AppendEntries_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftRpcServer).AppendEntries(ctx, req.(*AppendEntriesReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _RaftRpc_InstallSnapshot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InstallSnapshotReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftRpcServer).InstallSnapshot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RaftRpc_InstallSnapshot_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftRpcServer).InstallSnapshot(ctx, req.(*InstallSnapshotReq))
	}
	return interceptor(ctx, in, info, handler)
}

// RaftRpc_ServiceDesc is the grpc.ServiceDesc for RaftRpc service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RaftRpc_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.RaftRpc",
	HandlerType: (*RaftRpcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RequestVote",
			Handler:    _RaftRpc_RequestVote_Handler,
		},
		{
			MethodName: "AppendEntries",
			Handler:    _RaftRpc_AppendEntries_Handler,
		},
		{
			MethodName: "InstallSnapshot",
			Handler:    _RaftRpc_InstallSnapshot_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "raft/rpc/raft.proto",
}
