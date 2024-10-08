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

	//Generar el archivo Excel usando el helper
	if file, err := helpers.BuildReporteGeneral(); err == nil {
		//Enviar el archivo codificado en Base64 al Cliente
		c.Data["json"] = map[string]interface{}{
			"Success": true,
			"Status":  201,
			"Message": "Excel 'Reporte General' generado correctamente.",
			"Data":    file, //Archivo codificado en Base64
		}

		c.Ctx.Output.SetStatus(201)
	} else {
		//Manejar errores al generar el reporte
		c.Data["json"] = map[string]interface{}{
			"Success": false,
			"Status":  404,
			"Message": "Error al generar el Excel de Reporte Solicitud.",
			"Data":    err.Error(),
		}

		c.Ctx.Output.SetStatus(404)
		panic(err)
	}
	c.ServeJSON()
}
