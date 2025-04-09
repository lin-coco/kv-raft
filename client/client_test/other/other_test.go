package main

import (
	"encoding/json"
	"kv-raft/client"
	"testing"
)

var kvClient *client.Client

func init() {
	var err error
	servers := []string{"127.0.0.1:8000", "127.0.0.1:8001", "127.0.0.1:8002"}
	apiKey := "9jfAuiOc9L0dNXSJ36Cy5jh3ICnn9LOO"
	kvClient, err = client.NewClient(servers, apiKey)
	if err != nil {
		panic(err)
	}
}

/*
acl操作
*/

func TestAclGetMy(t *testing.T) {
	resp, err := kvClient.AclGet()
	if err != nil {
		t.Fatal("TestAclGet failed",err)
	}
	bytes, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("TestAclGet result: \n%s\n", string(bytes))
}

func TestAclGetOther(t *testing.T) {
	resp, err := kvClient.AclGet("lincoco")
	if err != nil {
		t.Error("TestAclGetOther failed",err)
	}
	bytes, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("TestAclGetOther result: \n%s\n", string(bytes))
}

func TestAclOtherCmd(t *testing.T) {
	_, err := kvClient.AclAll()
	if err != nil {
		t.Error("TestAclOtherCmd failed",err)
	}
	_, err = kvClient.AclAdd("test", "test:*,rw")
	if err != nil {
		t.Error("TestAclOtherCmd failed",err)
	}
	_, err = kvClient.AclUpdate("test", "testing:*,rw")
	if err != nil {
		t.Error("TestAclOtherCmd failed",err)
	}
	_, err = kvClient.AclDel("test")
	if err != nil {
		t.Error("TestAclOtherCmd failed",err)
	}
}


/*
键值操作
*/

func TestPut(t *testing.T) {
	resp, err := kvClient.Put("xys:computer", "Mac")
	if err != nil {
		t.Fatal("TestPut failed",err)
	}
	bytes, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("TestPut result: \n%s\n", string(bytes))
}
func TestPutNoPermission(t *testing.T) {
	resp, err := kvClient.Put("mango", "mango")
	if err != nil {
		t.Fatal("TestPutNoPermission failed",err)
	}
	bytes, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("TestPutNoPermission result: \n%s\n", string(bytes))
}

func TestGet(t *testing.T) {
	resp, err := kvClient.Get("xys:computer")
	if err != nil {
		t.Fatal("TestGet failed",err)
	}
	bytes, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("TestGet result: \n%s\n", string(bytes))
}

func TestGetNoPermission(t *testing.T) {
	resp, err := kvClient.Get("mango")
	if err != nil {
		t.Fatal("TestGetNoPermission failed",err)
	}
	bytes, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("TestGetNoPermission result: \n%s\n", string(bytes))
}

func TestExist(t *testing.T) {
	resp, err := kvClient.Exist("xys:computer")
	if err != nil {
		t.Fatal("TestExist failed",err)
	}
	bytes, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("TestExist result: \n%s\n", string(bytes))
}
func TestExistNoPermission(t *testing.T) {
	resp, err := kvClient.Exist("mango")
	if err != nil {
		t.Fatal("TestExistNoPermission failed",err)
	}
	bytes, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("TestExistNoPermission result: \n%s\n", string(bytes))
}
func TestRename(t *testing.T) {
	resp, err := kvClient.Rename("xys:computer", "xys:machine")
	if err != nil {
		t.Fatal("TestRename failed",err)
	}
	bytes, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("TestRename result: \n%s\n", string(bytes))
}
func TestRenameNoPermission(t *testing.T) {
	resp, err := kvClient.Rename("mango", "xys:mango")
	if err != nil {
		t.Fatal("TestRenameNoPermission failed",err)
	}
	bytes, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("TestRenameNoPermission result: \n%s\n", string(bytes))
}

func TestDel(t *testing.T) {
	resp, err := kvClient.Del("xys:computer")
	if err != nil {
		t.Fatal("TestDel failed",err)
	}
	bytes, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("TestDel result: \n%s\n", string(bytes))

}
func TestDelNoPermission(t *testing.T) {
	resp, err := kvClient.Del("mango")
	if err != nil {
		t.Fatal("TestDelNoPermission failed",err)
	}
	bytes, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("TestDelNoPermission result: \n%s\n", string(bytes))
}

/*
获取节点信息操作
*/

func TestNodeAll(t *testing.T) {
	resp, err := kvClient.NodeAll()
	if err != nil {
		t.Fatal("TestNodeAll failed",err)
	}
	bytes, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("TestNodeAll result: \n%s\n", string(bytes))
}
func TestNodeGet(t *testing.T) {
	resp, err := kvClient.NodeGet()
	if err != nil {
		t.Fatal("TestNodeAll failed",err)
	}
	bytes, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("TestNodeAll result: \n%s\n", string(bytes))
}
func TestNodeLeader(t *testing.T) {
	resp, err := kvClient.NodeLeader()
	if err != nil {
		t.Fatal("TestNodeLeader failed",err)
	}
	bytes, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("TestNodeLeader result: \n%s\n", string(bytes))
}