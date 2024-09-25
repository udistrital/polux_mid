package helpers

import (
	"fmt"
	"strconv"

	"github.com/udistrital/polux_mid/models"
)

// Obtener modalidades de parametros
func ObtenerModalidad(idModalidad models.CantidadEvaluadoresModalidad) (modalidad models.Parametro, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("ERROR ", err)
			panic(DeferHelpers("AddTransaccionSolicitud", err))
		}
	}()
	url := "parametro/" + strconv.Itoa(idModalidad.Modalidad)
	if err := GetRequestNew("UrlCrudParametros", url, &modalidad); err != nil {
		//logs.Error(err.Error())
		panic(err.Error())
	}
	return
}
