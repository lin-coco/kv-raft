# raft-core

## TODO

1. 在通常情况下，一条指令能在集群的大多数响应一轮远程过程调用（RPC）后完成；少数的较慢的服务器不会影响整个系统的性能。

   不要使用waitgroup等待所有rpc响应

2. 

## Raft拆分

- 领导选举
- 日志复制
- 安全性
- 变更集群成员

## 规则

### All Servers
- 如果commitIndex > lastApplied，自增lastApplied，应用lastApplied LogEntry到状态机
- 如果rpc请求或者响应包含term > currentTerm，设置currentTerm = term，更新身份状态到follower
### Followers
- 响应leader和candidate的rpc请求
- 如果electionTimeout没有收到appendEntries或者投票给一个candidate，转化成candidate
### Candidates
- 转换成candidate后，开始选举
    - 自增当前term
    - 投票给自己
    - 重置选举计时器
    - 发送投票请求向所有其他server
- 如果从一半以上的server获得票，成为leader
- 如果从新leader接收到appendEntries，成为follower
- 如果选举超时，开始新一轮选举

### Leaders
- 选举后：向每个server发送空的appendEntries；在空闲期间重复防止选举超时
- 如果从client接收到command，添加entry到local logs，在entry应用到状态机之后响应
- 如果lastLogIndex >= nextIndex[i]，发送带有从nextIndex开始的logEntries的appendEntries rpc
    - 如果成功，为follower更新nextIndex和matchIndex
    - 如果AppendEntries由于日志不一致而失败：递减nextIndex并重试
    - 如果存在一个N，使得N > commitIndex && 大多数matchIndex[i] >= N && logs[N].term == currentTerm，设置commitIndex = N

## 浓缩总结

![img.png](https://typora-img-xue.oss-cn-beijing.aliyuncs.com/img/img.png)

## 关键性质

<img src="https://typora-img-xue.oss-cn-beijing.aliyuncs.com/img/figure-3-20240323203650508.png" alt="image" style="zoom:38%;" />

## 日志压缩

快照仅覆盖已经提交的log

<img src="https://typora-img-xue.oss-cn-beijing.aliyuncs.com/img/figure-13-20240323203741063.png" alt="image" style="zoom:33%;" />

