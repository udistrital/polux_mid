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
	Contrato	   *DocumentoEscrito
	Carta		   *DocumentoEscrito
	HojaVidaDE	   *DocumentoEscrito
	DTG_Contrato   *DocumentoTrabajoGrado
	DTG_Carta	   *DocumentoTrabajoGrado
	DTG_HojaVida   *DocumentoTrabajoGrado
}