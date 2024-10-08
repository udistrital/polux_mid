// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/polux_mid/controllers"
)

func init() {

	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/verificarRequisitos",
			beego.NSInclude(
				&controllers.VerificarRequisitosController{},
			),
		),
		beego.NSNamespace("/cupos",
			beego.NSInclude(
				&controllers.CuposController{},
			),
		),
		beego.NSNamespace("/fechas",
			beego.NSInclude(
				&controllers.FechasController{},
			),
		),
		beego.NSNamespace("/creditos",
			beego.NSInclude(
				&controllers.CreditosController{},
			),
		),
		beego.NSNamespace("/evaluadores",
			beego.NSInclude(
				&controllers.EvaluadoresController{},
			),
		),
		beego.NSNamespace("/creditos_materias",
			beego.NSInclude(
				&controllers.CreditosMateriasController{},
			),
		),
		beego.NSNamespace("/tr_respuesta_solicitud",
			beego.NSInclude(
				&controllers.TrRespuestaSolicitudController{},
			),
		),
		beego.NSNamespace("/tr_solicitud",
			beego.NSInclude(
				&controllers.TrSolicitudController{},
			),
		),
		beego.NSNamespace("/tr_registrar_revision_tg",
			beego.NSInclude(
				&controllers.TrRegistrarRevisionTgController{},
			),
		),
		beego.NSNamespace("/tr_vinculado_registrar_nota",
			beego.NSInclude(
				&controllers.TrVinculadoRegistrarNotaController{},
			),
		),
		beego.NSNamespace("/tr_subir_arl",
			beego.NSInclude(
				&controllers.TrSubirArl{},
			),
		),
		beego.NSNamespace("/tr_registrar_acta_seguimiento",
			beego.NSInclude(
				&controllers.TrRegistrarActaSeguimiento{},
			),
		),
		beego.NSNamespace("/reporte_general",
			beego.NSInclude(
				&controllers.ReporteGeneralController{},
			),
		),
		beego.NSNamespace("/reporte_solicitud",
			beego.NSInclude(
				&controllers.ReporteSolicitudController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
