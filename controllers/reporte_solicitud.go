package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/polux_mid/helpers"
)

// Reporte_generalController operations for Reporte_general
type ReporteSolicitudController struct {
	beego.Controller
}

// URLMapping ...
func (c *ReporteSolicitudController) URLMapping() {
	c.Mapping("Post", c.Post)
}

// Post ...
// @Title Create
// @Description create reporte_general
// @Success 201
// @Failure 403
// @router / [post]
func (c *ReporteSolicitudController) Post() {
	if err := helpers.BuildReporteSolicitud(); err == nil {
		c.Ctx.Output.SetStatus(201)
		c.Data["json"] = "Reporte generado correctamente."
	} else {
		c.Data["json"] = err.Error()
		c.Ctx.Output.SetStatus(403)
	}
	c.ServeJSON()
}
