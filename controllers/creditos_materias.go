package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/udistrital/polux_mid/golog"
	"github.com/udistrital/polux_mid/models"
	"github.com/udistrital/utils_oas/ruler"
	"strconv"
)

// CreditosMateriasController operations for CreditosMaterias
type CreditosMateriasController struct {
	beego.Controller
}

// URLMapping ...
func (c *CreditosMateriasController) URLMapping() {
	c.Mapping("ObtenerCreditos", c.ObtenerCreditos)
}

// ObtenerCreditos ...
// Get ...
// @Title ObtenerCreditos Materias
// @Description Obtener el número de créditos minimos que se pueden cursar en la modalidad de materias de posgrado o profundización
// @Success 200 {object} models.CreditosMaterias
// @Failure 400 the request contains incorrect syntax
// @router /ObtenerCreditos [get]
func (c *CreditosMateriasController) ObtenerCreditos() {

	var creditosMaterias models.CreditosMaterias
	var comprobacion string
	//consultar las reglas
	reglasBase := ruler.CargarReglasBase("MateriasPosgrado")
	if reglasBase != "" {
		//obtener minimo de creditos para  materias de posgrado
		comprobacion = "min_creditos_asignaturas_posgrado(Y)."
		r := golog.Obtener(reglasBase, comprobacion)
		fmt.Println(r)
		creditosMaterias.MateriasPosgrado, _ = strconv.Atoi(r)

		comprobacion = "min_creditos_asignaturas_profundizacion(Y)."
		r = golog.Obtener(reglasBase, comprobacion)
		fmt.Println(r)
		creditosMaterias.MateriasProfundizacion, _ = strconv.Atoi(r)

		c.Data["json"] = creditosMaterias
	} else {
		beego.Error("Sin reglas base")
		c.Abort("400")
	}
	c.ServeJSON()

}
