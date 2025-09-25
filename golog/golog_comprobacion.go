package golog

import (
	"fmt"

	golog "github.com/mndrix/golog"

)

// Comprobar ...
func Comprobar(reglas string, reglaInyectada string) (rest string) {

	fmt.Println("=== INICIO Comprobar ===")
	fmt.Println("Reglas recibidas:")
	fmt.Println(reglas)
	fmt.Println("Regla inyectada:", reglaInyectada)

	//******QUITAR ARREGLO, DEJAR UNA SOLA VARIABLE PARA LAS REGLAS ******
	/*m := golog.NewMachine().Consult(reglas)
	if m.CanProve(reglaInyectada) {
		rest = "true"
	} else {
		rest = "false"
	}*/

	// Crear máquina con reglas
    m := golog.NewMachine().Consult(reglas)

    // Hechos que deberían cumplirse para validacion_requisitos
    pruebas := []string{
        "estado(20212015031, activo).",
        "porcentaje_tg(80).",                  // depende de lo que cargues
        "nivel_carrera(pregrado).",            // idem
        "cursado(20212015031, 75).",
        "nivel(20212015031, pregrado).",
    }

    for _, p := range pruebas {
        ok := m.CanProve(p)
        fmt.Printf("¿Se cumple %-40s ? -> %v\n", p, ok)
    }

    // Probar la regla completa
    if m.CanProve(reglaInyectada) {
        fmt.Println("✔ La regla se cumple:", reglaInyectada)
        rest = "true"
    } else {
        fmt.Println("❌ La regla NO se cumple:", reglaInyectada)
        rest = "false"
    }

	fmt.Println("=== FIN Comprobar ===")

	return

}

// Obtener ...
func Obtener(reglas string, reglaInyectada string) (rest string) {

	var res string
	m := golog.NewMachine().Consult(reglas)

	resultados := m.ProveAll(reglaInyectada)
	for _, solution := range resultados {
		res = fmt.Sprintf("%s", solution.ByName_("Y"))
	}

	return res

}
