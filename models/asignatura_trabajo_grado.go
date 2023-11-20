package models

import (
)

type AsignaturaTrabajoGrado struct {
	Id                           int
	CodigoAsignatura             int
	Periodo                      float64
	Anio                         float64
	Calificacion                 float64
	TrabajoGrado                 *TrabajoGrado
	EstadoAsignaturaTrabajoGrado *EstadoAsignaturaTrabajoGrado
	Activo                       bool
	FechaCreacion                string
	FechaModificacion            string
}