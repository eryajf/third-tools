package main

import (
	"context"
	"fmt"

	"github.com/chyroc/lark"
)

const (
	AppId     = "xxxxxxx"
	AppSecret = "xxxxxxxxxxxx"
)

func FeishuClient() *lark.Lark {
	return lark.New(lark.WithAppCredential(AppId, AppSecret))
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

	ids, err := GetLeaveUserIds()
	if err != nil {
		fmt.Printf("获取离职人员ID列表失败: %v\n", err)
	}
	for _, id := range ids {
		fmt.Println(id)
	}
}

func GetDepts() (depts []*lark.GetDepartmentListRespItem, err error) {
	var (
		fetchChild bool  = true
		pageSize   int64 = 50
	)

	req := lark.GetDepartmentListReq{
		FetchChild:   &fetchChild,
		PageSize:     &pageSize,
		DepartmentID: "0",
	}

	for {
		res, _, err := FeishuClient().Contact.GetDepartmentList(context.TODO(), &req)
		if err != nil {
			fmt.Printf("GetDepartmentList error: %v\n", err)
		}
		depts = append(depts, res.Items...)
		if !res.HasMore {
			break
		}
		req.PageToken = &res.PageToken
	}
	return
}

func GetUsers() (users []*lark.GetUserListRespItem, err error) {
	var (
		pageSize int64 = 50
	)
	depts, err := GetDepts()
	if err != nil {
		fmt.Printf(" get all depts failed, err:%v\n", err)
	}
	for _, dept := range depts {

		req := lark.GetUserListReq{
			PageSize:     &pageSize,
			PageToken:    new(string),
			DepartmentID: dept.OpenDepartmentID,
		}

		for {
			res, _, err := FeishuClient().Contact.GetUserList(context.Background(), &req)
			if err != nil {
				return nil, err
			}
			users = append(users, res.Items...)
			if !res.HasMore {
				break
			}
			req.PageToken = &res.PageToken
		}
	}
	return
}

// 官方文档： https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/ehr/ehr-v1/employee/list
// GetLeaveUserIds 获取离职人员ID列表
func GetLeaveUserIds() ([]string, error) {
	var ids []string
	users, _, err := FeishuClient().EHR.GetEHREmployeeList(context.TODO(), &lark.GetEHREmployeeListReq{
		Status:     []int64{5},
		UserIDType: lark.IDTypePtr(lark.IDTypeUnionID), // 只查询unionID
	})
	if err != nil {
		return nil, err
	}
	for _, user := range users.Items {
		ids = append(ids, user.UserID)
	}
	return ids, nil
}
