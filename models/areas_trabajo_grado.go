package models

import (
)

type AreasTrabajoGrado struct {
	Id               int
	AreaConocimiento int
	TrabajoGrado     *TrabajoGrado
	Activo           bool
}