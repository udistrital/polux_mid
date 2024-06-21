package helpers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/polux_mid/models"

	//request "github.com/udistrital/utils_oas/blob/master/request"
	request "github.com/udistrital/utils_oas/request"
)

func AddTransaccionRegistrarRevisionTg(transaccion *models.TrRegistrarRevisionTg) (alerta []string, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "AddTransaccionRegistrarRevisionTg", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
	alerta = append(alerta, "Success")
	var correccion models.Correccion

	rollbackPost := false
	if transaccion.RevisionTrabajoGrado.Id == 0 {
		rollbackPost = true
		url := "parametro?query=CodigoAbreviacion:FINALIZADA_PLX,TipoParametroId__CodigoAbreviacion:ESTREV_TRG"
		var parametroEstadoRevision []models.Parametro
		if err := request.GetRequestNew("UrlCrudParametros", url, &parametroEstadoRevision); err != nil {
			logs.Error(err.Error())
			panic(err.Error())
		}
		var parametroEstadoTrabajoGrado []models.Parametro
		url = "parametro?query=CodigoAbreviacion:EC_PLX,TipoParametroId__CodigoAbreviacion:EST_TRG"
		if err := request.GetRequestNew("UrlCrudParametros", url, &parametroEstadoTrabajoGrado); err != nil {
			logs.Error(err.Error())
			panic(err.Error())
		}
		transaccion.RevisionTrabajoGrado.EstadoRevisionTrabajoGrado = parametroEstadoRevision[0].Id
		transaccion.RevisionTrabajoGrado.VinculacionTrabajoGrado.TrabajoGrado.EstadoTrabajoGrado = parametroEstadoTrabajoGrado[0].Id

		var revisiones []models.RevisionTrabajoGrado
		url = beego.AppConfig.String("PoluxCrudUrl") + "/v1/revision_trabajo_grado?query=DocumentoTrabajoGrado__TrabajoGrado__Id:" + strconv.Itoa(transaccion.RevisionTrabajoGrado.DocumentoTrabajoGrado.TrabajoGrado.Id)
		if err := request.GetJson(url, &revisiones); err != nil {
			logs.Error(err.Error())
			panic(err.Error())
		}

		transaccion.RevisionTrabajoGrado.NumeroRevision = len(revisiones) + 1
		transaccion.RevisionTrabajoGrado.FechaRecepcion = time.Now()
		transaccion.RevisionTrabajoGrado.FechaRevision = &transaccion.RevisionTrabajoGrado.FechaRecepcion

		url = "/v1/revision_trabajo_grado"
		var resRevisionTrabajoGrado map[string]interface{}
		if err := request.SendRequestNew("PoluxCrudUrl", url, "POST", &resRevisionTrabajoGrado, &transaccion.RevisionTrabajoGrado); err == nil {
			transaccion.RevisionTrabajoGrado.Id = int(resRevisionTrabajoGrado["Id"].(float64))
		} else {
			logs.Error(err)
			panic(err.Error())
		}

		correccion.RevisionTrabajoGrado = transaccion.RevisionTrabajoGrado
		correccion.Documento = false

		url = "/v1/correccion"
		var resCorreccion map[string]interface{}
		if err := request.SendRequestNew("PoluxCrudUrl", url, "POST", &resCorreccion, &correccion); err == nil {
			correccion.Id = int(resCorreccion["Id"].(float64))
		} else {
			logs.Error(err)
			rollbackPostRevisionTrabajoGrado(transaccion)
		}

		url = "/v1/comentario"
		var comentarios = make([]map[string]interface{}, 0)
		for i, v := range *transaccion.Comentarios {
			var resComentario map[string]interface{}
			v.Correccion = &correccion
			v.Fecha = time.Now()
			if err := request.SendRequestNew("PoluxCrudUrl", url, "POST", &resComentario, &v); err == nil {
				(*transaccion.Comentarios)[i].Id = int(resComentario["Id"].(float64))
				comentarios = append(comentarios, resComentario)
			} else {
				logs.Error(err)
				if len(comentarios) > 0 {
					rollbackPostComentario(transaccion, &correccion)
				} else {
					rollbackPostCorreccion(transaccion, &correccion)
				}
			}
		}

		url = "/v1/trabajo_grado/" + strconv.Itoa(transaccion.RevisionTrabajoGrado.VinculacionTrabajoGrado.TrabajoGrado.Id)
		var resTrabajoGrado map[string]interface{}
		if err := request.SendRequestNew("PoluxCrudUrl", url, "PUT", &resTrabajoGrado, &transaccion.RevisionTrabajoGrado.VinculacionTrabajoGrado.TrabajoGrado); err != nil {
			logs.Error(err)
			panic(err.Error())
		}

		return alerta, nil
	}

	url := "/v1/revision_trabajo_grado/" + strconv.Itoa(transaccion.RevisionTrabajoGrado.Id)
	var resRevisionTrabajoGrado map[string]interface{}
	if err := request.SendRequestNew("PoluxCrudUrl", url, "PUT", &resRevisionTrabajoGrado, &transaccion.RevisionTrabajoGrado); err != nil {
		logs.Error(err)
		if rollbackPost {
			rollbackPostComentario(transaccion, &correccion)
		} else {
			panic(err.Error())
		}
	}

	var correcciones = make([]map[string]interface{}, 0)
	var comentarios = make([]map[string]interface{}, 0)
	for _, v := range *transaccion.Comentarios {
		url = "/v1/correccion"
		var resCorreccion map[string]interface{}
		if err := request.SendRequestNew("PoluxCrudUrl", url, "POST", &resCorreccion, &v.Correccion); err == nil {
			v.Correccion.Id = int(resCorreccion["Id"].(float64))
			correcciones = append(correcciones, resCorreccion)
		} else {
			logs.Error(err)
			if len(correcciones) > 0 {
				rollbackPostCorreccion2(transaccion)
			}
			if len(comentarios) > 0 {
				rollbackPostComentario2(transaccion)
			}
			rollbackPutRevisionTrabajoGrado(transaccion, &correccion, rollbackPost)
		}
		url = "/v1/comentario"
		var resComentario map[string]interface{}
		if err := request.SendRequestNew("PoluxCrudUrl", url, "POST", &resComentario, &v); err == nil {
			v.Id = int(resComentario["Id"].(float64))
			comentarios = append(comentarios, resComentario)
		} else {
			logs.Error(err)
			if len(correcciones) > 0 {
				rollbackPostCorreccion2(transaccion)
			}
			if len(comentarios) > 0 {
				rollbackPostComentario2(transaccion)
			}
			rollbackPutRevisionTrabajoGrado(transaccion, &correccion, rollbackPost)
		}

	}

	return alerta, outputError
}

func rollbackPostRevisionTrabajoGrado(transaccion *models.TrRegistrarRevisionTg) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK POST REVISION TRABAJO GRADO ")
	var respuesta map[string]interface{}
	url := "/v1/revision_trabajo_grado/" + strconv.Itoa(transaccion.RevisionTrabajoGrado.Id)
	if err := request.SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, &transaccion.RevisionTrabajoGrado); err != nil {
		panic("Rollback post revision trabajo grado " + err.Error())
	}
	return nil
}

func rollbackPostComentario(transaccion *models.TrRegistrarRevisionTg, correccion *models.Correccion) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK POST COMENTARIO")
	var respuesta map[string]interface{}
	for _, v := range *transaccion.Comentarios {
		fmt.Println("V ", v)
		if v.Id != 0 {
			url := "/v1/comentario/" + strconv.Itoa(v.Id)
			if err := request.SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil {
				panic("Rollback post comentario" + err.Error())
			}
		}
	}
	rollbackPostCorreccion(transaccion, correccion)
	return nil
}

func rollbackPostCorreccion(transaccion *models.TrRegistrarRevisionTg, correccion *models.Correccion) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK POST CORRECCION")
	var respuesta map[string]interface{}
	url := "/v1/correccion/" + strconv.Itoa(correccion.Id)
	if err := request.SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, &correccion); err != nil {
		panic("Rollback post revision trabajo grado " + err.Error())
	}
	rollbackPostRevisionTrabajoGrado(transaccion)
	return nil
}

func rollbackPutRevisionTrabajoGrado(transaccion *models.TrRegistrarRevisionTg, correccion *models.Correccion, rollbackPost bool) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK PUT REVISION TRABAJO GRADO")
	var respuesta map[string]interface{}

	url := "parametro?query=CodigoAbreviacion:PENDIENTE_PLX,TipoParametroId__CodigoAbreviacion:ESTREV_TRG"
	var parametroEstadoRevision []models.Parametro
	if err := request.GetRequestNew("UrlCrudParametros", url, &parametroEstadoRevision); err != nil {
		logs.Error(err.Error())
		panic(err.Error())
	}

	transaccion.RevisionTrabajoGrado.EstadoRevisionTrabajoGrado = parametroEstadoRevision[0].Id
	transaccion.RevisionTrabajoGrado.FechaRevision = nil
	url = "/v1/revision_trabajo_grado/" + strconv.Itoa(transaccion.RevisionTrabajoGrado.Id)
	if err := request.SendRequestNew("PoluxCrudUrl", url, "PUT", &respuesta, &transaccion.RevisionTrabajoGrado); err != nil {
		panic("Rollback put revision trabajo grado" + err.Error())
	} else if rollbackPost {
		rollbackPostComentario(transaccion, correccion)
	}
	return nil
}

func rollbackPostCorreccion2(transaccion *models.TrRegistrarRevisionTg) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK POST CORRECCION 2")
	var respuesta map[string]interface{}
	for _, v := range *transaccion.Comentarios {
		url := "/v1/correccion/" + strconv.Itoa(v.Correccion.Id)
		if err := request.SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, &v.Correccion); err != nil {
			panic("Rollback corrección 2 " + err.Error())
		}
	}
	return nil
}

func rollbackPostComentario2(transaccion *models.TrRegistrarRevisionTg) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK POST COMENTARIO 2")
	var respuesta map[string]interface{}
	for _, v := range *transaccion.Comentarios {
		url := "/v1/comentario/" + strconv.Itoa(v.Id)
		if err := request.SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, &v); err != nil {
			panic("Rollback comentario 2 " + err.Error())
		}
	}
	return nil
}
