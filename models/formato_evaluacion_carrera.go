package models

import "time"

type FormatoEvaluacionCarrera struct {
	Id             int       `orm:"column(id);pk;auto"`
	Activo         bool      `orm:"column(activo)"`
	CodigoProyecto int       `orm:"column(codigo_proyecto)"`
	FechaInicio    time.Time `orm:"column(fecha_inicio);type(timestamp without time zone)"`
	FechaFin       time.Time `orm:"column(fecha_fin);type(timestamp without time zone);null"`
	Modalidad      int       `orm:"column(modalidad);"`
	Formato        *Formato  `orm:"column(formato);rel(fk)"`
}
