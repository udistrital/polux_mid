package models

type ModalidadTipoSolicitud struct {
	Id            int
	TipoSolicitud *TipoSolicitud
	Modalidad     *Modalidad
}
