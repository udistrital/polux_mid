package models

type TrSolicitud struct {
	Solicitud        *SolicitudTrabajoGrado
	Respuesta        *RespuestaSolicitud
	DetalleSolicitud *[]DetalleSolicitud
	UsuarioSolicitud *[]UsuarioSolicitud
}
