package models

type DocumentosMaterias struct {
	SolicitudEscrita 		*DocumentoEscrito
    Justificacion 			*DocumentoEscrito
    CartaAceptacion			*DocumentoEscrito
    SabanaNotas				*DocumentoEscrito
	DTG_SolicitudEscrita 	*DocumentoTrabajoGrado
    DTG_Justificacion 		*DocumentoTrabajoGrado
    DTG_CartaAceptacion		*DocumentoTrabajoGrado
    DTG_SabanaNotas			*DocumentoTrabajoGrado
}