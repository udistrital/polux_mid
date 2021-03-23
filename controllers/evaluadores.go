package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/udistrital/polux_mid/golog"
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

		var modalidad string
		switch v.Modalidad {
		case 1:
			modalidad = "pasantia"
		case 2:
			modalidad = "posgrado"
		case 3:
			modalidad = "profundizacion"
		case 4:
			modalidad = "monografia"
		case 5:
			modalidad = "investigacion"
		case 6:
			modalidad = "creacion"
		case 7:
			modalidad = "emprendimiento"
		case 8:
			modalidad = "articulo"
		}

		comprobacion = "numero_evaluadores(" + modalidad + ",Y)."
		r := golog.Obtener(reglasBase, comprobacion)
		var m = make(map[string]string)
		m["cantidad_evaluadores"] = r

		fmt.Println(m)

		c.Data["json"] = m
	} else {
		beego.Error("Sin reglas base")
		c.Abort("400")
	}
	c.ServeJSON()

}
