package routers

import (
	"github.com/astaxie/beego"
)

func init() {

		beego.GlobalControllerRouter["polux_api_mid/controllers:DisponibilidadController"] = append(beego.GlobalControllerRouter["polux_api_mid/controllers:DisponibilidadController"],
			beego.ControllerComments{
				Method:           "Registrar",
				Router:           `/Registrar`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})

			beego.GlobalControllerRouter["polux_api_mid/controllers:SeleccionController"] = append(beego.GlobalControllerRouter["polux_api_mid/controllers:SeleccionController"],
				beego.ControllerComments{
					Method:           "Seleccionar",
					Router:           `/Seleccionar`,
					AllowHTTPMethods: []string{"post"},
					Params:           nil})

}
