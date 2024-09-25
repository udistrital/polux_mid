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
	//Generar el archivo Excel usando el helper
	if file, err := helpers.BuildReporteSolicitud(); err == nil {
		//Enviar el archivo codificado en Base64 al Cliente
		c.Data["json"] = map[string]interface{}{
			"file":     file,                    //Archivo codificado en Base64
			"filename": "ReporteSolicitud.xlsx", //Nombre del archivo
		}
		c.Ctx.Output.SetStatus(201)
	} else {
		//Manejar errores al generar el reporte
		c.Data["json"] = err.Error()
		c.Ctx.Output.SetStatus(403)
	}
	c.ServeJSON()
}
