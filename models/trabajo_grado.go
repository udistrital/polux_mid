package models

import (

)

type TrabajoGrado struct {
	Id          int
	IdModalidad *Modalidad
	Titulo      string
	Distincion  string
	Etapa       string
}
