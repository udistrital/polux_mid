package models

import "time"

type Socializacion struct {
	Id           int           `orm:"column(id);pk;auto"`
	Fecha        time.Time     `orm:"column(fecha);type(timestamp without time zone)"`
	Lugar        int           `orm:"column(lugar)"`
	TrabajoGrado *TrabajoGrado `orm:"column(trabajo_grado);rel(fk)"`
}
