package helpers

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/udistrital/polux_mid/models"
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
	if err := GetRequestNew("UrlCrudParametros", url, &parametro); err != nil {
		logs.Error(err.Error())
		panic(err.Error())
	}
	url = "parametro/" + strconv.Itoa(transaccion.RespuestaAnterior.EstadoSolicitud)
	if err := GetRequestNew("UrlCrudParametros", url, &parametroEstadoSolicitud); err != nil {
		logs.Error(err.Error())
		panic(err.Error())
	}
	fmt.Println("PARAMETRO ", parametro)
	var resRespuestaAnterior map[string]interface{}
	url = "/v1/respuesta_solicitud/" + strconv.Itoa(transaccion.RespuestaAnterior.Id)
	//payload := "/" + strconv.Itoa(transaccion.RespuestaAnterior.Id)
	//if outputError = Post(new(models.RespuestaSolicitud), payload, &resRespuestaAnterior); outputError == nil {
	if err := SendRequestNew("PoluxCrudUrl", url, "PUT", &resRespuestaAnterior, &transaccion.RespuestaAnterior); err == nil {
		fmt.Println("RES ANTERIOR ", resRespuestaAnterior)
		url = "/v1/respuesta_solicitud"
		var resRespuestaNueva map[string]interface{}
		if err := SendRequestNew("PoluxCrudUrl", url, "POST", &resRespuestaNueva, &transaccion.RespuestaNueva); err == nil {
			fmt.Println("RESPUESTA NUEVA ", resRespuestaNueva)
			transaccion.RespuestaNueva.Id = int(resRespuestaNueva["Id"].(float64))
			url = "/v1/documento_solicitud"
			var resDocumentoSolicitud map[string]interface{}
			if err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoSolicitud, &transaccion.DocumentoSolicitud); err == nil {
				fmt.Println("RESPUESTA DOCUMENTO ", resDocumentoSolicitud)
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
			fmt.Println("TRABAJO GRADO ", resTrabajoGrado)
			var idTrabajoGrado = int(resTrabajoGrado["Id"].(float64))
			transaccion.TrTrabajoGrado.TrabajoGrado.Id = idTrabajoGrado
			transaccion.SolicitudTrabajoGrado.TrabajoGrado.Id = idTrabajoGrado
			var resSolicitudTrabajoGrado map[string]interface{}
			url = "/v1/solicitud_trabajo_grado/" + strconv.Itoa(transaccion.SolicitudTrabajoGrado.Id)
			if err := SendRequestNew("PoluxCrudUrl", url, "PUT", &resSolicitudTrabajoGrado, &transaccion.SolicitudTrabajoGrado); err == nil {
				fmt.Println("SOLICITUD TRABAJO GRADO ", resSolicitudTrabajoGrado)
				url = "/v1/asignatura_trabajo_grado"
				var materias = make([]map[string]interface{}, 0)
				for i, v := range *transaccion.TrTrabajoGrado.AsignaturasTrabajoGrado {
					var resAsignaturaTrabajoGrado map[string]interface{}
					v.TrabajoGrado.Id = idTrabajoGrado
					if err := SendRequestNew("PoluxCrudUrl", url, "POST", &resAsignaturaTrabajoGrado, &v); err == nil {
						fmt.Println("ASIGNATURA TRABAJO GRADO ", resAsignaturaTrabajoGrado)
						(*transaccion.TrTrabajoGrado.AsignaturasTrabajoGrado)[i].Id = int(resAsignaturaTrabajoGrado["Id"].(float64))
						materias = append(materias, resAsignaturaTrabajoGrado)
						fmt.Println(transaccion.TrTrabajoGrado.AsignaturasTrabajoGrado)
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
						fmt.Println("ESTUDIANTE TRABAJO GRADO ", resEstudianteTrabajoGrado)
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
							fmt.Println("AREAS TRABAJO GRADO ", resAreasTrabajoGrado)
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
							fmt.Println("VINCULACION TRABAJO GRADO ", resVinculacionTrabajoGrado)
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
					url = "/v1/documento_escrito"
					var resDocumentoEscrito map[string]interface{}
					if err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoEscrito, &transaccion.TrTrabajoGrado.DocumentoEscrito); err == nil {
						fmt.Println("DOCUMENTO ESCRITO ", resDocumentoEscrito)
						transaccion.TrTrabajoGrado.DocumentoTrabajoGrado.TrabajoGrado.Id = idTrabajoGrado
						transaccion.TrTrabajoGrado.DocumentoTrabajoGrado.DocumentoEscrito.Id = int(resDocumentoEscrito["Id"].(float64))
						transaccion.TrTrabajoGrado.DocumentoEscrito.Id = int(resDocumentoEscrito["Id"].(float64))
						var resDocumentoTrabajoGrado map[string]interface{}
						url = "/v1/documento_trabajo_grado"
						if err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoTrabajoGrado, &transaccion.TrTrabajoGrado.DocumentoTrabajoGrado); err == nil {
							fmt.Println("DOCUMENTO TRABAJO GRADO ", resDocumentoTrabajoGrado)
							transaccion.TrTrabajoGrado.DocumentoTrabajoGrado.Id = int(resDocumentoTrabajoGrado["Id"].(float64))
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
							fmt.Println("ESPACIO ACADEMICO INSCRITO ", resEspaciosAcademicosInscritos)
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
					url = "/v1/detalle_pasantia"
					var resDetallePasantia map[string]interface{}
					if err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDetallePasantia, &transaccion.DetallesPasantia); err == nil {
						fmt.Println("DETALLE PASANTIA ", resDetallePasantia)
						transaccion.DetallesPasantia.Id = int(resDetallePasantia["Id"].(float64))
					} else {
						logs.Error(err)
						rollbackDetallesPasantia(transaccion)
					}
				}
				if transaccion.DetallesPasantiaExterna != nil && parametro.CodigoAbreviacion == "PASEX_PLX" {
					url = "/v1/detalle_trabajo_grado"
					var detalles_pasantia_externa = make([]map[string]interface{}, 0)
					for i, v := range *transaccion.DetallesPasantiaExterna {
						var resDetallesPasantiaExterna map[string]interface{}
						v.TrabajoGrado.Id = idTrabajoGrado
						if err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDetallesPasantiaExterna, &v); err == nil {
							fmt.Println("DETALLE PASANTIA EXTERNA ", resDetallesPasantiaExterna)
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
		fmt.Println("V ", v)
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
	if transaccion.DetallesPasantia != nil {
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
