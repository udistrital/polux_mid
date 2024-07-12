package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/udistrital/polux_mid/helpers"
	"github.com/udistrital/polux_mid/models"
)

type TrRegistrarActaSeguimiento struct {
	beego.Controller
}

// URLMapping ...
func (c *TrRegistrarActaSeguimiento) URLMapping() {
	c.Mapping("Post", c.Post)
}

// @Title PostTrRegistrarActaSeguimiento
// @Description create the TrRegistrarActaSeguimiento
// @Param	body		body 	models.TrRegistrarActaSeguimiento	true	"body for TrRegistrarActaSeguimiento content"
// @Success 201 {int} models.TrRegistrarActaSeguimiento
// @Failure 400 the request contains incorrect syntax
// @router / [post]
func (c *TrRegistrarActaSeguimiento) Post() {
	var v models.TrRegistrarActaSeguimiento
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if response, err := helpers.AddTransaccionRegistrarActaSeguimiento(&v); err == nil {
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