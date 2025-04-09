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
	apiKey := "qGCEUffJMrKQqeUn7x28GYcIcZj07GBy"
	kvClient, err = client.NewClient(servers, apiKey)
	if err != nil {
		panic(err)
	}
}
/*
键值操作
*/

func TestPut(t *testing.T) {
	resp, err := kvClient.Put("apple", "1")
	if err != nil {
		t.Fatal("TestPut failed",err)
	}
	bytes, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("TestPut result: \n%s\n", string(bytes))
}

func TestPutBatch(t *testing.T) {
	resp, err := kvClient.Put("pear", "2", "banana", "3", "mango", "4")
	if err != nil {
		t.Fatal("TestPutBatch failed",err)
	}
	bytes, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("TestPutBatch result: \n%s\n\n", string(bytes))
}

func TestGet(t *testing.T) {
	resp, err := kvClient.Get("apple")
	if err != nil {
		t.Fatal("TestGet failed",err)
	}
	bytes, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("TestGet result: \n%s\n", string(bytes))
}

func TestGetBatch(t *testing.T) {
	resp, err := kvClient.Get("pear", "banana", "mango")
	if err != nil {
		t.Fatal("TestGet failed",err)
	}
	bytes, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("TestGet result: \n%s\n", string(bytes))
}

func TestDel(t *testing.T) {
	resp, err := kvClient.Del("apple")
	if err != nil {
		t.Fatal("TestDel failed",err)
	}
	bytes, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("TestDel result: \n%s\n", string(bytes))

}
func TestDelBatch(t *testing.T) {
	resp, err := kvClient.Del("pear", "banana")
	if err != nil {
		t.Fatal("TestDelBatch failed",err)
	}
	bytes, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("TestDelBatch result: \n%s\n", string(bytes))

}
func TestExistKey(t *testing.T) {
	resp, err := kvClient.Exist("mango")
	if err != nil {
		t.Fatal("TestExist failed",err)
	}
	bytes, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("TestExist result: \n%s\n", string(bytes))
}
func TestNotExistKey(t *testing.T) {
	resp, err := kvClient.Exist("apple")
	if err != nil {
		t.Fatal("TestNotExistKey failed",err)
	}
	bytes, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("TestNotExistKey result: \n%s\n", string(bytes))
}
func TestRenameExistKey(t *testing.T) {
	resp, err := kvClient.Rename("mango", "mongo")
	if err != nil {
		t.Fatal("TestRenameExistKey failed",err)
	}
	bytes, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("TestRenameExistKey result: \n%s\n", string(bytes))
}
func TestRenameNotExistKey(t *testing.T) {
	resp, err := kvClient.Rename("apple", "apples")
	if err != nil {
		t.Fatal("TestRenameNotExistKey failed",err)
	}
	bytes, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("TestRenameNotExistKey result: \n%s\n", string(bytes))
}

/*
导入导出操作
*/
func TestImportCSV(t *testing.T) {
	resp, err := kvClient.Import("import1.csv")
	if err != nil {
		t.Fatal("TestImportCSV failed",err)
	}
	bytes, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("TestImportCSV result: \n%s\n", string(bytes))
}
func TestImportJSON(t *testing.T) {
	resp, err := kvClient.Import("import2.json")
	if err != nil {
		t.Fatal("TestImportJSON failed",err)
	}
	bytes, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("TestImportJSON result: \n%s\n", string(bytes))
}
func TestExportCSV(t *testing.T) {
	resp, err := kvClient.Export("export.csv")
	if err != nil {
		t.Fatal("TestExportCSV failed",err)
	}
	bytes, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("TestExportCSV result: \n%s\n", string(bytes))
}
func TestExportJSON(t *testing.T) {
	resp, err := kvClient.Export("export.json")
	if err != nil {
		t.Fatal("TestExportJSON failed",err)
	}
	bytes, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("TestExportJSON result: \n%s\n", string(bytes))
}

/*
acl操作
*/

func TestAclAdd(t *testing.T) {
	resp, err := kvClient.AclAdd("xueyeshang", "xys:*,rw")
	if err != nil {
		t.Fatal("TestAclAdd failed",err)
	}
	bytes, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("TestAclAdd result: \n%s\n", string(bytes))
	resp, err = kvClient.AclAdd("lincoco", "coco:*,rw")
	if err != nil {
		t.Fatal("TestAclAdd failed",err)
	}
	bytes, _ = json.MarshalIndent(resp, "", "  ")
	t.Logf("TestAclAdd result: \n%s\n", string(bytes))
}
func TestAclAll(t *testing.T) {
	resp, err := kvClient.AclAll()
	if err != nil {
		t.Fatal("TestAclAdd failed",err)
	}
	bytes, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("TestAclAdd result: \n%s\n", string(bytes))
}
func TestAclDel(t *testing.T) {
	resp, err := kvClient.AclDel("lincoco")
	if err != nil {
		t.Fatal("TestAclDel failed",err)
	}
	bytes, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("TestAclDel result: \n%s\n", string(bytes))
}
func TestAclGet(t *testing.T) {
	resp, err := kvClient.AclGet("xueyeshang")
	if err != nil {
		t.Fatal("TestAclGet failed",err)
	}
	bytes, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("TestAclGet result: \n%s\n", string(bytes))
}
func TestAclGetNotExist(t *testing.T) {
	resp, err := kvClient.AclGet("lincoco")
	if err != nil {
		t.Fatal("TestAclGetNotExist failed",err)
	}
	bytes, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("TestAclGetNotExist result: \n%s\n", string(bytes))
}
func TestAclUpdateExist(t *testing.T) {
	resp, err := kvClient.AclUpdate("xueyeshang", "xys:*,rw", "xueyeshang,rw")
	if err != nil {
		t.Fatal("TestAclUpdateExist failed",err)
	}
	bytes, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("TestAclUpdateExist result: \n%s\n", string(bytes))
}
func TestAclUpdateNotExist(t *testing.T) {
	resp, err := kvClient.AclUpdate("lincoco", "coco:*,rw", "lincoco,rw")
	if err != nil {
		t.Fatal("TestAclUpdateExist failed",err)
	}
	bytes, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("TestAclUpdateExist result: \n%s\n", string(bytes))
}



/*
获取节点信息操作
*/


// 获取节点信息，演示宕机效果、重启效果
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

// 获取节点信息，演示宕机效果、重启效果
func TestNodeLeader(t *testing.T) {
	resp, err := kvClient.NodeLeader()
	if err != nil {
		t.Fatal("TestNodeLeader failed",err)
	}
	bytes, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("TestNodeLeader result: \n%s\n", string(bytes))
}


