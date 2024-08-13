package models

import (
	"time"
)

type ReporteSolicitud struct {
	Id                 int       `orm:"column(id)"`
	TrabajoGrado       int       `orm:"column(trabajo_grado)"`
	Titulo             string    `orm:"column(titulo)"`
	Modalidad          string    `orm:"column(modalidad)"`
	EstadoTrabajoGrado string    `orm:"column(estado_trabajo_grado)"`
	IdEstudiante       string    `orm:"column(id_estudiante)"`
	IdCoestudiante     string    `orm:"column(id_coestudiante)"`
	ProgramaAcademico  string    `orm:"column(programa_academico)"`
	Coordinador        int       `orm:"column(coordinador)"`
	DocenteDirector    int       `orm:"column(docente_director)"`
	DocenteCodirector  int       `orm:"column(docente_codirector)"`
	Evaluador          int       `orm:"column(evaluador)"`
	FechaSolicitud     time.Time `orm:"column(fecha_solicitud);type(timestamp without time zone)"`
	FechaRevision      time.Time `orm:"column(fecha_revision);type(timestamp without time zone);null"`
	Solicitud          string    `orm:"column(tipo_solicitud)"`
	Observacion        string    `orm:"column(justificacion)"`
	Respuesta          string    `orm:"column(estado_solicitud)"`
}
