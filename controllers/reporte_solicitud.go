package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/polux_mid/helpers"
)

// Reporte_solicitudController operations for Reporte_solicitud
type ReporteSolicitudController struct {
	beego.Controller
}

// URLMapping ...
func (c *ReporteSolicitudController) URLMapping() {
	c.Mapping("Post", c.Post)
}

// Post ...
// @Title Create
// @Description create reporte_solicitud
// @Success 201
// @Failure 403
// @router / [post]
func (c *ReporteSolicitudController) Post() {
	defer helpers.ErrorController(c.Controller, "ReporteSolicitudController")

	//Generar el archivo Excel usando el helper
	if file, err := helpers.BuildReporteSolicitud(); err == nil {
		//Enviar el archivo codificado en Base64 al Cliente
		c.Data["json"] = map[string]interface{}{
			"Success": true,
			"Status":  201,
			"Message": "Excel 'Reporte Solicitud' generado correctamente.",
			"Data":    file, //Archivo codificado en Base64
		}

		c.Ctx.Output.SetStatus(201)
	} else {
		//Manejar errores al generar el reporte
		c.Data["json"] = map[string]interface{}{
			"Success": false,
			"Status":  404,
			"Message": "Error al generar el Excel de Reporte General.",
			"Data":    err.Error(),
		}

		c.Ctx.Output.SetStatus(404)
		panic(err)
	}
	c.ServeJSON()
}
