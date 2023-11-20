package models

import (
	"time"
)

type SolicitudTrabajoGrado struct {
	Id                     int
	Fecha                  time.Time
	ModalidadTipoSolicitud *ModalidadTipoSolicitud
	TrabajoGrado           *TrabajoGrado
	PeriodoAcademico       string
}