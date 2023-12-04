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
	var v models.TrRegistrarRevisionTg
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if response, err := helpers.AddTransaccionRegistrarRevisionTg(&v); err == nil {
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
