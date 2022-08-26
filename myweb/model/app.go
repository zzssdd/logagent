package model

import (
	"database/sql"
	"test/logagent/myweb/conf"

	"github.com/astaxie/beego/logs"
)

func GetAppList() (applist []conf.AppData, err error) {
	err = Db.Select(&applist, "select aid,aname,atype,create_time,develop_path from appdata")
	if err != nil {
		logs.Warn("Get All data failed:", err)
		return
	}
	return
}

func CreateApp(createlist conf.AppData) (err error) {
	var ret sql.Result
	ret, err = Db.Exec("insert into appdata(aname,atype,develop_path,ip)values(?,?,?,?)", createlist.AppName, createlist.AppType, createlist.DevelopPath, createlist.Ip)
	if err != nil {
		logs.Warn("Create App failed,err:", err)
		return
	}
	_, err = ret.LastInsertId()
	if err != nil {
		logs.Warn("Create App failed,err:", err)
		return
	}
	return
}
