package helpers

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/polux_mid/models"
	"github.com/udistrital/utils_oas/time_bogota"
)

func AddTransaccionRespuestaSolicitud(transaccion *models.TrRespuestaSolicitud) (alerta []string, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "AddTransaccionRespuestaSolicitud", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
	alerta = append(alerta, "Success")

	url := "parametro/" + strconv.Itoa(transaccion.ModalidadTipoSolicitud.Modalidad)
	var parametro models.Parametro
	var parametroEstadoSolicitud models.Parametro

	var TrAnteproyecto = transaccion

	if err := GetRequestNew("UrlCrudParametros", url, &parametro); err != nil {
		logs.Error(err.Error())
		panic(err.Error())
	}
	url = "parametro/" + strconv.Itoa(transaccion.RespuestaAnterior.EstadoSolicitud)
	if err := GetRequestNew("UrlCrudParametros", url, &parametroEstadoSolicitud); err != nil {
		logs.Error(err.Error())
		panic(err.Error())
	}
	var resRespuestaAnterior map[string]interface{}
	url = "/v1/respuesta_solicitud/" + strconv.Itoa(transaccion.RespuestaAnterior.Id)
	//payload := "/" + strconv.Itoa(transaccion.RespuestaAnterior.Id)
	//if outputError = Post(new(models.RespuestaSolicitud), payload, &resRespuestaAnterior); outputError == nil {
	if err := SendRequestNew("PoluxCrudUrl", url, "PUT", &resRespuestaAnterior, &transaccion.RespuestaAnterior); err == nil {
		url = "/v1/respuesta_solicitud"
		var resRespuestaNueva map[string]interface{}
		if err := SendRequestNew("PoluxCrudUrl", url, "POST", &resRespuestaNueva, &transaccion.RespuestaNueva); err == nil {
			transaccion.RespuestaNueva.Id = int(resRespuestaNueva["Id"].(float64))
			url = "/v1/documento_solicitud"
			var resDocumentoSolicitud map[string]interface{}
			if err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoSolicitud, &transaccion.DocumentoSolicitud); err == nil {
				transaccion.DocumentoSolicitud.Id = int(resDocumentoSolicitud["Id"].(float64))
			} else {
				logs.Error(err)
				rollbackResNueva(transaccion)
			}
		} else {
			logs.Error(err)
			rollbackResAnterior(transaccion)
		}
	} else {
		logs.Error(err)
		panic(err.Error())
	}

	if transaccion.TrTrabajoGrado != nil && (parametro.CodigoAbreviacion != "EAPOS_PLX" || parametroEstadoSolicitud.CodigoAbreviacion == "ACPR_PLX") {
		url = "/v1/trabajo_grado"
		var resTrabajoGrado map[string]interface{}
		if err := SendRequestNew("PoluxCrudUrl", url, "POST", &resTrabajoGrado, &transaccion.TrTrabajoGrado.TrabajoGrado); err == nil {
			var idTrabajoGrado = int(resTrabajoGrado["Id"].(float64))
			transaccion.TrTrabajoGrado.TrabajoGrado.Id = idTrabajoGrado
			transaccion.SolicitudTrabajoGrado.TrabajoGrado.Id = idTrabajoGrado
			var resSolicitudTrabajoGrado map[string]interface{}
			url = "/v1/solicitud_trabajo_grado/" + strconv.Itoa(transaccion.SolicitudTrabajoGrado.Id)
			if err := SendRequestNew("PoluxCrudUrl", url, "PUT", &resSolicitudTrabajoGrado, &transaccion.SolicitudTrabajoGrado); err == nil {
				url = "/v1/asignatura_trabajo_grado"
				var materias = make([]map[string]interface{}, 0)
				for i, v := range *transaccion.TrTrabajoGrado.AsignaturasTrabajoGrado {
					var resAsignaturaTrabajoGrado map[string]interface{}
					v.TrabajoGrado.Id = idTrabajoGrado
					if err := SendRequestNew("PoluxCrudUrl", url, "POST", &resAsignaturaTrabajoGrado, &v); err == nil {
						(*transaccion.TrTrabajoGrado.AsignaturasTrabajoGrado)[i].Id = int(resAsignaturaTrabajoGrado["Id"].(float64))
						materias = append(materias, resAsignaturaTrabajoGrado)
					} else {
						logs.Error(err)
						if len(materias) > 0 {
							rollbackAsignaturasTrabajoGrado(transaccion)
						} else {
							rollbackSolicitudTrabajoGrado(transaccion)
						}
					}
				}
				url = "/v1/estudiante_trabajo_grado"
				var estudiantes = make([]map[string]interface{}, 0)
				for i, v := range *transaccion.TrTrabajoGrado.EstudianteTrabajoGrado {
					var resEstudianteTrabajoGrado map[string]interface{}
					v.TrabajoGrado.Id = idTrabajoGrado
					if err := SendRequestNew("PoluxCrudUrl", url, "POST", &resEstudianteTrabajoGrado, &v); err == nil {
						(*transaccion.TrTrabajoGrado.EstudianteTrabajoGrado)[i].Id = int(resEstudianteTrabajoGrado["Id"].(float64))
						estudiantes = append(estudiantes, resEstudianteTrabajoGrado)
					} else {
						logs.Error(err)
						if len(estudiantes) > 0 {
							rollbackEstudianteTrabajoGrado(transaccion)
						} else {
							rollbackAsignaturasTrabajoGrado(transaccion)
						}
					}
				}
				if transaccion.TrTrabajoGrado.AreasTrabajoGrado != nil {
					url = "/v1/areas_trabajo_grado"
					var areas = make([]map[string]interface{}, 0)
					for i, v := range *transaccion.TrTrabajoGrado.AreasTrabajoGrado {
						var resAreasTrabajoGrado map[string]interface{}
						v.TrabajoGrado.Id = idTrabajoGrado
						if err := SendRequestNew("PoluxCrudUrl", url, "POST", &resAreasTrabajoGrado, &v); err == nil {
							(*transaccion.TrTrabajoGrado.AreasTrabajoGrado)[i].Id = int(resAreasTrabajoGrado["Id"].(float64))
							areas = append(areas, resAreasTrabajoGrado)
						} else {
							logs.Error(err)
							if len(areas) > 0 {
								rollbackAreasTrabajoGrado(transaccion)
							} else {
								rollbackEstudianteTrabajoGrado(transaccion)
							}
						}
					}
				}
				if transaccion.TrTrabajoGrado.VinculacionTrabajoGrado != nil {
					url = "/v1/vinculacion_trabajo_grado"
					var vinculaciones = make([]map[string]interface{}, 0)
					for i, v := range *transaccion.TrTrabajoGrado.VinculacionTrabajoGrado {
						var resVinculacionTrabajoGrado map[string]interface{}
						v.TrabajoGrado.Id = idTrabajoGrado
						if err := SendRequestNew("PoluxCrudUrl", url, "POST", &resVinculacionTrabajoGrado, &v); err == nil {
							(*transaccion.TrTrabajoGrado.VinculacionTrabajoGrado)[i].Id = int(resVinculacionTrabajoGrado["Id"].(float64))
							vinculaciones = append(vinculaciones, resVinculacionTrabajoGrado)
						} else {
							logs.Error(err)
							if len(vinculaciones) > 0 {
								rollbackVinculacionTrabajoGrado(transaccion)
							} else {
								rollbackAreasTrabajoGrado(transaccion)
							}
						}
					}
				}
				if transaccion.TrTrabajoGrado.DocumentoEscrito != nil {

					TrAnteproyecto = transaccion

					url = "/v1/documento_escrito"
					var resDocumentoEscrito map[string]interface{}
					if err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoEscrito, &transaccion.TrTrabajoGrado.DocumentoEscrito); err == nil { // Se guarda el documento con tipo_documento = DTR_PLX
						transaccion.TrTrabajoGrado.DocumentoTrabajoGrado.TrabajoGrado.Id = idTrabajoGrado
						transaccion.TrTrabajoGrado.DocumentoTrabajoGrado.DocumentoEscrito.Id = int(resDocumentoEscrito["Id"].(float64))
						transaccion.TrTrabajoGrado.DocumentoEscrito.Id = int(resDocumentoEscrito["Id"].(float64))
						var resDocumentoTrabajoGrado map[string]interface{}
						url = "/v1/documento_trabajo_grado"
						if err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoTrabajoGrado, &transaccion.TrTrabajoGrado.DocumentoTrabajoGrado); err == nil { //Se asocia el documento con el trabajo de grado
							transaccion.TrTrabajoGrado.DocumentoTrabajoGrado.Id = int(resDocumentoTrabajoGrado["Id"].(float64))

							//Funcionalidad para almacenar el Anteproyecto

							//Se recupera el id del tipo de documento ANP_PLX
							var tipoDocumento []models.TipoDocumento
							url = beego.AppConfig.String("UrlDocumentos") + "tipo_documento?query=CodigoAbreviacion:ANP_PLX"
							if err := GetJson(url, &tipoDocumento); err != nil {
								logs.Error(err.Error())
								panic(err.Error())
							}

							//Se prepara el documento, reseteando los ids, estableciendo el tipo de documento y cambiando el resumen y el titulo del documento
							TrAnteproyecto.TrTrabajoGrado.DocumentoTrabajoGrado.TrabajoGrado.Id = 0
							TrAnteproyecto.TrTrabajoGrado.DocumentoTrabajoGrado.DocumentoEscrito.Id = 0
							TrAnteproyecto.TrTrabajoGrado.DocumentoEscrito.Id = 0
							transaccion.TrTrabajoGrado.DocumentoTrabajoGrado.Id = 0
							TrAnteproyecto.TrTrabajoGrado.DocumentoEscrito.TipoDocumentoEscrito = tipoDocumento[0].Id
							TrAnteproyecto.TrTrabajoGrado.DocumentoEscrito.Resumen = "Anteproyecto del trabajo de grado con ID: " + strconv.Itoa(TrAnteproyecto.TrTrabajoGrado.TrabajoGrado.Id) + " Nombre: " + TrAnteproyecto.TrTrabajoGrado.TrabajoGrado.Titulo
							TrAnteproyecto.TrTrabajoGrado.DocumentoEscrito.Titulo = "Anteproyecto " + TrAnteproyecto.TrTrabajoGrado.DocumentoEscrito.Titulo

							url = "/v1/documento_escrito"
							var resDocumentoEscrito map[string]interface{}
							if err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoEscrito, &TrAnteproyecto.TrTrabajoGrado.DocumentoEscrito); err == nil { // Se guarda el documento con tipo_documento = ANP_PLX
								TrAnteproyecto.TrTrabajoGrado.DocumentoTrabajoGrado.TrabajoGrado.Id = idTrabajoGrado
								TrAnteproyecto.TrTrabajoGrado.DocumentoTrabajoGrado.DocumentoEscrito.Id = int(resDocumentoEscrito["Id"].(float64))
								TrAnteproyecto.TrTrabajoGrado.DocumentoEscrito.Id = int(resDocumentoEscrito["Id"].(float64))

								var resDocumentoTrabajoGradoAnt map[string]interface{}
								url = "/v1/documento_trabajo_grado"
								if err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoTrabajoGradoAnt, &TrAnteproyecto.TrTrabajoGrado.DocumentoTrabajoGrado); err == nil { //Se asocia el anteproyecto con el trabajo de grado
									TrAnteproyecto.TrTrabajoGrado.DocumentoTrabajoGrado.Id = int(resDocumentoTrabajoGradoAnt["Id"].(float64))
								} else {
									logs.Error(err)
									rollbackDocumentoTrabajoGrado(TrAnteproyecto)
								}
							} else {
								logs.Error(err)
								rollbackDocumentoEscrito(TrAnteproyecto)
							}

						} else {
							logs.Error(err)
							rollbackDocumentoTrabajoGrado(transaccion)

						}
					} else {
						logs.Error(err)
						rollbackDocumentoEscrito(transaccion)
					}
				}
				if transaccion.EspaciosAcademicosInscritos != nil {
					url = "/v1/espacio_academico_inscrito"
					var espacios_academicos_inscritos = make([]map[string]interface{}, 0)
					for i, v := range *transaccion.EspaciosAcademicosInscritos {
						var resEspaciosAcademicosInscritos map[string]interface{}
						v.TrabajoGrado.Id = idTrabajoGrado
						if err := SendRequestNew("PoluxCrudUrl", url, "POST", &resEspaciosAcademicosInscritos, &v); err == nil {
							(*transaccion.EspaciosAcademicosInscritos)[i].Id = int(resEspaciosAcademicosInscritos["Id"].(float64))
							espacios_academicos_inscritos = append(espacios_academicos_inscritos, resEspaciosAcademicosInscritos)
						} else {
							logs.Error(err)
							if len(espacios_academicos_inscritos) > 0 {
								rollbackEspaciosAcademicosInscritos(transaccion)
							} else {
								rollbackDocumentoTrabajoGrado(transaccion)
							}
						}
					}
				}
				if transaccion.DetallesPasantia != nil {
					transaccion.DetallesPasantia.TrabajoGrado.Id = idTrabajoGrado
					url = "/v1/detalle_pasantia"
					var resDetallePasantia map[string]interface{}
					if err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDetallePasantia, &transaccion.DetallesPasantia); err == nil {
						transaccion.DetallesPasantia.Id = int(resDetallePasantia["Id"].(float64))

						//Se guardan los documentos asociados a la pasantia

						transaccion.DetallesPasantia.Contrato.Titulo = transaccion.DetallesPasantia.Contrato.Titulo + " del trabajo de grado con id: " + strconv.Itoa(idTrabajoGrado)

						transaccion.DetallesPasantia.Carta.Titulo = transaccion.DetallesPasantia.Carta.Titulo + " del trabajo de grado con id: " + strconv.Itoa(idTrabajoGrado)

						//Se envia el contrato
						url = "/v1/documento_escrito"
						var resDocumentoEscrito map[string]interface{}
						if err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoEscrito, &transaccion.DetallesPasantia.Contrato); err == nil { // Se guarda el contrato
							transaccion.DetallesPasantia.DTG_Contrato.TrabajoGrado.Id = idTrabajoGrado
							transaccion.DetallesPasantia.Contrato.Id = int(resDocumentoEscrito["Id"].(float64))
							transaccion.DetallesPasantia.DTG_Contrato.DocumentoEscrito.Id = int(resDocumentoEscrito["Id"].(float64))

							var resDocumentoTrabajoGrado map[string]interface{}
							url = "/v1/documento_trabajo_grado"
							if err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoTrabajoGrado, &transaccion.DetallesPasantia.DTG_Contrato); err == nil { //Se asocia el contrato con el trabajo de grado
								transaccion.DetallesPasantia.DTG_Contrato.Id = int(resDocumentoTrabajoGrado["Id"].(float64))

								//Se envia la Carta
								url = "/v1/documento_escrito"
								var resDocumentoEscrito map[string]interface{}
								if err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoEscrito, &transaccion.DetallesPasantia.Carta); err == nil { // Se guarda la carta
									transaccion.DetallesPasantia.Carta.Id = int(resDocumentoEscrito["Id"].(float64))
									transaccion.DetallesPasantia.DTG_Carta.TrabajoGrado.Id = idTrabajoGrado
									transaccion.DetallesPasantia.DTG_Carta.DocumentoEscrito.Id = int(resDocumentoEscrito["Id"].(float64))

									var resDocumentoTrabajoGrado map[string]interface{}
									url = "/v1/documento_trabajo_grado"
									if err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoTrabajoGrado, &transaccion.DetallesPasantia.DTG_Carta); err == nil { //Se asocia la carta con el trabajo de grado
										transaccion.DetallesPasantia.DTG_Carta.Id = int(resDocumentoTrabajoGrado["Id"].(float64))

									} else {
										logs.Error(err)
										rollbackDocumentoTrabajoGradoPasantia(transaccion)
									}
								} else {
									logs.Error(err)
									rollbackDocumentoEscritoPasantia(transaccion)
								}

							} else {
								logs.Error(err)
								rollbackDocumentoTrabajoGradoPasantia(transaccion)
							}
						} else {
							logs.Error(err)
							rollbackDocumentoEscritoPasantia(transaccion)
						}

					} else {
						logs.Error(err)
						rollbackDetallesPasantia(transaccion)
					}
				}
				if transaccion.DetallesPasantiaExterna != nil && parametro.CodigoAbreviacion == "PASIN_PLX" {
					url = "/v1/detalle_trabajo_grado"
					var detalles_pasantia_externa = make([]map[string]interface{}, 0)
					for i, v := range *transaccion.DetallesPasantiaExterna {
						var resDetallesPasantiaExterna map[string]interface{}
						v.TrabajoGrado.Id = idTrabajoGrado
						if err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDetallesPasantiaExterna, &v); err == nil {
							(*transaccion.DetallesPasantiaExterna)[i].Id = int(resDetallesPasantiaExterna["Id"].(float64))
							detalles_pasantia_externa = append(detalles_pasantia_externa, resDetallesPasantiaExterna)
						} else {
							logs.Error(err)
							if len(detalles_pasantia_externa) > 0 {
								rollbackDetallesPasantiaExterna(transaccion)
							} else {
								rollbackDocumentoTrabajoGrado(transaccion)
							}
						}
					}
				}
			} else {
				logs.Error(err)
				rollbackTrabajoGrado(transaccion)
			}
		} else {
			logs.Error(err)
			rollbackDocumentoSolicitud(transaccion)
		}
	}

	// Solicitud de cambio de docente evaluador o docente director
	if transaccion.Vinculaciones != nil {
		var idVinculadoAntiguo int
		var idVinculadoNuevo int64
		var vinculaciones_trabajo_grado = make([]map[string]interface{}, 0)
		var vinculaciones_originales_trabajo_grado []models.VinculacionTrabajoGrado
		var vinculaciones_trabajo_grado_post = make([]map[string]interface{}, 0)
		var vinculaciones_trabajo_grado_canceladas = make([]map[string]interface{}, 0)
		for _, v := range *transaccion.Vinculaciones {
			//Si esta activo es nuevo y se inserta sino se actualiza la fecha de fin y el activo
			if v.Activo {
				// Se buscar si el docente ya estuvo vinculado y se actualiza
				var vinculado []models.VinculacionTrabajoGrado
				url = beego.AppConfig.String("PoluxCrudUrl") + "/v1/vinculacion_trabajo_grado?query=TrabajoGrado__Id:" + strconv.Itoa(v.TrabajoGrado.Id) +
					",Usuario:" + strconv.Itoa(v.Usuario) + ",RolTrabajoGrado:" + strconv.Itoa(v.RolTrabajoGrado) + "&limit=1"
				fmt.Println("URL ", url)
				if err := GetJson(url, &vinculado); err != nil {
					logs.Error(err.Error())
					panic(err.Error())
				}
				if vinculado[0].Id != 0 {
					var vinculadoAux = vinculado[0]
					idVinculadoNuevo = int64(vinculado[0].Id)
					vinculado[0].Activo = v.Activo
					vinculado[0].FechaFin = v.FechaFin
					vinculado[0].FechaInicio = v.FechaInicio
					var resVinculacionTrabajoGrado map[string]interface{}
					url = "/v1/vinculacion_trabajo_grado/" + strconv.Itoa(vinculado[0].Id)
					if err := SendRequestNew("PoluxCrudUrl", url, "PUT", &resVinculacionTrabajoGrado, &vinculado[0]); err == nil {
						vinculaciones_originales_trabajo_grado = append(vinculaciones_originales_trabajo_grado, vinculadoAux)
						vinculaciones_trabajo_grado = append(vinculaciones_trabajo_grado, resVinculacionTrabajoGrado)
					} else {
						logs.Error(err)
						if len(vinculaciones_trabajo_grado) > 0 || len(vinculaciones_trabajo_grado_post) > 0 {
							rollbackVinculacionTrabajoGradoRS(transaccion, vinculaciones_originales_trabajo_grado)
							rollbackVinculacionTrabajoGradoRSPost(transaccion, vinculaciones_trabajo_grado_post)
						} else {
							rollbackDocumentoSolicitud(transaccion)
						}
					}
				} else {
					var resVinculacionTrabajoGrado map[string]interface{}
					url = "/v1/vinculacion_trabajo_grado"
					if err := SendRequestNew("PoluxCrudUrl", url, "POST", &resVinculacionTrabajoGrado, &v); err == nil {
						idVinculadoNuevo = int64(resVinculacionTrabajoGrado["Id"].(float64))
						vinculaciones_trabajo_grado_post = append(vinculaciones_trabajo_grado_post, resVinculacionTrabajoGrado)
					} else {
						logs.Error(err)
						if len(vinculaciones_trabajo_grado) > 0 || len(vinculaciones_trabajo_grado_post) > 0 {
							rollbackVinculacionTrabajoGradoRS(transaccion, vinculaciones_originales_trabajo_grado)
							rollbackVinculacionTrabajoGradoRSPost(transaccion, vinculaciones_trabajo_grado_post)
						} else {
							rollbackDocumentoSolicitud(transaccion)
						}
					}
				}
			} else {
				idVinculadoAntiguo = v.Id
				var resVinculacionTrabajoGrado map[string]interface{}
				url = "/v1/vinculacion_trabajo_grado/" + strconv.Itoa(v.Id)
				if err := SendRequestNew("PoluxCrudUrl", url, "PUT", &resVinculacionTrabajoGrado, &v); err == nil {
					vinculaciones_trabajo_grado_canceladas = append(vinculaciones_trabajo_grado_canceladas, resVinculacionTrabajoGrado)
				} else {
					if len(vinculaciones_trabajo_grado) > 0 || len(vinculaciones_trabajo_grado_post) > 0 {
						rollbackVinculacionTrabajoGradoRS(transaccion, vinculaciones_originales_trabajo_grado)
						rollbackVinculacionTrabajoGradoRSPost(transaccion, vinculaciones_trabajo_grado_post)
					} else {
						rollbackDocumentoSolicitud(transaccion)
					}
				}
			}
		}

		//Se busca si el vinculado antiguo tiene una revision pendiente
		var revisionTrabajoGrado []models.RevisionTrabajoGrado
		url := "parametro?query=CodigoAbreviacion:PENDIENTE_PLX,TipoParametroId__CodigoAbreviacion:ESTREV_TRG"
		var parametroEstadoRevision []models.Parametro
		if err := GetRequestNew("UrlCrudParametros", url, &parametroEstadoRevision); err != nil {
			logs.Error(err.Error())
			panic(err.Error())
		}
		url = beego.AppConfig.String("PoluxCrudUrl") + "/v1/revision_trabajo_grado?query=VinculacionTrabajoGrado:" + strconv.Itoa(idVinculadoAntiguo) +
			",EstadoRevisionTrabajoGrado:" + strconv.Itoa(parametroEstadoRevision[0].Id) + "&limit=1"
		fmt.Println("URL ", url)
		if err := GetJson(url, &revisionTrabajoGrado); err != nil {
			logs.Error(err.Error())
			panic(err.Error())
		}

		var vinc_orig = revisionTrabajoGrado[0]

		fmt.Println("revisitonTrabajoGrado", revisionTrabajoGrado)

		// Verificación adicional para asegurar que la revisión encontrada es válida
		if len(revisionTrabajoGrado) > 0 && revisionTrabajoGrado[0].Id != 0 {
			revisionTrabajoGrado[0].VinculacionTrabajoGrado.Id = int(idVinculadoNuevo)
			var resRevisionTrabajoGrado map[string]interface{}
			url = "/v1/revision_trabajo_grado/" + strconv.Itoa(revisionTrabajoGrado[0].Id)
			if err := SendRequestNew("PoluxCrudUrl", url, "PUT", &resRevisionTrabajoGrado, &revisionTrabajoGrado[0]); err != nil {
				rollbackRevisionTrabajoGrado(&vinc_orig, vinculaciones_originales_trabajo_grado, vinculaciones_trabajo_grado_post)
			}
		}

		// Si  el cambio es de director externo, se recibe la data del detalle de la pasantia y
		// se actualiza
		if transaccion.DetallesPasantia != nil {
			// Se busca el detalle de la pasantia asociado al tg
			var detallePasantia *models.DetallePasantia
			url = beego.AppConfig.String("PoluxCrudUrl") + "/v1/detalle_pasantia?query=TrabajoGrado__Id:" + strconv.Itoa(transaccion.DetallesPasantia.TrabajoGrado.Id) + "&limit=1"
			fmt.Println("URL ", url)
			if err := GetJson(url, &detallePasantia); err != nil {
				logs.Error(err.Error())
				panic(err.Error())
			}
			detallePasantia.Observaciones = strings.Split(detallePasantia.Observaciones, " y dirigida por ")[0]
			detallePasantia.Observaciones += transaccion.DetallesPasantia.Observaciones
			var resDetallePasantia map[string]interface{}
			url = "/v1/detalle_pasantia/" + strconv.Itoa(detallePasantia.Id)
			if err := SendRequestNew("PoluxCrudUrl", url, "PUT", &resDetallePasantia, &detallePasantia); err != nil {
				rollbackRevisionTrabajoGrado(&vinc_orig, vinculaciones_originales_trabajo_grado, vinculaciones_trabajo_grado_post)
			}
		}
	}

	//Solicitud de cambio de nombre del trabajo de grado
	if transaccion.TrabajoGrado != nil {
		var resTrabajoGrado map[string]interface{}
		url := "/v1/trabajo_grado/" + strconv.Itoa(transaccion.TrabajoGrado.Id)
		if err := SendRequestNew("PoluxCrudUrl", url, "PUT", &resTrabajoGrado, &transaccion.TrabajoGrado); err != nil {
			logs.Error(err)
			rollbackDocumentoSolicitud(transaccion)
		}
	}

	var parametroTipoSolicitud models.Parametro
	url = "parametro/" + strconv.Itoa(transaccion.TipoSolicitud.Id)
	if err := GetRequestNew("UrlCrudParametros", url, &parametroTipoSolicitud); err != nil {
		logs.Error(err.Error())
		panic(err.Error())
	}

	// Solicitud de prorroga
	if parametroTipoSolicitud.CodigoAbreviacion == "SPR_PLX" {
		if transaccion.CausaProrroga != nil {
			for _, data := range *transaccion.CausaProrroga {
				data.Activo = true
				data.FechaCreacion = time_bogota.TiempoBogotaFormato()
				data.FechaModificacion = time_bogota.TiempoBogotaFormato()
				url = "/v1/detalle_trabajo_grado"
				var resDetalleTrabajoGrado map[string]interface{}
				if err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDetalleTrabajoGrado, &data); err != nil {
					logs.Error(err)
					rollbackDocumentoSolicitud(transaccion)
				}
			}
		}
	}

	// Solicitud de cancelación de modalidad
	if transaccion.EstudianteTrabajoGrado != nil {

		var parametroEstadoEstudianteTrGr []models.Parametro
		url = "parametro?query=CodigoAbreviacion:EST_ACT_PLX,TipoParametroId__CodigoAbreviacion:EST_ESTU_TRG"
		if err := GetRequestNew("UrlCrudParametros", url, &parametroEstadoEstudianteTrGr); err != nil {
			logs.Error(err.Error())
			panic(err.Error())
		}

		//Se busca al estudiante con el trabajo de grado
		var estudianteTrabajoGrado []models.EstudianteTrabajoGrado
		url = beego.AppConfig.String("PoluxCrudUrl") + "/v1/estudiante_trabajo_grado?query=TrabajoGrado__Id:" + strconv.Itoa(transaccion.EstudianteTrabajoGrado.TrabajoGrado.Id) +
			",Estudiante:" + transaccion.EstudianteTrabajoGrado.Estudiante + ",EstadoEstudianteTrabajoGrado:" + strconv.Itoa(parametroEstadoEstudianteTrGr[0].Id) + "&limit=1"
		fmt.Println("URL ", url)
		if err := GetJson(url, &estudianteTrabajoGrado); err != nil {
			logs.Error(err.Error())
			panic(err.Error())
		}
		transaccion.EstudianteTrabajoGrado.Id = estudianteTrabajoGrado[0].Id
		var resEstudianteTrabajoGrado map[string]interface{}
		url := "/v1/estudiante_trabajo_grado/" + strconv.Itoa(estudianteTrabajoGrado[0].Id)
		if err := SendRequestNew("PoluxCrudUrl", url, "PUT", &resEstudianteTrabajoGrado, &transaccion.EstudianteTrabajoGrado); err == nil {
			var estudianteTrabajoGradoAux []models.EstudianteTrabajoGrado
			url = beego.AppConfig.String("PoluxCrudUrl") + "/v1/estudiante_trabajo_grado?query=TrabajoGrado__Id:" + strconv.Itoa(transaccion.EstudianteTrabajoGrado.TrabajoGrado.Id) +
				",EstadoEstudianteTrabajoGrado:" + strconv.Itoa(parametroEstadoEstudianteTrGr[0].Id)
			fmt.Println("URL ", url)
			if err := GetJson(url, &estudianteTrabajoGradoAux); err != nil {
				logs.Error(err.Error())
				panic(err.Error())
			}
			// si no hay estudiantes vinculados se inactivan las vinculaciones y se cancela el tg
			if estudianteTrabajoGradoAux[0].Id == 0 {
				// Se inactivan las vinculaciones
				var vinculaciones_tr_gr = make([]map[string]interface{}, 0)
				for _, v := range *transaccion.VinculacionesCancelacion {
					v.Activo = false
					var resVinculacionTrabajoGrado map[string]interface{}
					url := "/v1/vinculacion_trabajo_grado/" + strconv.Itoa(v.Id)
					if err := SendRequestNew("PoluxCrudUrl", url, "PUT", &resVinculacionTrabajoGrado, v); err == nil {
						vinculaciones_tr_gr = append(vinculaciones_tr_gr, resVinculacionTrabajoGrado)
					} else {
						if len(vinculaciones_tr_gr) > 0 {
							logs.Error(err)
							rollbackVincTrGrCanc(transaccion)
						} else {
							rollbackEstTrGrCanc(transaccion)
						}
					}
				}
				var parametroEstadoTrabajoGrado []models.Parametro
				url = "parametro?query=CodigoAbreviacion:CNC_PLX,TipoParametroId__CodigoAbreviacion:EST_TRG"
				if err := GetRequestNew("UrlCrudParametros", url, &parametroEstadoTrabajoGrado); err != nil {
					logs.Error(err.Error())
					panic(err.Error())
				}
				// Se cancela el trabajo de grado
				tg := transaccion.EstudianteTrabajoGrado.TrabajoGrado
				tg.EstadoTrabajoGrado = parametroEstadoTrabajoGrado[0].Id
				var resTrabajoGrado map[string]interface{}
				url := "/v1/trabajo_grado/" + strconv.Itoa(tg.Id)
				if err := SendRequestNew("PoluxCrudUrl", url, "PUT", &resTrabajoGrado, tg); err != nil {
					logs.Error(err)
					rollbackVincTrGrCanc(transaccion)
				}

				// Actualizar asignaturas trabajo de grado a cancelado
				var asignaturasTrabajoGrado []models.AsignaturaTrabajoGrado
				// Se busca asignaturas trabajo grado
				url = beego.AppConfig.String("PoluxCrudUrl") + "/v1/asignatura_trabajo_grado?query=TrabajoGrado__Id:" + strconv.Itoa(transaccion.EstudianteTrabajoGrado.TrabajoGrado.Id)
				fmt.Println("URL ", url)
				if err := GetJson(url, &asignaturasTrabajoGrado); err != nil {
					logs.Error(err.Error())
					panic(err.Error())
				}
				var parametroEstAsTrGr []models.Parametro
				url = "parametro?query=CodigoAbreviacion:CNC_PLX,TipoParametroId__CodigoAbreviacion:EST_ASIG_TRG"
				if err := GetRequestNew("UrlCrudParametros", url, &parametroEstAsTrGr); err != nil {
					logs.Error(err.Error())
					panic(err.Error())
				}
				var asignaturas_tr_gr = make([]map[string]interface{}, 0)
				for _, v := range asignaturasTrabajoGrado {
					//Id de la asignatura cancelada
					v.EstadoAsignaturaTrabajoGrado = parametroEstAsTrGr[0].Id
					var resAsignaturaTrabajoGrado map[string]interface{}
					url := "/v1/asignatura_trabajo_grado/" + strconv.Itoa(v.Id)
					if err := SendRequestNew("PoluxCrudUrl", url, "PUT", &resAsignaturaTrabajoGrado, v); err == nil {
						asignaturas_tr_gr = append(asignaturas_tr_gr, resAsignaturaTrabajoGrado)
					} else {
						logs.Error(err)
						if len(asignaturas_tr_gr) > 0 {
							rollbackAsTrGr(transaccion, &asignaturasTrabajoGrado)
						} else {
							rollbackTrGrCanc(transaccion)
						}
					}
				}

				var parametroEspAcadIns []models.Parametro
				url = "parametro?query=CodigoAbreviacion:ESP_ACT_PLX,TipoParametroId__CodigoAbreviacion:EST_ESP"
				if err := GetRequestNew("UrlCrudParametros", url, &parametroEspAcadIns); err != nil {
					logs.Error(err.Error())
					panic(err.Error())
				}
				//Actualizar espacios academicos inscritos
				var espaciosAcademicosInscritos []models.EspacioAcademicoInscrito
				// Se buscan espacios academicos inscritos activos
				url = beego.AppConfig.String("PoluxCrudUrl") + "/v1/espacio_academico_inscrito?query=TrabajoGrado__Id:" + strconv.Itoa(transaccion.EstudianteTrabajoGrado.TrabajoGrado.Id) +
					",EstadoEspacioAcademicoInscrito:" + strconv.Itoa(parametroEspAcadIns[0].Id)
				fmt.Println("URL ", url)
				if err := GetJson(url, &espaciosAcademicosInscritos); err != nil {
					logs.Error(err.Error())
					panic(err.Error())
				}
				if espaciosAcademicosInscritos[0].Id != 0 {
					var parametroEspAcadInsAux []models.Parametro
					url = "parametro?query=CodigoAbreviacion:ESP_CAN_PLX,TipoParametroId__CodigoAbreviacion:EST_ESP"
					if err := GetRequestNew("UrlCrudParametros", url, &parametroEspAcadInsAux); err != nil {
						logs.Error(err.Error())
						panic(err.Error())
					}
					var espacios_acad_insc = make([]map[string]interface{}, 0)
					for _, v := range espaciosAcademicosInscritos {
						// Id del espacio cancelado
						v.EstadoEspacioAcademicoInscrito = parametroEspAcadInsAux[0].Id
						var resEspacioAcademicoInscrito map[string]interface{}
						url := "/v1/espacio_academico_inscrito/" + strconv.Itoa(v.Id)
						if err := SendRequestNew("PoluxCrudUrl", url, "PUT", &resEspacioAcademicoInscrito, v); err == nil {
							espacios_acad_insc = append(espacios_acad_insc, resEspacioAcademicoInscrito)
						} else {
							logs.Error(err)
							if len(espacios_acad_insc) > 0 {
								rollbackEsAcadInsc(transaccion, &asignaturasTrabajoGrado, &espaciosAcademicosInscritos)
							} else {
								rollbackAsTrGr(transaccion, &asignaturasTrabajoGrado)
							}
						}
					}
				}
			}
		} else {
			logs.Error(err)
			rollbackDocumentoSolicitud(transaccion)
		}
	}

	// Solicitud de revisión del trabajo de grado
	if transaccion.TrRevision != nil {
		// Se actualiza el trabajo de grado
		var resTrabajoGrado map[string]interface{}
		url := "/v1/trabajo_grado/" + strconv.Itoa(transaccion.TrRevision.TrabajoGrado.Id)
		if err := SendRequestNew("PoluxCrudUrl", url, "PUT", &resTrabajoGrado, &transaccion.TrRevision.TrabajoGrado); err != nil {
			logs.Error(err)
			rollbackDocumentoSolicitud(transaccion)
		}

		//INSERTA EN LA TABLA DETALLE TRABAJO GRADO
		if transaccion.TrRevision.DetalleTrabajoGrado != nil {
			var detalles_trabajo_grado = make([]map[string]interface{}, 0)
			for _, data := range *transaccion.TrRevision.DetalleTrabajoGrado {
				data.Activo = true
				data.FechaCreacion = time_bogota.TiempoBogotaFormato()
				data.FechaModificacion = time_bogota.TiempoBogotaFormato()
				var resDetalleTrabajoGrado map[string]interface{}
				url = "/v1/detalle_trabajo_grado"
				if err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDetalleTrabajoGrado, &data); err == nil {
					data.Id = int(resDetalleTrabajoGrado["Id"].(float64))
					detalles_trabajo_grado = append(detalles_trabajo_grado, resDetalleTrabajoGrado)
				} else {
					if len(detalles_trabajo_grado) > 0 {
						rollbackDetTrGrRev(transaccion)
					} else {
						rollbackTrGrRev(transaccion)
					}
				}
			}
		}

		// Se inserta el documento final de la revisión y se relaciona con el trabajo de grado
		var resDocumentoEscrito map[string]interface{}
		url = "/v1/documento_escrito"
		if err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoEscrito, &transaccion.TrRevision.DocumentoEscrito); err == nil {
			transaccion.TrRevision.DocumentoEscrito.Id = int(resDocumentoEscrito["Id"].(float64))
			transaccion.TrRevision.DocumentoTrabajoGrado.DocumentoEscrito.Id = int(resDocumentoEscrito["Id"].(float64))
			var resDocumentoTrabajoGrado map[string]interface{}
			url = "/v1/documento_trabajo_grado"
			if err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoTrabajoGrado, &transaccion.TrRevision.DocumentoTrabajoGrado); err == nil {
				transaccion.TrRevision.DocumentoEscrito.Id = int(resDocumentoTrabajoGrado["Id"].(float64))
			} else {
				rollbackDocEscrRev(transaccion)
			}
		}

		// Se actualizan las vinculaciones
		var vinculaciones_trabajo_grado = make([]map[string]interface{}, 0)
		var vinculaciones_originales_trabajo_grado []models.VinculacionTrabajoGrado
		var vinculaciones_trabajo_grado_post = make([]map[string]interface{}, 0)
		for _, v := range *transaccion.TrRevision.Vinculaciones {
			// Si esta activo es nuevo y se inserta sino se actualiza la fecha de fin y el activo
			if v.Activo {
				// Se buscar si el docente ya estuvo vinculado y se actualiza
				var vinculado []models.VinculacionTrabajoGrado
				url = beego.AppConfig.String("PoluxCrudUrl") + "/v1/vinculacion_trabajo_grado?query=TrabajoGrado__Id:" + strconv.Itoa(v.TrabajoGrado.Id) +
					",Usuario:" + strconv.Itoa(v.Usuario) + ",RolTrabajoGrado:" + strconv.Itoa(v.RolTrabajoGrado) + "&limit=1"
				fmt.Println("URL ", url)
				if err := GetJson(url, &vinculado); err != nil {
					logs.Error(err.Error())
					panic(err.Error())
				}
				if vinculado[0].Id != 0 {
					var vinculadoAux = vinculado[0]
					vinculado[0].Activo = v.Activo
					vinculado[0].FechaFin = v.FechaFin
					vinculado[0].FechaInicio = v.FechaInicio
					var resVinculacionTrabajoGrado map[string]interface{}
					url = "/v1/vinculacion_trabajo_grado/" + strconv.Itoa(vinculado[0].Id)
					if err := SendRequestNew("PoluxCrudUrl", url, "PUT", &resVinculacionTrabajoGrado, &vinculado[0]); err == nil {
						vinculaciones_originales_trabajo_grado = append(vinculaciones_originales_trabajo_grado, vinculadoAux)
						vinculaciones_trabajo_grado = append(vinculaciones_trabajo_grado, resVinculacionTrabajoGrado)
					} else {
						logs.Error(err)
						if len(vinculaciones_trabajo_grado) > 0 || len(vinculaciones_trabajo_grado_post) > 0 {
							rollbackVincTrGrRev(transaccion, vinculaciones_originales_trabajo_grado)
							rollbackVincTrGrPostRev(transaccion, vinculaciones_trabajo_grado_post)
							rollbackDocTrGrRev(transaccion)
						} else {
							rollbackDocTrGrRev(transaccion)
						}
					}
				} else {
					var resVinculacionTrabajoGrado map[string]interface{}
					url = "/v1/vinculacion_trabajo_grado"
					if err := SendRequestNew("PoluxCrudUrl", url, "POST", &resVinculacionTrabajoGrado, &v); err == nil {
						vinculaciones_trabajo_grado_post = append(vinculaciones_trabajo_grado_post, resVinculacionTrabajoGrado)
					} else {
						logs.Error(err)
						if len(vinculaciones_trabajo_grado) > 0 || len(vinculaciones_trabajo_grado_post) > 0 {
							rollbackVincTrGrRev(transaccion, vinculaciones_originales_trabajo_grado)
							rollbackVincTrGrPostRev(transaccion, vinculaciones_trabajo_grado_post)
							rollbackDocTrGrRev(transaccion)
						} else {
							rollbackDocTrGrRev(transaccion)
						}
					}
				}
			}
		}
	}
	return alerta, outputError
}

func rollbackResAnterior(transaccion *models.TrRespuestaSolicitud) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK RES ANTERIOR ")
	var respuesta map[string]interface{}
	transaccion.RespuestaAnterior.Activo = true
	url := "/v1/respuesta_solicitud/" + strconv.Itoa(transaccion.RespuestaAnterior.Id)
	if err := SendRequestNew("PoluxCrudUrl", url, "PUT", &respuesta, &transaccion.RespuestaAnterior); err != nil {
		panic("Rollback respuesta anteror " + err.Error())
	}
	return nil
}

func rollbackResNueva(transaccion *models.TrRespuestaSolicitud) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK RES NUEVA")
	var respuesta map[string]interface{}
	url := "/v1/respuesta_solicitud/" + strconv.Itoa(transaccion.RespuestaNueva.Id)
	if err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil {
		panic("Rollback respuesta nueva" + err.Error())
	} else {
		rollbackResAnterior(transaccion)
	}
	return nil
}

func rollbackDocumentoSolicitud(transaccion *models.TrRespuestaSolicitud) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK DOCUMENTO SOLICITUD")
	var respuesta map[string]interface{}
	url := "/v1/documento_solicitud/" + strconv.Itoa(transaccion.DocumentoSolicitud.Id)
	if err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil {
		panic("Rollback documento solicitud" + err.Error())
	} else {
		rollbackResNueva(transaccion)
	}
	return nil
}

func rollbackTrabajoGrado(transaccion *models.TrRespuestaSolicitud) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK TRABAJO GRADO")
	var respuesta map[string]interface{}
	url := "/v1/trabajo_grado/" + strconv.Itoa(transaccion.TrTrabajoGrado.TrabajoGrado.Id)
	if err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil {
		panic("Rollback trabajo grado" + err.Error())
	} else {
		rollbackDocumentoSolicitud(transaccion)
	}
	return nil
}

func rollbackSolicitudTrabajoGrado(transaccion *models.TrRespuestaSolicitud) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK SOLICITUD TRABAJO GRADO")
	var respuesta map[string]interface{}
	transaccion.SolicitudTrabajoGrado.TrabajoGrado = nil
	url := "/v1/solicitud_trabajo_grado/" + strconv.Itoa(transaccion.SolicitudTrabajoGrado.Id)
	if err := SendRequestNew("PoluxCrudUrl", url, "PUT", &respuesta, &transaccion.SolicitudTrabajoGrado); err != nil {
		panic("Rollback solicitud trabajo grado" + err.Error())
	} else {
		rollbackTrabajoGrado(transaccion)
	}
	return nil
}

func rollbackAsignaturasTrabajoGrado(transaccion *models.TrRespuestaSolicitud) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK ASIGNATURAS TRABAJO GRADO")
	var respuesta map[string]interface{}
	for _, v := range *transaccion.TrTrabajoGrado.AsignaturasTrabajoGrado {
		if v.Id != 0 {
			url := "/v1/asignatura_trabajo_grado/" + strconv.Itoa(v.Id)
			if err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil {
				panic("Rollback asignaturas trabajo grado" + err.Error())
			}
		}
	}
	rollbackSolicitudTrabajoGrado(transaccion)
	return nil
}

func rollbackEstudianteTrabajoGrado(transaccion *models.TrRespuestaSolicitud) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK ESTUDIANTE TRABAJO GRADO")
	var respuesta map[string]interface{}
	for _, v := range *transaccion.TrTrabajoGrado.EstudianteTrabajoGrado {
		if v.Id != 0 {
			url := "/v1/estudiante_trabajo_grado/" + strconv.Itoa(v.Id)
			if err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil {
				panic("Rollback estudiante trabajo grado" + err.Error())
			}
		}
	}
	rollbackAsignaturasTrabajoGrado(transaccion)
	return nil
}

func rollbackAreasTrabajoGrado(transaccion *models.TrRespuestaSolicitud) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK AREAS TRABAJO GRADO")
	var respuesta map[string]interface{}
	if transaccion.TrTrabajoGrado.AreasTrabajoGrado != nil {
		for _, v := range *transaccion.TrTrabajoGrado.AreasTrabajoGrado {
			if v.Id != 0 {
				url := "/v1/areas_trabajo_grado/" + strconv.Itoa(v.Id)
				if err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil {
					panic("Rollback areas trabajo grado" + err.Error())
				}
			}
		}
	}
	rollbackEstudianteTrabajoGrado(transaccion)
	return nil
}

func rollbackVinculacionTrabajoGrado(transaccion *models.TrRespuestaSolicitud) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK VINCULACION TRABAJO GRADO")
	var respuesta map[string]interface{}
	if transaccion.TrTrabajoGrado.VinculacionTrabajoGrado != nil {
		for _, v := range *transaccion.TrTrabajoGrado.VinculacionTrabajoGrado {
			if v.Id != 0 {
				url := "/v1/vinculacion_trabajo_grado/" + strconv.Itoa(v.Id)
				if err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil {
					panic("Rollback vinculacion trabajo grado" + err.Error())
				}
			}
		}
	}
	rollbackAreasTrabajoGrado(transaccion)
	return nil
}

func rollbackDocumentoEscrito(transaccion *models.TrRespuestaSolicitud) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK DOCUMENTO ESCRITO")
	var respuesta map[string]interface{}
	if transaccion.TrTrabajoGrado.DocumentoEscrito != nil {
		url := "/v1/documento_escrito/" + strconv.Itoa(transaccion.TrTrabajoGrado.DocumentoEscrito.Id)
		if err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil {
			panic("Rollback documento escrito " + err.Error())
		}
	}
	rollbackVinculacionTrabajoGrado(transaccion)
	return nil
}

func rollbackDocumentoTrabajoGrado(transaccion *models.TrRespuestaSolicitud) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK DOCUMENTO TRABAJO GRADO")
	var respuesta map[string]interface{}
	if transaccion.TrTrabajoGrado.DocumentoTrabajoGrado != nil {
		url := "/v1/documento_trabajo_grado/" + strconv.Itoa(transaccion.TrTrabajoGrado.DocumentoEscrito.Id)
		if err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil {
			panic("Rollback documento trabajo grado " + err.Error())
		}
	}
	rollbackDocumentoEscrito(transaccion)
	return nil
}

func rollbackDocumentoEscritoPasantia(transaccion *models.TrRespuestaSolicitud) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK DOCUMENTO ESCRITO PASANTIA")
	var respuesta map[string]interface{}
	if transaccion.DetallesPasantia.Contrato.Id != 0 {
		url := "/v1/documento_escrito/" + strconv.Itoa(transaccion.DetallesPasantia.Contrato.Id)
		if err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil {
			panic("Rollback documento escrito " + err.Error())
		}
	}
	if transaccion.DetallesPasantia.Carta.Id != 0 {
		url := "/v1/documento_escrito/" + strconv.Itoa(transaccion.DetallesPasantia.Carta.Id)
		if err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil {
			panic("Rollback documento escrito " + err.Error())
		}
	}
	rollbackDocumentoEscrito(transaccion)
	return nil
}

func rollbackDocumentoTrabajoGradoPasantia(transaccion *models.TrRespuestaSolicitud) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK DOCUMENTO TRABAJO GRADO PASANTIA")
	var respuesta map[string]interface{}
	if transaccion.DetallesPasantia.DTG_Contrato.Id != 0 {
		url := "/v1/documento_trabajo_grado/" + strconv.Itoa(transaccion.DetallesPasantia.DTG_Contrato.Id)
		if err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil {
			panic("Rollback documento trabajo grado " + err.Error())
		}
	}
	if transaccion.DetallesPasantia.DTG_Carta.Id != 0 {
		url := "/v1/documento_trabajo_grado/" + strconv.Itoa(transaccion.DetallesPasantia.DTG_Carta.Id)
		if err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil {
			panic("Rollback documento trabajo grado " + err.Error())
		}
	}
	rollbackDocumentoEscritoPasantia(transaccion)
	return nil
}

func rollbackEspaciosAcademicosInscritos(transaccion *models.TrRespuestaSolicitud) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK ESPACIOS ACADEMICOS INSCRITOS")
	var respuesta map[string]interface{}
	if transaccion.EspaciosAcademicosInscritos != nil {
		for _, v := range *transaccion.EspaciosAcademicosInscritos {
			if v.Id != 0 {
				url := "/v1/espacio_academico_inscrito/" + strconv.Itoa(v.Id)
				if err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil {
					panic("Rollback espacios academicos inscritos " + err.Error())
				}
			}
		}
	}
	rollbackDocumentoTrabajoGrado(transaccion)
	return nil
}

func rollbackDetallesPasantia(transaccion *models.TrRespuestaSolicitud) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK DETALLES PASANTIA")
	var respuesta map[string]interface{}
	if transaccion.DetallesPasantia.Id != 0 {
		url := "/v1/detalle_pasantia/" + strconv.Itoa(transaccion.DetallesPasantia.Id)
		if err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil {
			panic("Rollback detalles pasantia " + err.Error())
		}
	}
	rollbackEspaciosAcademicosInscritos(transaccion)
	return nil
}

func rollbackDetallesPasantiaExterna(transaccion *models.TrRespuestaSolicitud) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK DETALLES PASANTIA EXTERNA")
	var respuesta map[string]interface{}
	if transaccion.DetallesPasantiaExterna != nil {
		for _, v := range *transaccion.DetallesPasantiaExterna {
			if v.Id != 0 {
				url := "/v1/detalle_trabajo_grado/" + strconv.Itoa(v.Id)
				if err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil {
					panic("Rollback detalles pasantia externa " + err.Error())
				}
			}
		}
	}
	rollbackDetallesPasantia(transaccion)
	return nil
}

func rollbackVinculacionTrabajoGradoRS(transaccion *models.TrRespuestaSolicitud, vinculacionesOriginales []models.VinculacionTrabajoGrado) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK VINCULACION TRABAJO GRADO RS")
	var respuesta map[string]interface{}
	for _, v := range vinculacionesOriginales {
		url := "/v1/vinculacion_trabajo_grado/" + strconv.Itoa(v.Id)
		if err := SendRequestNew("PoluxCrudUrl", url, "PUT", &respuesta, &v); err != nil {
			panic("Rollback vinculacion trabajo grado rs" + err.Error())
		}
	}
	rollbackDocumentoSolicitud(transaccion)
	return nil
}

func rollbackVinculacionTrabajoGradoRSPost(transaccion *models.TrRespuestaSolicitud, vinculacionesNuevas []map[string]interface{}) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK VINCULACION TRABAJO GRADO RS POST")
	var respuesta map[string]interface{}
	for _, v := range vinculacionesNuevas {
		var vinculacionNueva models.VinculacionTrabajoGrado
		jsonData, err := json.Marshal(v)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(jsonData, &vinculacionNueva)
		if err != nil {
			log.Fatal(err)
		}
		url := "/v1/vinculacion_trabajo_grado/" + strconv.Itoa(vinculacionNueva.Id)
		if err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, &vinculacionNueva); err != nil {
			panic("Rollback vinculacion trabajo grado rs post" + err.Error())
		}
	}
	return nil
}

func rollbackRevisionTrabajoGrado(revisionAnterior *models.RevisionTrabajoGrado, vinc_orig []models.VinculacionTrabajoGrado, vinc_post []map[string]interface{}) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK REVISON TRABAJO GRADO")
	var respuesta map[string]interface{}
	url := "/v1/revision_trabajo_grado/" + strconv.Itoa(revisionAnterior.Id)
	if err := SendRequestNew("PoluxCrudUrl", url, "PUT", &respuesta, &revisionAnterior); err != nil {
		panic("Rollback revision trabajo grado" + err.Error())
	}
	rollbackVinculacionTrabajoGradoRS(nil, vinc_orig)
	rollbackVinculacionTrabajoGradoRSPost(nil, vinc_post)
	return nil
}

func rollbackEstTrGrCanc(transaccion *models.TrRespuestaSolicitud) (outputError map[string]interface{}) {
	var parametroEstadoEstudianteTrGr []models.Parametro
	url := "parametro?query=CodigoAbreviacion:EST_ACT_PLX,TipoParametroId__CodigoAbreviacion:EST_ESTU_TRG"
	if err := GetRequestNew("UrlCrudParametros", url, &parametroEstadoEstudianteTrGr); err != nil {
		logs.Error(err.Error())
		panic(err.Error())
	}
	fmt.Println("ROLLBACK ESTUDIANTE TRABAJO GRADO CANCELACIÓN")
	var respuesta map[string]interface{}
	url = "/v1/estudiante_trabajo_grado/" + strconv.Itoa(transaccion.EstudianteTrabajoGrado.Id)
	if err := SendRequestNew("PoluxCrudUrl", url, "PUT", &respuesta, &transaccion.EstudianteTrabajoGrado); err != nil {
		panic("Rollback estudiante trabajo grado cancelación " + err.Error())
	} else {
		rollbackDocumentoSolicitud(transaccion)
	}
	return nil
}

func rollbackVincTrGrCanc(transaccion *models.TrRespuestaSolicitud) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK VINCULACIÓN TRABAJO GRADO CANCELACIÓN")
	for _, v := range *transaccion.VinculacionesCancelacion {
		v.Activo = true
		var respuesta map[string]interface{}
		url := "/v1/vinculacion_trabajo_grado/" + strconv.Itoa(v.Id)
		if err := SendRequestNew("PoluxCrudUrl", url, "PUT", &respuesta, v); err != nil {
			panic("Rollback vinculación trabajo grado cancelación" + err.Error())
		}
	}
	rollbackEstTrGrCanc(transaccion)
	return nil
}

func rollbackTrGrCanc(transaccion *models.TrRespuestaSolicitud) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK TRABAJO GRADO CANCELACIÓN")
	var respuesta map[string]interface{}
	url := "/v1/trabajo_grado/" + strconv.Itoa(transaccion.EstudianteTrabajoGrado.TrabajoGrado.Id)
	if err := SendRequestNew("PoluxCrudUrl", url, "PUT", &respuesta, transaccion.EstudianteTrabajoGrado.TrabajoGrado); err != nil {
		panic("Rollback vinculación trabajo grado cancelación" + err.Error())
	} else {
		rollbackVincTrGrCanc(transaccion)
	}
	return nil
}

func rollbackAsTrGr(transaccion *models.TrRespuestaSolicitud, asignaturas *[]models.AsignaturaTrabajoGrado) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK ASIGNATURAS TRABAJO GRADO")
	var parametroEstAsTrGr []models.Parametro
	url := "parametro?query=CodigoAbreviacion:CND_PLX,TipoParametroId__CodigoAbreviacion:EST_ASIG_TRG"
	if err := GetRequestNew("UrlCrudParametros", url, &parametroEstAsTrGr); err != nil {
		logs.Error(err.Error())
		panic(err.Error())
	}
	for _, v := range *asignaturas {
		var respuesta map[string]interface{}
		v.EstadoAsignaturaTrabajoGrado = parametroEstAsTrGr[0].Id
		url := "/v1/asignatura_trabajo_grado/" + strconv.Itoa(v.Id)
		if err := SendRequestNew("PoluxCrudUrl", url, "PUT", &respuesta, v); err != nil {
			logs.Error(err.Error())
			panic(err.Error())
		}
	}
	rollbackTrGrCanc(transaccion)
	return nil
}

func rollbackEsAcadInsc(transaccion *models.TrRespuestaSolicitud, asignaturas *[]models.AsignaturaTrabajoGrado, espaciosAcad *[]models.EspacioAcademicoInscrito) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK ESPACIO ACADEMICO INSCRITO")
	var parametroEspAcadIns []models.Parametro
	url := "parametro?query=CodigoAbreviacion:ESP_ACT_PLX,TipoParametroId__CodigoAbreviacion:EST_ESP"
	if err := GetRequestNew("UrlCrudParametros", url, &parametroEspAcadIns); err != nil {
		logs.Error(err.Error())
		panic(err.Error())
	}
	for _, v := range *espaciosAcad {
		var respuesta map[string]interface{}
		v.EstadoEspacioAcademicoInscrito = parametroEspAcadIns[0].Id
		url := "/v1/espacio_academico_inscrito/" + strconv.Itoa(v.Id)
		if err := SendRequestNew("PoluxCrudUrl", url, "PUT", &respuesta, v); err != nil {
			logs.Error(err.Error())
			panic(err.Error())
		}
	}
	rollbackAsTrGr(transaccion, asignaturas)
	return nil
}

func rollbackTrGrRev(transaccion *models.TrRespuestaSolicitud) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK TRABAJO GRADO REVISION")
	var parametroEstTrGr []models.Parametro
	url := "parametro?query=CodigoAbreviacion:RDE_PLX,TipoParametroId__CodigoAbreviacion:EST_TRG"
	if err := GetRequestNew("UrlCrudParametros", url, &parametroEstTrGr); err != nil {
		logs.Error(err.Error())
		panic(err.Error())
	}
	var respuesta map[string]interface{}
	url = "/v1/trabajo_grado/" + strconv.Itoa(transaccion.TrRevision.TrabajoGrado.Id)
	transaccion.TrRevision.TrabajoGrado.EstadoTrabajoGrado = parametroEstTrGr[0].Id
	if err := SendRequestNew("PoluxCrudUrl", url, "PUT", &respuesta, transaccion.TrRevision.TrabajoGrado); err != nil {
		panic("Rollback trabajo grado revision" + err.Error())
	} else {
		rollbackDocumentoSolicitud(transaccion)
	}
	return nil
}

func rollbackDetTrGrRev(transaccion *models.TrRespuestaSolicitud) (outputError map[string]interface{}) {
	if transaccion.TrRevision.DetalleTrabajoGrado != nil {
		fmt.Println("ROLLBACK DETALLE TRABAJO GRADO REVISION")
		for _, data := range *transaccion.TrRevision.DetalleTrabajoGrado {
			var respuesta map[string]interface{}
			url := "/v1/detalle_trabajo_grado/" + strconv.Itoa(data.Id)
			if err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil {
				panic("Rollback detalle trabajo grado revision" + err.Error())
			}
		}
	}
	rollbackTrGrRev(transaccion)
	return nil
}

func rollbackDocEscrRev(transaccion *models.TrRespuestaSolicitud) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK DOCUMENTO ESCRITO REVISION")
	var respuesta map[string]interface{}
	url := "/v1/documento_escrito/" + strconv.Itoa(transaccion.TrRevision.DocumentoEscrito.Id)
	if err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil {
		panic("Rollback documento escrito revision" + err.Error())
	} else {
		rollbackDetTrGrRev(transaccion)
	}
	return nil
}

func rollbackDocTrGrRev(transaccion *models.TrRespuestaSolicitud) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK DOCUMENTO TRABAJO GRADO REVISION")
	var respuesta map[string]interface{}
	url := "/v1/documento_trabajo_grado/" + strconv.Itoa(transaccion.TrRevision.DocumentoTrabajoGrado.Id)
	if err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil {
		panic("Rollback documento trabajo grado revision" + err.Error())
	} else {
		rollbackDocEscrRev(transaccion)
	}
	return nil
}

func rollbackVincTrGrRev(transaccion *models.TrRespuestaSolicitud, vinculacionesOriginales []models.VinculacionTrabajoGrado) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK VINCULACION TRABAJO GRADO REVISION")
	var respuesta map[string]interface{}
	for _, v := range vinculacionesOriginales {
		url := "/v1/vinculacion_trabajo_grado/" + strconv.Itoa(v.Id)
		if err := SendRequestNew("PoluxCrudUrl", url, "PUT", &respuesta, &v); err != nil {
			panic("Rollback vinculacion trabajo grado revision" + err.Error())
		}
	}
	return nil
}

func rollbackVincTrGrPostRev(transaccion *models.TrRespuestaSolicitud, vinculacionesNuevas []map[string]interface{}) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK VINCULACION TRABAJO GRADO POST REVISION")
	var respuesta map[string]interface{}
	for _, v := range vinculacionesNuevas {
		var vinculacionNueva models.VinculacionTrabajoGrado
		jsonData, err := json.Marshal(v)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(jsonData, &vinculacionNueva)
		if err != nil {
			log.Fatal(err)
		}
		url := "/v1/vinculacion_trabajo_grado/" + strconv.Itoa(vinculacionNueva.Id)
		if err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, &vinculacionNueva); err != nil {
			panic("Rollback vinculacion trabajo grado post revision" + err.Error())
		}
	}
	return nil
}
