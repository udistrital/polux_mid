package models

import (
	"time"
)

type SolicitudMaterias struct {
	Solicitud   int
	Fecha       time.Time
	Estudiante  string
	Nombre      string
	Promedio    string
	Rendimiento string
	Estado      *EstadoSolicitud
	Respuesta   string
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

type RespuestaSolicitud struct {
	Id                    int
	Fecha                 time.Time
	Justificacion         string
	EnteResponsable       int
	Usuario               int
	EstadoSolicitud       *EstadoSolicitud
	SolicitudTrabajoGrado *SolicitudTrabajoGrado
}

type EstadoSolicitud struct {
	Id                int
	Nombre            string
	Descripcion       string
	CodigoAbreviacion string
	Activo            bool
}

type SolicitudTrabajoGrado struct {
	Id                     int
	Fecha                  time.Time
	ModalidadTipoSolicitud *ModalidadTipoSolicitud
	TrabajoGrado           *TrabajoGrado
}

type ModalidadTipoSolicitud struct {
	Id            int
	TipoSolicitud *TipoSolicitud
	Modalidad     *Modalidad
}

type TipoSolicitud struct {
	Id                int
	Nombre            string
	Descripcion       string
	CodigoAbreviacion string
	Activo            bool
}

type TrSolicitud struct {
	NumAdmitidos *Cupos
	Solicitudes  *[]SolicitudMaterias
}

type Vals []SolicitudMaterias
