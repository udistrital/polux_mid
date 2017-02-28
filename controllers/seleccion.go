package controllers

import (

	"fmt"
	"Polux_API_mid/models"
	"strconv"
	"github.com/astaxie/beego"
	//"Polux_API_mid/golog"
	"encoding/json"
	"sort"
)

type SeleccionController struct {
	beego.Controller
}

func (c *SeleccionController) URLMapping() {
	c.Mapping("Seleccionar", c.Seleccionar)
}

func (c *SeleccionController) Seleccionar() {
	var v models.TrSolicitud

  if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		o := make(models.Vals, len(*v.Solicitudes))
		for i, x := range *v.Solicitudes {
				o[i]=x
    }
		//ordenar arreglo
		sorts := []sort.Interface{
			models.ValsAscByC{o},
		}
		for _, s := range sorts {
			sort.Sort(s)
		}


		if(v.NumAdmitidos.Cupos_excelencia>0 && len(*v.Solicitudes)>0){
			var filas int
			if (v.NumAdmitidos.Cupos_excelencia<=len(*v.Solicitudes)){
				filas=v.NumAdmitidos.Cupos_excelencia
			} else if (v.NumAdmitidos.Cupos_excelencia > len(*v.Solicitudes)){
				filas=len(*v.Solicitudes)
			}else{
				filas=0
			}
			var rta models.Solicitud

			for i := 0; i < filas; i++ {
				if err :=getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/solicitud_materias/"+strconv.Itoa(o[i].Solicitud), &rta); err == nil{
					rta.Estado="aprobado"
					//cambiar estado de la solicitud
					var respuesta interface{}
					if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/solicitud_materias/"+strconv.Itoa(o[i].Solicitud),"PUT",&respuesta ,&rta); err == nil{
						c.Data["json"] = "Solicitudes Aceptadas"
					}
				}else{
					c.Data["json"] = err.Error()
				}
			}
		}

		if(v.NumAdmitidos.Cupos_excelencia+v.NumAdmitidos.Cupos_adicionales)<=len(*v.Solicitudes){
			var rta2 models.Solicitud

			for i := v.NumAdmitidos.Cupos_excelencia; i < v.NumAdmitidos.Cupos_excelencia+v.NumAdmitidos.Cupos_adicionales; i++ {
				if err :=getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/solicitud_materias/"+strconv.Itoa(o[i].Solicitud), &rta2); err == nil{
					rta2.Estado="aprobado con pago"
					//cambiar estado de la solicitud
					var respuesta2 interface{}
					if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/solicitud_materias/"+strconv.Itoa(o[i].Solicitud),"PUT",&respuesta2 ,&rta2); err == nil{
						fmt.Println(respuesta2)
						c.Data["json"] = "Solicitudes Aceptadas"
					}
				}else{
					c.Data["json"] = err.Error()
				}
			}
		}


  } else {
    c.Data["json"] = err.Error()
  }
    c.ServeJSON()

}
