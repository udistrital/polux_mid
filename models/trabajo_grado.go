package models

type TrabajoGrado struct {
	Id                     int
	Titulo                 string
	Modalidad              *Modalidad
	EstadoTrabajoGrado     *EstadoTrabajoGrado
	DistincionTrabajoGrado *DistincionTrabajoGrado
	PeriodoAcademico 	   int
}

type EstadoTrabajoGrado struct {
	Id                int
	Nombre            string
	Descripcion       string
	CodigoAbreviacion string
	Activo            bool
}

type DistincionTrabajoGrado struct {
	Id                int
	Nombre            string
	Descripcion       string
	CodigoAbreviacion string
	Activo            bool
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
