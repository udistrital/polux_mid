package models

import (
)

type DetallePasantia struct {
	Id             int
	Empresa        int
	Horas          int
	ObjetoContrato string
	Observaciones  string
	TrabajoGrado   *TrabajoGrado
}