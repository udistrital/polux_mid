package models

type Detalle struct {
	Id                int
	Nombre            string
	Enunciado         string
	Descripcion       string
	CodigoAbreviacion string
	Activo            bool
	TipoDetalle       int
}
