package helpers

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/polux_mid/models"
	"github.com/udistrital/utils_oas/time_bogota"
)

func AddTransaccionVinculadoRegistrarNota(transaccion *models.TrVinculadoRegistrarNota) (alerta []string, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "AddTransaccionVinculadoRegistrarNota", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
	alerta = append(alerta, "Success")

	url := "/v1/evaluacion_trabajo_grado"
	var resEvaluacionTrabajoGrado map[string]interface{}

	if err := SendRequestNew("PoluxCrudUrl", url, "POST", &resEvaluacionTrabajoGrado, &transaccion.EvaluacionTrabajoGrado); err == nil {
		var idDocumentoTrabajoGrado = 0
		var idEvaluacionTrabajoGrado = int(resEvaluacionTrabajoGrado["Id"].(float64))
		transaccion.EvaluacionTrabajoGrado.Id = idEvaluacionTrabajoGrado

		if transaccion.DocumentoEscrito != nil {
			url := "/v1/documento_escrito"
			var resDocumentoEscrito map[string]interface{}
			if err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoEscrito, &transaccion.DocumentoEscrito); err == nil {
				var idDocumentoEscrito = int(resDocumentoEscrito["Id"].(float64))
				transaccion.DocumentoEscrito.Id = idDocumentoEscrito
				url := "/v1/documento_trabajo_grado"
				var resDocumentoTrabajoGrado map[string]interface{}
				var documentoTrabajoGrado *models.DocumentoTrabajoGrado

				documentoTrabajoGrado.Id = 0
				documentoTrabajoGrado.DocumentoEscrito = transaccion.DocumentoEscrito
				documentoTrabajoGrado.TrabajoGrado = transaccion.TrabajoGrado

				if err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoTrabajoGrado, documentoTrabajoGrado); err == nil {
					idDocumentoTrabajoGrado = int(resDocumentoTrabajoGrado["Id"].(float64))
					fmt.Println(idDocumentoTrabajoGrado)
				} else {
					rollbackEvaluacionTrabajoGrado(transaccion)
					rollbackDocEscrito(transaccion)
					logs.Error(err.Error())
					panic(err.Error())
				}
			} else {
				rollbackEvaluacionTrabajoGrado(transaccion)
			}
		}

		var parametrosRolTrabajoGrado []models.Parametro
		url := "parametro?query=CodigoAbreviacion.in:DIRECTOR_PLX|EVALUADOR_PLX,TipoParametroId__CodigoAbreviacion:ROL_TRG"
		if err := GetRequestNew("UrlCrudParametros", url, &parametrosRolTrabajoGrado); err != nil {
			logs.Error(err.Error())
			panic(err.Error())
		}

		var actualizarNotasTg bool
		var promedio float64
		var notaDirector float64

		var vinculacionesTrabajoGrado []models.VinculacionTrabajoGrado
		url = beego.AppConfig.String("PoluxCrudUrl") + "/v1/vinculacion_trabajo_grado?query=Activo:True,RolTrabajoGrado.in:" + strconv.Itoa(parametrosRolTrabajoGrado[0].Id) + "|" + strconv.Itoa(parametrosRolTrabajoGrado[1].Id) + ",TrabajoGrado__Id:" + strconv.Itoa(transaccion.TrabajoGrado.Id)
		if err := GetJson(url, &vinculacionesTrabajoGrado); err == nil {
			if len(vinculacionesTrabajoGrado) == 1 {
				actualizarNotasTg = true
				promedio = transaccion.EvaluacionTrabajoGrado.Nota
				notaDirector = transaccion.EvaluacionTrabajoGrado.Nota
				fmt.Println(actualizarNotasTg, promedio, notaDirector)
			} else {
				var notasRegistradas []models.EvaluacionTrabajoGrado
				for _, vinculacion := range vinculacionesTrabajoGrado {
					var nota []models.EvaluacionTrabajoGrado
					url = beego.AppConfig.String("PoluxCrudUrl") + "/v1/evaluacion_trabajo_grado?query=VinculacionTrabajoGrado__Id:" + strconv.Itoa(vinculacion.Id)
					print("URL:", url)
					if err := GetJson(url, &nota); err == nil {
						notasRegistradas = append(notasRegistradas, nota[0])
					} else {
						logs.Error(err.Error())
						panic(err.Error())
					}
				}
				if len(notasRegistradas) != len(vinculacionesTrabajoGrado) {
					actualizarNotasTg = true
				} else {
					var promedioTemp float64
					promedioTemp = 0
					var parametroRolTrabajoGrado []models.Parametro
					url := "parametro?query=CodigoAbreviacion:DIRECTOR_PLX,TipoParametroId__CodigoAbreviacion:ROL_TRG"
					if err := GetRequestNew("UrlCrudParametros", url, &parametroRolTrabajoGrado); err != nil {
						logs.Error(err.Error())
						panic(err.Error())
					}
					for _, data := range notasRegistradas {
						promedioTemp += data.Nota
						if (data.Id != 0) && (data.VinculacionTrabajoGrado.RolTrabajoGrado == parametroRolTrabajoGrado[0].Id) {
							notaDirector = data.Nota
						}
					}
					promedio = promedioTemp / float64(len(notasRegistradas))
					actualizarNotasTg = true
				}
			}
		} else {
			rollbackEvaluacionTrabajoGrado(transaccion)
			rollbackDocEscrito(transaccion)
			rollbackDocumentoTrGr(idDocumentoTrabajoGrado)
			logs.Error(err.Error())
			panic(err.Error())
		}
		if actualizarNotasTg == true {
			var asignaturasTrabajoGrado []models.AsignaturaTrabajoGrado
			url := beego.AppConfig.String("PoluxCrudUrl") + "/v1/asignatura_trabajo_grado?query=TrabajoGrado__Id:" + strconv.Itoa(transaccion.TrabajoGrado.Id)

			if err := GetJson(url, &asignaturasTrabajoGrado); err == nil {
				for _, asignatura := range asignaturasTrabajoGrado {
					if asignatura.CodigoAsignatura == 1 {
						asignatura.Calificacion = notaDirector
					} else {
						asignatura.Calificacion = promedio
					}

					var parametroEstadoAsignaturaTrabajoGrado []models.Parametro
					url := "parametro?query=CodigoAbreviacion.in:CDO_PLX,TipoParametroId__CodigoAbreviacion:EST_ASIG_TRG"
					if err := GetRequestNew("UrlCrudParametros", url, &parametroEstadoAsignaturaTrabajoGrado); err != nil {
						logs.Error(err.Error())
						panic(err.Error())
					}

					asignatura.EstadoAsignaturaTrabajoGrado = parametroEstadoAsignaturaTrabajoGrado[0].Id
					asignatura.FechaModificacion = time_bogota.TiempoBogotaFormato()

					url = "/v1/asignatura_trabajo_grado/" + strconv.Itoa(asignatura.Id)
					var resAsignaturaTrabajoGrado map[string]interface{}
					if err := SendRequestNew("PoluxCrudUrl", url, "UPDATE", &resAsignaturaTrabajoGrado, &asignatura); err == nil {
						//CONTINUAR AQUÍ ......
						fmt.Println("Actualizó")
					} else {
						fmt.Println("ERROR AQUÍ")
						rollbackEvaluacionTrabajoGrado(transaccion)
						rollbackDocEscrito(transaccion)
						rollbackDocumentoTrGr(idDocumentoTrabajoGrado)
						logs.Error(err.Error())
						panic(err.Error())
					}
				}
			} else {
				logs.Error(err.Error())
				panic(err.Error())
			}
		}
	} else {
		logs.Error(err)
		panic(err.Error())
	}
	return alerta, outputError
}

func rollbackEvaluacionTrabajoGrado(transaccion *models.TrVinculadoRegistrarNota) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK EVALUACION TRABAJO GRADO")
	var respuesta map[string]interface{}
	url := "/v1/evaluacion_trabajo_grado/" + strconv.Itoa(transaccion.EvaluacionTrabajoGrado.Id)
	if err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil {
		panic("Rollback evaluacion trabajo grado" + err.Error())
	}
	return nil
}

func rollbackDocEscrito(transaccion *models.TrVinculadoRegistrarNota) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK DOCUMENTO ESCRITO")
	var respuesta map[string]interface{}
	url := "/v1/documento_escrito/" + strconv.Itoa(transaccion.DocumentoEscrito.Id)
	if err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil {
		panic("Rollback documento escrito" + err.Error())
	}
	return nil
}

func rollbackDocumentoTrGr(ID int) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK DOCUMENTO TRABAJO GRADO")
	var respuesta map[string]interface{}
	url := "/v1/documento_trabajo_grado/" + strconv.Itoa(ID)
	if err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil {
		panic("Rollback documento trabajo grado" + err.Error())
	}
	return nil
}
