package main

import (
	"fmt"
	"sync"
	"test/logagent/mylog/config"
	"test/logagent/mylog/es"
	"test/logagent/mylog/etcd"
	"test/logagent/mylog/tail"
	"test/logagent/mylog/utils"

	"test/logagent/mylog/kafka"

	"gopkg.in/ini.v1"
)

var wg sync.WaitGroup

func main() {
	var cfg config.AppConf
	err := ini.MapTo(&cfg, "./config/config.ini")
	if err != nil {
		fmt.Println("Decode Map failed!", err)
	}
	err = kafka.Init([]string{cfg.KafkaConf.Address}, cfg.KafkaConf.Max_size)
	if err != nil {
		fmt.Println("init kafka failed", err)
		return
	}
	fmt.Println("init kafka success!")

	err = etcd.Init([]string{cfg.EtcdConf.Address}, cfg.EtcdConf.Timeout)
	var path config.LogEntryConf
	ipStr, err := utils.GetOutboundIP()
	if err != nil {
		panic(err)
	}
	etcdConfKey := fmt.Sprintf(cfg.Log_key, ipStr)
	fmt.Printf("etcdConfKey:%s\n", etcdConfKey)
	path, err = etcd.GetConf(etcdConfKey)
	fmt.Println(path)
	if err != nil {
		return
	}
	tail.Init(path)
	es.Init(cfg.ESConf.Address, cfg.ESConf.Max_size)
	for index, value := range path {
		fmt.Printf("index:%v value:%v topic:%v\n", index, value, value.Topic)
		kafka.Consumer(value.Topic)
	}
	wg.Add(1)
	etcd.WatchConf(etcdConfKey)
	wg.Done()
}
