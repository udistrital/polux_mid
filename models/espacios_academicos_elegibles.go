package models

type EspaciosAcademicosElegibles struct {
	Id               int
	CodigoAsignatura int
	Activo           bool
	CarreraElegible  *CarreraElegible
}
