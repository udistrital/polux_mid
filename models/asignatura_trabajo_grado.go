package models

type AsignaturaTrabajoGrado struct {
	Id                           int
	CodigoAsignatura             int
	Periodo                      float64
	Anio                         float64
	Calificacion                 float64
	TrabajoGrado                 *TrabajoGrado
	EstadoAsignaturaTrabajoGrado int
	Activo                       bool
	FechaCreacion                string
	FechaModificacion            string
}
