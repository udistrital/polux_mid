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
func (v Vals) Len() int      { return len(v) }
func (v Vals) Swap(i, j int) { v[i], v[j] = v[j], v[i] }

type ValsAscByC struct{ Vals }
func (v ValsAscByC) Less(i, j int) bool { return v.Vals[i].Promedio > v.Vals[j].Promedio }
