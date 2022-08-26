package main

import (
	"github.com/astaxie/beego"

	_ "test/logagent/myweb/router"
)

func main() {
	beego.Run()
}
