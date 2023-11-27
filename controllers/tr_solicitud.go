package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/udistrital/polux_mid/helpers"
	"github.com/udistrital/polux_mid/models"
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
	var v models.TrSolicitud
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if response, err := helpers.AddTransaccionSolicitud(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = response
		} else {
			beego.Error(err)
			c.Abort("400")
		}
	} else {
		beego.Error(err)
		c.Abort("400")
	}
	c.ServeJSON()
}
