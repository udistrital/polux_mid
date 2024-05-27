package models

type EstudianteTrabajoGrado struct {
	Id                           int
	Estudiante                   string
	TrabajoGrado                 *TrabajoGrado
	EstadoEstudianteTrabajoGrado int
}
