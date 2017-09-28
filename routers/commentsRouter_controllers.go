package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/udistrital/Polux_API_mid/controllers:CreditosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/Polux_API_mid/controllers:CreditosController"],
		beego.ControllerComments{
			Method: "ObtenerMinimo",
			Router: `/ObtenerMinimo`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/Polux_API_mid/controllers:CuposController"] = append(beego.GlobalControllerRouter["github.com/udistrital/Polux_API_mid/controllers:CuposController"],
		beego.ControllerComments{
			Method: "Obtener",
			Router: `/Obtener`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/Polux_API_mid/controllers:FechasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/Polux_API_mid/controllers:FechasController"],
		beego.ControllerComments{
			Method: "ObtenerFechas",
			Router: `/ObtenerFechas`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/Polux_API_mid/controllers:SeleccionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/Polux_API_mid/controllers:SeleccionController"],
		beego.ControllerComments{
			Method: "Seleccionar",
			Router: `/Seleccionar`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/Polux_API_mid/controllers:VerificarRequisitosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/Polux_API_mid/controllers:VerificarRequisitosController"],
		beego.ControllerComments{
			Method: "CantidadModalidades",
			Router: `/CantidadModalidades`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/Polux_API_mid/controllers:VerificarRequisitosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/Polux_API_mid/controllers:VerificarRequisitosController"],
		beego.ControllerComments{
			Method: "Registrar",
			Router: `/Registrar`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

}
