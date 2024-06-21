package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/udistrital/polux_mid/helpers"
	"github.com/udistrital/polux_mid/models"
)

type TrSubirArl struct {
	beego.Controller
}

//URLMapping ...
func (c *TrSubirArl) URLMapping() {
	c.Mapping("Post", c.Post)
}

// @Title PostTrSubirArl
// @Description create the TrSubirArl
// @Param	body		body 	models.TrSubirArl	true	"body for TrSubirArl content"
// @Success 201 {int} models.TrSubirArl
// @Failure 400 the request contains incorrect syntax
// @router / [post]
func (c *TrSubirArl) Post() {
	var v models.TrSubirArl
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if response, err := helpers.AddTransaccionSubirArl(&v); err == nil {
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