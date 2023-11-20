package models

import (
)

type DetalleTrabajoGrado struct {
	Id                int
	Parametro         string
	Valor             string
	TrabajoGrado      *TrabajoGrado
	Activo            bool
	FechaCreacion     string
	FechaModificacion string
}