package models

import (
)

type EspaciosAcademicosElegibles struct {
	Id               int
	CodigoAsignatura int
	Activo           bool
	CarreraElegible  *CarreraElegible
}