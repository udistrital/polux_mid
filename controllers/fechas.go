package controllers

import (

	"fmt"
	"github.com/astaxie/beego"
	"github.com/udistrital/Polux_API_mid/golog"
)

type FechasController struct {
	beego.Controller
}

func (c *FechasController) URLMapping() {
	c.Mapping("ObtenerFechas", c.ObtenerFechas)
}

// Get ...
// @Title ObtenerFechas
// @Description get fechas
// @Success 200 {object} make(map[string]string)
// @Failure 403
// @router /ObtenerFechas [get]
func (c *FechasController) ObtenerFechas() {
  var comprobacion string = ""
	//consultar las reglas
	reglasBase := CargarReglasBase("FechasSeleccion")
	fmt.Println(reglasBase)

  comprobacion="fecha_inicio_proceso_seleccion(Y)."
  r:=golog.Obtener(reglasBase,comprobacion)
	var m =make(map[string]string)
	m["inicio_proceso"] = r

	comprobacion="segunda_fecha_proceso_seleccion(Y)."
  r=golog.Obtener(reglasBase,comprobacion)
	m["segunda_fecha"] =r

	comprobacion="fecha_fin_proceso_seleccion(Y)."
  r=golog.Obtener(reglasBase,comprobacion)
	m["fecha_fin"] =r

	fmt.Println(m)

  c.Data["json"] = m
  c.ServeJSON()

}
