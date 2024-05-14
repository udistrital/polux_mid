package models

import "time"

type Formato struct {
	Id               int       `orm:"column(id);pk;auto"`
	Nombre           string    `orm:"column(nombre)"`
	Introduccion     string    `orm:"column(introduccion);null"`
	FechaRealizacion time.Time `orm:"column(fecha_realizacion);type(timestamp without time zone)"`
}
