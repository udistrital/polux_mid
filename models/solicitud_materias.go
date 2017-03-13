package models

import (
	"time"
)

type SolicitudMaterias struct {
	Solicitud      int
	Fecha          time.Time
	Estudiante     string
	Nombre         string
	Promedio       string
	Rendimiento    string
	Estado         string
}

type Solicitud struct {
	Id             int
	IdTrabajoGrado *TrabajoGrado
	Fecha          time.Time
	Estado         string
	Formalizacion  string
	CodigoCarrera  float64
	Periodo        string
	Anio           float64
}

type TrSolicitud struct {
	NumAdmitidos *Cupos
	Solicitudes *[]SolicitudMaterias
}

type Vals []SolicitudMaterias
