# 功能测试

在三个终端中分别源码运行raft节点

```bash
go run server/main.go -f server8000.json
go run server/main.go -f server8001.json
go run server/main.go -f server8002.json
```

运行指定用例

```bash
cd client/testing/root
go test -v -run [TestName] 运行指定用例
```

## root账号测试

## other账号测试
