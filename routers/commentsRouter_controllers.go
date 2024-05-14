package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/polux_mid/controllers:CreditosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/polux_mid/controllers:CreditosController"],
        beego.ControllerComments{
            Method: "ObtenerMinimo",
            Router: "/ObtenerMinimo",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/polux_mid/controllers:CreditosMateriasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/polux_mid/controllers:CreditosMateriasController"],
        beego.ControllerComments{
            Method: "ObtenerCreditos",
            Router: "/ObtenerCreditos",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/polux_mid/controllers:CuposController"] = append(beego.GlobalControllerRouter["github.com/udistrital/polux_mid/controllers:CuposController"],
        beego.ControllerComments{
            Method: "Obtener",
            Router: "/Obtener",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/polux_mid/controllers:EvaluadoresController"] = append(beego.GlobalControllerRouter["github.com/udistrital/polux_mid/controllers:EvaluadoresController"],
        beego.ControllerComments{
            Method: "ObtenerEvaluadores",
            Router: "/ObtenerEvaluadores",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/polux_mid/controllers:FechasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/polux_mid/controllers:FechasController"],
        beego.ControllerComments{
            Method: "ObtenerFechas",
            Router: "/ObtenerFechas",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/polux_mid/controllers:TrRegistrarRevisionTgController"] = append(beego.GlobalControllerRouter["github.com/udistrital/polux_mid/controllers:TrRegistrarRevisionTgController"],
        beego.ControllerComments{
            Method: "Post",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/polux_mid/controllers:TrRespuestaSolicitudController"] = append(beego.GlobalControllerRouter["github.com/udistrital/polux_mid/controllers:TrRespuestaSolicitudController"],
        beego.ControllerComments{
            Method: "Post",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/polux_mid/controllers:TrSolicitudController"] = append(beego.GlobalControllerRouter["github.com/udistrital/polux_mid/controllers:TrSolicitudController"],
        beego.ControllerComments{
            Method: "Post",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/polux_mid/controllers:TrVinculadoRegistrarNotaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/polux_mid/controllers:TrVinculadoRegistrarNotaController"],
        beego.ControllerComments{
            Method: "Post",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/polux_mid/controllers:VerificarRequisitosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/polux_mid/controllers:VerificarRequisitosController"],
        beego.ControllerComments{
            Method: "CantidadModalidades",
            Router: "/CantidadModalidades",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/polux_mid/controllers:VerificarRequisitosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/polux_mid/controllers:VerificarRequisitosController"],
        beego.ControllerComments{
            Method: "Registrar",
            Router: "/Registrar",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
