package main

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type Config struct {
	Jenkins struct {
		User string `json:"user"`
		Pass string `json:"pass"`
	} `json:"jenkins"`
	Rancher struct {
		User string `json:"user"`
		Pass string `json:"pass"`
	} `json:"rancher"`
}

func NacosClient() {
	// 至少一个ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: "localhost",
			Port:   8848,
		},
	}

	// 创建clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         "1e096d0a-00c9-4fe2-b188-1065a44d3228", // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		LogLevel:            "debug",
	}
	// 创建动态配置客户端的另一种方式 (推荐)
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		panic(err)
	}
	//获取配置信息
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: "dev",
		Group:  "DEFAULT_GROUP"})
	if err != nil {
		fmt.Println("GetConfig err: ", err)
	}

	fmt.Println(content)

	var pro Config
	err = jsoniter.Unmarshal([]byte(content), &pro)
	if err != nil {
		fmt.Printf("json unmarshal err: %v", err)
	}
	fmt.Println(pro.Jenkins.User)

	//监听配置
	// err = configClient.ListenConfig(vo.ConfigParam{
	// 	DataId: "dev",
	// 	Group:  "DEFAULT_GROUP",
	// 	OnChange: func(namespace, group, dataId, data string) {
	// 		fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)
	// 	},
	// })
	// if err != nil {
	// 	return
	// }
	// time.Sleep(time.Second * 1000)

}
