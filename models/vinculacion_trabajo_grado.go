package models

type VinculacionTrabajoGrado struct {
	Id              int
	Usuario         int
	Activo          bool
	FechaInicio     *string
	FechaFin        *string
	RolTrabajoGrado int
	TrabajoGrado    *TrabajoGrado
}
