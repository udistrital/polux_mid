package models

import "time"

type RevisionTrabajoGrado struct {
	Id                         int
	NumeroRevision             int
	FechaRecepcion             time.Time
	FechaRevision              *time.Time
	EstadoRevisionTrabajoGrado int
	DocumentoTrabajoGrado      *DocumentoTrabajoGrado
	VinculacionTrabajoGrado    *VinculacionTrabajoGrado
}
