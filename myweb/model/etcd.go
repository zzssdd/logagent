package model

import (
	"context"
	"encoding/json"
	"fmt"
	"test/logagent/myweb/conf"
	"test/referen/log/etcd"
	"time"

	"github.com/astaxie/beego/logs"
)

func SetLogConfToEtcd(etcdKey string, info *conf.EtcdData) (err error) {

	var logConfArr []*etcd.LogEntry
	logConfArr = append(logConfArr, &etcd.LogEntry{
		Path:  info.Path,
		Topic: info.Topic,
	},
	)

	data, err := json.Marshal(logConfArr)
	if err != nil {
		logs.Warn("marshal failed, err:%v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//cli.Delete(ctx, EtcdKey)
	//return
	fmt.Println(etcdKey)
	_, err = cli.Put(ctx, etcdKey, string(data))
	//value := `[{"path":"d:/tmp/nginx.log","topic":"web_log"},{"path":"d:/xxx/redis.log","topic":"redis_log"},{"path":"d:/xxx/mysql.log","topic":"mysql_log"},{"path":"d:/xxx/mysql.log","topic":"kafka_log"}]`
	//_, err = cli.Put(ctx, "/log/192.168.1.7/collect_config", value)
	cancel()
	if err != nil {
		logs.Warn("Put failed, err:%v", err)
		return
	}

	logs.Debug("put etcd succ, data:%v", string(data))
	return
}
