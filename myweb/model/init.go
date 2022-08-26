package model

import (
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/coreos/etcd/clientv3"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var Db *sqlx.DB
var cli *clientv3.Client

func init() {
	var err error
	Db, err = sqlx.Open("mysql", "root:*******@tcp(127.0.0.1:3306)/logdata")
	if err != nil {
		logs.Warn("open mysql failed", err)
		return
	}
	err = Db.Ping()
	if err != nil {
		return
	}

	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"}, //节点
		DialTimeout: 5 * time.Second,                                                  //超过5秒钟连不上超时
	})
	if err != nil {
		logs.Warn("connect to etcd failed:", err)
		return
	}
}
