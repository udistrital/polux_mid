package models

import (
)

type DocumentoTrabajoGrado struct {
	Id               int
	TrabajoGrado     *TrabajoGrado
	DocumentoEscrito *DocumentoEscrito
}