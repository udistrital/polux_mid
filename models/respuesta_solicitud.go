package models

import (
	"time"

	"github.com/astaxie/beego"
)

type RespuestaSolicitud struct {
	Id                    int
	Fecha                 time.Time
	Justificacion         string
	EnteResponsable       int
	Usuario               int
	EstadoSolicitud       int
	SolicitudTrabajoGrado *SolicitudTrabajoGrado
	Activo                bool
}

func (*RespuestaSolicitud) BasePath() string {
	return beego.AppConfig.String("PoluxCrudUrl")
}

func (*RespuestaSolicitud) Endpoint() string {
	return "respuesta_solicitud"
}

type RespuestaSolicitudRevisar struct {
	RespuestaSolicitud
	Revisar bool
}
