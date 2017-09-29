package models

type TrabajoGrado struct {
	Id          int
	IdModalidad *Modalidad
	Titulo      string
	Distincion  string
	Etapa       string
}

type CantidadModalidad struct {
	Modalidad string
	Cantidad  string
}

type Datos struct {
	Codigo            string
	Nombre            string
	Tipo              string
	Modalidad         int
	PorcentajeCursado string
	Promedio          string
	Rendimiento       string
	Estado            string
	Nivel             string
	TipoCarrera       string
}
