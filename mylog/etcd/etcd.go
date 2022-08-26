package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"test/logagent/mylog/config"
	"test/logagent/mylog/tail"

	"github.com/coreos/etcd/clientv3"
)

var (
	client  *clientv3.Client
	logdata config.LogEntryConf
)

func Init(address []string, timeout int) (err error) {
	client, err = clientv3.New(clientv3.Config{
		Endpoints:   address,
		DialTimeout: time.Duration(timeout) * time.Second,
	})
	if err != nil {
		fmt.Println("connect to etcd failed,err:\n", err)
		return
	}
	fmt.Println("connect to etcd success!")
	return
}

func GetConf(key string) (logconf config.LogEntryConf, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	var resp *clientv3.GetResponse
	resp, err = client.Get(ctx, key)
	cancel()
	if err != nil {
		fmt.Println("get from etcd failed,err:", err)
		return
	}
	for _, ev := range resp.Kvs {
		fmt.Println("values:", ev.Value)
		err = json.Unmarshal(ev.Value, &logconf)
		if err != nil {
			fmt.Println("Unmarshal json failed:", err)
			return
		}
	}
	return
}

func WatchConf(topic string) {
	rch := client.Watch(context.Background(), topic)
	channel := tail.Get_chan()
	for wresp := range rch {
		for _, ev := range wresp.Events {
			err := json.Unmarshal(ev.Kv.Value, &logdata)
			if err != nil {
				fmt.Println("Update conf failed,err:", err)
				return
			}
			fmt.Println("update config success:", ev.Kv.Value)
			channel <- logdata
		}
	}
}
