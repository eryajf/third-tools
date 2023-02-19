package main

import (
	"encoding/json"
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
)

const (
	// consulAddress = "10.6.6.17:8500"
	consulAddress = "10.6.6.66:8500"
)

func InitConsulCli() *consulapi.Client {
	config := consulapi.DefaultConfig()
	config.Address = consulAddress
	client, err := consulapi.NewClient(config)
	if err != nil {
		msg := fmt.Sprintf("init consul client failed,err: %v\n", err)
		panic(msg)
	}
	return client
}

type NgConsulKey struct {
	Env     string `json:"env"`
	Service string `json:"service"`
	Backend string `json:"backend"`
}
type NgConsulValue struct {
	Weight      int `json:"weight"`
	MaxFails    int `json:"max_fails"`
	FailTimeout int `json:"fail_timeout"`
	Down        int `json:"down"`
}

func Addkv(k NgConsulKey, v NgConsulValue) (*consulapi.WriteMeta, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	kv := InitConsulCli().KV()
	p := &consulapi.KVPair{Key: fmt.Sprintf("%s/nginx/upstreams/%s/%s", k.Env, k.Service, k.Backend), Value: []byte(b)}
	wm, err := kv.Put(p, nil)
	if err != nil {
		return nil, err
	}
	return wm, nil
}
func TestKV() {
	key := NgConsulKey{
		Env:     "prod",
		Service: "eryajf-api-upsync",
		Backend: "10.6.6.66:9902",
	}
	value := NgConsulValue{
		Weight:      10,
		MaxFails:    3,
		FailTimeout: 10,
		Down:        0,
	}
	_, err := Addkv(key, value)
	if err != nil {
		fmt.Printf("add kv failed:%v\n", err)
	}
	fmt.Println("add kv success")
}
func main() {
	// TestKV()
	TestService()
}

func (u NgConsulValue) GetNgConsulValue(weight, max_fails, fail_timeout int) *NgConsulValue {
	return &NgConsulValue{
		Weight:      weight,
		MaxFails:    max_fails,
		FailTimeout: fail_timeout,
		Down:        0,
	}
}

func DelKey() {
	kv := InitConsulCli().KV()
	pair, err := kv.Delete("upstreams/eryajf_test/10.6.6.66:9901", nil)
	if err != nil {
		fmt.Printf("del key failed,err :%v\n", err)
	}
	fmt.Println(pair.RequestTime.Hours(), pair.RequestTime.Microseconds())
}
func GetKey() {
	kv := InitConsulCli().KV()
	// Lookup the pair
	pair, _, err := kv.Get("upstreams/eryajf_test/10.6.6.66:9901", nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("KV: %v %s\n", pair.Key, pair.Value)
}
func GetAllKey() {
	kv := InitConsulCli().KV()
	pair, _, err := kv.List("prod/nginx/upstreams/eryajf-test-upsync", nil)
	// Lookup the pair
	// pair, _, err := kv.Get("upstreams/eryajf_test/10.6.6.66:9901", nil)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("KV: %v %s\n", pair.Key, pair.Value)
	for _, v := range pair {

		fmt.Println(v.Key, string(v.Value))
	}

}

func TestService() {
	// _ = ConsulRegister()
	ConsulDeRegister()
}

// 注册服务到consul
func ConsulRegister() error {
	// 创建注册到consul的服务到
	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = "ddd"
	registration.Name = "go-consul-test"
	registration.Port = 2379
	registration.Tags = []string{"go-consul-test"}
	registration.Address = "10.6.6.66"

	// 增加consul健康检查回调函数
	check := new(consulapi.AgentServiceCheck)
	check.HTTP = fmt.Sprintf("http://%s:%d", registration.Address, registration.Port)
	check.Timeout = "5s"
	check.Interval = "5s"
	// check.DeregisterCriticalServiceAfter = "30s" // 故障检查失败30s后 consul自动将注册服务删除
	registration.Check = check

	// 注册服务到consul
	return InitConsulCli().Agent().ServiceRegister(registration)
}

// 取消consul注册的服务
func ConsulDeRegister() {
	// // 创建连接consul服务配置
	// config := consulapi.DefaultConfig()
	// config.Address = "172.16.242.129:8500"
	// client, err := consulapi.NewClient(config)
	// if err != nil {
	// 	log.Fatal("consul client error : ", err)
	// }

	// client.Agent().ServiceDeregister("111")
	err := InitConsulCli().Agent().ServiceDeregister("ddd")
	if err != nil {
		fmt.Printf("deregister failed,err:%v\n", err)
	}
}
