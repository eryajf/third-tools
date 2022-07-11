package main

import (
	"fmt"

	"github.com/zhaoyunxing92/dingtalk/v2"
	dingreq "github.com/zhaoyunxing92/dingtalk/v2/request"
)

const (
	// dingding
	DingTalkAppKey    = "xxxxxx"
	DingTalkAppSecret = "xxxxxxxxxx"
)

func DingClient() *dingtalk.DingTalk {
	client, err := dingtalk.NewClient(DingTalkAppKey, DingTalkAppSecret)
	if err != nil {
		fmt.Printf("init dingding client failed, err:%v\n", err)
	}
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
func GetDepts() (result []*DingTalkDept, err error) {
	rsp, err := DingClient().FetchDeptList(1, true, "zh_CN")
	if err != nil {
		fmt.Printf("get dept failed:%v\n", err)
	}
	for _, rs := range rsp.Dept {
		result = append(result, &DingTalkDept{
			Id:       rs.Id,
			Name:     rs.Name,
			Remark:   rs.Name,
			ParentId: rs.ParentId,
		})
	}
	return
}

func GetUsers() ([]*DeptDetailUserInfo, error) {
	depts, err := GetDepts()
	if err != nil {
		return nil, err
	}
	var users []*DeptDetailUserInfo
	for _, dept := range depts {
		r := dingreq.DeptDetailUserInfo{
			DeptId:   dept.Id,
			Cursor:   0,
			Size:     1,
			Language: "zh_CN",
		}
		for {
			rsp, err := DingClient().GetDeptDetailUserInfo(&r)
			if err != nil {
				return nil, err
			}
			for _, user := range rsp.Page.List {
				users = append(users, &DeptDetailUserInfo{
					UserId:               user.UserId,
					UnionId:              user.UnionId,
					Name:                 user.Name,
					Avatar:               user.Avatar,
					StateCode:            user.StateCode,
					ManagerUserId:        user.ManagerUserId,
					Mobile:               user.Mobile,
					HideMobile:           user.HideMobile,
					Telephone:            user.Telephone,
					JobNumber:            user.JobNumber,
					Title:                user.Title,
					WorkPlace:            user.WorkPlace,
					Remark:               user.Remark,
					LoginId:              user.LoginId,
					DeptIds:              user.DeptIds,
					DeptOrder:            user.DeptOrder,
					Extension:            user.Extension,
					HiredDate:            user.HiredDate,
					Active:               user.Active,
					Admin:                user.Admin,
					Boss:                 user.Boss,
					ExclusiveAccount:     user.ExclusiveAccount,
					Leader:               user.Leader,
					ExclusiveAccountType: user.ExclusiveAccountType,
					OrgEmail:             user.OrgEmail,
					Email:                user.Email,
				})
			}
			if !rsp.Page.HasMore {
				break
			}
			r.Cursor = rsp.Page.NextCursor
		}
	}
	return users, nil
}

type DingTalkDept struct {
	Id       int    `json:"dept_id"`
	Name     string `json:"name"`
	Remark   string `json:"remark"`
	ParentId int    `json:"parent_id"`
}

type DeptDetailUserInfo struct {
	UserId string `json:"userid"`
	// 员工在当前开发者企业账号范围内的唯一标识
	UnionId string `json:"unionid"`
	// 员工名称
	Name string `json:"name"`
	// 头像
	Avatar string `json:"avatar"`
	// 国际电话区号
	StateCode string `json:"state_code"`
	// 员工的直属主管
	ManagerUserId string `json:"manager_userid"`
	// 手机号码
	Mobile string `json:"mobile"`
	// 是否号码隐藏
	HideMobile bool `json:"hide_mobile"`
	// 分机号
	Telephone string `json:"telephone"`
	// 员工工号
	JobNumber string `json:"job_number"`
	// 职位
	Title string `json:"title"`
	// 办公地点
	WorkPlace string `json:"work_place"`
	// 备注
	Remark string `json:"remark"`
	// 专属帐号登录名
	LoginId string `json:"loginId"`
	// 所属部门ID列表
	DeptIds []int `json:"dept_id_list"`
	// 员工在部门中的排序
	DeptOrder int `json:"dept_order"`
	// 员工在对应的部门中的排序
	Extension string `json:"extension"`
	// 入职时间
	HiredDate int `json:"hired_date"`
	// 是否激活了钉钉
	Active bool `json:"active"`
	//是否为企业的管理员：
	//
	//true：是
	//
	//false：不是
	Admin bool `json:"admin"`
	// 是否为企业的老板
	Boss bool `json:"boss"`
	// 是否专属帐号
	ExclusiveAccount bool `json:"exclusive_account"`
	// 是否是部门的主管
	Leader bool `json:"leader"`
	//专属帐号类型：
	//
	//sso：企业自建专属帐号
	//
	//dingtalk：钉钉自建专属帐号
	ExclusiveAccountType string `json:"exclusive_account_type"`
	//员工的企业邮箱
	//
	//如果员工的企业邮箱没有开通，返回信息中不包含该数据
	OrgEmail string `json:"org_email"`
	//员工邮箱
	//
	//企业内部应用如果没有返回该字段，需要检查当前应用通讯录权限中邮箱等个人信息权限是否开启
	//
	//员工信息面板中有邮箱字段值才返回该字段
	//
	//第三方企业应用不返回该参数
	Email string `json:"email"`
}
