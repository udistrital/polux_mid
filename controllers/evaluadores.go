package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/polux_mid/golog"
	"github.com/udistrital/polux_mid/helpers"
	"github.com/udistrital/polux_mid/models"
	"github.com/udistrital/utils_oas/ruler"
)

// EvaluadoresController operations for Evaluadores
type EvaluadoresController struct {
	beego.Controller
}

// URLMapping ...
func (c *EvaluadoresController) URLMapping() {
	c.Mapping("ObtenerEvaluadores", c.ObtenerEvaluadores)
}

// ObtenerEvaluadores ...
// Get ...
// @Title ObtenerEvaluadores
// @Description get Evaluadores
// @Param	body		body 	models.CantidadEvaluadoresModalidad	true		"body for Registrar content"
// @Success 200 {object} make(map[string]string)
// @Failure 400 the request contains incorrect syntax
// @router /ObtenerEvaluadores [post]
func (c *EvaluadoresController) ObtenerEvaluadores() {
	var comprobacion string
	//consultar las reglas
	reglasBase := ruler.CargarReglasBase("RequisitosModalidades")
	if reglasBase != "" {
		fmt.Println(reglasBase)

		var v models.CantidadEvaluadoresModalidad
		if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
			fmt.Println(v)
		} else {
			beego.Error("Sin modalidad valida")
			c.Abort("400")
		}
		if modalidadParam, err2 := helpers.ObtenerModalidad(v); err2 == nil {
			var modalidad string
			switch modalidadParam.CodigoAbreviacion {
			case "PASEX_PLX":
				modalidad = "pasantia"
			case "EAPOS_PLX":
				modalidad = "posgrado"
			case "EAPRO_PLX":
				modalidad = "profundizacion"
			case "MONO_PLX":
				modalidad = "monografia"
			case "INV_PLX":
				modalidad = "investigacion"
			case "CRE_PLX":
				modalidad = "creacion"
			case "PEMP_PLX":
				modalidad = "emprendimiento"
			case "PACAD_PLX":
				modalidad = "articulo"
			case "PASIN_PLX":
				modalidad = "pasantia"
			}
			comprobacion = "numero_evaluadores(" + modalidad + ",Y)."
			r := golog.Obtener(reglasBase, comprobacion)
			var m = make(map[string]string)
			m["cantidad_evaluadores"] = r
			c.Data["json"] = m
		}
	} else {
		beego.Error("Sin reglas base")
		c.Abort("400")
	}
	c.ServeJSON()

}
