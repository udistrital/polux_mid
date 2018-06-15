package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/udistrital/Polux_API_mid/golog"
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

// Get ...
// @Title ObtenerEvaluadores
// @Description get Evaluadores
// @Param	body		body 	int	true		"body for Registrar content"
// @Success 200 {object} make(map[string]string)
// @Failure 400 the request contains incorrect syntax
// @router /ObtenerEvaluadores [post]
func (this *EvaluadoresController) ObtenerEvaluadores() {
	var comprobacion string = ""
	//consultar las reglas
	reglasBase := ruler.CargarReglasBase("RequisitosModalidades")
	fmt.Println(reglasBase)

	var idmodalidad int
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &idmodalidad); err == nil {
		fmt.Println(idmodalidad)
	} else {
		fmt.Println(err)
	}

	var modalidad string
	switch idmodalidad {
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

	this.Data["json"] = m
	this.ServeJSON()

}
