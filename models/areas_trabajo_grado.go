package models

import (
)

type AreasTrabajoGrado struct {
	Id               int
	AreaConocimiento *AreaConocimiento
	TrabajoGrado     *TrabajoGrado
	Activo           bool
}