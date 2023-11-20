package models

import (
)

type EstudianteTrabajoGrado struct {
	Id                           int
	Estudiante                   string
	TrabajoGrado                 *TrabajoGrado
	EstadoEstudianteTrabajoGrado *EstadoEstudianteTrabajoGrado
}