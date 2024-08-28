package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/beego/beego/logs"
	"github.com/udistrital/polux_mid/helpers"
	"github.com/udistrital/polux_mid/models"
	//errorControl "github.com/udistrital/utils_oas/errorctrl"
)

type TrSolicitudController struct {
	beego.Controller
}

// URLMapping ...
func (c *TrSolicitudController) URLMapping() {
	c.Mapping("Post", c.Post)
}

// @Title PostTrSolicitud
// @Description create the TrSolicitud
// @Param	body		body 	models.TrSolicitud	true	"body for TrSolicitud content"
// @Success 201 {int} models.TrSolicitud
// @Failure 400 the request contains incorrect syntax
// @router / [post]
func (c *TrSolicitudController) Post() {
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		logs.Error(err)
	// 		fmt.Println("Defer ERR", err)
	// 		localError := err.(map[string]interface{})
	// 		fmt.Println("LOCAL ERROR", localError)
	// 		c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "TrSolicitudController" + "/" + (localError["funcion"]).(string))
	// 		c.Data["data"] = (localError["err"])
	// 		fmt.Println("C.DATA", c.Data)

	// 		if status, ok := localError["status"]; ok {
	// 			c.Abort(status.(string))
	// 		} else {
	// 			c.Abort("500")
	// 		}
	// 	}
	// }()
	defer helpers.ErrorController(c.Controller, "GestionResolucionesController")
	var v models.TrSolicitud
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if response, err := helpers.AddTransaccionSolicitud(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": 201, "Message": "Solicitud realizada con exito", "Data": response}
		} else {
			logs.Error(err)
			fmt.Println("ERROR", err)
			//c.Data["mesaage"] = "Error service POST: The request contains an incorrect data type or an invalid parameter"
			//c.Abort("400")
			panic(err)
		}
	} else {
		panic(map[string]interface{}{"funcion": "Post", "err": err.Error(), "status": "400"})
	}
	c.ServeJSON()
}
