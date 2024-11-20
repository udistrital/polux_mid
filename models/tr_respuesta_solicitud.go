package models

type TrRevision struct {
	TrabajoGrado          *TrabajoGrado
	Vinculaciones         *[]VinculacionTrabajoGrado //Cambio de director o evaluador
	DocumentoEscrito      *[]DocumentoEscrito
	DocumentoTrabajoGrado *DocumentoTrabajoGrado
	DetalleTrabajoGrado   *[]DetalleTrabajoGrado //Solicitud inicial de TrabajoGrado
}

type TrRespuestaSolicitud struct {
	RespuestaAnterior           *RespuestaSolicitud
	RespuestaNueva              *RespuestaSolicitud
	DocumentoSolicitud          *DocumentoSolicitud
	TipoSolicitud               *TipoSolicitud
	Vinculaciones               *[]VinculacionTrabajoGrado  //Cambio de director o evaluador..//Cancelación TG
	EstudianteTrabajoGrado      *EstudianteTrabajoGrado     //Cancelación trabajo grado
	VinculacionesCancelacion    *[]VinculacionTrabajoGrado  //Vinculaciones para cancelacion de trabajo de grado
	TrTrabajoGrado              *TrTrabajoGrado             //Solictudes iniciales
	ModalidadTipoSolicitud      *ModalidadTipoSolicitud     //Para saber el tipo de solicitud inicial
	TrabajoGrado                *TrabajoGrado               //Cambio Titulo
	SolicitudTrabajoGrado       *SolicitudTrabajoGrado      //solicitud inicial
	EspaciosAcademicos          *[]EspacioAcademicoInscrito //Solicitud de cambio de asignaturas
	DetallesPasantia            *DetallePasantia            //SOlicitud inicial de pasantia
	TrRevision                  *TrRevision                 //Solicitud de revisión
	EspaciosAcademicosInscritos *[]EspacioAcademicoInscrito //Espacios academicos inscritos
	CausaProrroga               *[]DetalleTrabajoGrado      //Solicitud de Prorroga
	DetallesPasantiaExterna     *[]DetalleTrabajoGrado      //Solicitud Inicial PASANTÍA EXTERNA
}
