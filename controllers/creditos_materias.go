package controllers

import (
	"fmt"
	"github.com/udistrital/Polux_API_mid/models"
	"strconv"
	"github.com/astaxie/beego"
	"github.com/udistrital/Polux_API_mid/golog"
  	"github.com/udistrital/utils_oas/ruler"
)

// CreditosMateriasController operations for CreditosMaterias
type CreditosMateriasController struct {
	beego.Controller
}

func (c *CreditosMateriasController) URLMapping() {
	c.Mapping("ObtenerCreditos", c.ObtenerCreditos)
}

// Get ...
// @Title ObtenerCreditos Materias
// @Description Obtener el número de créditos minimos que se pueden cursar en la modalidad de materias de posgrado o profundización
// @Success 200 {object} models.CreditosMaterias
// @Failure 403
// @router /ObtenerCreditos [get]
func (c *CreditosMateriasController) ObtenerCreditos() {

  var creditosMaterias models.CreditosMaterias
  var comprobacion string = ""
  //consultar las reglas
  reglasBase := ruler.CargarReglasBase("MateriasPosgrado")

  //obtener minimo de creditos para  materias de posgrado
  comprobacion="min_creditos_asignaturas_posgrado(Y)."
  r:=golog.Obtener(reglasBase,comprobacion)
  fmt.Println(r)
  creditosMaterias.MateriasPosgrado , _ = strconv.Atoi(r)

  comprobacion="min_creditos_asignaturas_profundizacion(Y)."
  r=golog.Obtener(reglasBase,comprobacion)
  fmt.Println(r)
  creditosMaterias.MateriasProfundizacion , _ = strconv.Atoi(r)

  c.Data["json"] = creditosMaterias
  c.ServeJSON()

}


