package helpers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/udistrital/polux_mid/models"
	"github.com/udistrital/utils_oas/time_bogota"
	//request "github.com/udistrital/utils_oas/blob/master/request"
)

func AddTransaccionRegistrarRevisionTg(transaccion *models.TrRegistrarRevisionTg) (response map[string]interface{}, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("ERROR ", err)
			panic(DeferHelpers("AddTransaccionSolicitud", err))
		}
	}()
	//alerta = append(alerta, "Success")
	var correccion models.Correccion

	rollbackPost := false
	if transaccion.RevisionTrabajoGrado.Id == 0 {
		rollbackPost = true
		url := "parametro?query=CodigoAbreviacion:FINALIZADA_PLX,TipoParametroId__CodigoAbreviacion:ESTREV_TRG"
		var parametroEstadoRevision []models.Parametro
		fmt.Println("UrlCrudParametros", url)
		if err := GetRequestNew("UrlCrudParametros", url, &parametroEstadoRevision); err != nil {
			fmt.Println("ENTRA A ERROR PARAMETROS")
			//logs.Error(err.Error())
			panic(err.Error())
		}
		var parametroEstadoTrabajoGrado []models.Parametro
		url = "parametro?query=CodigoAbreviacion:EC_PLX,TipoParametroId__CodigoAbreviacion:EST_TRG"
		if err := GetRequestNew("UrlCrudParametros", url, &parametroEstadoTrabajoGrado); err != nil {
			//logs.Error(err.Error())
			panic(err.Error())
		}
		transaccion.RevisionTrabajoGrado.EstadoRevisionTrabajoGrado = parametroEstadoRevision[0].Id
		transaccion.RevisionTrabajoGrado.VinculacionTrabajoGrado.TrabajoGrado.EstadoTrabajoGrado = parametroEstadoTrabajoGrado[0].Id

		var revisiones []models.RevisionTrabajoGrado
		url = "/v1/revision_trabajo_grado?query=DocumentoTrabajoGrado__TrabajoGrado__Id:" + strconv.Itoa(transaccion.RevisionTrabajoGrado.DocumentoTrabajoGrado.TrabajoGrado.Id)
		if err := GetRequestNew("PoluxCrudUrl", url, &revisiones); err != nil {
			//logs.Error(err.Error())
			panic(err.Error())
		}

		transaccion.RevisionTrabajoGrado.NumeroRevision = len(revisiones) + 1

		url = "/v1/revision_trabajo_grado"
		var resRevisionTrabajoGrado map[string]interface{}
		if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resRevisionTrabajoGrado, &transaccion.RevisionTrabajoGrado); err == nil && status == "201" {
			transaccion.RevisionTrabajoGrado.Id = int(resRevisionTrabajoGrado["Id"].(float64))
		} else {
			//logs.Error(err)
			panic(err.Error())
		}

		correccion.RevisionTrabajoGrado = transaccion.RevisionTrabajoGrado
		correccion.Documento = false

		url = "/v1/correccion"
		var resCorreccion map[string]interface{}
		if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resCorreccion, &correccion); err == nil && status == "201" {
			fmt.Println("Entra a cORRECCIÓN!", correccion)
			correccion.Id = int(resCorreccion["Id"].(float64))
		} else {
			logs.Error(err)
			rollbackPostRevisionTrabajoGrado(transaccion)
		}

		url = "/v1/comentario"
		var comentarios = make([]map[string]interface{}, 0)
		fmt.Println("TRANSACCIÓN COMENTARIOS", transaccion.Comentarios)
		for i, v := range *transaccion.Comentarios {
			var resComentario map[string]interface{}
			v.Correccion = &correccion
			v.Fecha = time.Now()
			fmt.Println("pASA pOR aCÁ!!!")
			if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resComentario, &v); err == nil && status == "201" {
				fmt.Println("Entra a Comentario!", v)
				fmt.Println("Responde Comentario!", resComentario)
				(*transaccion.Comentarios)[i].Id = int(resComentario["Id"].(float64))
				comentarios = append(comentarios, resComentario)
			} else {
				fmt.Println("Entra a ELSE!", resComentario)
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
		if status, err := SendRequestNew("PoluxCrudUrl", url, "PUT", &resTrabajoGrado, &transaccion.RevisionTrabajoGrado.VinculacionTrabajoGrado.TrabajoGrado); err != nil && status != "200" {
			//logs.Error(err)
			panic(err.Error())
		}

		response = map[string]interface{}{
			"RevisionTrabajoGrado": resRevisionTrabajoGrado,
			"Correccion":           resCorreccion,
			"comentarios":          comentarios,
		}
		return response, nil
	}

	url := "/v1/revision_trabajo_grado/" + strconv.Itoa(transaccion.RevisionTrabajoGrado.Id)
	var resRevisionTrabajoGrado map[string]interface{}
	transaccion.RevisionTrabajoGrado.FechaRecepcion = strings.Replace(transaccion.RevisionTrabajoGrado.FechaRecepcion, " +0000 +0000", " +0000", 1)
	transaccion.RevisionTrabajoGrado.FechaRevision = time_bogota.TiempoBogotaFormato()
	if status, err := SendRequestNew("PoluxCrudUrl", url, "PUT", &resRevisionTrabajoGrado, &transaccion.RevisionTrabajoGrado); err != nil && status != "200" {
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
		if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resCorreccion, &v.Correccion); err == nil && status == "201" {
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
		if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resComentario, &v); err == nil && status == "201" {
			fmt.Println("Entra a Comentario!", v)
			fmt.Println("Responde Comentario!", resComentario["Id"])
			v.Id = int(resComentario["Id"].(float64))
			fmt.Println("FORMATO COMENTRAIO!", resComentario)
			comentarios = append(comentarios, resComentario)
			fmt.Println("APPEND COMENTARIOS!", comentarios)
			fmt.Println("ERROR!", err)
		} else {
			fmt.Println("Entra a ELSE!", resComentario)
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

	response = map[string]interface{}{
		"RevisionTrabajoGrado": resRevisionTrabajoGrado,
		"Correccion":           correcciones,
		"comentarios":          comentarios,
	}
	return response, nil
}

func rollbackPostRevisionTrabajoGrado(transaccion *models.TrRegistrarRevisionTg) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK POST REVISION TRABAJO GRADO ")
	var respuesta map[string]interface{}
	url := "/v1/revision_trabajo_grado/" + strconv.Itoa(transaccion.RevisionTrabajoGrado.Id)
	if status, err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, &transaccion.RevisionTrabajoGrado); err != nil && status != "200" {
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
			if status, err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil && status != "200" {
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
	if status, err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, &correccion); err != nil && status != "200" {
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
	if err := GetRequestNew("UrlCrudParametros", url, &parametroEstadoRevision); err != nil {
		logs.Error(err.Error())
		panic(err.Error())
	}

	transaccion.RevisionTrabajoGrado.EstadoRevisionTrabajoGrado = parametroEstadoRevision[0].Id
	transaccion.RevisionTrabajoGrado.FechaRevision = ""
	url = "/v1/revision_trabajo_grado/" + strconv.Itoa(transaccion.RevisionTrabajoGrado.Id)
	if status, err := SendRequestNew("PoluxCrudUrl", url, "PUT", &respuesta, &transaccion.RevisionTrabajoGrado); err != nil && status != "200" {
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
		if status, err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, &v.Correccion); err != nil && status != "200" {
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
		if status, err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, &v); err != nil && status != "200" {
			panic("Rollback comentario 2 " + err.Error())
		}
	}
	return nil
}
