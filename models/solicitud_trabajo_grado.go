package models

import (
	"time"

	"github.com/astaxie/beego"
)

type SolicitudTrabajoGrado struct {
	Id                     int
	Fecha                  time.Time
	ModalidadTipoSolicitud *ModalidadTipoSolicitud
	TrabajoGrado           *TrabajoGrado
	PeriodoAcademico       string
}

func (*SolicitudTrabajoGrado) BasePath() string {
	return beego.AppConfig.String("PoluxCrud")
}

func (*SolicitudTrabajoGrado) Endpoint() string {
	return "solicitud_trabajo_grado"
}
