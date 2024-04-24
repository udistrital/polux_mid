package models

type AreasTrabajoGrado struct {
	Id               int
	AreaConocimiento int
	TrabajoGrado     *TrabajoGrado
	Activo           bool
}
