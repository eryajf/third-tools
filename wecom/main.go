package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/wenerme/go-wecom/wecom"
)

const (
	CorpID     = "wwc6656f03c222fd11"
	AgentID    = 1000003
	CorpSecret = "vhgdQfGS7pM46O40oSK8GBFTcDH_N3gkN5ZJJfOMpog"
)

func WeComClient() *wecom.Client {
	// token store - 默认内存 Map - 可以使用数据库实现
	store := &wecom.SyncMapStore{}
	// 加载缓存 - 复用之前的 Token
	if bytes, err := os.ReadFile("wecom-cache.json"); err == nil {
		_ = store.Restore(bytes)
	}
	// 当 Token 变化时生成缓存文件
	store.OnChange = func(s *wecom.SyncMapStore) {
		_ = os.WriteFile("wecom-cache.json", s.Dump(), 0o600)
	}

	client := wecom.NewClient(wecom.Conf{
		CorpID:     CorpID,
		AgentID:    AgentID,
		CorpSecret: CorpSecret,
		// 不配置默认使用 内存缓存
		TokenProvider: &wecom.TokenCache{
			Store: store,
		},
	})
	return client
}

func main() {
	groups, err := GetDepts()
	if err != nil {
		fmt.Printf("get all group failed: %v\n", err)
	}
	for _, group := range groups {
		fmt.Println("分组信息:", group)
	}
	users, err := GetUsers()
	if err != nil {
		fmt.Printf("get all user failed: %v\n", err)
	}
	for _, user := range users {
		fmt.Println("分组信息:", user)
	}
}

func GetDepts() ([]wecom.ListDepartmentResponseItem, error) {
	depts, err := WeComClient().ListDepartment(
		&wecom.ListDepartmentRequest{},
	)
	if err != nil {
		return nil, err
	}
	return depts.Department, nil
}

func GetUsers() ([]wecom.ListUserResponseItem, error) {
	depts, err := GetDepts()
	if err != nil {
		return nil, err
	}
	var us []wecom.ListUserResponseItem
	for _, dept := range depts {
		users, err := WeComClient().ListUser(
			&wecom.ListUserRequest{
				DepartmentID: strconv.Itoa(dept.ID),
				FetchChild:   "1",
			},
		)
		if err != nil {
			return nil, err
		}
		us = append(us, users.UserList...)
	}
	return us, nil
}
