package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/udistrital/Polux_API_mid/golog"
	"github.com/udistrital/Polux_API_mid/models"
	"github.com/udistrital/utils_oas/ruler"
	"strconv"
)

// CuposController operations for Cupos
type CuposController struct {
	beego.Controller
}

func (c *CuposController) URLMapping() {
	c.Mapping("Obtener", c.Obtener)
}

// Get ...
// @Title Obtener
// @Description Obtener el número máximo de cupos de excelencia academica y adicionales que se pueden aceptar en la modalidad de materias de posgrado y profundización
// @Success 200 {object} models.Cupos
// @Failure 400 the request contains incorrect syntax
// @router /Obtener [get]
func (c *CuposController) Obtener() {

	var NumAdmitidos models.Cupos
	var comprobacion string = ""
	//consultar las reglas
	reglasBase := ruler.CargarReglasBase("MateriasPosgrado")

	//obtener máximo de cupos por excelencia académica
	comprobacion = "max_cupos_excelencia_academica(Y)."
	r := golog.Obtener(reglasBase, comprobacion)
	fmt.Println(r)
	NumAdmitidos.Cupos_excelencia, _ = strconv.Atoi(r)

	//obtener máximo de cupos adicionales
	comprobacion = "max_cupos_adicionales(Y)."
	r2 := golog.Obtener(reglasBase, comprobacion)
	fmt.Println(r2)
	NumAdmitidos.Cupos_adicionales, _ = strconv.Atoi(r2)

	c.Data["json"] = NumAdmitidos
	c.ServeJSON()
	///////////////////////////////////////////////////////////////////////////

}
