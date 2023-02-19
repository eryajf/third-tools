package main

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/go-resty/resty/v2"
)

func InitHBCli() *resty.Client {
	return resty.New().SetTimeout(3*time.Second).SetBasicAuth("admin", "123465")
}

type Tags struct {
	Digest        string    `json:"digest"`
	Name          string    `json:"name"`
	Size          int       `json:"size"`
	Architecture  string    `json:"architecture"`
	Os            string    `json:"os"`
	OsVersion     string    `json:"os.version"`
	DockerVersion string    `json:"docker_version"`
	Author        string    `json:"author"`
	Created       time.Time `json:"created"`
	Config        struct {
		Labels interface{} `json:"labels"`
	} `json:"config"`
	Signature interface{}   `json:"signature"`
	Labels    []interface{} `json:"labels"`
	PushTime  time.Time     `json:"push_time"`
	PullTime  time.Time     `json:"pull_time"`
}

type tags []Tags

func (s tags) Len() int {
	return len(s)
}
func (s tags) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type ByTime struct {
	tags
}

func (b ByTime) Less(i, j int) bool {
	return b.tags[i].Created.Before(b.tags[j].Created)
}

func GetProTags(url string) ([]Tags, error) {
	var data []Tags
	resp, err := InitHBCli().R().Get(url)
	if err != nil {
		fmt.Printf("get err:%v\n", err)
		return nil, err
	}
	json.Unmarshal(resp.Body(), &data)
	return data, nil
}

func main() {
	url := "https://reg.eryajf.net/api/repositories/multienv/admin/tags"
	a, err := GetProTags(url)
	if err != nil {
		fmt.Printf("get err:%v\n", err)
	}
	fmt.Println(len(a))
	sort.Sort(ByTime{a})
	for _, v := range a {
		size := v.Size / 1000 / 1000
		fmt.Println(v.Name, size, v.DockerVersion, v.Created)
	}
}
func SlicePage(page, pageSize, nums int) (sliceStart, sliceEnd int) {
	if page < 0 {
		page = 1
	}

	if pageSize < 0 {
		pageSize = 20
	}

	if pageSize > nums {
		return 0, nums
	}

	// 总页数
	pageCount := int(math.Ceil(float64(nums) / float64(pageSize)))
	if page > pageCount {
		return 0, 0
	}
	sliceStart = (page - 1) * pageSize
	sliceEnd = sliceStart + pageSize

	if sliceEnd > nums {
		sliceEnd = nums
	}
	return sliceStart, sliceEnd
}
