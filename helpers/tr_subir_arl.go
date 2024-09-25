package helpers

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/polux_mid/models"
)

func AddTransaccionSubirArl(transaccion *models.TrSubirArl) (alerta []string, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("ERROR ", err)
			panic(DeferHelpers("AddTransaccionSolicitud", err))
		}
	}()
	alerta = append(alerta, "Success")

	transaccion.DocumentoEscrito.Id = 0
	transaccion.DocumentoEscrito.Titulo = "ARL de la pasantía con titulo: " + transaccion.DocumentoEscrito.Titulo
	transaccion.DocumentoEscrito.Resumen = transaccion.DocumentoEscrito.Titulo

	url := "/v1/documento_escrito"
	var resDocumentoEscrito map[string]interface{}
	if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoEscrito, &transaccion.DocumentoEscrito); err == nil && status == "201" { //Se guarda la ARL en Documento Escrito. Desde el Cliente ya viene con el tipo "Documentos adicionales asociados a las pasantias"

		var estadoRechazado []models.Parametro
		url = "parametro?query=CodigoAbreviacion:ARC_PLX"
		if err := GetRequestNew("UrlCrudParametros", url, &estadoRechazado); err != nil { //Se busca el Estado Trabajo de Grado "ARL Rechazada"
			logs.Error(err.Error())
			panic(err.Error())
		}

		if transaccion.TrabajoGrado.EstadoTrabajoGrado == estadoRechazado[0].Id { //Si el estado del trabajo de grado es "ARL Rechazado"

			var tipoDocumento []models.TipoDocumento
			url = beego.AppConfig.String("UrlDocumentos") + "tipo_documento?query=CodigoAbreviacion:DPAS_PLX"
			if err := GetJson(url, &tipoDocumento); err != nil { //Se busca el Tipo de Documento asociado a la ARL
				logs.Error(err.Error())
				panic(err.Error())
			}

			var documentosTG []models.DocumentoTrabajoGrado
			url = "/v1/documento_trabajo_grado?query=trabajo_grado__Id:" + strconv.Itoa(transaccion.TrabajoGrado.Id) + ",documento_escrito__tipo_documento_escrito:" + strconv.Itoa(tipoDocumento[0].Id)
			if err := GetRequestNew("PoluxCrudUrl", url, &documentosTG); err != nil { //Se busca el registro en documento_trabajo_grado en la que relacione la ARL con el trabajo de grado
				logs.Error(err.Error())
				panic(err.Error())
			}

			transaccion.DocumentoTrabajoGrado.Id = documentosTG[0].Id
			transaccion.DocumentoEscrito.Id = int(resDocumentoEscrito["Id"].(float64))
			transaccion.DocumentoTrabajoGrado.DocumentoEscrito.Id = int(resDocumentoEscrito["Id"].(float64))

			url := "/v1/documento_trabajo_grado/" + strconv.Itoa(documentosTG[0].Id)
			var resDocumentoTrabajoGrado map[string]interface{}
			if status, err := SendRequestNew("PoluxCrudUrl", url, "PUT", &resDocumentoTrabajoGrado, &transaccion.DocumentoTrabajoGrado); err == nil && status == "200" { //Se actualiza la relación entre la ARL y el Trabajo de Grado

				transaccion.DocumentoTrabajoGrado.Id = int(resDocumentoTrabajoGrado["Id"].(float64))

				var parametro []models.Parametro
				url = "parametro?query=CodigoAbreviacion:ACEA_PLX"
				if err := GetRequestNew("UrlCrudParametros", url, &parametro); err != nil { //Se busca el Estado Trabajo de Grado "ARL Cargada, en espera de aprobación"
					logs.Error(err.Error())
					panic(err.Error())
				}

				transaccion.TrabajoGrado.EstadoTrabajoGrado = parametro[0].Id

				url := "/v1/trabajo_grado/" + strconv.Itoa(transaccion.TrabajoGrado.Id)
				var resTrabajoGrado map[string]interface{}
				if status, err := SendRequestNew("PoluxCrudUrl", url, "PUT", &resTrabajoGrado, &transaccion.TrabajoGrado); err != nil && status != "200" { //Se actualiza el Trabajo de Grado con el nuevo Estado
					rollbackDocumentoTrabajoGradoARL(transaccion)
					//logs.Error(err)
					panic(err.Error())
				}

			} else {
				rollbackDocumentoEscritoARL(transaccion)
				//logs.Error(err)
				panic(err.Error())
			}

		} else { //Si es la primera ARL que sube el estudiante

			transaccion.DocumentoTrabajoGrado.Id = 0
			transaccion.DocumentoEscrito.Id = int(resDocumentoEscrito["Id"].(float64))
			transaccion.DocumentoTrabajoGrado.DocumentoEscrito.Id = int(resDocumentoEscrito["Id"].(float64))

			url := "/v1/documento_trabajo_grado"
			var resDocumentoTrabajoGrado map[string]interface{}
			if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoTrabajoGrado, &transaccion.DocumentoTrabajoGrado); err == nil && status == "201" { //Se guarda la relación entre Documento Escrito y el Trabajo de Grado

				transaccion.DocumentoTrabajoGrado.Id = int(resDocumentoTrabajoGrado["Id"].(float64))

				var parametro []models.Parametro
				url = "parametro?query=CodigoAbreviacion:ACEA_PLX"
				if err := GetRequestNew("UrlCrudParametros", url, &parametro); err != nil { //Se busca el Estado Trabajo de Grado "ARL Cargada, en espera de aprobación"
					logs.Error(err.Error())
					panic(err.Error())
				}

				transaccion.TrabajoGrado.EstadoTrabajoGrado = parametro[0].Id

				url := "/v1/trabajo_grado/" + strconv.Itoa(transaccion.TrabajoGrado.Id)
				var resTrabajoGrado map[string]interface{}
				if status, err := SendRequestNew("PoluxCrudUrl", url, "PUT", &resTrabajoGrado, &transaccion.TrabajoGrado); err != nil && status != "200" { //Se actualiza el Trabajo de Grado con el nuevo Estado
					rollbackDocumentoTrabajoGradoARL(transaccion)
					//logs.Error(err)
					panic(err.Error())
				}

			} else {
				rollbackDocumentoEscritoARL(transaccion)
				//logs.Error(err)
				panic(err.Error())
			}
		}
	} else {
		//logs.Error(err)
		panic(err.Error())
	}
	return alerta, outputError
}

func rollbackDocumentoTrabajoGradoARL(transaccion *models.TrSubirArl) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK DOCUMENTO TRABAAJO GRADO")
	var respuesta map[string]interface{}
	url := "/v1/documento_trabajo_grado/" + strconv.Itoa(transaccion.DocumentoTrabajoGrado.Id)
	if status, err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil && status != "200" {
		panic("Rollback solicitud trabajo grado" + err.Error())
	}
	rollbackDocumentoEscritoARL(transaccion)
	return nil
}

func rollbackDocumentoEscritoARL(transaccion *models.TrSubirArl) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK DOCUMENTO ESCRITO")
	var respuesta map[string]interface{}
	url := "/v1/documento_escrito/" + strconv.Itoa(transaccion.DocumentoEscrito.Id)
	if status, err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil && status != "200" {
		panic("Rollback solicitud trabajo grado" + err.Error())
	}
	return nil
}
