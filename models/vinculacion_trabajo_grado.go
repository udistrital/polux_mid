package models

import (
	"time"

	"github.com/astaxie/beego"
)

type VinculacionTrabajoGrado struct {
	Id              int
	Usuario         int
	Activo          bool
	FechaInicio     time.Time
	FechaFin        time.Time
	RolTrabajoGrado *RolTrabajoGrado
	TrabajoGrado    *TrabajoGrado
}

func (s *VinculacionTrabajoGrado) BasePath() string {
	return beego.AppConfig.String("PoluxCrud")
}

func (s *VinculacionTrabajoGrado) Endpoint() string {
	return "vinculacion_trabajo_grado"
}
