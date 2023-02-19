package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"go.etcd.io/etcd/clientv3"
)

func main() {
	addr := ""
	AllBack(addr)
}

func AllBack(addr string) {
	// 可以先获取所有的key
	fmt.Printf("env:%s\n", addr)
	a, err := Get(addr, "/")
	if err != nil {
		fmt.Printf("get from etcd failed, err:%v\n", err)
	}
	for k, v := range a {
		// fmt.Printf("key:%s,value:%s\n", k, v)
		WriteToFile("/backup/"+"/"+k, []byte(v))
	}
}

func Put(Address, path, value string) error {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{Address},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return err
	}
	defer cli.Close()

	// put
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = cli.Put(ctx, path, value)
	cancel()
	if err != nil {
		fmt.Printf("put to etcd failed, err:%v\n", err)
		return err
	}
	return nil
}

func Delete(Address, path string) error {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{Address},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return err
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	_, err = cli.Delete(ctx, path, clientv3.WithPrefix())
	cancel()
	if err != nil {
		return err
	}
	return nil
}

func Get(Address, prefix string) (map[string]string, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{Address},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	defer cli.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	// 如果后边的超时时间设置的比较短，而获取的key数量又比较多的时候，就可能会报如下错误
	// {"level":"warn","ts":"2021-11-25T09:53:08.726+0800","caller":"clientv3/retry_interceptor.go:61","msg":"retrying of unary invoker failed","target":"endpoint://client-161a587c-ca69-48ab-82cc-aa2e7d4c912c/10.6.6.66:2379","attempt":0,"error":"rpc error: code = DeadlineExceeded desc = context deadline exceeded"}
	resp, err := cli.Get(ctx, prefix, clientv3.WithPrefix())
	defer cancel()
	if err != nil {
		fmt.Printf("get from etcd failed, err:%v\n", err)
		return nil, err
	}
	kvs := make(map[string]string)

	for _, ev := range resp.Kvs {
		if string(ev.Value) != "" {
			kvs[string(ev.Key)] = string(ev.Value)
		}
	}

	return kvs, nil
}

func WriteToFile(path string, data []byte) error {
	pa := "./muban-config" + path
	tmp := strings.Split(pa, "/")
	if len(tmp) > 0 {
		tmp = tmp[:len(tmp)-1]
	}
	err := os.MkdirAll(strings.Join(tmp, "/"), os.ModePerm)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(pa, data, 0755)
	if err != nil {
		return err
	}
	return nil
}
