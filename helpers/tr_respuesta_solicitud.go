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
	request "github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/time_bogota"
)

func AddTransaccionRespuestaSolicitud(transaccion *models.TrRespuestaSolicitud) (response map[string]interface{}, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "AddTransaccionRespuestaSolicitud", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
	//alerta = append(alerta, "Success")

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
	if status, err := SendRequestNew("PoluxCrudUrl", url, "PUT", &resRespuestaAnterior, &transaccion.RespuestaAnterior); err == nil && status == "200" {
		fmt.Println("RESPUESTA ANTERIOR", resRespuestaAnterior)
		url = "/v1/respuesta_solicitud"
		var resRespuestaNueva map[string]interface{}
		fmt.Println("TR _RESPUESTA NUEVA", transaccion.RespuestaNueva)
		if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resRespuestaNueva, &transaccion.RespuestaNueva); err == nil && status == "201" {
			fmt.Println("TR _RESPUESTA NUEVA", transaccion.RespuestaNueva)
			fmt.Println("RESPUESTA NUEVA", resRespuestaNueva)
			transaccion.RespuestaNueva.Id = int(resRespuestaNueva["Id"].(float64))
			url = "/v1/documento_solicitud"
			var resDocumentoSolicitud map[string]interface{}
			if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoSolicitud, &transaccion.DocumentoSolicitud); err == nil && status == "201" {
				transaccion.DocumentoSolicitud.Id = int(resDocumentoSolicitud["Id"].(float64))
			} else {
				fmt.Println("ERR3", err)
				logs.Error(err)
				rollbackResNueva(transaccion)
				logs.Error(err)
				panic(err.Error())
			}
		} else {
			logs.Error(err)
			fmt.Println("ERR", err)
			rollbackResAnterior(transaccion)
			logs.Error(err)
			panic(err.Error())
		}
	} else {
		fmt.Println("ERR2", err)
		logs.Error(err)
		panic(err.Error())
	}

	if transaccion.TrTrabajoGrado != nil && ((parametro.CodigoAbreviacion != "EAPOS_PLX" && parametroEstadoSolicitud.CodigoAbreviacion == "ACPR_PLX") || transaccion.MateriasProPos) {
		url = "/v1/trabajo_grado"
		var resTrabajoGrado map[string]interface{}
		if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resTrabajoGrado, &transaccion.TrTrabajoGrado.TrabajoGrado); err == nil && status == "201" {
			var idTrabajoGrado = int(resTrabajoGrado["Id"].(float64))
			transaccion.TrTrabajoGrado.TrabajoGrado.Id = idTrabajoGrado
			transaccion.SolicitudTrabajoGrado.TrabajoGrado.Id = idTrabajoGrado
			var resSolicitudTrabajoGrado map[string]interface{}
			url = "/v1/solicitud_trabajo_grado/" + strconv.Itoa(transaccion.SolicitudTrabajoGrado.Id)
			if status, err := SendRequestNew("PoluxCrudUrl", url, "PUT", &resSolicitudTrabajoGrado, &transaccion.SolicitudTrabajoGrado); err == nil && status == "200" {
				url = "/v1/asignatura_trabajo_grado"
				var materias = make([]map[string]interface{}, 0)
				for i, v := range *transaccion.TrTrabajoGrado.AsignaturasTrabajoGrado {
					var resAsignaturaTrabajoGrado map[string]interface{}
					v.TrabajoGrado.Id = idTrabajoGrado
					if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resAsignaturaTrabajoGrado, &v); err == nil && status == "201" {
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
					if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resEstudianteTrabajoGrado, &v); err == nil && status == "201" {
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
						if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resAreasTrabajoGrado, &v); err == nil && status == "201" {
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
						if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resVinculacionTrabajoGrado, &v); err == nil && status == "201" {
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
					if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoEscrito, &transaccion.TrTrabajoGrado.DocumentoEscrito); err == nil && status == "201" { // Se guarda el documento con tipo_documento = DTR_PLX
						transaccion.TrTrabajoGrado.DocumentoTrabajoGrado.TrabajoGrado.Id = idTrabajoGrado
						transaccion.TrTrabajoGrado.DocumentoTrabajoGrado.DocumentoEscrito.Id = int(resDocumentoEscrito["Id"].(float64))
						transaccion.TrTrabajoGrado.DocumentoEscrito.Id = int(resDocumentoEscrito["Id"].(float64))
						var resDocumentoTrabajoGrado map[string]interface{}
						url = "/v1/documento_trabajo_grado"
						if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoTrabajoGrado, &transaccion.TrTrabajoGrado.DocumentoTrabajoGrado); err == nil && status == "201" { //Se asocia el documento con el trabajo de grado
							transaccion.TrTrabajoGrado.DocumentoTrabajoGrado.Id = int(resDocumentoTrabajoGrado["Id"].(float64))

							//Funcionalidad para almacenar el Anteproyecto

							//Se recupera el id del tipo de documento ANP_PLX
							var tipoDocumento []models.TipoDocumento
							url = beego.AppConfig.String("UrlDocumentos") + "tipo_documento?query=CodigoAbreviacion:ANP_PLX"
							fmt.Println("RUTACOMPLETA!", url)
							if err := request.GetJson(url, &tipoDocumento); err != nil {
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
							if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoEscrito, &TrAnteproyecto.TrTrabajoGrado.DocumentoEscrito); err == nil && status == "201" { // Se guarda el documento con tipo_documento = ANP_PLX
								TrAnteproyecto.TrTrabajoGrado.DocumentoTrabajoGrado.TrabajoGrado.Id = idTrabajoGrado
								TrAnteproyecto.TrTrabajoGrado.DocumentoTrabajoGrado.DocumentoEscrito.Id = int(resDocumentoEscrito["Id"].(float64))
								TrAnteproyecto.TrTrabajoGrado.DocumentoEscrito.Id = int(resDocumentoEscrito["Id"].(float64))

								var resDocumentoTrabajoGradoAnt map[string]interface{}
								url = "/v1/documento_trabajo_grado"
								if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoTrabajoGradoAnt, &TrAnteproyecto.TrTrabajoGrado.DocumentoTrabajoGrado); err == nil && status == "201" { //Se asocia el anteproyecto con el trabajo de grado
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
					var espaciosAcademicosInscritos = make([]map[string]interface{}, 0)
					for i, v := range *transaccion.EspaciosAcademicosInscritos {
						var resEspaciosAcademicosInscritos map[string]interface{}
						v.TrabajoGrado.Id = idTrabajoGrado
						if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resEspaciosAcademicosInscritos, &v); err == nil && status == "201" {
							(*transaccion.EspaciosAcademicosInscritos)[i].Id = int(resEspaciosAcademicosInscritos["Id"].(float64))
							espaciosAcademicosInscritos = append(espaciosAcademicosInscritos, resEspaciosAcademicosInscritos)
						} else {
							logs.Error(err)
							if len(espaciosAcademicosInscritos) > 0 {
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
					if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDetallePasantia, &transaccion.DetallesPasantia); err == nil && status == "201" {
						transaccion.DetallesPasantia.Id = int(resDetallePasantia["Id"].(float64))

						//Se guardan los documentos asociados a la pasantia

						transaccion.DetallesPasantia.Contrato.Titulo = transaccion.DetallesPasantia.Contrato.Titulo + " del trabajo de grado con id: " + strconv.Itoa(idTrabajoGrado)

						transaccion.DetallesPasantia.Carta.Titulo = transaccion.DetallesPasantia.Carta.Titulo + " del trabajo de grado con id: " + strconv.Itoa(idTrabajoGrado)

						//Se envia el contrato
						url = "/v1/documento_escrito"
						var resDocumentoEscrito map[string]interface{}
						if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoEscrito, &transaccion.DetallesPasantia.Contrato); err == nil && status == "201" { // Se guarda el contrato
							transaccion.DetallesPasantia.DTG_Contrato.TrabajoGrado.Id = idTrabajoGrado
							transaccion.DetallesPasantia.Contrato.Id = int(resDocumentoEscrito["Id"].(float64))
							transaccion.DetallesPasantia.DTG_Contrato.DocumentoEscrito.Id = int(resDocumentoEscrito["Id"].(float64))

							var resDocumentoTrabajoGrado map[string]interface{}
							url = "/v1/documento_trabajo_grado"
							if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoTrabajoGrado, &transaccion.DetallesPasantia.DTG_Contrato); err == nil && status == "201" { //Se asocia el contrato con el trabajo de grado
								transaccion.DetallesPasantia.DTG_Contrato.Id = int(resDocumentoTrabajoGrado["Id"].(float64))

								//Se envia la Carta
								url = "/v1/documento_escrito"
								var resDocumentoEscrito map[string]interface{}
								if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoEscrito, &transaccion.DetallesPasantia.Carta); err == nil && status == "201" { // Se guarda la carta
									transaccion.DetallesPasantia.Carta.Id = int(resDocumentoEscrito["Id"].(float64))
									transaccion.DetallesPasantia.DTG_Carta.TrabajoGrado.Id = idTrabajoGrado
									transaccion.DetallesPasantia.DTG_Carta.DocumentoEscrito.Id = int(resDocumentoEscrito["Id"].(float64))

									var resDocumentoTrabajoGrado map[string]interface{}
									url = "/v1/documento_trabajo_grado"
									if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoTrabajoGrado, &transaccion.DetallesPasantia.DTG_Carta); err == nil && status == "201" { //Se asocia la carta con el trabajo de grado
										transaccion.DetallesPasantia.DTG_Carta.Id = int(resDocumentoTrabajoGrado["Id"].(float64))

										//Se envia la hoja de vida del director externo

										transaccion.DetallesPasantia.HojaVidaDE.Resumen = transaccion.DetallesPasantia.HojaVidaDE.Resumen + " con id: " + strconv.Itoa(idTrabajoGrado)

										url = "/v1/documento_escrito"
										var resDocumentoEscrito map[string]interface{}
										if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoEscrito, &transaccion.DetallesPasantia.HojaVidaDE); err == nil && status == "201" { // Se guarda la hoja de vida
											transaccion.DetallesPasantia.HojaVidaDE.Id = int(resDocumentoEscrito["Id"].(float64))
											transaccion.DetallesPasantia.DTG_HojaVida.TrabajoGrado.Id = idTrabajoGrado
											transaccion.DetallesPasantia.DTG_HojaVida.DocumentoEscrito.Id = int(resDocumentoEscrito["Id"].(float64))

											var resDocumentoTrabajoGrado map[string]interface{}
											url = "/v1/documento_trabajo_grado"
											if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoTrabajoGrado, &transaccion.DetallesPasantia.DTG_HojaVida); err == nil && status == "201" { //Se asocia la hv con el trabajo de grado
												transaccion.DetallesPasantia.DTG_HojaVida.Id = int(resDocumentoTrabajoGrado["Id"].(float64))

											} else {
												logs.Error(err)
												rollbackDocumentoTrabajoGradoPasantia(transaccion)
												panic(err.Error())
											}
										} else {
											logs.Error(err)
											rollbackDocumentoEscritoPasantia(transaccion)
											panic(err.Error())
										}
									} else {
										logs.Error(err)
										rollbackDocumentoTrabajoGradoPasantia(transaccion)
										panic(err.Error())
									}
								} else {
									logs.Error(err)
									rollbackDocumentoEscritoPasantia(transaccion)
									panic(err.Error())
								}

							} else {
								logs.Error(err)
								rollbackDocumentoTrabajoGradoPasantia(transaccion)
								panic(err.Error())
							}
						} else {
							logs.Error(err)
							rollbackDocumentoEscritoPasantia(transaccion)
							panic(err.Error())
						}

					} else {
						logs.Error(err)
						rollbackDetallesPasantia(transaccion)
						panic(err.Error())
					}
				}
				if transaccion.DetallesPasantiaExterna != nil && parametro.CodigoAbreviacion == "PAS_PLX" {
					url = "/v1/detalle_trabajo_grado"
					var detallesPasantiaExterna = make([]map[string]interface{}, 0)
					for i, v := range *transaccion.DetallesPasantiaExterna {
						var resDetallesPasantiaExterna map[string]interface{}
						v.TrabajoGrado.Id = idTrabajoGrado
						if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDetallesPasantiaExterna, &v); err == nil && status == "201" {
							(*transaccion.DetallesPasantiaExterna)[i].Id = int(resDetallesPasantiaExterna["Id"].(float64))
							detallesPasantiaExterna = append(detallesPasantiaExterna, resDetallesPasantiaExterna)
						} else {
							logs.Error(err)
							if len(detallesPasantiaExterna) > 0 {
								rollbackDetallesPasantiaExterna(transaccion)
								logs.Error(err)
								panic(err.Error())
							} else {
								rollbackDocumentoTrabajoGrado(transaccion)
								logs.Error(err)
								panic(err.Error())
							}
						}
					}
				}
				if transaccion.TrTrabajoGrado.DocumentosMaterias != nil {
					//se almacenan los documentos de la S.I de Materias de Posgrado
					var resDocumentoEscrito map[string]interface{}

					//Se almacena la Solicitud Escrita
					url = "/v1/documento_escrito"
					if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoEscrito, &transaccion.TrTrabajoGrado.DocumentosMaterias.SolicitudEscrita); err == nil && status == "201" {

						transaccion.TrTrabajoGrado.DocumentosMaterias.DTG_SolicitudEscrita.TrabajoGrado.Id = idTrabajoGrado
						transaccion.TrTrabajoGrado.DocumentosMaterias.DTG_SolicitudEscrita.DocumentoEscrito.Id = int(resDocumentoEscrito["Id"].(float64))
						transaccion.TrTrabajoGrado.DocumentosMaterias.SolicitudEscrita.Id = int(resDocumentoEscrito["Id"].(float64))

						//Se almacena DTG de Solicitud Escrita
						var resDocumentoTrabajoGrado map[string]interface{}
						url = "/v1/documento_trabajo_grado"
						if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoTrabajoGrado, &transaccion.TrTrabajoGrado.DocumentosMaterias.DTG_SolicitudEscrita); err == nil && status == "201" {
							transaccion.TrTrabajoGrado.DocumentosMaterias.DTG_SolicitudEscrita.Id = int(resDocumentoTrabajoGrado["Id"].(float64))
						}
					}

					//Se almacena la Justificación
					url = "/v1/documento_escrito"
					if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoEscrito, &transaccion.TrTrabajoGrado.DocumentosMaterias.Justificacion); err == nil && status == "201" {

						transaccion.TrTrabajoGrado.DocumentosMaterias.DTG_Justificacion.TrabajoGrado.Id = idTrabajoGrado
						transaccion.TrTrabajoGrado.DocumentosMaterias.DTG_Justificacion.DocumentoEscrito.Id = int(resDocumentoEscrito["Id"].(float64))
						transaccion.TrTrabajoGrado.DocumentosMaterias.Justificacion.Id = int(resDocumentoEscrito["Id"].(float64))

						//Se almacena DTG de Solicitud Escrita
						var resDocumentoTrabajoGrado map[string]interface{}
						url = "/v1/documento_trabajo_grado"
						if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoTrabajoGrado, &transaccion.TrTrabajoGrado.DocumentosMaterias.DTG_Justificacion); err == nil && status == "201" {
							transaccion.TrTrabajoGrado.DocumentosMaterias.DTG_Justificacion.Id = int(resDocumentoTrabajoGrado["Id"].(float64))
						}
					}

					//Se almacena la Carta Aceptación
					url = "/v1/documento_escrito"
					if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoEscrito, &transaccion.TrTrabajoGrado.DocumentosMaterias.CartaAceptacion); err == nil && status == "201" {

						transaccion.TrTrabajoGrado.DocumentosMaterias.DTG_CartaAceptacion.TrabajoGrado.Id = idTrabajoGrado
						transaccion.TrTrabajoGrado.DocumentosMaterias.DTG_CartaAceptacion.DocumentoEscrito.Id = int(resDocumentoEscrito["Id"].(float64))
						transaccion.TrTrabajoGrado.DocumentosMaterias.CartaAceptacion.Id = int(resDocumentoEscrito["Id"].(float64))

						//Se almacena DTG de Solicitud Escrita
						var resDocumentoTrabajoGrado map[string]interface{}
						url = "/v1/documento_trabajo_grado"
						if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoTrabajoGrado, &transaccion.TrTrabajoGrado.DocumentosMaterias.DTG_CartaAceptacion); err == nil && status == "201" {
							transaccion.TrTrabajoGrado.DocumentosMaterias.DTG_CartaAceptacion.Id = int(resDocumentoTrabajoGrado["Id"].(float64))
						}
					}

					//Se almacena la Sábana de Notas
					url = "/v1/documento_escrito"
					if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoEscrito, &transaccion.TrTrabajoGrado.DocumentosMaterias.SabanaNotas); err == nil && status == "201" {

						transaccion.TrTrabajoGrado.DocumentosMaterias.DTG_SabanaNotas.TrabajoGrado.Id = idTrabajoGrado
						transaccion.TrTrabajoGrado.DocumentosMaterias.DTG_SabanaNotas.DocumentoEscrito.Id = int(resDocumentoEscrito["Id"].(float64))
						transaccion.TrTrabajoGrado.DocumentosMaterias.SabanaNotas.Id = int(resDocumentoEscrito["Id"].(float64))

						//Se almacena DTG de Solicitud Escrita
						var resDocumentoTrabajoGrado map[string]interface{}
						url = "/v1/documento_trabajo_grado"
						if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoTrabajoGrado, &transaccion.TrTrabajoGrado.DocumentosMaterias.DTG_SabanaNotas); err == nil && status == "201" {
							transaccion.TrTrabajoGrado.DocumentosMaterias.DTG_SabanaNotas.Id = int(resDocumentoTrabajoGrado["Id"].(float64))
						}
					}
				}
			} else {
				logs.Error(err)
				rollbackTrabajoGrado(transaccion)
				logs.Error(err)
				panic(err.Error())
			}
		} else {
			logs.Error(err)
			rollbackDocumentoSolicitud(transaccion)
			logs.Error(err)
			panic(err.Error())
		}
	}

	// Solicitud de cambio de docente evaluador o docente director
	if transaccion.Vinculaciones != nil {
		var idVinculadoAntiguo int
		var idVinculadoNuevo int64
		var vinculacionesTrabajoGrado = make([]map[string]interface{}, 0)
		var vinculacionesOriginalesTrabajoGrado []models.VinculacionTrabajoGrado
		var vinculacionesTrabajoGradoPost = make([]map[string]interface{}, 0)
		var vinculacionesTrabajoGradoCanceladas []models.VinculacionTrabajoGrado
		for _, v := range *transaccion.Vinculaciones {
			//Si esta activo es nuevo y se inserta sino se actualiza la fecha de fin y el activo
			if v.Activo {
				// Se buscar si el docente ya estuvo vinculado y se actualiza
				var vinculado []models.VinculacionTrabajoGrado
				url = "/v1/vinculacion_trabajo_grado?query=TrabajoGrado__Id:" + strconv.Itoa(v.TrabajoGrado.Id) +
					",Usuario:" + strconv.Itoa(v.Usuario) + ",RolTrabajoGrado:" + strconv.Itoa(v.RolTrabajoGrado) + "&limit=1"
				fmt.Println("URL ", url)
				if err := GetRequestNew("PoluxCrudUrl", url, &vinculado); err != nil {
					logs.Error(err.Error())
					panic(err.Error())
				}

				if len(vinculado) > 0 {
					var vinculadoAux = vinculado[0]
					idVinculadoNuevo = int64(vinculado[0].Id)
					vinculado[0].Activo = v.Activo
					vinculado[0].FechaFin = v.FechaFin
					vinculado[0].FechaInicio = v.FechaInicio
					var resVinculacionTrabajoGrado map[string]interface{}
					url = "/v1/vinculacion_trabajo_grado/" + strconv.Itoa(vinculado[0].Id)
					if status, err := SendRequestNew("PoluxCrudUrl", url, "PUT", &resVinculacionTrabajoGrado, &vinculado[0]); err == nil && status == "200" {
						vinculacionesOriginalesTrabajoGrado = append(vinculacionesOriginalesTrabajoGrado, vinculadoAux)
						vinculacionesTrabajoGrado = append(vinculacionesTrabajoGrado, resVinculacionTrabajoGrado)
					} else {
						logs.Error(err)
						if len(vinculacionesTrabajoGrado) > 0 || len(vinculacionesTrabajoGradoPost) > 0 {
							rollbackVinculacionTrabajoGradoRS(transaccion, vinculacionesOriginalesTrabajoGrado)
							rollbackVinculacionTrabajoGradoRSPost(transaccion, vinculacionesTrabajoGradoPost)
							logs.Error(err)
							panic(err.Error())
						} else {
							rollbackDocumentoSolicitud(transaccion)
							logs.Error(err)
							panic(err.Error())
						}
					}
				} else {
					var resVinculacionTrabajoGrado map[string]interface{}
					url = "/v1/vinculacion_trabajo_grado"
					if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resVinculacionTrabajoGrado, &v); err == nil && status == "201" {
						idVinculadoNuevo = int64(resVinculacionTrabajoGrado["Id"].(float64))
						vinculacionesTrabajoGradoPost = append(vinculacionesTrabajoGradoPost, resVinculacionTrabajoGrado)
					} else {
						logs.Error(err)
						if len(vinculacionesTrabajoGrado) > 0 || len(vinculacionesTrabajoGradoPost) > 0 {
							rollbackVinculacionTrabajoGradoRS(transaccion, vinculacionesOriginalesTrabajoGrado)
							rollbackVinculacionTrabajoGradoRSPost(transaccion, vinculacionesTrabajoGradoPost)
							logs.Error(err)
							panic(err.Error())
						} else {
							rollbackDocumentoSolicitud(transaccion)
							logs.Error(err)
							panic(err.Error())
						}
					}
				}

			} else {
				idVinculadoAntiguo = v.Id
				var resVinculacionTrabajoGrado models.VinculacionTrabajoGrado
				url = "/v1/vinculacion_trabajo_grado/" + strconv.Itoa(v.Id)
				if status, err := SendRequestNew("PoluxCrudUrl", url, "PUT", &resVinculacionTrabajoGrado, &v); err == nil && status == "200" {
					vinculacionesTrabajoGradoCanceladas = append(vinculacionesTrabajoGradoCanceladas, resVinculacionTrabajoGrado)
				} else {
					if len(vinculacionesTrabajoGrado) > 0 || len(vinculacionesTrabajoGradoPost) > 0 {
						rollbackVinculacionTrabajoGradoRS(transaccion, vinculacionesOriginalesTrabajoGrado)
						rollbackVinculacionTrabajoGradoRSPost(transaccion, vinculacionesTrabajoGradoPost)
						logs.Error(err)
						panic(err.Error())
					} else {
						rollbackDocumentoSolicitud(transaccion)
						logs.Error(err)
						panic(err.Error())
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
		url = "/v1/revision_trabajo_grado?query=VinculacionTrabajoGrado:" + strconv.Itoa(idVinculadoAntiguo) +
			",EstadoRevisionTrabajoGrado:" + strconv.Itoa(parametroEstadoRevision[0].Id) + "&limit=1"
		fmt.Println("URL ", url)
		if err := GetRequestNew("PoluxCrudUrl", url, &revisionTrabajoGrado); err != nil {
			logs.Error(err.Error())
			panic(err.Error())
		}

		var vincOrig models.RevisionTrabajoGrado

		if len(revisionTrabajoGrado) > 0 {
			vincOrig = revisionTrabajoGrado[0]
		}

		// Verificación adicional para asegurar que la revisión encontrada es válida
		if len(revisionTrabajoGrado) > 0 && revisionTrabajoGrado[0].Id != 0 {
			revisionTrabajoGrado[0].VinculacionTrabajoGrado.Id = int(idVinculadoNuevo)
			var resRevisionTrabajoGrado map[string]interface{}
			url = "/v1/revision_trabajo_grado/" + strconv.Itoa(revisionTrabajoGrado[0].Id)
			if status, err := SendRequestNew("PoluxCrudUrl", url, "PUT", &resRevisionTrabajoGrado, &revisionTrabajoGrado[0]); err != nil && status == "200" {
				rollbackRevisionTrabajoGrado(transaccion, &vincOrig, vinculacionesOriginalesTrabajoGrado, vinculacionesTrabajoGradoPost, vinculacionesTrabajoGradoCanceladas)
				logs.Error(err)
				panic(err.Error())
			}
		}

		// Si  el cambio es de director externo, se recibe la data del detalle de la pasantia y
		// se actualiza
		if transaccion.DetallesPasantia != nil {
			// Se busca el detalle de la pasantia asociado al tg
			var detallePasantia []models.DetallePasantia
			url = "/v1/detalle_pasantia?query=TrabajoGrado__Id:" + strconv.Itoa(transaccion.DetallesPasantia.TrabajoGrado.Id) + "&limit=1"
			fmt.Println("URL ", url)
			if err := GetRequestNew("PoluxCrudUrl", url, &detallePasantia); err != nil {
				logs.Error(err.Error())
				panic(err.Error())
			}

			var detallePre = detallePasantia[0].Observaciones

			detallePasantia[0].Observaciones = strings.Split(detallePasantia[0].Observaciones, " y dirigida por ")[0]
			detallePasantia[0].Observaciones += transaccion.DetallesPasantia.Observaciones

			var resDetallePasantia map[string]interface{}
			url = "/v1/detalle_pasantia/" + strconv.Itoa(detallePasantia[0].Id)
			if status, err := SendRequestNew("PoluxCrudUrl", url, "PUT", &resDetallePasantia, &detallePasantia[0]); err != nil && status != "200" {
				rollbackRevisionTrabajoGrado(transaccion, &vincOrig, vinculacionesOriginalesTrabajoGrado, vinculacionesTrabajoGradoPost, vinculacionesTrabajoGradoCanceladas)
				logs.Error(err)
				panic(err.Error())
			}

			//Se guarda la hoja de vida del nuevo director externo
			var idTrabajoGrado = transaccion.DetallesPasantia.TrabajoGrado.Id
			transaccion.DetallesPasantia.HojaVidaDE.Resumen = transaccion.DetallesPasantia.HojaVidaDE.Resumen + " con id: " + strconv.Itoa(idTrabajoGrado)

			url := "/v1/documento_escrito"
			var resDocumentoEscrito map[string]interface{}
			if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoEscrito, &transaccion.DetallesPasantia.HojaVidaDE); err == nil && status == "201" { //Se guarda la ARL en Documento Escrito. Desde el Cliente ya viene con el tipo "Documentos adicionales asociados a las pasantias"

				transaccion.DetallesPasantia.HojaVidaDE.Id = int(resDocumentoEscrito["Id"].(float64))

				//se actualiza la hoja de vida del Director Externo

				var tipoDocumento []models.TipoDocumento
				url = beego.AppConfig.String("UrlDocumentos") + "tipo_documento?query=CodigoAbreviacion:HVDE_PLX"
				if err := GetJson(url, &tipoDocumento); err != nil { //Se busca el Tipo de Documento asociado a la Hoja de Vida del Director Externo
					logs.Error(err.Error())
					panic(err.Error())
				}

				var documentosTG []models.DocumentoTrabajoGrado
				url = "/v1/documento_trabajo_grado?query=trabajo_grado__Id:" + strconv.Itoa(idTrabajoGrado) + ",documento_escrito__tipo_documento_escrito:" + strconv.Itoa(tipoDocumento[0].Id)
				if err := GetRequestNew("PoluxCrudUrl", url, &documentosTG); err != nil { //Se busca el registro en documento_trabajo_grado en la que relacione la HV con el trabajo de grado
					logs.Error(err.Error())
					panic(err.Error())
				}

				transaccion.DetallesPasantia.DTG_HojaVida.Id = documentosTG[0].Id
				transaccion.DetallesPasantia.DTG_HojaVida.DocumentoEscrito.Id = int(resDocumentoEscrito["Id"].(float64))
				transaccion.DetallesPasantia.DTG_HojaVida.TrabajoGrado.Id = idTrabajoGrado
				transaccion.DetallesPasantia.DTG_HojaVida.Activo = true;
				transaccion.DetallesPasantia.DTG_HojaVida.FechaCreacion = documentosTG[0].FechaCreacion
				transaccion.DetallesPasantia.DTG_HojaVida.FechaModificacion = documentosTG[0].FechaModificacion

				url := "/v1/documento_trabajo_grado/" + strconv.Itoa(documentosTG[0].Id)
				var resDocumentoTrabajoGrado map[string]interface{}
				if status, err := SendRequestNew("PoluxCrudUrl", url, "PUT", &resDocumentoTrabajoGrado, &transaccion.DetallesPasantia.DTG_HojaVida); err == nil && status == "200" { //Se actualiza la relación entre la HV y el Trabajo de Grado
					transaccion.DetallesPasantia.DTG_HojaVida.Id = int(resDocumentoTrabajoGrado["Id"].(float64))
				} else {
					rollbackDocEsHv(transaccion, detallePre, &vincOrig, vinculacionesOriginalesTrabajoGrado, vinculacionesTrabajoGradoPost, vinculacionesTrabajoGradoCanceladas)
				}
			} else {
				rollbackDetallesPasantiaCambioDirector(transaccion, detallePre, &vincOrig, vinculacionesOriginalesTrabajoGrado, vinculacionesTrabajoGradoPost, vinculacionesTrabajoGradoCanceladas)
			}
		}
	}

	//Solicitud de cambio de nombre del trabajo de grado
	if transaccion.TrabajoGrado != nil {
		var resTrabajoGrado map[string]interface{}
		url := "/v1/trabajo_grado/" + strconv.Itoa(transaccion.TrabajoGrado.Id)
		if status, err := SendRequestNew("PoluxCrudUrl", url, "PUT", &resTrabajoGrado, &transaccion.TrabajoGrado); err != nil && status != "200" {
			logs.Error(err)
			rollbackDocumentoSolicitud(transaccion)
			logs.Error(err)
			panic(err.Error())
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
				if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDetalleTrabajoGrado, &data); err != nil && status != "201" {
					logs.Error(err)
					rollbackDocumentoSolicitud(transaccion)
					logs.Error(err)
					panic(err.Error())
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
		url = "/v1/estudiante_trabajo_grado?query=TrabajoGrado__Id:" + strconv.Itoa(transaccion.EstudianteTrabajoGrado.TrabajoGrado.Id) +
			",Estudiante:" + transaccion.EstudianteTrabajoGrado.Estudiante + ",EstadoEstudianteTrabajoGrado:" + strconv.Itoa(parametroEstadoEstudianteTrGr[0].Id) + "&limit=1"
		fmt.Println("URL ", url)
		if err := GetRequestNew("PoluxCrudUrl", url, &estudianteTrabajoGrado); err != nil {
			logs.Error(err.Error())
			panic(err.Error())
		}
		transaccion.EstudianteTrabajoGrado.Id = estudianteTrabajoGrado[0].Id
		var resEstudianteTrabajoGrado map[string]interface{}
		url := "/v1/estudiante_trabajo_grado/" + strconv.Itoa(estudianteTrabajoGrado[0].Id)
		if status, err := SendRequestNew("PoluxCrudUrl", url, "PUT", &resEstudianteTrabajoGrado, &transaccion.EstudianteTrabajoGrado); err == nil && status == "200" {
			var estudianteTrabajoGradoAux []models.EstudianteTrabajoGrado
			url = "/v1/estudiante_trabajo_grado?query=TrabajoGrado__Id:" + strconv.Itoa(transaccion.EstudianteTrabajoGrado.TrabajoGrado.Id) +
				",EstadoEstudianteTrabajoGrado:" + strconv.Itoa(parametroEstadoEstudianteTrGr[0].Id)
			fmt.Println("URL ", url)
			if err := GetRequestNew("PoluxCrudUrl", url, &estudianteTrabajoGradoAux); err != nil {
				logs.Error(err.Error())
				panic(err.Error())
			}

			fmt.Println("acá se soluciona esto!!!", len(estudianteTrabajoGradoAux))
			// si no hay estudiantes vinculados se inactivan las vinculaciones y se cancela el tg
			if len(estudianteTrabajoGradoAux) == 0 || estudianteTrabajoGradoAux[0].Id == 0 {
				fmt.Println("acá se soluciona esto 2222!!!", len(estudianteTrabajoGradoAux))
				// Se inactivan las vinculaciones
				var vinculacionesTrGr = make([]map[string]interface{}, 0)
				for _, v := range *transaccion.VinculacionesCancelacion {
					v.Activo = false
					var resVinculacionTrabajoGrado map[string]interface{}
					url := "/v1/vinculacion_trabajo_grado/" + strconv.Itoa(v.Id)
					if status, err := SendRequestNew("PoluxCrudUrl", url, "PUT", &resVinculacionTrabajoGrado, v); err == nil && status == "200" {
						vinculacionesTrGr = append(vinculacionesTrGr, resVinculacionTrabajoGrado)
					} else {
						if len(vinculacionesTrGr) > 0 {
							logs.Error(err)
							rollbackVincTrGrCanc(transaccion)
							logs.Error(err)
							panic(err.Error())
						} else {
							rollbackEstTrGrCanc(transaccion)
							logs.Error(err)
							panic(err.Error())
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
				if status, err := SendRequestNew("PoluxCrudUrl", url, "PUT", &resTrabajoGrado, tg); err != nil && status != "200" {
					logs.Error(err)
					rollbackVincTrGrCanc(transaccion)
					logs.Error(err)
					panic(err.Error())
				}

				// Actualizar asignaturas trabajo de grado a cancelado
				var asignaturasTrabajoGrado []models.AsignaturaTrabajoGrado
				// Se busca asignaturas trabajo grado
				url = "/v1/asignatura_trabajo_grado?query=TrabajoGrado__Id:" + strconv.Itoa(transaccion.EstudianteTrabajoGrado.TrabajoGrado.Id)
				fmt.Println("URL ", url)
				if err := GetRequestNew("PoluxCrudUrl", url, &asignaturasTrabajoGrado); err != nil {
					logs.Error(err.Error())
					panic(err.Error())
				}
				var parametroEstAsTrGr []models.Parametro
				url = "parametro?query=CodigoAbreviacion:CNC_PLX,TipoParametroId__CodigoAbreviacion:EST_ASIG_TRG"
				if err := GetRequestNew("UrlCrudParametros", url, &parametroEstAsTrGr); err != nil {
					logs.Error(err.Error())
					panic(err.Error())
				}
				var asignaturasTrGr = make([]map[string]interface{}, 0)
				for _, v := range asignaturasTrabajoGrado {
					//Id de la asignatura cancelada
					v.EstadoAsignaturaTrabajoGrado = parametroEstAsTrGr[0].Id
					var resAsignaturaTrabajoGrado map[string]interface{}
					url := "/v1/asignatura_trabajo_grado/" + strconv.Itoa(v.Id)
					if status, err := SendRequestNew("PoluxCrudUrl", url, "PUT", &resAsignaturaTrabajoGrado, v); err == nil && status == "200" {
						asignaturasTrGr = append(asignaturasTrGr, resAsignaturaTrabajoGrado)
					} else {
						logs.Error(err)
						if len(asignaturasTrGr) > 0 {
							rollbackAsTrGr(transaccion, &asignaturasTrabajoGrado)
							logs.Error(err)
							panic(err.Error())
						} else {
							rollbackTrGrCanc(transaccion)
							logs.Error(err)
							panic(err.Error())
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
				url = "/v1/espacio_academico_inscrito?query=TrabajoGrado__Id:" + strconv.Itoa(transaccion.EstudianteTrabajoGrado.TrabajoGrado.Id) +
					",EstadoEspacioAcademicoInscrito:" + strconv.Itoa(parametroEspAcadIns[0].Id)
				fmt.Println("URL ", url)
				if err := GetRequestNew("PoluxCrudUrl", url, &espaciosAcademicosInscritos); err != nil {
					logs.Error(err.Error())
					panic(err.Error())
				}
				if len(espaciosAcademicosInscritos) > 0 {
					if espaciosAcademicosInscritos[0].Id != 0 {
						var parametroEspAcadInsAux []models.Parametro
						url = "parametro?query=CodigoAbreviacion:ESP_CAN_PLX,TipoParametroId__CodigoAbreviacion:EST_ESP"
						if err := GetRequestNew("UrlCrudParametros", url, &parametroEspAcadInsAux); err != nil {
							logs.Error(err.Error())
							panic(err.Error())
						}
						var espaciosAcadInsc = make([]map[string]interface{}, 0)
						for _, v := range espaciosAcademicosInscritos {
							// Id del espacio cancelado
							v.EstadoEspacioAcademicoInscrito = parametroEspAcadInsAux[0].Id
							var resEspacioAcademicoInscrito map[string]interface{}
							url := "/v1/espacio_academico_inscrito/" + strconv.Itoa(v.Id)
							if status, err := SendRequestNew("PoluxCrudUrl", url, "PUT", &resEspacioAcademicoInscrito, v); err == nil && status == "200" {
								espaciosAcadInsc = append(espaciosAcadInsc, resEspacioAcademicoInscrito)
							} else {
								logs.Error(err)
								if len(espaciosAcadInsc) > 0 {
									rollbackEsAcadInsc(transaccion, &asignaturasTrabajoGrado, &espaciosAcademicosInscritos)
									logs.Error(err)
									panic(err.Error())
								} else {
									rollbackAsTrGr(transaccion, &asignaturasTrabajoGrado)
									logs.Error(err)
									panic(err.Error())
								}
							}
						}
					}
				}

			}
		} else {
			logs.Error(err)
			rollbackDocumentoSolicitud(transaccion)
			logs.Error(err)
			panic(err.Error())
		}
	}

	// Solicitud de revisión del trabajo de grado
	if transaccion.TrRevision != nil {
		// Se actualiza el trabajo de grado
		var resTrabajoGrado map[string]interface{}
		trabajoGradoId := strconv.Itoa(transaccion.TrRevision.TrabajoGrado.Id)
		url := "/v1/trabajo_grado/" + trabajoGradoId
		if status, err := SendRequestNew("PoluxCrudUrl", url, "PUT", &resTrabajoGrado, &transaccion.TrRevision.TrabajoGrado); err != nil && status != "200" {
			logs.Error(err)
			rollbackDocumentoSolicitud(transaccion)
			logs.Error(err)
			panic(err.Error())
		}

		//INSERTA EN LA TABLA DETALLE TRABAJO GRADO
		if transaccion.TrRevision.DetalleTrabajoGrado != nil {
			var detallesTrabajoGrado = make([]map[string]interface{}, 0)
			for _, data := range *transaccion.TrRevision.DetalleTrabajoGrado {
				data.Activo = true
				data.FechaCreacion = time_bogota.TiempoBogotaFormato()
				data.FechaModificacion = time_bogota.TiempoBogotaFormato()
				var resDetalleTrabajoGrado map[string]interface{}
				url = "/v1/detalle_trabajo_grado"
				if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDetalleTrabajoGrado, &data); err == nil && status == "201" {
					data.Id = int(resDetalleTrabajoGrado["Id"].(float64))
					detallesTrabajoGrado = append(detallesTrabajoGrado, resDetalleTrabajoGrado)
				} else {
					if len(detallesTrabajoGrado) > 0 {
						rollbackDetTrGrRev(transaccion)
						logs.Error(err)
						panic(err.Error())
					} else {
						rollbackTrGrRev(transaccion)
						logs.Error(err)
						panic(err.Error())
					}
				}
			}
		}

		// Se obtiene ID de tipo documento anexo para la obtención de documentos activos en documento_trabajo_grado
		var tipoDocumentoAux []models.TipoDocumento
		urlTipoDocumento := beego.AppConfig.String("UrlDocumentos") + "tipo_documento?query=CodigoAbreviacion:ANX_PLX"
		// Consulta el tipo de documento con CodigoAbreviacion ANX_PLX
		if err := GetJson(urlTipoDocumento, &tipoDocumentoAux); err != nil {
			logs.Error(err.Error())
			panic(err.Error())
		}
		fmt.Println("CONSULTA TIPO DOCUMENTO!!!", tipoDocumentoAux)

		tipoDocumentoId := tipoDocumentoAux[0].Id

		// Realiza la segunda consulta en documento_trabajo_grado con el Id de tipo de documento anexo
		var documentosTrabajoGradoAux []models.DocumentoTrabajoGrado
		urlDocumentosTrabajoGrado := "/v1/documento_trabajo_grado?query=DocumentoEscrito.TipoDocumentoEscrito:" +
			strconv.Itoa(tipoDocumentoId) + ",TrabajoGrado.Id:" + trabajoGradoId + ",Activo:true&sortby=id&order=desc"
		fmt.Println("URL!!!", urlDocumentosTrabajoGrado)
		if err := GetRequestNew("PoluxCrudUrl", urlDocumentosTrabajoGrado, &documentosTrabajoGradoAux); err != nil {
			logs.Error(err.Error())
			panic(err.Error())
		}
		if len(documentosTrabajoGradoAux) > 0 {
			// Para cada documento activo, se actualiza el campo `Activo` a `false`
			fmt.Println("LLEGA AL FOR", documentosTrabajoGradoAux)
			for _, documento := range documentosTrabajoGradoAux {
				fmt.Println("REGISTRO DOCUMENTO TRABAJO GRADO", documento)
				documento.Activo = false
				urlUpdate := "/v1/documento_trabajo_grado/" + strconv.Itoa(documento.Id)
				var resUpdate map[string]interface{}

				// Realizamos la solicitud PUT para desactivar el documento
				if status, err := SendRequestNew("PoluxCrudUrl", urlUpdate, "PUT", &resUpdate, &documento); err != nil || status != "200" {
					logs.Error("Error al actualizar documento: ", err)
					rollbackDocumentoSolicitud(transaccion) // Rollback en caso de error
					//return nil, err
				}
			}
		}

		// Se inserta el documento final de la revisión y se relaciona con el trabajo de grado
		// Itera sobre el array de documentos escritos
		for _, documentoEscrito := range *transaccion.TrRevision.DocumentoEscrito {
			var resDocumentoEscrito map[string]interface{}
			url := "/v1/documento_escrito"

			// Envía la solicitud para cada documento escrito
			if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoEscrito, &documentoEscrito); err == nil && status == "201" {
				// Asigna el ID devuelto al documento escrito actual
				documentoEscrito.Id = int(resDocumentoEscrito["Id"].(float64))

				// Relaciona el documento escrito con el trabajo de grado
				transaccion.TrRevision.DocumentoTrabajoGrado.DocumentoEscrito.Id = int(resDocumentoEscrito["Id"].(float64))

				var resDocumentoTrabajoGrado map[string]interface{}
				url = "/v1/documento_trabajo_grado"

				// Envía la solicitud para asociar el documento con el trabajo de grado
				if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoTrabajoGrado, &transaccion.TrRevision.DocumentoTrabajoGrado); err == nil && status == "201" {
					documentoEscrito.Id = int(resDocumentoTrabajoGrado["Id"].(float64))
				} else {
					// En caso de error, realizar el rollback y detener el proceso
					rollbackDocEscrRev(transaccion)
					logs.Error(err)
					panic(err.Error())
				}
			} else {
				// En caso de error, realizar el rollback y detener el proceso
				rollbackDocEscrRev(transaccion)
				logs.Error(err)
				panic(err.Error())
			}
		}

		// Se actualizan las vinculaciones
		var vinculacionesTrabajoGrado = make([]map[string]interface{}, 0)
		var vinculacionesOriginalesTrabajoGrado []models.VinculacionTrabajoGrado
		var vinculacionesTrabajoGradoPost = make([]map[string]interface{}, 0)
		for _, v := range *transaccion.TrRevision.Vinculaciones {
			// Si esta activo es nuevo y se inserta sino se actualiza la fecha de fin y el activo
			if v.Activo {
				// Se buscar si el docente ya estuvo vinculado y se actualiza
				var vinculado []models.VinculacionTrabajoGrado
				url = "/v1/vinculacion_trabajo_grado?query=TrabajoGrado__Id:" + strconv.Itoa(v.TrabajoGrado.Id) +
					",Usuario:" + strconv.Itoa(v.Usuario) + ",RolTrabajoGrado:" + strconv.Itoa(v.RolTrabajoGrado) + "&limit=1"
				fmt.Println("URL ", url)
				if err := GetRequestNew("PoluxCrudUrl", url, &vinculado); err != nil {
					logs.Error(err.Error())
					panic(err.Error())
				}
				if len(vinculado) > 0 {
					if vinculado[0].Id != 0 {
						var vinculadoAux = vinculado[0]
						vinculado[0].Activo = v.Activo
						vinculado[0].FechaFin = v.FechaFin
						vinculado[0].FechaInicio = v.FechaInicio
						var resVinculacionTrabajoGrado map[string]interface{}
						url = "/v1/vinculacion_trabajo_grado/" + strconv.Itoa(vinculado[0].Id)
						if status, err := SendRequestNew("PoluxCrudUrl", url, "PUT", &resVinculacionTrabajoGrado, &vinculado[0]); err == nil && status == "200" {
							vinculacionesOriginalesTrabajoGrado = append(vinculacionesOriginalesTrabajoGrado, vinculadoAux)
							vinculacionesTrabajoGrado = append(vinculacionesTrabajoGrado, resVinculacionTrabajoGrado)
						} else {
							logs.Error(err)
							if len(vinculacionesTrabajoGrado) > 0 || len(vinculacionesTrabajoGradoPost) > 0 {
								rollbackVincTrGrRev(transaccion, vinculacionesOriginalesTrabajoGrado)
								rollbackVincTrGrPostRev(transaccion, vinculacionesTrabajoGradoPost)
								rollbackDocTrGrRev(transaccion)
								logs.Error(err)
								panic(err.Error())
							} else {
								rollbackDocTrGrRev(transaccion)
								logs.Error(err)
								panic(err.Error())
							}
						}
					} else {
						var resVinculacionTrabajoGrado map[string]interface{}
						url = "/v1/vinculacion_trabajo_grado"
						if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resVinculacionTrabajoGrado, &v); err == nil && status == "201" {
							vinculacionesTrabajoGradoPost = append(vinculacionesTrabajoGradoPost, resVinculacionTrabajoGrado)
						} else {
							logs.Error(err)
							if len(vinculacionesTrabajoGrado) > 0 || len(vinculacionesTrabajoGradoPost) > 0 {
								rollbackVincTrGrRev(transaccion, vinculacionesOriginalesTrabajoGrado)
								rollbackVincTrGrPostRev(transaccion, vinculacionesTrabajoGradoPost)
								rollbackDocTrGrRev(transaccion)
								logs.Error(err)
								panic(err.Error())
							} else {
								rollbackDocTrGrRev(transaccion)
								logs.Error(err)
								panic(err.Error())
							}
						}
					}
				}

			}
		}
	}
	response = map[string]interface{}{
		"RespuestaAnterior": resRespuestaAnterior,
		//"RespuestaNueva":    resRespuestaNueva,
		//"DetalleSolicitud":  detalleSolicitud,
		//"VinculacionTrabajoGrado":  resVinculacionTrabajoGrado,
	}
	return response, outputError
}

func rollbackResAnterior(transaccion *models.TrRespuestaSolicitud) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK RES ANTERIOR")
	var respuesta map[string]interface{}
	transaccion.RespuestaAnterior.Activo = true
	url := "/v1/respuesta_solicitud/" + strconv.Itoa(transaccion.RespuestaAnterior.Id)
	if status, err := SendRequestNew("PoluxCrudUrl", url, "PUT", &respuesta, &transaccion.RespuestaAnterior); err != nil && status != "200" {
		panic("Rollback respuesta anteror " + err.Error())
	}
	return nil
}

func rollbackResNueva(transaccion *models.TrRespuestaSolicitud) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK RES NUEVA")
	var respuesta map[string]interface{}
	if transaccion.RespuestaNueva.Id != 0 {
		url := "/v1/respuesta_solicitud/" + strconv.Itoa(transaccion.RespuestaNueva.Id)
		if status, err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil && status != "200" {
			panic("Rollback respuesta nueva" + err.Error())
		}
	}
	rollbackResAnterior(transaccion)
	return nil
}

func rollbackDocumentoSolicitud(transaccion *models.TrRespuestaSolicitud) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK DOCUMENTO SOLICITUD")
	var respuesta map[string]interface{}
	if transaccion.DocumentoSolicitud.Id != 0 {
		url := "/v1/documento_solicitud/" + strconv.Itoa(transaccion.DocumentoSolicitud.Id)
		if status, err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil && status != "200" {
			panic("Rollback documento solicitud" + err.Error())
		}
	}
	rollbackResNueva(transaccion)
	return nil
}

func rollbackDocumentoSolicitudCambioDirEx(transaccion *models.TrRespuestaSolicitud) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK DOCUMENTO SOLICITUD")
	var respuesta map[string]interface{}

	if transaccion.DocumentoSolicitud.DocumentoEscrito.Id != 0 {
		url := "/v1/documento_solicitud/" + strconv.Itoa(transaccion.DocumentoSolicitud.Id)
		if status, err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil && status != "200" {
			panic("Rollback documento solicitud" + err.Error())
		}
	}
	rollbackResNueva(transaccion)

	return nil
}

func rollbackTrabajoGrado(transaccion *models.TrRespuestaSolicitud) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK TRABAJO GRADO")
	var respuesta map[string]interface{}
	url := "/v1/trabajo_grado/" + strconv.Itoa(transaccion.TrTrabajoGrado.TrabajoGrado.Id)
	if status, err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil && status != "200" {
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
	if status, err := SendRequestNew("PoluxCrudUrl", url, "PUT", &respuesta, &transaccion.SolicitudTrabajoGrado); err != nil && status != "200" {
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
			if status, err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil && status != "200" {
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
			if status, err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil && status != "200" {
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
				if status, err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil && status != "200" {
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
				if status, err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil && status != "200" {
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
		if status, err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil && status != "200" {
			panic("Rollback documento escrito " + err.Error())
		}
	}
	rollbackVinculacionTrabajoGrado(transaccion)
	return nil
}

func rollbackDocuEscrito(id int) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK DOCUMENTO ESCRITO")
	var respuesta map[string]interface{}
	url := "/v1/documento_escrito/" + strconv.Itoa(id)
	if status, err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil && status != "200" {
		panic("Rollback documento escrito " + err.Error())
	}

	return nil
}

func rollbackDocumentoTrabajoGrado(transaccion *models.TrRespuestaSolicitud) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK DOCUMENTO TRABAJO GRADO")

	var documentostg []models.DocumentoTrabajoGrado
	var url = "/v1/documento_trabajo_grado?query=TrabajoGrado__Id:" + strconv.Itoa(transaccion.TrTrabajoGrado.TrabajoGrado.Id)
	if err := GetRequestNew("PoluxCrudUrl", url, &documentostg); err != nil {
		logs.Error(err.Error())
		panic(err.Error())
	}

	if len(documentostg) > 0 {
		if documentostg[0].Id != 0 {
			for _, v := range documentostg {
				var respuesta map[string]interface{}
				if transaccion.TrTrabajoGrado.DocumentoTrabajoGrado != nil {
					url := "/v1/documento_trabajo_grado/" + strconv.Itoa(v.Id)
					if status, err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil && status != "200" {
						panic("Rollback documento trabajo grado " + err.Error())
					}
				}
				rollbackDocuEscrito(v.DocumentoEscrito.Id)
			}
			rollbackVinculacionTrabajoGrado(transaccion)
		}
	}

	return nil
}

func rollbackDocumentoEscritoPasantia(transaccion *models.TrRespuestaSolicitud) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK DOCUMENTO ESCRITO PASANTIA")
	var respuesta map[string]interface{}
	if transaccion.DetallesPasantia.Contrato.Id != 0 {
		url := "/v1/documento_escrito/" + strconv.Itoa(transaccion.DetallesPasantia.Contrato.Id)
		if status, err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil && status != "200" {
			panic("Rollback documento escrito " + err.Error())
		}
	}
	if transaccion.DetallesPasantia.Carta.Id != 0 {
		url := "/v1/documento_escrito/" + strconv.Itoa(transaccion.DetallesPasantia.Carta.Id)
		if status, err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil && status != "200" {
			panic("Rollback documento escrito " + err.Error())
		}
	}
	if transaccion.DetallesPasantia.HojaVidaDE.Id != 0 {
		url := "/v1/documento_escrito/" + strconv.Itoa(transaccion.DetallesPasantia.HojaVidaDE.Id)
		if status, err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil && status != "200" {
			panic("Rollback documento escrito " + err.Error())
		}
	}
	rollbackDetallesPasantia(transaccion)
	return nil
}

func rollbackDocumentoTrabajoGradoPasantia(transaccion *models.TrRespuestaSolicitud) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK DOCUMENTO TRABAJO GRADO PASANTIA")
	var respuesta map[string]interface{}
	if transaccion.DetallesPasantia.DTG_Contrato.Id != 0 {
		url := "/v1/documento_trabajo_grado/" + strconv.Itoa(transaccion.DetallesPasantia.DTG_Contrato.Id)
		if status, err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil && status != "200" {
			panic("Rollback documento trabajo grado " + err.Error())
		}
	}
	if transaccion.DetallesPasantia.DTG_Carta.Id != 0 {
		url := "/v1/documento_trabajo_grado/" + strconv.Itoa(transaccion.DetallesPasantia.DTG_Carta.Id)
		if status, err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil && status != "200" {
			panic("Rollback documento trabajo grado " + err.Error())
		}
	}
	if transaccion.DetallesPasantia.DTG_HojaVida.Id != 0 {
		url := "/v1/documento_trabajo_grado/" + strconv.Itoa(transaccion.DetallesPasantia.DTG_HojaVida.Id)
		if status, err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil && status != "200" {
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
				if status, err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil && status != "200" {
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
		if status, err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil && status != "200" {
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
				if status, err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil && status != "200" {
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
		if status, err := SendRequestNew("PoluxCrudUrl", url, "PUT", &respuesta, &v); err != nil && status != "200" {
			panic("Rollback vinculacion trabajo grado rs" + err.Error())
		}
	}
	rollbackDocumentoSolicitudCambioDirEx(transaccion)
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
		if status, err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, &vinculacionNueva); err != nil && status != "200" {
			panic("Rollback vinculacion trabajo grado rs post" + err.Error())
		}
	}
	return nil
}

func rollbackVinculacionTrabajoGradoCan(transaccion *models.TrRespuestaSolicitud, vinculacionesCanceladas []models.VinculacionTrabajoGrado) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK VINCULACION TRABAJO GRADO CAN")
	var respuesta map[string]interface{}
	for _, v := range vinculacionesCanceladas {
		v.Activo = true
		url := "/v1/vinculacion_trabajo_grado/" + strconv.Itoa(v.Id)
		if status, err := SendRequestNew("PoluxCrudUrl", url, "PUT", &respuesta, &v); err != nil && status != "200" {
			panic("Rollback vinculacion trabajo grado can" + err.Error())
		}
	}
	return nil
}

func rollbackRevisionTrabajoGrado(transaccion *models.TrRespuestaSolicitud, revisionAnterior *models.RevisionTrabajoGrado, vincOrig []models.VinculacionTrabajoGrado, vincPost []map[string]interface{}, vinCan []models.VinculacionTrabajoGrado) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK REVISON TRABAJO GRADO")
	var respuesta map[string]interface{}
	url := "/v1/revision_trabajo_grado/" + strconv.Itoa(revisionAnterior.Id)
	if status, err := SendRequestNew("PoluxCrudUrl", url, "PUT", &respuesta, &revisionAnterior); err != nil && status != "200" {
		panic("Rollback revision trabajo grado" + err.Error())
	}

	rollbackVinculacionTrabajoGradoRSPost(transaccion, vincPost)
	rollbackVinculacionTrabajoGradoCan(transaccion, vinCan)
	rollbackVinculacionTrabajoGradoRS(transaccion, vincOrig)
	return nil
}

func rollbackDetallesPasantiaCambioDirector(transaccion *models.TrRespuestaSolicitud, detalle string, revisionAnterior *models.RevisionTrabajoGrado, vincOrig []models.VinculacionTrabajoGrado, vincPost []map[string]interface{}, vincCan []models.VinculacionTrabajoGrado) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK DETALLE PASANTÍA")

	var detallePasantia []models.DetallePasantia
	var url = "/v1/detalle_pasantia?query=TrabajoGrado__Id:" + strconv.Itoa(transaccion.DetallesPasantia.TrabajoGrado.Id) + "&limit=1"
	fmt.Println("URL ", url)
	if err := GetRequestNew("PoluxCrudUrl", url, &detallePasantia); err != nil {
		logs.Error(err.Error())
		panic(err.Error())
	}

	detallePasantia[0].Observaciones = detalle

	var respuesta map[string]interface{}
	url = "/v1/detalle_pasantia/" + strconv.Itoa(detallePasantia[0].Id)
	if status, err := SendRequestNew("PoluxCrudUrl", url, "PUT", &respuesta, &detallePasantia[0]); err != nil && status != "200" {
		panic("Rollback detalle pasantia " + err.Error())
	} else {
		rollbackRevisionTrabajoGrado(transaccion, revisionAnterior, vincOrig, vincPost, vincCan)
	}
	return nil
}

func rollbackDocEsHv(transaccion *models.TrRespuestaSolicitud, detalle string, revisionAnterior *models.RevisionTrabajoGrado, vincOrig []models.VinculacionTrabajoGrado, vincPost []map[string]interface{}, vincCan []models.VinculacionTrabajoGrado) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK DOCUMENTO ESCRITO PARA HOJA DE VIDA")

	var respuesta map[string]interface{}
	url := "/v1/documento_escrito/" + strconv.Itoa(transaccion.DetallesPasantia.HojaVidaDE.Id)
	if status, err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, &transaccion.DetallesPasantia.HojaVidaDE); err != nil && status != "200" {
		panic("Rollback documento escrito hoja de vida " + err.Error())
	} else {
		rollbackDetallesPasantiaCambioDirector(transaccion, detalle, revisionAnterior, vincOrig, vincPost, vincCan)
	}
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
	if status, err := SendRequestNew("PoluxCrudUrl", url, "PUT", &respuesta, &transaccion.EstudianteTrabajoGrado); err != nil && status != "200" {
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
		if status, err := SendRequestNew("PoluxCrudUrl", url, "PUT", &respuesta, v); err != nil && status != "200" {
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
	if status, err := SendRequestNew("PoluxCrudUrl", url, "PUT", &respuesta, transaccion.EstudianteTrabajoGrado.TrabajoGrado); err != nil && status != "200" {
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
		if status, err := SendRequestNew("PoluxCrudUrl", url, "PUT", &respuesta, v); err != nil && status != "200" {
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
		if status, err := SendRequestNew("PoluxCrudUrl", url, "PUT", &respuesta, v); err != nil && status != "200" {
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
	if status, err := SendRequestNew("PoluxCrudUrl", url, "PUT", &respuesta, transaccion.TrRevision.TrabajoGrado); err != nil && status != "200" {
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
			if status, err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil && status != "200" {
				panic("Rollback detalle trabajo grado revision" + err.Error())
			}
		}
	}
	rollbackTrGrRev(transaccion)
	return nil
}

func rollbackDocEscrRev(transaccion *models.TrRespuestaSolicitud) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK DOCUMENTO ESCRITO REVISION")
	// Iterar sobre el array de documentos escritos y hacer rollback para cada uno
	for _, documentoEscrito := range *transaccion.TrRevision.DocumentoEscrito {
		var respuesta map[string]interface{}
		url := "/v1/documento_escrito/" + strconv.Itoa(documentoEscrito.Id)
		if status, err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil && status != "200" {
			panic("Rollback documento escrito revision: " + err.Error())
		}
	}

	// Llamar al siguiente rollback una vez que se hayan eliminado todos los documentos
	rollbackDetTrGrRev(transaccion)
	return nil
}

func rollbackDocTrGrRev(transaccion *models.TrRespuestaSolicitud) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK DOCUMENTO TRABAJO GRADO REVISION")
	var respuesta map[string]interface{}
	url := "/v1/documento_trabajo_grado/" + strconv.Itoa(transaccion.TrRevision.DocumentoTrabajoGrado.Id)
	if status, err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil && status != "200" {
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
		if status, err := SendRequestNew("PoluxCrudUrl", url, "PUT", &respuesta, &v); err != nil && status != "200" {
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
		if status, err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, &vinculacionNueva); err != nil && status != "200" {
			panic("Rollback vinculacion trabajo grado post revision" + err.Error())
		}
	}
	return nil
}
