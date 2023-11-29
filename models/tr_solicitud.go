package models

type TrSolicitud struct {
	Solicitud         *SolicitudTrabajoGrado
	Respuesta         *RespuestaSolicitud
	DetallesSolicitud *[]DetalleSolicitud
	UsuariosSolicitud *[]UsuarioSolicitud
}
