package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/udistrital/auditoria"
	_ "github.com/udistrital/polux_mid/routers"
	apistatus "github.com/udistrital/utils_oas/apiStatusLib"
	"github.com/udistrital/utils_oas/customerrorv2"
	"github.com/udistrital/utils_oas/xray"
)

func main() {
	//orm.Debug = true
	logPath := "{\"filename\":\""
	logPath += beego.AppConfig.String("logPath")
	logPath += "\"}"
	if err := logs.SetLogger(logs.AdapterFile, logPath); err != nil {
		if err := logs.SetLogger("console", ""); err != nil {
			logs.Warn("logPath not set")
		}
	}
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowHeaders: []string{"Origin", "x-requested-with",
			"content-type",
			"accept",
			"origin",
			"authorization",
			"x-csrftoken"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	beego.ErrorController(&customerrorv2.CustomErrorController{})
	xray.InitXRay()
	apistatus.Init()
	auditoria.InitMiddleware()
	beego.Run()
}
