package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/plugins/cors"
	_ "github.com/udistrital/Polux_API_mid/routers"
	"github.com/udistrital/utils_oas/customerror"
)

func main() {
	//orm.Debug = true
	logPath := "{\"filename\":\""
	logPath += beego.AppConfig.String("logPath")
	logPath += "\"}"
	if err:= logs.SetLogger(logs.AdapterFile, logPath); err != nil{
		if err:= logs.SetLogger("console", ""); err != nil {
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
	beego.ErrorController(&customerror.CustomErrorController{})
	beego.Run()
}
