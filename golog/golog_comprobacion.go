package golog

import (
	"fmt"

	. "github.com/mndrix/golog"
)

func Comprobar(reglas string, regla_inyectada string) (rest string) {

	//******QUITAR ARREGLO, DEJAR UNA SOLA VARIABLE PARA LAS REGLAS ******
	var res string
	m := NewMachine().Consult(reglas)

	if m.CanProve(regla_inyectada) {
		res = "true"
	} else {
		res = "false"
	}

	return res

}

func Obtener(reglas string, regla_inyectada string) (rest string) {

	var res string
	m := NewMachine().Consult(reglas)

	resultados := m.ProveAll(regla_inyectada)
	for _, solution := range resultados {
		res = fmt.Sprintf("%s", solution.ByName_("Y"))
	}

	return res

}
