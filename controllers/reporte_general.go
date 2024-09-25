package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/polux_mid/helpers"
)

// Reporte_generalController operations for Reporte_general
type ReporteGeneralController struct {
	beego.Controller
}

// URLMapping ...
func (c *ReporteGeneralController) URLMapping() {
	c.Mapping("Post", c.Post)
}

// Post ...
// @Title Create
// @Description create reporte_general
// @Success 201
// @Failure 403
// @router / [post]
func (c *ReporteGeneralController) Post() {
	defer helpers.ErrorController(c.Controller, "ReporteGeneralController")
	if err := helpers.BuildReporteGeneral(); err == nil {
		//c.Ctx.Output.SetStatus(201)
		//c.Data["json"] = "Reporte generado correctamente."
		c.Ctx.Output.SetStatus(201)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": 201, "Message": "Reporte generado correctamente.", "Data": "Reporte generado correctamente."}
	} else {
		//c.Data["json"] = err.Error()
		//c.Ctx.Output.SetStatus(403)
		panic(err)
	}
	c.ServeJSON()
}
