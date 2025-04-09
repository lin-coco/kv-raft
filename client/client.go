package client

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"kv-raft/raft"
	"kv-raft/server/command"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Client struct {
	Servers []string
	ApiKey string
	leaderAddr string
	httpClient *http.Client
}

func NewClient(servers []string, apiKey string) (*Client,error) {
	if len(servers) == 0 || apiKey == "" {
		return nil,errors.New("must request servers and apikey")
	}
	return &Client{
		Servers: servers,
		ApiKey: apiKey,
		leaderAddr: servers[0],
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}, nil
}

func (c *Client) Get(subList ...string) (command.GetResponse, error) {
	var resp command.GetResponse
	userCommand := "get " + strings.Join(subList, " ")
	result, err := c.SendCommand(userCommand)
	if err != nil {
		return resp, err
	}
	json.Unmarshal([]byte(result), &resp)
	return resp, nil
}

func (c *Client) Del(subList ...string) (command.DelResponse, error) {
	var resp command.DelResponse
	userCommand := "del " + strings.Join(subList, " ")
	result, err := c.SendCommand(userCommand)
	if err != nil {
		return resp, err
	}
	json.Unmarshal([]byte(result), &resp)
	return resp, nil
}

func (c *Client) Put(subList ...string) (command.PutResponse, error) {
	var resp command.PutResponse
	userCommand := "put " + strings.Join(subList, " ")
	result, err := c.SendCommand(userCommand)
	if err != nil {
		return resp, err
	}
	json.Unmarshal([]byte(result), &resp)
	return resp, nil
}

func (c *Client) Exist(subList ...string) (command.ExistResponse, error) {
	var resp command.ExistResponse
	userCommand := "exist " + strings.Join(subList, " ")
	result, err := c.SendCommand(userCommand)
	if err != nil {
		return resp, err
	}
	json.Unmarshal([]byte(result), &resp)
	return resp, nil
}

func (c *Client) Rename(subList ...string) (command.RenameResponse, error) {
	var resp command.RenameResponse
	userCommand := "rename " + strings.Join(subList, " ")
	result, err := c.SendCommand(userCommand)
	if err != nil {
		return resp, err
	}
	json.Unmarshal([]byte(result), &resp)
	return resp, nil
}

func (c *Client) AclAdd(subList ...string) (command.AclAddResponse, error) {
	var resp command.AclAddResponse
	userCommand := "acl-add " + strings.Join(subList, " ")
	result, err := c.SendCommand(userCommand)
	if err != nil {
		return resp, err
	}
	json.Unmarshal([]byte(result), &resp)
	return resp, nil
}

func (c *Client) AclUpdate(subList ...string) (command.AclUpdateResponse, error) {
	var resp command.AclUpdateResponse
	userCommand := "acl-update " + strings.Join(subList, " ")
	result, err := c.SendCommand(userCommand)
	if err != nil {
		return resp, err
	}
	json.Unmarshal([]byte(result), &resp)
	return resp, nil
}

func (c *Client) AclDel(subList ...string) (command.AclDelResponse, error) {
	var resp command.AclDelResponse
	userCommand := "acl-del " + strings.Join(subList, " ")
	result, err := c.SendCommand(userCommand)
	if err != nil {
		return resp, err
	}
	json.Unmarshal([]byte(result), &resp)
	return resp, nil
}

func (c *Client) AclGet(subList ...string) (command.AclGetResponse, error) {
	var resp command.AclGetResponse
	userCommand := "acl-get " + strings.Join(subList, " ")
	result, err := c.SendCommand(userCommand)
	if err != nil {
		return resp, err
	}
	json.Unmarshal([]byte(result), &resp)
	return resp, nil
}

func (c *Client) AclAll() (command.AclAllResponse, error) {
	var resp command.AclAllResponse
	userCommand := "acl-all"
	result, err := c.SendCommand(userCommand)
	if err != nil {
		return resp, err
	}
	json.Unmarshal([]byte(result), &resp)
	return resp, nil
}

func (c *Client) NodeAll() ([]raft.NodeInfo, error) {
	var resp []raft.NodeInfo
	userCommand := "node-all"
	result, err := c.SendCommand(userCommand)
	if err != nil {
		return resp, err
	}
	json.Unmarshal([]byte(result), &resp)
	return resp, nil
}

func (c *Client) NodeLeader() (raft.NodeInfo, error) {
	var resp raft.NodeInfo
	userCommand := "node-leader"
	result, err := c.SendCommand(userCommand)
	if err != nil {
		return resp, err
	}
	json.Unmarshal([]byte(result), &resp)
	return resp, nil
}

func (c *Client) NodeGet() (raft.NodeInfo, error) {
	var resp raft.NodeInfo
	userCommand := "node-get"
	result, err := c.SendCommand(userCommand)
	if err != nil {
		return resp, err
	}
	json.Unmarshal([]byte(result), &resp)
	return resp, nil
}

type ImportResp struct {
	ImportNums int `json:"import_nums"`
}

/*
csv: key,value
json: {"key1": "value1","key2": "value2"}
*/
func (c *Client) Import(f string) (ImportResp, error) {
	var resp ImportResp
	// 读取文件
	file, err := os.Open(f)
	if err != nil {
		return resp, err
	}
	defer file.Close()
	// 解析发送
	if strings.HasSuffix(f,".csv") {
		reader := csv.NewReader(file)
		// 跳过首行标题
		if _, err := reader.Read(); err != nil {
			return resp, fmt.Errorf("no records err:%v", err)
		}
		// 读取key,value
		var end bool
		for !end {
			// 每读取10行就发送
			subList := make([]string, 0, 20)
			for i := 0;i < 10;i++ {
				record, err := reader.Read()
				if err != nil {
					if err == io.EOF { // 文件结束
						end = true
						err = nil
						break
					}
					return resp, err
				}
				if len(record) != 2 {
					continue
				}
				subList = append(subList, record[0], record[1])
			}
			if len(subList) == 0 {
				continue
			}
			_, err := c.Put(subList...)
			if err != nil {
				return resp, err
			}
			resp.ImportNums += len(subList) / 2
		}
	} else if strings.HasSuffix(f,".json") {
		bytes, err := io.ReadAll(file)
		if err != nil {
			return resp, err
		}
		var data map[string]string
		if err := json.Unmarshal(bytes, &data); err != nil {
			return resp, err
		}
		for k,v := range data {
			_, err := c.Put(k,v)
			if err != nil {
				return resp, err
			}
			resp.ImportNums++
		}
	} else {
		return resp, errors.New("不支持的文件格式，目前支持.json .csv")
	}
	return resp, nil
}

type ExportResp struct {
	ExportNums int `json:"export_nums"`
}

// 导出所有键
func (c *Client) Export(f string) (ExportResp, error) {
	var resp ExportResp
	// 校验文件类型
	if !strings.HasSuffix(f,".csv") && !strings.HasSuffix(f,".json") {
		i := strings.Index(f, ".")
		return resp, fmt.Errorf("不支持存储此文件格式: %s", string(f[i+1:]))
	}
	// 创建或打开文件:使用 os.Create 自动处理，存在则清空，不存在则创建
	file, err := os.Create(f)
	if err != nil {
		return resp, err
	}
	defer file.Close()
	// 执行命令
	result, err := c.SendCommand("keys")
	if err != nil {
		return resp, err
	}
	var keysResp command.KeysResponse
	json.Unmarshal([]byte(result), &keysResp)
	data := make(map[string]string)
	for i := 0;i < len(keysResp.Kvs);i++ {
		data[keysResp.Kvs[i].Key] = keysResp.Kvs[i].Value
	}
	if strings.HasSuffix(f,".csv") {
		writer := csv.NewWriter(file)
		defer writer.Flush()
		writer.Write([]string{"key", "value"})
		for key,value := range data {
			writer.Write([]string{key, value})
		}
	} else if strings.HasSuffix(f, ".json") {
		bytes,_ := json.MarshalIndent(data, "", "  ")
		file.WriteString(string(bytes))
	}
	resp.ExportNums = len(data)
	return resp, nil
}

func (c *Client) SendCommand(userCommand string) (string,error) {
	requestId := uuid.New().String()
	path := "receive"
	if strings.HasPrefix(userCommand, "debug-allkey") {
		path = "debug-allkey"
	} else if strings.HasPrefix(userCommand, "debug") {
		path = "debug"
		userCommand  = userCommand[6:]
	}
	// 尝试连接到一个节点
	// log.Printf("current leader: %s\n", c.leaderAddr)
	req := c.newRequest(c.leaderAddr, path, requestId, userCommand)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Printf("ping pong failed: %s", c.leaderAddr)
		// 拒绝连接，可能宕机，尝试其他地址
		var success bool
		resp, success = c.tryOtherAddr(path, requestId, userCommand)
		if !success {
			return "",fmt.Errorf("所有地址都发送失败 %v", err)
		}
	}
	// 可能要转发到leader节点
	if resp.StatusCode == http.StatusSeeOther {
		log.Printf("leader node %s change to %s\n", c.leaderAddr, resp.Header.Get("leader-id"))
		c.leaderAddr = resp.Header.Get("leader-id")
		req = c.newRequest(c.leaderAddr, path, requestId, userCommand)
		resp, err = c.httpClient.Do(req)
		if err != nil {
			return "",fmt.Errorf("leader节点连接失败 %v", err)
		}
	}
	// leader节点的响应处理
	if resp.StatusCode == http.StatusBadRequest {
		// 命令错误
		return "",fmt.Errorf("命令错误或请求id重复 %v", http.StatusBadRequest)
	} else if resp.StatusCode == http.StatusForbidden {
		// 禁止访问
		return "",fmt.Errorf("命令或键无权限访问 %v", http.StatusForbidden)
	} else if resp.StatusCode == http.StatusUnauthorized {
		// 认证失败
		return "",fmt.Errorf("apikey无效或不存在 %v", http.StatusUnauthorized)
	} else if resp.StatusCode == http.StatusInternalServerError {
		// 服务端错误
		return "",fmt.Errorf("服务端错误，请调试服务端 %v: %v", http.StatusInternalServerError, resp.Header.Get("message"))
	} else if resp.StatusCode == http.StatusOK {
		// 成功响应
		body, _ := io.ReadAll(resp.Body)
		return strings.TrimSpace(string(body)), nil
	} else {
		// 未知错误
		return "",fmt.Errorf("未知错误 %v", resp.StatusCode)
	}
}

func (c *Client) tryOtherAddr(path string,requestId string, userCommand string) (*http.Response, bool) {
	for i := 0;i < len(c.Servers);i++ {
		req := c.newRequest(c.Servers[i], path, requestId, userCommand)
		resp, err := c.httpClient.Do(req)
		if err == nil {
			log.Printf("ping pong:%s", c.Servers[i])
			return resp, true
		}
	}
	return nil, false
}

func (c *Client) newRequest(addr string, path string,requestId string, userCommand string) *http.Request {
	req, _ := http.NewRequest("POST", "http://"+addr+"/" + path, bytes.NewBuffer([]byte(userCommand)))
	req.Header.Set("Content-Type", "application/stream+json")
	req.Header.Set("Authorization", c.ApiKey)
	req.Header.Set("kv-raft-request-id", requestId)
	return req
}

