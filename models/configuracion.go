package models

import "github.com/astaxie/beego"

type PerfilXMenuOpcion struct {
	Id int
}

func (*PerfilXMenuOpcion) BasePath() string {
	return beego.AppConfig.String("ConfiguracionCrud")
}

func (*PerfilXMenuOpcion) Endpoint() string {
	return "perfil_x_menu_opcion"
}
