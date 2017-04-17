package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"github.com/udistrital/Polux_API_mid/golog"
	"github.com/udistrital/Polux_API_mid/models"

	"github.com/astaxie/beego"
)

type DisponibilidadController struct {
	beego.Controller
}

func (c *DisponibilidadController) URLMapping() {
	c.Mapping("Registrar", c.Registrar)
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

// Registrar ...
// @Title Registrar
// @Description get requisitos
// @Param	body		body 	models.Datos	true		"body for Registrar content"
// @Success 200 {bool}
// @Failure 403 body is empty
// @router /Registrar [post]
func (this *DisponibilidadController) Registrar() {
	//var predicados []models.Predicado
	//var postdominio string = ""
	var comprobacion string = ""
	var reglasbase string = ""

	reglasBase := CargarReglasBase("RequisitosModalidades")

	var v models.Datos
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &v); err == nil {
		fmt.Println(v)
	}else{
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

		codigo:=v.Codigo
		modalidad:=v.Modalidad
		//estado in (J, A, ...)
		//estado:=v.Estado
		estado:=""
		porcentaje:=strconv.FormatFloat(v.PorcentajeCursado, 'f', -1, 64)
		promedio:=v.Promedio
		nivel:=strings.ToLower(v.Nivel)
		tipo_carrera:=strings.ToLower(v.TipoCarrera)

		estados := []string{"A", "B", "V", "T", "J"}
		modalidades := []int{1, 2, 3, 4, 5} //Modalidades que solo necesitan el Porcentaje cursado y el Estado del estudiante
		if (stringInSlice2(v.Estado, estados)){
			estado="activo"
		}
		reglasbase = reglasBase +"estado("+codigo+", "+estado+").cursado("+codigo+", "+porcentaje+").nivel("+codigo+", "+nivel+")."

		if (stringInSlice(modalidad, modalidades)){
			comprobacion="validacion_requisitos("+codigo+")."
		}else if(modalidad==6){
			reglasbase = reglasbase +"promedio("+codigo+", "+promedio+").tipo_carrera("+codigo+", "+tipo_carrera+")."
			comprobacion="validacion_posgrado("+codigo+")."
		}else if(modalidad==7){
			reglasbase = reglasbase +"tipo_carrera("+codigo+", "+tipo_carrera+")."
			comprobacion="validacion_profundizacion("+codigo+")."
		}else if(modalidad==8){
			reglasbase = reglasbase +"tipo_carrera("+codigo+", "+tipo_carrera+")."
			comprobacion="validacion_creacion("+codigo+")."
		}
		fmt.Println(reglasbase)

		r:=golog.Comprobar(reglasbase,comprobacion)

		this.Data["json"] = r
		this.ServeJSON()


}
