package helpers

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/udistrital/polux_mid/models"
)

func AddTransaccionSolicitud(transaccion *models.TrSolicitud) (response map[string]interface{}, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "AddTransaccionSolicitud", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	url := "/v1/solicitud_trabajo_grado"
	var resSolicitudTrabajoGrado map[string]interface{}
	if err := SendRequestNew("PoluxCrudUrl", url, "POST", &resSolicitudTrabajoGrado, &transaccion.Solicitud); err == nil {
		var idSolicitudTrabajoGrado = int(resSolicitudTrabajoGrado["Id"].(float64))
		transaccion.Respuesta.SolicitudTrabajoGrado.Id = idSolicitudTrabajoGrado
		transaccion.Solicitud.Id = idSolicitudTrabajoGrado

		url = "/v1/respuesta_solicitud"
		var resRespuestaSolicitud map[string]interface{}
		if err := SendRequestNew("PoluxCrudUrl", url, "POST", &resRespuestaSolicitud, &transaccion.Respuesta); err == nil {
			transaccion.Respuesta.Id = int(resRespuestaSolicitud["Id"].(float64))
		} else {
			fmt.Println("ENTRA A ROLLBACKSOLICITUDTABAJOGRADO", transaccion)
			rollbackSolicitudTrabajoGradoSol(transaccion)
			//return nil, fmt.Errorf("Error en Respuesta Solicitud: %v", err)
			//return response, outputError
		}

		url = "/v1/detalle_solicitud"
		var detalleSolicitud = make([]map[string]interface{}, 0)
		for i, v := range *transaccion.DetallesSolicitud {
			var resDetalleSolicitudSol map[string]interface{}
			v.SolicitudTrabajoGrado.Id = idSolicitudTrabajoGrado
			if err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDetalleSolicitudSol, &v); err == nil {
				(*transaccion.DetallesSolicitud)[i].Id = int(resDetalleSolicitudSol["Id"].(float64))
				detalleSolicitud = append(detalleSolicitud, resDetalleSolicitudSol)
			} else {
				logs.Error(err)
				if len(detalleSolicitud) > 0 {
					rollbackDetalleSolicitudSol(transaccion)
				}
				rollbackRespuestaSolicitudSol(transaccion)
				//return nil, fmt.Errorf("Error en Detalle Solicitud: %v", err)
				//return response, outputError
			}
		}
		url = "/v1/usuario_solicitud"
		var usuarioSolicitud = make([]map[string]interface{}, 0)
		for i, v := range *transaccion.UsuariosSolicitud {
			var resUsuarioSolicitud map[string]interface{}
			v.SolicitudTrabajoGrado.Id = idSolicitudTrabajoGrado
			if err := SendRequestNew("PoluxCrudUrl", url, "POST", &resUsuarioSolicitud, &v); err == nil {
				(*transaccion.UsuariosSolicitud)[i].Id = int(resUsuarioSolicitud["Id"].(float64))
				usuarioSolicitud = append(usuarioSolicitud, resUsuarioSolicitud)
			} else {
				logs.Error(err)
				if len(detalleSolicitud) > 0 {
					rollbackUsuarioSolicitudSol(transaccion)
				}
				rollbackDetalleSolicitudSol(transaccion)
				//return nil, fmt.Errorf("Error en Usuario Solicitud: %v", err)
				//return response, outputError
			}
		}

		response = map[string]interface{}{
			"SolicitudTrabajoGrado": resSolicitudTrabajoGrado,
			"RespuestaSolicitud":    resRespuestaSolicitud,
			"DetalleSolicitud":      detalleSolicitud,
			"UsuarioSolicitud":      usuarioSolicitud,
		}
		return response, outputError
	} else {
		//logs.Error(err)
		//return nil, fmt.Errorf("Error en Solicitud Trabajo Grado: %v", err)
		//localError := map[string]interface{}{"Success": false, "Status": 401, "Message": "Error en los datos de solicitud_trabajo_grado ", "Data": response}
		panic("Error en los datos de solicitud_trabajo_grado " + err.Error())
		//return response, outputError
	}
}

func rollbackSolicitudTrabajoGradoSol(transaccion *models.TrSolicitud) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK SOLICITUD TRABAJO GRADO SOL")
	var respuesta map[string]interface{}
	url := "/v1/respuesta_solicitud/" + strconv.Itoa(transaccion.Solicitud.Id)
	if err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil {
		panic("Rollback solicitud trabajo grado" + err.Error())
	}
	return nil
}

func rollbackRespuestaSolicitudSol(transaccion *models.TrSolicitud) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK RESPUESTA SOLICITUD SOL")
	var respuesta map[string]interface{}
	url := "/v1/respuesta_solicitud/" + strconv.Itoa(transaccion.Respuesta.Id)
	if err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil {
		panic("Rollback respuesta solicitud" + err.Error())
	}
	rollbackSolicitudTrabajoGradoSol(transaccion)
	return nil
}

func rollbackDetalleSolicitudSol(transaccion *models.TrSolicitud) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK DETALLE SOLICITUD SOL")
	var respuesta map[string]interface{}
	if transaccion.DetallesSolicitud != nil {
		for _, v := range *transaccion.DetallesSolicitud {
			if v.Id != 0 {
				url := "/v1/detalle_solicitud/" + strconv.Itoa(v.Id)
				if err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil {
					panic("Rollback detalle solicitud " + err.Error())
				}
			}
		}
	}
	rollbackRespuestaSolicitudSol(transaccion)
	return nil
}

func rollbackUsuarioSolicitudSol(transaccion *models.TrSolicitud) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK USUARIO SOLICITUD SOL")
	var respuesta map[string]interface{}
	if transaccion.UsuariosSolicitud != nil {
		for _, v := range *transaccion.UsuariosSolicitud {
			if v.Id != 0 {
				url := "/v1/usuario_solicitud/" + strconv.Itoa(v.Id)
				if err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil {
					panic("Rollback usuario solicitud " + err.Error())
				}
			}
		}
	}
	rollbackDetalleSolicitudSol(transaccion)
	return nil
}
