package models

type DetalleTipoSolicitud struct {
	Id                     int
	Detalle                *Detalle
	ModalidadTipoSolicitud *ModalidadTipoSolicitud
	Activo                 bool
	Requerido              bool
	NumeroOrden            int
}
