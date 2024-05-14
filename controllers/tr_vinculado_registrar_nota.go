package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/udistrital/polux_mid/helpers"
	"github.com/udistrital/polux_mid/models"
)

type TrVinculadoRegistrarNotaController struct {
	beego.Controller
}

//URLMapping ...
func (c *TrVinculadoRegistrarNotaController) URLMapping() {
	c.Mapping("Post", c.Post)
}

// @Title PostTrVinculadoRegistrarNota
// @Description create the TrVinculadoRegistrarNota
// @Param	body		body 	models.TrVinculadoRegistrarNota	true	"body for TrVinculadoRegistrarNota content"
// @Success 201 {int} models.TrVinculadoRegistrarNota
// @Failure 400 the request contains incorrect syntax
// @router / [post]
func (c *TrVinculadoRegistrarNotaController) Post() {
	var v models.TrVinculadoRegistrarNota
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if response, err := helpers.AddTransaccionVinculadoRegistrarNota(&v); err == nil {
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
