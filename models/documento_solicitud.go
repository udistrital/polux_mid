package models

import (
)

type DocumentoSolicitud struct {
	Id                    int
	DocumentoEscrito      *DocumentoEscrito
	SolicitudTrabajoGrado *SolicitudTrabajoGrado
}