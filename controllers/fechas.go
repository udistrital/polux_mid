package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/udistrital/Polux_API_mid/golog"
	"github.com/udistrital/utils_oas/ruler"
)

// FechasController operations for Fechas
type FechasController struct {
	beego.Controller
}

// URLMapping ...
func (c *FechasController) URLMapping() {
	c.Mapping("ObtenerFechas", c.ObtenerFechas)
}

// ObtenerFechas ...
// Get ...
// @Title ObtenerFechas
// @Description Obtener fechas para el procso de selección de admitidos
// @Success 200 {object} make(map[string]string)
// @Failure 400 the request contains incorrect syntax
// @router /ObtenerFechas [get]
func (c *FechasController) ObtenerFechas() {
	var comprobacion string
	//consultar las reglas
	fmt.Println("http://" + beego.AppConfig.String("Urlruler") + "/" + beego.AppConfig.String("Nsruler") + "/predicado?limit=0&query=Dominio.Nombre:" + "FechasSeleccion")
	fmt.Println("http://" + beego.AppConfig.String("Urlruler") + ":" + beego.AppConfig.String("Portruler") + "/" + beego.AppConfig.String("Nsruler") + "/predicado?limit=0&query=Dominio.Nombre:" + "FechasSeleccion")
	reglasBase := ruler.CargarReglasBase("FechasSeleccion")
	fmt.Println(reglasBase)

	comprobacion = "fecha_inicio_proceso_seleccion(Y)."
	r := golog.Obtener(reglasBase, comprobacion)
	var m = make(map[string]string)
	m["inicio_proceso"] = r

	comprobacion = "segunda_fecha_proceso_seleccion(Y)."
	r = golog.Obtener(reglasBase, comprobacion)
	m["segunda_fecha"] = r

	comprobacion = "fecha_fin_proceso_seleccion(Y)."
	r = golog.Obtener(reglasBase, comprobacion)
	m["fecha_fin"] = r

	fmt.Println(m)

	c.Data["json"] = m
	c.ServeJSON()

}
