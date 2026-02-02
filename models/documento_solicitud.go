package models

type DocumentoSolicitud struct {
	Id                    int
	DocumentoEscrito      *DocumentoEscrito
	SolicitudTrabajoGrado *SolicitudTrabajoGrado
}
