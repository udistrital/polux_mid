package models

type TrTrabajoGrado struct {
	TrabajoGrado            *TrabajoGrado
	EstudianteTrabajoGrado  *[]EstudianteTrabajoGrado
	DocumentoEscrito        *DocumentoEscrito
	DocumentoTrabajoGrado   *DocumentoTrabajoGrado
	AreasTrabajoGrado       *[]AreasTrabajoGrado
	VinculacionTrabajoGrado *[]VinculacionTrabajoGrado
	AsignaturasTrabajoGrado *[]AsignaturaTrabajoGrado
	DocumentosMaterias		*DocumentosMaterias
}
