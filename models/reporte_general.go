package models

import (
	"time"
)

type ReporteGeneral struct {
	Id                 int
	TrabajoGrado       int       `orm:"column(trabajo_grado)"`
	Titulo             string    `orm:"column(titulo)"`
	Modalidad          string    `orm:"column(modalidad)"`
	EstadoTrabajoGrado string    `orm:"column(estado)"`
	IdEstudiante       string    `orm:"column(id_estudiante)"`
	IdCoestudiante     string    `orm:"column(id_coestudiante)"`
	AreaConocimiento   string    `orm:"column(area_conocimiento)"`
	DocenteDirector    int       `orm:"column(docente_director)"`
	DocenteCodirector  int       `orm:"column(docente_codirector)"`
	Evaluador          int       `orm:"column(evaluador)"`
	FechaInicio        time.Time `orm:"column(fecha_inicio);type(timestamp without time zone)"`
	FechaFin           time.Time `orm:"column(fecha_fin);type(timestamp without time zone);null"`
	CalificacionUno    float32   `orm:"column(calificacion_1)"`
	CalificacionDos    float32   `orm:"column(calificacion_2)"`
}
