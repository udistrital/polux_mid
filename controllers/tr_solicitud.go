package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/polux_mid/helpers"
	"github.com/udistrital/polux_mid/models"
	//errorControl "github.com/udistrital/utils_oas/errorctrl"
)

type TrSolicitudController struct {
	beego.Controller
}

// URLMapping ...
func (c *TrSolicitudController) URLMapping() {
	c.Mapping("Post", c.Post)
}

// @Title PostTrSolicitud
// @Description create the TrSolicitud
// @Param	body		body 	models.TrSolicitud	true	"body for TrSolicitud content"
// @Success 201 {int} models.TrSolicitud
// @Failure 400 the request contains incorrect syntax
// @router / [post]
func (c *TrSolicitudController) Post() {
	defer helpers.ErrorController(c.Controller, "TrSolicitudController")
	//Se realiza tratamiento y deserialización inicial del cuerpo para realizar el tratamiento de datosPersonalesArl que pertenece al modelo solicitudTrabajoGrado
	rawBody := c.Ctx.Input.RequestBody
	var bodyMap map[string]interface{}
	if err := json.Unmarshal(rawBody, &bodyMap); err != nil {
		panic(map[string]interface{}{"funcion": "Post", "err": err.Error(), "status": "400"})
	}
	//tratamiento del dato
	if solicitud, ok := bodyMap["Solicitud"].(map[string]interface{}); ok {
		if datosPersonales, exists := solicitud["DatosPersonalesArl"]; exists {
			serializedDatosPersonales, err := json.Marshal(datosPersonales)
			if err != nil {
				panic(map[string]interface{}{"funcion": "Post", "err": err.Error(), "status": "400"})
			}
			//Se reemplaza el campo DatosPersonalesArl con el string JSON serializado
			solicitud["DatosPersonalesArl"] = string(serializedDatosPersonales)
		}
	}
	//Se vuelve a serializar para seguir el flujo
	processedBody, err := json.Marshal(bodyMap)
	if err != nil {
		panic(map[string]interface{}{"funcion": "Post", "err": err.Error(), "status": "400"})
	}

	var v models.TrSolicitud
	if err := json.Unmarshal(processedBody, &v); err == nil {
		if response, err := helpers.AddTransaccionSolicitud(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": 201, "Message": "Solicitud realizada con exito", "Data": response}
		} else {
			// logs.Error(err)
			// fmt.Println("ERROR", err)
			// c.Data["mesaage"] = "Error service POST: The request contains an incorrect data type or an invalid parameter"
			// c.Abort("400")
			panic(err)
		}
	} else {
		// logs.Error(err)
		// c.Data["mesaage"] = "Error service POST: The request contains an incorrect data type or an invalid parameter"
		// c.Abort("400")
		fmt.Println("Se rompe en el último ELSE")
		panic(map[string]interface{}{"funcion": "Post", "err": err.Error(), "status": "400"})
	}
	c.ServeJSON()
}
