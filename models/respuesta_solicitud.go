package models

import (
	"time"
)

type RespuestaSolicitud struct {
	Id                    int
	Fecha                 time.Time
	Justificacion         string
	EnteResponsable       int
	Usuario               int
	EstadoSolicitud       *EstadoSolicitud
	SolicitudTrabajoGrado *SolicitudTrabajoGrado
	Activo                bool
}