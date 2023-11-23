package models

type EspacioAcademicoInscrito struct {
	Id                             int
	Nota                           float64
	EspaciosAcademicosElegibles    *EspaciosAcademicosElegibles
	EstadoEspacioAcademicoInscrito int
	TrabajoGrado                   *TrabajoGrado
}
