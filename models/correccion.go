package models

type Correccion struct {
	Id                   int
	Observacion          string
	Pagina               float64
	RevisionTrabajoGrado *RevisionTrabajoGrado
	Documento            bool
}
