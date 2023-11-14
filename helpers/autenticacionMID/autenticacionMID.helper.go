package autenticacionMID

import (
	"github.com/udistrital/polux_mid/helpers"
	"github.com/udistrital/polux_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

func DataUsuario(usuarioWSO2 string) (dataUsuario models.UsuarioAutenticacion, outputError map[string]interface{}) {

	funcion := "DataUsuario - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	payload := models.UsuarioDataRequest{
		User: usuarioWSO2,
	}

	outputError = helpers.Post(new(models.UsuarioAutenticacion), payload, &dataUsuario)

	return
}
