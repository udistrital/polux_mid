package controllers

import (
	"fmt"

	"encoding/json"
	"sort"

	"github.com/udistrital/Polux_API_mid/models"

	"github.com/astaxie/beego"
)

type SeleccionController struct {
	beego.Controller
}

func (c *SeleccionController) URLMapping() {
	c.Mapping("Seleccionar", c.Seleccionar)
}

// Seleccionar ...
// @Title Seleccionar
// @Description post admitidos
// @Param	body		body 	models.TrSolicitud	true		"body for Seleccionar content"
// @Success 200 {string}
// @Failure 403 body is empty
// @router /Seleccionar [post]
func (c *SeleccionController) Seleccionar() {

	var v models.TrSolicitud

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {

		o := make(models.Vals, len(*v.Solicitudes))
		for i, x := range *v.Solicitudes {
			o[i] = x

		}

		sort.SliceStable(o, func(i, j int) bool {
			return o[i].Rendimiento > o[j].Rendimiento
		})

		sort.SliceStable(o, func(i, j int) bool {
			return o[i].Promedio > o[j].Promedio
		})

		if v.NumAdmitidos.Cupos_excelencia > 0 && len(*v.Solicitudes) > 0 {
			var filas int
			if v.NumAdmitidos.Cupos_excelencia <= len(*v.Solicitudes) {
				filas = v.NumAdmitidos.Cupos_excelencia
			} else if v.NumAdmitidos.Cupos_excelencia > len(*v.Solicitudes) {
				filas = len(*v.Solicitudes)
			} else {
				filas = 0
			}
			var rta models.RespuestaSolicitud

			for i := 0; i < filas; i++ {
				fmt.Println(o[i].Solicitud)

				if err := getJson(beego.AppConfig.String("Urlcrud")+"/"+beego.AppConfig.String("Nscrud")+"/respuesta_solicitud/"+o[i].Respuesta, &rta); err == nil {

					rta.EstadoSolicitud.Id = 7
					//cambiar estado de la solicitud
					var respuesta interface{}

					if err := sendJson(beego.AppConfig.String("Urlcrud")+"/"+beego.AppConfig.String("Nscrud")+"/respuesta_solicitud/"+o[i].Respuesta, "PUT", &respuesta, &rta); err == nil {
						c.Data["json"] = "Solicitudes Aceptadas"
					}
				} else {
					c.Data["json"] = err.Error()
				}
			}
		}

		var filas2 = 0
		if v.NumAdmitidos.Cupos_adicionales > 0 {

			if v.NumAdmitidos.Cupos_excelencia+v.NumAdmitidos.Cupos_adicionales >= len(*v.Solicitudes) {
				filas2 = len(*v.Solicitudes) - v.NumAdmitidos.Cupos_excelencia

			} else {
				filas2 = v.NumAdmitidos.Cupos_adicionales
			}

			var rta2 models.RespuestaSolicitud

			for i := v.NumAdmitidos.Cupos_excelencia; i < v.NumAdmitidos.Cupos_excelencia+filas2; i++ {
				fmt.Println(beego.AppConfig.String("Urlcrud") + "/" + beego.AppConfig.String("Nscrud") + "/respuesta_solicitud/" + o[i].Respuesta)
				if err := getJson(beego.AppConfig.String("Urlcrud")+"/"+beego.AppConfig.String("Nscrud")+"/respuesta_solicitud/"+o[i].Respuesta, &rta2); err == nil {
					rta2.EstadoSolicitud.Id = 8
					//cambiar estado de la solicitud
					var respuesta2 interface{}
					fmt.Println(beego.AppConfig.String("Urlcrud") + "/" + beego.AppConfig.String("Nscrud") + "/respuesta_solicitud/" + o[i].Respuesta)
					if err := sendJson(beego.AppConfig.String("Urlcrud")+"/"+beego.AppConfig.String("Nscrud")+"/respuesta_solicitud/"+o[i].Respuesta, "PUT", &respuesta2, &rta2); err == nil {
						fmt.Println(respuesta2)
						c.Data["json"] = "Solicitudes Aceptadas"
					}
				} else {
					c.Data["json"] = err.Error()
				}
			}
		}

	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()

}
