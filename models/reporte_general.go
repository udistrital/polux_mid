package models

import (
	"encoding/xml"
	"time"
)

type ReporteGeneral struct {
	Id                      int
	TrabajoGrado            int    `orm:"column(trabajo_grado)"`
	Titulo                  string `orm:"column(titulo)"`
	Modalidad               string `orm:"column(modalidad)"`
	EstadoTrabajoGrado      string `orm:"column(estado)"`
	IdEstudiante            string `orm:"column(id_estudiante)"`
	NombreEstudiante        string
	IdCoestudiante          string `orm:"column(id_coestudiante)"`
	NombreCoestudiante      string
	ProgramaAcademico       string
	AreaConocimiento        string `orm:"column(area_conocimiento)"`
	DocenteDirector         int    `orm:"column(docente_director)"`
	NombreDocenteDirector   string
	DocenteCodirector       int `orm:"column(docente_codirector)"`
	NombreDocenteCodirector string
	Evaluador               int `orm:"column(evaluador)"`
	NombreEvaluador         string
	FechaInicio             time.Time `orm:"column(fecha_inicio);type(timestamp without time zone)"`
	FechaFin                time.Time `orm:"column(fecha_fin);type(timestamp without time zone);null"`
	CalificacionUno         float32   `orm:"column(calificacion_1)"`
	CalificacionDos         float32   `orm:"column(calificacion_2)"`
}

type DatosBasicosEstudiante struct {
	Nombre  string `xml:"nombre"`
	Carrera string `xml:"carrera"`
}

type DatosEstudianteCollection struct {
	XMLName                xml.Name                 `xml:"datosEstudianteCollection"`
	DatosBasicosEstudiante []DatosBasicosEstudiante `xml:"datosBasicosEstudiante"`
}

type FiltrosReporte struct {
	ProyectoCurricular	string
	Estado				string
	FechaInicio			time.Time
	FechaFin			time.Time
	IdEstFinalizado		int
	IdEstCancelado		int
}