package models

import "time"

type Comentario struct {
	Id         int
	Comentario string
	Fecha      time.Time
	Autor      string
	Correccion *Correccion
}
