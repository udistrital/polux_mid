package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/udistrital/polux_mid/helpers"
	"github.com/udistrital/polux_mid/models"
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
// @Param	body		body 	models.FiltrosReporte	true	"body for FiltrosReporte content"
// @Success 201
// @Failure 403
// @router / [post]
func (c *ReporteSolicitudController) Post() {
	defer helpers.ErrorController(c.Controller, "ReporteSolicitudController")
	var v models.FiltrosReporte
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		//Generar el archivo Excel usando el helper
		if file, err := helpers.BuildReporteSolicitud(&v); err == nil {
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
				"Status":  400,
				"Message": "Error al generar el Excel de Reporte Solicitud.",
				"Data":    nil,
			}

			c.Ctx.Output.SetStatus(400)
			//panic(err)
		}
	} else {
		panic(map[string]interface{}{"funcion": "Post", "err": err.Error(), "status": "400"})
	}
	c.ServeJSON()
}
