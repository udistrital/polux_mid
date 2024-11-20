package helpers

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/udistrital/polux_mid/models"
	"github.com/udistrital/utils_oas/time_bogota"
)

func AddTransaccionVinculadoRegistrarNota(transaccion *models.TrVinculadoRegistrarNota) (response map[string]interface{}, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("ERROR ", err)
			panic(DeferHelpers("AddTransaccionSolicitud", err))
		}
	}()

	//alerta = append(alerta, "Success")

	url := "/v1/evaluacion_trabajo_grado"
	var resEvaluacionTrabajoGrado map[string]interface{}

	//Se registra la nota
	if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resEvaluacionTrabajoGrado, &transaccion.EvaluacionTrabajoGrado); err == nil && status == "201" {
		fmt.Println("Evaluacion_Trabajo_Trado", resEvaluacionTrabajoGrado)
		var idDocumentoTrabajoGrado = 0
		var idEvaluacionTrabajoGrado = int(resEvaluacionTrabajoGrado["Id"].(float64))
		transaccion.EvaluacionTrabajoGrado.Id = idEvaluacionTrabajoGrado

		//Si se recibió el documento del acta de sustentanción
		fmt.Println("DOCUMENTO ESCRITO", transaccion.DocumentoEscrito)
		if transaccion.DocumentoEscrito != nil {
			url := "/v1/documento_escrito"
			var resDocumentoEscrito map[string]interface{}

			//Se guarda el documento del acta de sustentanción
			if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoEscrito, &transaccion.DocumentoEscrito); err == nil && status == "201" {
				fmt.Println("Documento___Escrito", transaccion.DocumentoEscrito)
				var idDocumentoEscrito = int(resDocumentoEscrito["Id"].(float64))
				transaccion.DocumentoEscrito.Id = idDocumentoEscrito

				//Se envia el ID del acta de sustentacion y el ID del trabajo de grado de la monografía para asociarlos en la tabla documento_trabajo_grado

				url := "/v1/documento_trabajo_grado"

				var resDocumentoTrabajoGrado map[string]interface{}
				documentoTrabajoGrado := &models.DocumentoTrabajoGrado{}

				documentoTrabajoGrado.Id = 0
				documentoTrabajoGrado.DocumentoEscrito = transaccion.DocumentoEscrito
				documentoTrabajoGrado.TrabajoGrado = transaccion.TrabajoGrado

				if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoTrabajoGrado, documentoTrabajoGrado); err == nil && status == "201" {
					fmt.Println("Documento___Trabajo____Grado", documentoTrabajoGrado)
					idDocumentoTrabajoGrado = int(resDocumentoTrabajoGrado["Id"].(float64))
				} else {
					rollbackEvaluacionTrabajoGrado(transaccion)
					rollbackDocEscrito(transaccion)
					//logs.Error(err.Error())
					panic(err.Error())
				}
			} else {
				rollbackEvaluacionTrabajoGrado(transaccion)
			}
		}

		//Se obtiene la información de los roles de Docente Director y Docente Evaluador
		var parametrosRolTrabajoGrado []models.Parametro
		url := "parametro?query=CodigoAbreviacion.in:DIRECTOR_PLX|EVALUADOR_PLX,TipoParametroId__CodigoAbreviacion:ROL_TRG"
		if err := GetRequestNew("UrlCrudParametros", url, &parametrosRolTrabajoGrado); err != nil {
			fmt.Println("UrlCrudParametros!!!", parametrosRolTrabajoGrado)
			logs.Error(err.Error())
			panic(err.Error())
		}
		fmt.Println("parametrosRolTrabajoGrado", parametrosRolTrabajoGrado)
		var actualizarNotasTg bool
		var promedio float64
		var notaDirector float64

		//Se consultan los Docentes Vinculados en el trabajo de grado
		var vinculacionesTrabajoGrado []models.VinculacionTrabajoGrado
		url = "/v1/vinculacion_trabajo_grado?query=Activo:True,RolTrabajoGrado.in:" + strconv.Itoa(parametrosRolTrabajoGrado[0].Id) + "|" + strconv.Itoa(parametrosRolTrabajoGrado[1].Id) + ",TrabajoGrado__Id:" + strconv.Itoa(transaccion.TrabajoGrado.Id)
		fmt.Println("URL GET PoluxCrudUrl", url)
		if err := GetRequestNew("PoluxCrudUrl", url, &vinculacionesTrabajoGrado); err == nil {
			fmt.Println("RESPUESTA GET", vinculacionesTrabajoGrado)
			//Si la cantidad de evaluadores registrados es 1 entonces se actualiza la nota
			//fmt.Println("LenVinculaciones", len(vinculacionesTrabajoGrado))
			if len(vinculacionesTrabajoGrado) == 1 {
				actualizarNotasTg = true
				promedio = transaccion.EvaluacionTrabajoGrado.Nota
				notaDirector = transaccion.EvaluacionTrabajoGrado.Nota
			} else {
				//Si la cantidad de evaluadores registrados es mayor a 1 entonces se busca las notas registradas por cada docente vinculado
				var notasRegistradas []models.EvaluacionTrabajoGrado
				for _, vinculacion := range vinculacionesTrabajoGrado {
					var nota []models.EvaluacionTrabajoGrado
					url = "/v1/evaluacion_trabajo_grado?query=VinculacionTrabajoGrado__Id:" + strconv.Itoa(vinculacion.Id)
					if err := GetRequestNew("PoluxCrudUrl", url, &nota); err == nil {
						if len(nota) > 0 {
							if nota[0].Id != 0 {
								notasRegistradas = append(notasRegistradas, nota[0])
							}
						}

					} else {
						logs.Error(err.Error())
						panic(err.Error())
					}
				}
				//Si no se tienen todas las notas registradas por parte de los docentes, entonces no se actualiza la nota en las asignaturas
				if len(notasRegistradas) != len(vinculacionesTrabajoGrado) {
					actualizarNotasTg = false
				} else { //Si la cantidad de notas es la misma que la cantidad de vinculados, se actualiza las notas
					var promedioTemp float64
					promedioTemp = 0
					var parametroRolTrabajoGrado []models.Parametro
					url := "parametro?query=CodigoAbreviacion:DIRECTOR_PLX,TipoParametroId__CodigoAbreviacion:ROL_TRG"
					if err := GetRequestNew("UrlCrudParametros", url, &parametroRolTrabajoGrado); err != nil {
						logs.Error(err.Error())
						panic(err.Error())
					}
					//Se recorren las notas registradas y se calcula el promedio
					for _, data := range notasRegistradas {
						promedioTemp += data.Nota

						//Se trae el id del rol de trabajo para verificar si es el docente director
						var vinculacionTrabajoGrado []models.VinculacionTrabajoGrado
						url = "/v1/vinculacion_trabajo_grado?query=id:" + strconv.Itoa(data.VinculacionTrabajoGrado.Id)
						if err := GetRequestNew("PoluxCrudUrl", url, &vinculacionTrabajoGrado); err != nil {
							logs.Error(err.Error())
							panic(err.Error())
						}

						//Si el docente actual tiene el rol de Docente Director, almacena la nota para ponerla en la materia 1
						if (data.Id != 0) && (vinculacionTrabajoGrado[0].RolTrabajoGrado == parametroRolTrabajoGrado[0].Id) {
							notaDirector = data.Nota
						}
					}
					promedio = promedioTemp / float64(len(notasRegistradas))
					actualizarNotasTg = true
				}
			}
		} else {
			fmt.Println("ENTRA A ELSE", err)
			rollbackEvaluacionTrabajoGrado(transaccion)
			rollbackDocumentoTrGr(idDocumentoTrabajoGrado)
			rollbackDocEscrito(transaccion)
			//logs.Error(err.Error())
			panic(err.Error())
		}
		//Se actualizan las notas de TG teniendo en cuenta el número de evaluadores registrados y el tipo de vinculación
		if actualizarNotasTg {
			var asignaturasTrabajoGrado []models.AsignaturaTrabajoGrado

			//Se buscan las materias asociadas al trabajo de grado
			url := "/v1/asignatura_trabajo_grado?query=TrabajoGrado__Id:" + strconv.Itoa(transaccion.TrabajoGrado.Id)

			if err := GetRequestNew("PoluxCrudUrl", url, &asignaturasTrabajoGrado); err == nil {

				var estadoAsignaturas = asignaturasTrabajoGrado[0].EstadoAsignaturaTrabajoGrado
				var fechaAnterior = asignaturasTrabajoGrado[0].FechaModificacion

				for _, asignatura := range asignaturasTrabajoGrado {
					if asignatura.CodigoAsignatura == 1 { //Para la primera materia se registra la nota del Docente Director
						asignatura.Calificacion = notaDirector
					} else { //Para la segunda materia se registra el promedio de notas
						asignatura.Calificacion = promedio
					}

					//Trae el ID de "cursado" para las asignaturas de trabajo de grado
					var parametroEstadoAsignaturaTrabajoGrado []models.Parametro
					url := "parametro?query=CodigoAbreviacion.in:CDO_PLX,TipoParametroId__CodigoAbreviacion:EST_ASIG_TRG"
					if err := GetRequestNew("UrlCrudParametros", url, &parametroEstadoAsignaturaTrabajoGrado); err != nil {
						logs.Error(err.Error())
						panic(err.Error())
					}

					asignatura.EstadoAsignaturaTrabajoGrado = parametroEstadoAsignaturaTrabajoGrado[0].Id
					asignatura.FechaModificacion = time_bogota.TiempoBogotaFormato()

					//Envía los datos para actualizar el estado de las asignaturas
					url = "/v1/asignatura_trabajo_grado/" + strconv.Itoa(asignatura.Id)
					var resAsignaturaTrabajoGrado map[string]interface{}
					if status, err := SendRequestNew("PoluxCrudUrl", url, "PUT", &resAsignaturaTrabajoGrado, &asignatura); err == nil && status == "200" {
						//fmt.Println("Actualizó")
					} else {
						//fmt.Println("ERROR AQUÍ")
						rollbackEvaluacionTrabajoGrado(transaccion)
						rollbackDocumentoTrGr(idDocumentoTrabajoGrado)
						rollbackDocEscrito(transaccion)
						//logs.Error(err.Error())
						panic(err.Error())
					}
				}

				var trabajoGrado []models.TrabajoGrado

				//Se busca el proyecto de grado
				url = "/v1/trabajo_grado?query=Id:" + strconv.Itoa(transaccion.TrabajoGrado.Id)

				if err := GetRequestNew("PoluxCrudUrl", url, &trabajoGrado); err == nil {

					//Se busca el estado de "Notificado a Coordinación con calificación" para reemplazar en el trabajo de grado
					var parametroEstadoTrabajoGrado []models.Parametro
					url := "parametro?query=CodigoAbreviacion.in:NTF_PLX,TipoParametroId__CodigoAbreviacion:EST_TRG"
					if err := GetRequestNew("UrlCrudParametros", url, &parametroEstadoTrabajoGrado); err != nil {
						logs.Error(err.Error())
						panic(err.Error())
					}

					//Se inserta el ID del Estado Obtenido
					trabajoGrado[0].EstadoTrabajoGrado = parametroEstadoTrabajoGrado[0].Id

					url = "/v1/trabajo_grado/" + strconv.Itoa(trabajoGrado[0].Id)

					var resTrabajoGrado map[string]interface{}
					if status, err := SendRequestNew("PoluxCrudUrl", url, "PUT", &resTrabajoGrado, &trabajoGrado[0]); err != nil && status != "200" {
						rollbackEvaluacionTrabajoGrado(transaccion)
						rollbackDocumentoTrGr(idDocumentoTrabajoGrado)
						rollbackDocEscrito(transaccion)
						rollbackEstadoAsignaturas(asignaturasTrabajoGrado, estadoAsignaturas, fechaAnterior)
						//logs.Error(err.Error())
						panic(err.Error())
					}
				} else {
					logs.Error(err.Error())
					panic(err.Error())
				}
			} else {
				logs.Error(err.Error())
				panic(err.Error())
			}
		} else {
			//SOLO HAY UNA NOTA REGISTRADA

			var trabajoGrado []models.TrabajoGrado

			//Se busca el proyecto de grado
			url = "/v1/trabajo_grado?query=Id:" + strconv.Itoa(transaccion.TrabajoGrado.Id)

			if err := GetRequestNew("PoluxCrudUrl", url, &trabajoGrado); err == nil {

				//Se busca el estado de "Sustentado" para reemplazar en el trabajo de grado
				var parametroEstadoTrabajoGrado []models.Parametro
				url := "parametro?query=CodigoAbreviacion.in:STN_PLX,TipoParametroId__CodigoAbreviacion:EST_TRG"
				if err := GetRequestNew("UrlCrudParametros", url, &parametroEstadoTrabajoGrado); err != nil {
					logs.Error(err.Error())
					panic(err.Error())
				}

				//Se inserta el ID del Estado Obtenido
				trabajoGrado[0].EstadoTrabajoGrado = parametroEstadoTrabajoGrado[0].Id

				url = "/v1/trabajo_grado/" + strconv.Itoa(trabajoGrado[0].Id)

				var resTrabajoGrado map[string]interface{}
				if status, err := SendRequestNew("PoluxCrudUrl", url, "PUT", &resTrabajoGrado, &trabajoGrado[0]); err != nil && status != "200" {
					rollbackEvaluacionTrabajoGrado(transaccion)
					rollbackDocumentoTrGr(idDocumentoTrabajoGrado)
					rollbackDocEscrito(transaccion)
					//logs.Error(err.Error())
					panic(err.Error())
				}

			} else {
				logs.Error(err.Error())
				panic(err.Error())
			}
		}
	} else {
		//logs.Error(err)
		panic(err.Error())
	}
	response = map[string]interface{}{
		"EvaluacionTrabajoGrado": resEvaluacionTrabajoGrado,
		//"DocumentoEscrito":           resDocumentoEscrito,
		//"comentarios":          comentarios,
	}
	return response, outputError
}

func rollbackEvaluacionTrabajoGrado(transaccion *models.TrVinculadoRegistrarNota) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK EVALUACION TRABAJO GRADO")
	var respuesta map[string]interface{}
	url := "/v1/evaluacion_trabajo_grado/" + strconv.Itoa(transaccion.EvaluacionTrabajoGrado.Id)
	if status, err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil && status != "200" {
		panic("Rollback evaluacion trabajo grado" + err.Error())
	}
	return nil
}

func rollbackDocEscrito(transaccion *models.TrVinculadoRegistrarNota) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK DOCUMENTO ESCRITO")
	var respuesta map[string]interface{}
	url := "/v1/documento_escrito/" + strconv.Itoa(transaccion.DocumentoEscrito.Id)
	if status, err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil && status != "200" {
		panic("Rollback documento escrito" + err.Error())
	}
	return nil
}

func rollbackDocumentoTrGr(ID int) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK DOCUMENTO TRABAJO GRADO")
	var respuesta map[string]interface{}
	url := "/v1/documento_trabajo_grado/" + strconv.Itoa(ID)
	if status, err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil && status != "200" {
		panic("Rollback documento trabajo grado" + err.Error())
	}
	return nil
}

func rollbackEstadoAsignaturas(asignaturas []models.AsignaturaTrabajoGrado, Estado int, Fecha string) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK ESTADO ASIGNATURAS")
	var respuesta map[string]interface{}

	for _, asignatura := range asignaturas {
		asignatura.EstadoAsignaturaTrabajoGrado = Estado
		asignatura.FechaModificacion = Fecha

		url := "/v1/asignatura_trabajo_grado/" + strconv.Itoa(asignatura.Id)
		if status, err := SendRequestNew("PoluxCrudUrl", url, "PUT", &respuesta, &asignatura); err != nil && status != "200" {
			panic("Rollback estado asignatura" + err.Error())
		}
	}
	return nil
}
