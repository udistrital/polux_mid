package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
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
	defer helpers.ErrorController(c.Controller, "TrSolicitudController")
	var v models.TrSolicitud
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if response, err := helpers.AddTransaccionSolicitud(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": 201, "Message": "Solicitud realizada con exito", "Data": response}
		} else {
			// logs.Error(err)
			// fmt.Println("ERROR", err)
			// c.Data["mesaage"] = "Error service POST: The request contains an incorrect data type or an invalid parameter"
			// c.Abort("400")
			panic(err)
		}
	} else {
		// logs.Error(err)
		// c.Data["mesaage"] = "Error service POST: The request contains an incorrect data type or an invalid parameter"
		// c.Abort("400")
		fmt.Println("Se rompe en el último ELSE")
		panic(map[string]interface{}{"funcion": "Post", "err": err.Error(), "status": "400"})
	}
	c.ServeJSON()
}
