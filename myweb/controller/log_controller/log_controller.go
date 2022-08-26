package log_controller

import (
	"fmt"
	"test/logagent/myweb/conf"
	"test/logagent/myweb/model"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type LogController struct {
	beego.Controller
}

func (p *LogController) LogList() {
	logs.Debug("enter index controller")
	p.Layout = "layout/layout.html"

	loglist, err := model.GetLogList()
	if err != nil {
		p.Data["Error"] = "服务器繁忙"
		p.TplName = "log/error.html"
		logs.Warn("get loglist failed,err:", err)
		return
	}
	p.Data["loglist"] = loglist
	p.TplName = "log/index.html"
}

func (p *LogController) LogApply() {
	logs.Debug("enter index controller")
	p.Layout = "layout/layout.html"
	p.TplName = "log/apply.html"
}

func (p *LogController) LogCreate() {
	logs.Debug("enter create controller")
	appName := p.GetString("app_name")
	logPath := p.GetString("log_path")
	topic := p.GetString("topic")

	p.Layout = "layout/layout.html"
	if len(appName) == 0 {
		p.Data["Error"] = "项目名称不能为空"
		p.TplName = "log/error.html"

		logs.Warn("invalid parameter")
		return
	}
	if len(logPath) == 0 {
		p.Data["Error"] = "日志路径不能为空"
		p.TplName = "log/error.html"

		logs.Warn("invalid parameter")
		return
	}
	if len(topic) == 0 {
		p.Data["Error"] = "日志Topic不能为空"
		p.TplName = "log/error.html"

		logs.Warn("invalid parameter")
		return
	}

	logData := &conf.LogData{
		AppName: appName,
		LogPath: logPath,
		Topic:   topic,
	}
	fmt.Println(logData)
	err := model.CreateLog(logData)
	if err != nil {
		p.Data["Error"] = "创建Log失败"
		p.TplName = "log/error.html"
		logs.Warn("create log failed,err:", err)
		return
	}

	keyFormat := "/log/%s/collect_config"

	iplist, err := model.GetIpByName(appName)

	for _, ip := range iplist {
		key := fmt.Sprintf(keyFormat, ip)
		etcdData := &conf.EtcdData{
			Path:  logPath,
			Topic: topic,
		}
		err = model.SetLogConfToEtcd(key, etcdData)
		if err != nil {
			logs.Warn("Set log conf to etcd failed, err:%v", err)
			continue
		}
	}

	p.Redirect("/log/list", 302)
}
