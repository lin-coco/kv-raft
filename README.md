# kv-raft

## 前言

raft算法从开始学习到开发完成过程：

- 第一阶段：MIT6.824的实验，刚做完MapReduce，到了Raft感觉无从下手，放弃了
- 第二阶段：一段时间过后决定重新学习raft，这次变了学习方式没有做MIT，先看了大佬写的[分布式键值存储(c++)版](https://github.com/youngyangyang04/KVstorageBaseRaft-cpp)源码，看完了选举和心跳部分，感觉懵懵的，但是对raft变量和实现过程有了更深的认识
- 第三阶段：开始根据论文和之前读源码的印象自己写go语言版的raft，从选举、心跳、日志状态存储、上层状态机、客户端交互实现了，由于我是全部写完再测试，这个就太痛苦了，改了有几十个bug，但最终效果是好的，终于让我测不出bug了。

## 概述

目前已实现的内容：

- 领导选举：超时选举时间为随机150ms\~300ms，测试时为会改成15s~30s。
- 日志复制/心跳：心跳时间为固定25ms，测试时为会改成2.5s。
- 日志持久化：采用了Append Only File的方式持久化logEntry，在日志复制阶段多了个采用aof文件同步日志的选项
- 上层状态机：仅采用`map[string]string`的结构实现kv的存储，提供了get、put、del的命令
- 客户端交互：客户端发送或转发带有requestID的命令给leader节点，上层状态机会通过channel发送给raft，并创建requestID对应的channel，raft在commit并apply之后就会用此requestID对应的channel通知上层的状态机响应客户端

采用第三方库：

- [sirupsen/logrus](github.com/sirupsen/logrus) 日志组件
- [urfave/cli](github.com/urfave/cli) 终端工具
- [grpc](google.golang.org/grpc) raft通讯
- [uuid](github.com/google/uuid) requestID生成

有时间将要完善的内容

- 避免leader崩溃，命令被重复执行
- aof日志压缩
- 集群成员变更（读了论文感觉这个是不是有点难呀，我看MIT实验也没有这一部分实现）

另外我在raft层写了非常非常非常多debug日志，能完全能从日志中窥探每一时刻raft在干什么

## 运行

```bash
git clone https://github.com/lin-coco/kv-raft
```

在Goland上同时开启四个运行，一个客户端，三个raft节点，这三个raft节点，分别有不同的程序参数对应`-f server8000.json` `-f server8001.json` `-f server8002.json`

![image-20240405011101244](https://typora-img-xue.oss-cn-beijing.aliyuncs.com/img/image-20240405011101244.png)

![image-20240405011344797](https://typora-img-xue.oss-cn-beijing.aliyuncs.com/img/image-20240405011344797.png)

同时启动这四个程序，就可以在客户端自由的输入`get []` `put [] []` `del []`这三个命令了

日志存储和状态类似

![image-20240405011647883](https://typora-img-xue.oss-cn-beijing.aliyuncs.com/img/image-20240405011647883.png)





README待补充。。。