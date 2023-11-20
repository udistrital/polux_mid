package models

import (
)

type CarreraElegible struct {
	Id               int
	CodigoCarrera    int
	CuposExcelencia  float64
	CuposAdicionales float64
	Periodo          float64
	Anio             float64
	CodigoPensum     float64
	Nivel            string
}