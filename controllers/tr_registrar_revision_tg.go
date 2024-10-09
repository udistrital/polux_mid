package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/udistrital/polux_mid/helpers"
	"github.com/udistrital/polux_mid/models"
)

type TrRegistrarRevisionTgController struct {
	beego.Controller
}

// URLMapping ...
func (c *TrRegistrarRevisionTgController) URLMapping() {
	c.Mapping("Post", c.Post)
}

// @Title PostTrRegistrarRevisionTg
// @Description create the TrRegistrarRevisionTg
// @Param	body		body 	models.TrRegistrarRevisionTg	true	"body for TrRegistrarRevisionTg content"
// @Success 201 {int} models.TrRegistrarRevisionTg
// @Failure 400 the request contains incorrect syntax
// @router / [post]
func (c *TrRegistrarRevisionTgController) Post() {
	defer helpers.ErrorController(c.Controller, "TrRegistrarRevisionTgController")
	var v models.TrRegistrarRevisionTg
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if response, err := helpers.AddTransaccionRegistrarRevisionTg(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": 201, "Message": "Solicitud realizada con exito", "Data": response}
			//beego.Error(err)
			//c.Abort("400")
			//panic(err)
		}
	} else {
		//beego.Error(err)
		//c.Abort("400")
		panic(map[string]interface{}{"funcion": "Post", "err": err.Error(), "status": "400"})
	}
	c.ServeJSON()
}
