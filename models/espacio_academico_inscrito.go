package models

import (
)

type EspacioAcademicoInscrito struct {
	Id                             int
	Nota                           float64
	EspaciosAcademicosElegibles    *EspaciosAcademicosElegibles
	EstadoEspacioAcademicoInscrito *EstadoEspacioAcademicoInscrito
	TrabajoGrado                   *TrabajoGrado
}