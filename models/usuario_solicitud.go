package models

import "github.com/astaxie/beego"

type UsuarioSolicitud struct {
	Id                    int
	Usuario               string
	SolicitudTrabajoGrado *SolicitudTrabajoGrado
}

func (s *UsuarioSolicitud) BasePath() string {
	return beego.AppConfig.String("PoluxCrud")
}

func (s *UsuarioSolicitud) Endpoint() string {
	return "usuario_solicitud"
}
