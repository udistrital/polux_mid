package golog

import (
	"fmt"

	golog "github.com/mndrix/golog"
)

// Comprobar ...
func Comprobar(reglas string, regla_inyectada string) (rest string) {

	//******QUITAR ARREGLO, DEJAR UNA SOLA VARIABLE PARA LAS REGLAS ******
	m := golog.NewMachine().Consult(reglas)
	if m.CanProve(regla_inyectada) {
		rest = "true"
	} else {
		rest = "false"
	}

	return

}

// Obtener ...
func Obtener(reglas string, regla_inyectada string) (rest string) {

	var res string
	m := golog.NewMachine().Consult(reglas)

	resultados := m.ProveAll(regla_inyectada)
	for _, solution := range resultados {
		res = fmt.Sprintf("%s", solution.ByName_("Y"))
	}

	return res

}
