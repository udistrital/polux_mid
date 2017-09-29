package controllers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/udistrital/Polux_API_mid/golog"
	"github.com/udistrital/Polux_API_mid/models"

	"github.com/astaxie/beego"
)

type VerificarRequisitosController struct {
	beego.Controller
}

func (c *VerificarRequisitosController) URLMapping() {
	c.Mapping("Registrar", c.Registrar)
	c.Mapping("CantidadModalidades", c.CantidadModalidades)
}

//buscar elemento en arreglo
func stringInSlice(str int, list []int) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func stringInSlice2(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

// CantidadModalidades ...
// @Title CantidadModalidades
// @Description get CantidadModalidades
// @Param	body		body 	models.CantidadModalidad	true		"body for CantidadModalidades content"
// @Success 200 {bool}
// @Failure 403 body is empty
// @router /CantidadModalidades [post]
func (this *VerificarRequisitosController) CantidadModalidades() {

	reglasBase := CargarReglasBase("RequisitosModalidades")
	var v models.CantidadModalidad
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &v); err == nil {
		fmt.Println(v)
	} else {
		fmt.Println(err)
	}

	var modalidad string
	cantidad := v.Cantidad

	//modificar para que haga validacion de la modalidad aca! Switch
	switch os := v.Modalidad; os {
	case "1":
		modalidad = "pasantia"
	case "2":
		modalidad = "posgrado"
	case "3":
		modalidad = "profundizacion"
	case "4":
		modalidad = "monografia"
	case "5":
		modalidad = "investigacion"
	case "6":
		modalidad = "creacion"
	case "7":
		modalidad = "emprendimiento"
	case "8":
		modalidad = "articulo"
	}

	comprobacion := "validar_cantidad_estudiantes(" + modalidad + ", " + cantidad + ")."

	r := golog.Comprobar(reglasBase, comprobacion)

	fmt.Println(comprobacion)

	this.Data["json"] = r
	this.ServeJSON()
}

// Registrar ...
// @Title Registrar
// @Description get requisitos
// @Param	body		body 	models.Datos	true		"body for Registrar content"
// @Success 200 {bool}
// @Failure 403 body is empty
// @router /Registrar [post]
func (this *VerificarRequisitosController) Registrar() {
	//var predicados []models.Predicado
	//var postdominio string = ""
	var comprobacion string = ""
	var reglasbase string = ""

	reglasBase := CargarReglasBase("RequisitosModalidades")

	var v models.Datos
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &v); err == nil {
		fmt.Println(v)
	} else {
		fmt.Println(err)
	}

	/*
		Modalidad 1: Pasantía (Estado, Porcentaje, Nivel)
		Modalidad 2: Innovación-Investigación (Estado, Porcentaje, Nivel)
		Modalidad 3: Proyecto Emprendimiento (Estado, Porcentaje, Nivel)
		Modalidad 4: Producción Académica (Estado, Porcentaje, Nivel)
		Modalidad 5: Monografía (Estado, Porcentaje, Nivel)
		Modalidad 6: Materias de posgrado (Estado, Porcentaje, Nivel)+(Promedio, Tipo carrera)
		Modalidad 7: Materias de profundización (Estado, Porcentaje, Nivel)+(Tipo carrera)
		Modalidad 8: Creación o Interpretación (Estado, Porcentaje, Nivel)+(Tipo carrera)
	*/

	codigo := v.Codigo
	modalidad := v.Modalidad
	//estado in (J, A, ...)
	//estado:=v.Estado
	estado := ""
	porcentaje := v.PorcentajeCursado
	promedio := v.Promedio
	nivel := strings.ToLower(v.Nivel)
	tipo_carrera := strings.ToLower(v.TipoCarrera)

	estados := []string{"A", "B", "V", "T", "J"}
	modalidades := []int{1, 4, 5, 7, 8} //Modalidades que solo necesitan el Porcentaje cursado y el Estado del estudiante
	if stringInSlice2(v.Estado, estados) {
		estado = "activo"
	}
	reglasbase = reglasBase + "estado(" + codigo + ", " + estado + ").cursado(" + codigo + ", " + porcentaje + ").nivel(" + codigo + ", " + nivel + ")."

	if stringInSlice(modalidad, modalidades) {
		comprobacion = "validacion_requisitos(" + codigo + ")."
	} else if modalidad == 2 {
		reglasbase = reglasbase + "promedio(" + codigo + ", " + promedio + ").tipo_carrera(" + codigo + ", " + tipo_carrera + ")."
		comprobacion = "validacion_posgrado(" + codigo + ")."
	} else if modalidad == 3 {
		reglasbase = reglasbase + "tipo_carrera(" + codigo + ", " + tipo_carrera + ")."
		comprobacion = "validacion_profundizacion(" + codigo + ")."
	} else if modalidad == 6 {
		reglasbase = reglasbase + "tipo_carrera(" + codigo + ", " + tipo_carrera + ")."
		comprobacion = "validacion_creacion(" + codigo + ")."
	}
	fmt.Println(reglasbase)

	r := golog.Comprobar(reglasbase, comprobacion)

	this.Data["json"] = r
	this.ServeJSON()

}
