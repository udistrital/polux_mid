package models

import "github.com/astaxie/beego"

type DatosCoordinador struct {
	CoordinadorCollection struct {
		Coordinador []struct {
			NombreProyectoCurricular string `json:"nombre_proyecto_curricular"`
			CodigoProyectoCurricular string `json:"codigo_proyecto_curricular"`
			NombreCoordinador        string `json:"nombre_coordinador"`
		} `json:"coordinador"`
	} `json:"coordinadorCollection"`
}

func (*DatosCoordinador) BasePath() string {
	return beego.AppConfig.String("AcademicaJbpm")
}

func (*DatosCoordinador) Endpoint() string {
	return "coordinador_carrera"
}
