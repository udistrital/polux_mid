package models

// TrabajoGrado ...
type TrabajoGrado struct {
	Id                     int
	Titulo                 string
	Modalidad              *Modalidad
	EstadoTrabajoGrado     *EstadoTrabajoGrado
	DistincionTrabajoGrado *DistincionTrabajoGrado
	PeriodoAcademico 	   int
}

// EstadoTrabajoGrado ...
type EstadoTrabajoGrado struct {
	Id                int
	Nombre            string
	Descripcion       string
	CodigoAbreviacion string
	Activo            bool
}

// DistincionTrabajoGrado ...
type DistincionTrabajoGrado struct {
	Id                int
	Nombre            string
	Descripcion       string
	CodigoAbreviacion string
	Activo            bool
}

// CantidadModalidad ...
type CantidadModalidad struct {
	Modalidad string
	Cantidad  string
}

// Datos ...
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
