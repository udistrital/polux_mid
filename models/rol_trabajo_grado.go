package models

import "github.com/astaxie/beego"

type RolTrabajoGrado struct {
	Id                int
	Nombre            string
	Descripcion       string
	CodigoAbreviacion string
	Activo            bool
}

func (s *RolTrabajoGrado) BasePath() string {
	return beego.AppConfig.String("PoluxCrud")
}

func (s *RolTrabajoGrado) Endpoint() string {
	return "rol_trabajo_grado"
}
