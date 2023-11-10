package helpers

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/udistrital/polux_mid/models"
)

// Obtener modalidades de parametros
func ObtenerModalidad(idModalidad models.CantidadEvaluadoresModalidad) (modalidad models.Parametro, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "ObtenerModalidad", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
	url := "parametro/" + strconv.Itoa(idModalidad.Modalidad)
	fmt.Println("URL ", url)
	if err := GetRequestNew("UrlCrudParametros", url, &modalidad); err != nil {
		logs.Error(err.Error())
		panic(err.Error())
	}
	fmt.Println("ASDAS ", modalidad)
	return
}
