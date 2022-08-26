package config

import (
	"context"

	"github.com/hpcloud/tail"
)

type AppConf struct {
	KafkaConf `ini:"kafka"`
	EtcdConf  `ini:"etcd"`
	ESConf    `ini:"es"`
}

type KafkaConf struct {
	Address  string `ini:"address"`
	Max_size int    `ini:"chan_max_size"`
}

type EtcdConf struct {
	Address string `ini:"address"`
	Timeout int    `ini:"timeout"`
	Log_key string `ini:"log_key"`
}

type ESConf struct {
	Address  string `ini:"address"`
	Max_size int    `ini:"size"`
}

type LogConf struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

type LogEntryConf []*LogConf

type TailTask struct {
	Path     string
	Topic    string
	Instance *tail.Tail
	Ctx      context.Context
	CancelF  context.CancelFunc
}

type LogData struct {
	Topic string
	Data  string
}
