package models

type DetalleSolicitud struct {
	Id                    int
	Descripcion           string
	SolicitudTrabajoGrado *SolicitudTrabajoGrado
	DetalleTipoSolicitud  *DetalleTipoSolicitud
}
