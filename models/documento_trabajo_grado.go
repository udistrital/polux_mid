package models

type DocumentoTrabajoGrado struct {
	Id                int
	TrabajoGrado      *TrabajoGrado
	DocumentoEscrito  *DocumentoEscrito
	Activo            bool
	FechaCreacion     string
	FechaModificacion string
}
