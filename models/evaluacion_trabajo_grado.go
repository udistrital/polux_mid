package models

type EvaluacionTrabajoGrado struct {
	Id                       int                       `orm:"column(id);pk;auto"`
	Nota                     float64                   `orm:"column(nota)"`
	VinculacionTrabajoGrado  *VinculacionTrabajoGrado  `orm:"column(vinculacion_trabajo_grado);rel(fk)"`
	FormatoEvaluacionCarrera *FormatoEvaluacionCarrera `orm:"column(formato_evaluacion_carrera);rel(fk);null"`
	Socializacion            *Socializacion            `orm:"column(socializacion);rel(fk);null"`
}
