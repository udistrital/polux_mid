package controllers

import (
	"fmt"
	"github.com/udistrital/Polux_API_mid/models"
	"strconv"
	"github.com/astaxie/beego"
	"github.com/udistrital/Polux_API_mid/golog"
)

type CuposController struct {
	beego.Controller
}

func (c *CuposController) URLMapping() {
	c.Mapping("Obtener", c.Obtener)
}

func (c *CuposController) Obtener() {
  var NumAdmitidos models.Cupos
  //consultar las reglas
  var predicados []models.Predicado
  var postdominio string = ""
  var comprobacion string = ""
  if tdominio  := c.GetString("tdominio"); tdominio != "" {
      postdominio = postdominio +"&query=Dominio.Id:"+tdominio
  }else{
    c.Data["json"] = "no se especifico el domino del ruler"
    c.ServeJSON()
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

  //obtener máximo de cupos por excelencia académica
  comprobacion="max_cupos_excelencia_academica(Y)."
  r:=golog.Obtener(reglasbase,comprobacion)
  fmt.Println(r)
  NumAdmitidos.Cupos_excelencia, err = strconv.Atoi(r)
   if err != nil {
      fmt.Println(err)
   }

  //obtener máximo de cupos adicionales
  comprobacion="max_cupos_adicionales(Y)."
  r2:=golog.Obtener(reglasbase,comprobacion)
  fmt.Println(r2)

  NumAdmitidos.Cupos_adicionales, err = strconv.Atoi(r2)
   if err != nil {
      fmt.Println(err)
   }
  c.Data["json"] = NumAdmitidos
  c.ServeJSON()
  ///////////////////////////////////////////////////////////////////////////
  }
}
