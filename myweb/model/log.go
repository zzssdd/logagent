package model

import (
	"database/sql"
	"test/logagent/myweb/conf"

	"github.com/astaxie/beego/logs"
)

func GetLogList() (loglist []conf.LogData, err error) {
	err = Db.Select(&loglist, "SELECT l.lid,l.aname,l.create_time,log_path,topic FROM appdata a,logdata l WHERE a.aid=l.aid")
	if err != nil {
		logs.Warn("Get All data failed:", err)
		return
	}
	return
}

func CreateLog(logData *conf.LogData) (err error) {
	var tx *sql.Tx
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}

		tx.Commit()
	}()

	tx, err = Db.Begin()
	if err != nil {
		logs.Warn("Db begin failed,err", err)
		return
	}
	var id []int
	err = Db.Select(&id, "select aid from appdata where aname=?", logData.AppName)
	if err != nil {
		logs.Warn("Get aid failed,err:", err)
		return
	}

	var ret sql.Result
	ret, err = tx.Exec("insert into logdata(aname,log_path,topic,aid)Values(?,?,?,?)", logData.AppName, logData.LogPath, logData.Topic, id[0])
	if err != nil {
		logs.Warn("Exec failed,err:", err)
		return
	}
	_, err = ret.LastInsertId()
	if err != nil {
		logs.Warn("Insert failed,err:", err)
		return
	}

	return
}

func GetIpByName(appname string) (iplist []string, err error) {
	err = Db.Select(&iplist, "select ip from appdata where aname=?", appname)
	if err != nil {
		logs.Warn("Select ip failed,err:", err)
		return
	}
	return
}
