package controllers

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/polux_mid/golog"
	"github.com/udistrital/utils_oas/ruler"
)

// CreditosController operations for Creditos
type CreditosController struct {
	beego.Controller
}

// URLMapping ...
func (c *CreditosController) URLMapping() {
	c.Mapping("ObtenerMinimo", c.ObtenerMinimo)
}

// ObtenerMinimo ...
// Get ...
// @Title ObtenerMinimo
// @Description Obtener el numero de creditos minimos necesarios para solicitar materias de posgrado
// @Success 200 {object} make(map[string]string)
// @Failure 400 the request contains incorrect syntax
// @router /ObtenerMinimo [get]
func (c *CreditosController) ObtenerMinimo() {
	var comprobacion string 
	//consultar las reglas
	reglasBase := ruler.CargarReglasBase("RequisitosModalidades")
	if reglasBase != "" {
		fmt.Println(reglasBase)

		comprobacion = "minimo_numero_creditos_posgrado(Y)."
		r := golog.Obtener(reglasBase, comprobacion)
		var m = make(map[string]string)
		m["minimo_creditos_posgrado"] = r

		comprobacion = "minimo_numero_creditos_profundizacion(Y)."
		r = golog.Obtener(reglasBase, comprobacion)
		m["minimo_creditos_profundizacion"] = r

		fmt.Println(m)

		c.Data["json"] = m
	} else {
		beego.Error("Sin reglas base")
		c.Abort("400")
	}
	c.ServeJSON()

}
