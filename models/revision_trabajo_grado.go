package models

type RevisionTrabajoGrado struct {
	Id                         int
	NumeroRevision             int
	FechaRecepcion             string
	FechaRevision              string
	EstadoRevisionTrabajoGrado int
	DocumentoTrabajoGrado      *DocumentoTrabajoGrado
	VinculacionTrabajoGrado    *VinculacionTrabajoGrado
}
