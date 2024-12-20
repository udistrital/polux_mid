package controllers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/polux_mid/golog"
	"github.com/udistrital/polux_mid/helpers"
	"github.com/udistrital/polux_mid/models"
	"github.com/udistrital/utils_oas/ruler"
)

// VerificarRequisitosController operations for VerificarRequisitos
type VerificarRequisitosController struct {
	beego.Controller
}

// URLMapping ...
func (c *VerificarRequisitosController) URLMapping() {
	c.Mapping("Registrar", c.Registrar)
	c.Mapping("CantidadModalidades", c.CantidadModalidades)
}

// buscar elemento en arreglo
func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

// CantidadModalidades ...
// @Title CantidadModalidades
// @Description Validar si la cantidad de estudiantes solicitados es menor o igual a la cantidad de estudiantes que permite la modalidad
// @Param	body		body 	models.CantidadModalidad	true		"body for CantidadModalidades content"
// @Success 200 {object} make(map[string]bool)
// @Failure 400 the request contains incorrect syntax
// @router /CantidadModalidades [post]
func (c *VerificarRequisitosController) CantidadModalidades() {

	defer helpers.ErrorController(c.Controller, "VerificarRequisitosController")
	reglasBase := ruler.CargarReglasBase("RequisitosModalidades")
	if reglasBase != "" {
		var v models.CantidadModalidad
		if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
			fmt.Println(v)
			var modalidad string
			cantidad := v.Cantidad

			//modificar para que haga validacion de la modalidad aca! Switch
			switch os := v.Modalidad; os {
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
			case "PAS_PLX":
				modalidad = "pasantia"
			}

			comprobacion := "validar_cantidad_estudiantes(" + modalidad + ", " + cantidad + ")."

			r := golog.Comprobar(reglasBase, comprobacion)

			var m = make(map[string]bool)
			m["RequisitosModalidades"] = (r == "true")
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": 201, "Message": "Solicitud realizada con exito", "Data": m}
		} else {
			//beego.Error(err)
			//c.Abort("400")
			panic(map[string]interface{}{"funcion": "Post", "err": err.Error(), "status": "400"})
		}
	} else {
		beego.Error("Sin reglas base")
		c.Abort("400")
	}
	c.ServeJSON()
}

// Registrar ...
// @Title Registrar
// @Description Validar si un estudiante cumple con los requisitos para cursar una modalidad
// @Param	body		body 	models.Datos	true		"body for Registrar content"
// @Success 200 {object} make(map[string]bool)
// @Failure 400 the request contains incorrect syntax
// @router /Registrar [post]
func (c *VerificarRequisitosController) Registrar() {
	//var predicados []models.Predicado
	//var postdominio string = ""
	var comprobacion string
	var reglasbase string

	fmt.Println("CENTINELAAAAAAAAAAAA: ")

	fmt.Println(beego.AppConfig.String("Urlruler") + "predicado?limit=0&query=Dominio.Nombre:" + "RequisitosModalidades")
	reglasBase := ruler.CargarReglasBase("RequisitosModalidades")
	if reglasBase != "" {
		var v models.Datos
		if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
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
			//Realiza la lectura del estado
			estado := v.Estado
			porcentaje := v.PorcentajeCursado
			promedio := v.Promedio
			nivel := strings.ToLower(v.Nivel)
			tipoCarrera := strings.ToLower(v.TipoCarrera)

			estados := []string{"A", "B", "V", "T", "J"}
			modalidades := []string{"MONO_PLX", "INV_PLX", "PEMP_PLX", "PACAD_PLX", "PAS_PLX"} //Modalidades que solo necesitan el Porcentaje cursado y el Estado del estudiante
			if stringInSlice(v.Estado, estados) {
				estado = "activo"
			}
			reglasbase = reglasBase + "estado(" + codigo + ", " + estado + ").cursado(" + codigo + ", " + porcentaje + ").nivel(" + codigo + ", " + nivel + ")."
			if stringInSlice(modalidad, modalidades) {
				comprobacion = "validacion_requisitos(" + codigo + ")."
			} else if modalidad == "EAPOS_PLX" {
				reglasbase = reglasbase + "promedio(" + codigo + ", " + promedio + ").tipo_carrera(" + codigo + ", " + tipoCarrera + ")."
				comprobacion = "validacion_posgrado(" + codigo + ")."
			} else if modalidad == "EAPRO_PLX" {
				reglasbase = reglasbase + "tipo_carrera(" + codigo + ", " + tipoCarrera + ")."
				comprobacion = "validacion_profundizacion(" + codigo + ")."
			} else if modalidad == "CRE_PLX" {
				reglasbase = reglasbase + "tipo_carrera(" + codigo + ", " + tipoCarrera + ")."
				comprobacion = "validacion_creacion(" + codigo + ")."
			}
			fmt.Println(reglasbase)

			r := golog.Comprobar(reglasbase, comprobacion)

			var m = make(map[string]bool)
			m["RequisitosModalidades"] = (r == "true")
			fmt.Println("RESULTADO ", m)
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": 201, "Message": "Solicitud realizada con exito", "Data": m}
			fmt.Println("RESULTADO DE LA VARIABLE M: ", m)
		} else {
			beego.Error(err)
			c.Abort("400")
		}
	} else {
		beego.Error("Sin reglas base")
		c.Abort("400")
	}
	c.ServeJSON()

}
