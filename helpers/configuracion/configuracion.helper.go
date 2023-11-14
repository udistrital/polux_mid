package configuracionHelper

import (
	"strings"

	"github.com/udistrital/polux_mid/helpers"
	"github.com/udistrital/polux_mid/models"
)

func CheckPermisoRoles(permiso string, roles []string) (allowed bool, outputError map[string]interface{}) {

	permisos := make([]models.PerfilXMenuOpcion, 0)
	payload := "?limit=1" +
		"&query=Perfil__Aplicacion__Nombre:Polux" +
		",Opcion__Nombre:" + permiso +
		",Perfil__Nombre__in:" + strings.Join(roles, "|")
	outputError = helpers.Get(new(models.PerfilXMenuOpcion), payload, &permisos)
	allowed = len(permisos) == 1

	return
}
