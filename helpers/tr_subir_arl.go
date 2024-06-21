package helpers

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/udistrital/polux_mid/models"
)

func AddTransaccionSubirArl(transaccion *models.TrSubirArl) (alerta []string, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "AddTransaccionSubirArl", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
	alerta = append(alerta, "Success")

	transaccion.DocumentoEscrito.Id = 0
	transaccion.DocumentoEscrito.Titulo = "ARL de la pasantía con titulo: " + transaccion.DocumentoEscrito.Titulo
	transaccion.DocumentoEscrito.Resumen = transaccion.DocumentoEscrito.Titulo

	url := "/v1/documento_escrito"
	var resDocumentoEscrito map[string]interface{}
	if err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoEscrito, &transaccion.DocumentoEscrito); err == nil { //Se guarda la ARL en Documento Escrito. Desde el Cliente ya viene con el tipo "Documentos adicionales asociados a las pasantias"

		transaccion.DocumentoTrabajoGrado.Id = 0
		transaccion.DocumentoEscrito.Id = int(resDocumentoEscrito["Id"].(float64))
		transaccion.DocumentoTrabajoGrado.DocumentoEscrito.Id = int(resDocumentoEscrito["Id"].(float64))

		url := "/v1/documento_trabajo_grado"
		var resDocumentoTrabajoGrado map[string]interface{}
		if err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoTrabajoGrado, &transaccion.DocumentoTrabajoGrado); err == nil { //Se guarda la relación entre Documento Escrito y el Trabajo de Grado

			transaccion.DocumentoTrabajoGrado.Id = int(resDocumentoTrabajoGrado["Id"].(float64))

			var parametro []models.Parametro
			url = "parametro?query=CodigoAbreviacion:ACEA_PLX"
			if err := GetRequestNew("UrlCrudParametros", url, &parametro); err != nil { //Se busca el Estado Trabajo de Grado "El estudiante cargó el certificado de afiliación de ARL y es necesaria la aprobación por parte de la Oficina de Extención de Pasantías"
				logs.Error(err.Error())
				panic(err.Error())
			}

			transaccion.TrabajoGrado.EstadoTrabajoGrado = parametro[0].Id

			url := "/v1/trabajo_grado/"+ strconv.Itoa(transaccion.TrabajoGrado.Id)
			var resTrabajoGrado map[string]interface{}
			if err := SendRequestNew("PoluxCrudUrl", url, "PUT", &resTrabajoGrado, &transaccion.TrabajoGrado); err != nil { //Se actualiza el Trabajo de Grado con el nuevo Estado
				rollbackDocumentoTrabajoGradoARL(transaccion)
				logs.Error(err)
				panic(err.Error())
			}

		} else {
			rollbackDocumentoEscritoARL(transaccion)
			logs.Error(err)
			panic(err.Error())
		}

	} else {
		logs.Error(err)
		panic(err.Error())
	}
	return alerta, outputError
}

func rollbackUpdateTrabajoGrado(transaccion *models.TrSubirArl, EstadoAnterior int) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK TRABAJO GRADO")
	var respuesta map[string]interface{}
	transaccion.TrabajoGrado.EstadoTrabajoGrado = EstadoAnterior
	url := "/v1/trabajo_grado/" + strconv.Itoa(transaccion.TrabajoGrado.Id)
	if err := SendRequestNew("PoluxCrudUrl", url, "PUT", &respuesta, nil); err != nil {
		panic("Rollback solicitud trabajo grado" + err.Error())
	}
	return nil
}

func rollbackDocumentoTrabajoGradoARL(transaccion *models.TrSubirArl) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK DOCUMENTO TRABAAJO GRADO")
	var respuesta map[string]interface{}
	url := "/v1/documento_trabajo_grado/" + strconv.Itoa(transaccion.DocumentoTrabajoGrado.Id)
	if err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil {
		panic("Rollback solicitud trabajo grado" + err.Error())
	}
	rollbackDocumentoEscritoARL(transaccion)
	return nil
}

func rollbackDocumentoEscritoARL(transaccion *models.TrSubirArl) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK DOCUMENTO ESCRITO")
	var respuesta map[string]interface{}
	url := "/v1/documento_escrito/" + strconv.Itoa(transaccion.DocumentoEscrito.Id)
	if err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil {
		panic("Rollback solicitud trabajo grado" + err.Error())
	}
	return nil
}