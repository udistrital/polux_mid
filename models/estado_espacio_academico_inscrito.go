package models

import (
)

type EstadoEspacioAcademicoInscrito struct {
	Id                int
	Nombre            string
	Descripcion       string
	CodigoAbreviacion string
	Activo            bool
}