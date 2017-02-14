package controllers

import (

	"fmt"
	"polux_api_mid/models"
	//"strconv"
	"github.com/astaxie/beego"
	"polux_api_mid/golog"
)

type DisponibilidadController struct {
	beego.Controller
}

func (c *DisponibilidadController) URLMapping() {
	c.Mapping("Registrar", c.Registrar)
}

func stringInSlice(str int, list []int) bool {
 for _, v := range list {
	 if v == str {
		 return true
	 }
 }
 return false
}

func (this *DisponibilidadController) Registrar() {

	var predicados []models.Predicado
	var postdominio string = ""
	var comprobacion string = ""
	if tdominio  := this.GetString("tdominio"); tdominio != "" {
			postdominio = postdominio +"&query=Dominio.Id:"+tdominio
	}else{
		this.Data["json"] = "no se especifico el domino del ruler"
		this.ServeJSON()
	}

	if err := getJson("http://"+beego.AppConfig.String("Urlruler")+":"+beego.AppConfig.String("Portruler")+"/"+beego.AppConfig.String("Nsruler")+"/predicado?limit=0"+postdominio, &predicados); err == nil{
		var reglasbase string = ""

		var arregloReglas = make([]string, len(predicados))
		for i := 0; i < len(predicados); i++ {
			arregloReglas[i] = predicados[i].Nombre
		}

		for i := 0; i < len(arregloReglas); i++ {
			reglasbase = reglasbase + arregloReglas[i]
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

		codigo:="20102020007"
		estado:="activo"
		porcentaje:="60"
		promedio:="4.2"
		nivel:="pregrado"
		tipo_carrera:="artes"

		modalidad:=8
		modalidades := []int{1, 2, 3, 4, 5} //Modalidades que solo necesitan el Porcentaje cursado y el Estado del estudiante

		reglasbase = reglasbase +"estado("+codigo+", "+estado+").cursado("+codigo+", "+porcentaje+").nivel("+codigo+", "+nivel+")."

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
	//	reglasbase = reglasbase +"promedio(20102020007,4.6).tipo(20102020007,ARTES).estado(20102020007, activo).cursado(20102020007, 90).modalidad(20102020007, materias_posgrado)."
		r:=golog.Comprobar(reglasbase,comprobacion)
		if(r=="true"){
			this.Data["json"] = r
			this.ServeJSON()
		}else{
			this.Data["json"] = "No cumple con los requisitos exigidos para la Modalidad de TG"
			this.ServeJSON()
		}

	}


}
