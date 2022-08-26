package router

import (
	"test/logagent/myweb/controller/app_controller"
	"test/logagent/myweb/controller/log_controller"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/index", &app_controller.AppController{}, "*:AppList")
	beego.Router("/app/list", &app_controller.AppController{}, "*:AppList")
	beego.Router("/app/apply", &app_controller.AppController{}, "*:AppApply")

	beego.Router("/app/create", &app_controller.AppController{}, "*:AppCreate")

	beego.Router("/log/apply", &log_controller.LogController{}, "*:LogApply")
	beego.Router("/log/list", &log_controller.LogController{}, "*:LogList")
	beego.Router("/log/create", &log_controller.LogController{}, "*:LogCreate")
}
