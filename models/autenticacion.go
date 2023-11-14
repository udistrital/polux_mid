package models

import "github.com/astaxie/beego"

type UsuarioAutenticacion struct {
	Codigo             string
	Estado             string
	FamilyName         string
	Documento          string   `json:"documento"`
	DocumentoCompuesto string   `json:"documento_compuesto"`
	Email              string   `json:"email"`
	Role               []string `json:"role"`
}

func (*UsuarioAutenticacion) BasePath() string {
	return beego.AppConfig.String("AutenticacionMid")
}

func (*UsuarioAutenticacion) Endpoint() string {
	return "token/userRol"
}

type UsuarioDataRequest struct {
	User string `json:"user"`
}
