package models

import (
	"time"
)

type VinculacionTrabajoGrado struct {
	Id              int
	Usuario         int
	Activo          bool
	FechaInicio     time.Time
	FechaFin        time.Time
	RolTrabajoGrado int
	TrabajoGrado    *TrabajoGrado
}
