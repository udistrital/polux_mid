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

type Datos struct {
	Codigo        string
	Nombre        string
	Tipo   				string
	Modalidad			int
	PorcentajeCursado			float64
	Promedio			string
	Rendimiento		string
	Estado				string
	Nivel					string
	TipoCarrera		string
}
