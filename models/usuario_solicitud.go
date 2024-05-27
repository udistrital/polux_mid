package models

type UsuarioSolicitud struct {
	Id                    int
	Usuario               string
	SolicitudTrabajoGrado *SolicitudTrabajoGrado
}
