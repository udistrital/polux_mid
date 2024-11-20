package helpers

import (
	"fmt"
	"strconv"

	"github.com/udistrital/polux_mid/models"
)

func AddTransaccionRegistrarActaSeguimiento(transaccion *models.TrRegistrarActaSeguimiento) (alerta []string, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("ERROR ", err)
			panic(DeferHelpers("AddTransaccionSolicitud", err))
		}
	}()
	alerta = append(alerta, "Success")

	transaccion.DocumentoEscrito.Id = 0
	transaccion.DocumentoEscrito.Resumen = transaccion.DocumentoEscrito.Resumen + " del trabajo de grado con id: " + strconv.Itoa(transaccion.DocumentoTrabajoGrado.TrabajoGrado.Id)

	url := "/v1/documento_escrito"
	var resDocumentoEscrito map[string]interface{}
	if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoEscrito, &transaccion.DocumentoEscrito); err == nil && status == "201" { //Se guarda el acta en Documento Escrito

		transaccion.DocumentoTrabajoGrado.Id = 0
		transaccion.DocumentoEscrito.Id = int(resDocumentoEscrito["Id"].(float64))
		transaccion.DocumentoTrabajoGrado.DocumentoEscrito.Id = int(resDocumentoEscrito["Id"].(float64))

		url := "/v1/documento_trabajo_grado"
		var resDocumentoTrabajoGrado map[string]interface{}
		if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &resDocumentoTrabajoGrado, &transaccion.DocumentoTrabajoGrado); err == nil && status == "201" { //Se guarda la relación entre Documento Escrito y el Trabajo de Grado

			transaccion.DocumentoTrabajoGrado.Id = int(resDocumentoTrabajoGrado["Id"].(float64))

		} else {
			rollbackDocumentoEscritoActaSeguimiento(transaccion)
			// logs.Error(err)
			// panic(err.Error())
		}
	} else {
		//logs.Error(err)
		panic(err.Error())
	}
	return alerta, outputError
}

func rollbackDocumentoEscritoActaSeguimiento(transaccion *models.TrRegistrarActaSeguimiento) (outputError map[string]interface{}) {
	fmt.Println("ROLLBACK DOCUMENTO ESCRITO")
	var respuesta map[string]interface{}
	url := "/v1/documento_escrito/" + strconv.Itoa(transaccion.DocumentoEscrito.Id)
	if status, err := SendRequestNew("PoluxCrudUrl", url, "DELETE", &respuesta, nil); err != nil && status != "200" {
		panic("Rollback registrar acta de seguimiento" + err.Error())
	}
	return nil
}
