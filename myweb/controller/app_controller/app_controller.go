package app_controller

import (
	"test/logagent/myweb/conf"
	"test/logagent/myweb/model"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type AppController struct {
	beego.Controller
}

func (p *AppController) AppList() {
	logs.Debug("enter index controller")
	p.Layout = "layout/layout.html"
	appList, err := model.GetAppList()
	if err != nil {
		p.Data["Error"] = "服务器繁忙"
		p.TplName = "app/error.html"
		logs.Warn("get app list failed,err:", err)
	}
	logs.Debug("get app list success,data:%v", appList)
	p.Data["applist"] = appList
	p.TplName = "app/index.html"
}

func (p *AppController) AppApply() {
	logs.Debug("enter apply controller")
	p.Layout = "layout/layout.html"
	p.TplName = "app/apply.html"
}

func (p *AppController) AppCreate() {
	logs.Debug("enter create controller")
	appName := p.GetString("app_name")
	appType := p.GetString("app_type")
	developPath := p.GetString("develop_path")
	ipstr := p.GetString("iplist")

	p.Layout = "layout/layout.html"

	if len(appName) == 0 {
		p.Data["Error"] = "项目名称不能为空"
		p.TplName = "app/error.html"

		logs.Warn("invalid parameter")
		return
	}
	if len(appType) == 0 {
		p.Data["Error"] = "项目类型不能为空"
		p.TplName = "app/error.html"

		logs.Warn("invalid parameter")
		return
	}
	if len(developPath) == 0 {
		p.Data["Error"] = "部署路径不能为空"
		p.TplName = "app/error.html"

		logs.Warn("invalid parameter")
		return
	}
	if len(ipstr) == 0 {
		p.Data["Error"] = "IP地址不能为空"
		p.TplName = "app/error.html"

		logs.Warn("invalid parameter")
		return
	}

	appData := &conf.AppData{
		AppName:     appName,
		AppType:     appType,
		DevelopPath: developPath,
		Ip:          ipstr,
	}
	err := model.CreateApp(*appData)
	if err != nil {
		p.Data["Error"] = "创建项目失败"
		p.TplName = "app/error.html"

		logs.Warn("invalid parameter")
		return
	}
	p.Redirect("/app/list", 302)
}
