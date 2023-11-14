package solicitudesHelper

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/udistrital/polux_mid/helpers"
	"github.com/udistrital/polux_mid/helpers/autenticacionMID"
	configuracionHelper "github.com/udistrital/polux_mid/helpers/configuracion"
	"github.com/udistrital/polux_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

func GetSolicitudesByUser(user string) (solicitudes []models.RespuestaSolicitudRevisar, outputError map[string]interface{}) {

	funcion := "GetSolicitudesByUser - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	solicitudes = make([]models.RespuestaSolicitudRevisar, 0)
	dataUsuario, outputError := autenticacionMID.DataUsuario(user)
	if outputError != nil {
		return
	}

	permisoCoordinador, outputError := configuracionHelper.CheckPermisoRoles("Solicitudes", dataUsuario.Role)
	if outputError != nil {
		return
	}

	payload := "/user?"
	if dataUsuario.Codigo != "" {
		payload += "&codigoEstudiante=" + dataUsuario.Codigo
	}
	if dataUsuario.Documento != "" {
		payload += "&documento=" + dataUsuario.Documento
	}
	if permisoCoordinador {
		programas := new(models.DatosCoordinador)
		payload_ := "/" + dataUsuario.Documento + "/PREGRADO"
		outputError = helpers.GetXML(new(models.DatosCoordinador), payload_, &programas)
		if outputError != nil {
			return
		}

		if len(programas.CoordinadorCollection.Coordinador) > 0 {
			number, err := strconv.Atoi(programas.CoordinadorCollection.Coordinador[0].CodigoProyectoCurricular)
			if err != nil {
				logs.Error(err)
				eval := "strconv.Atoi(programas.CoordinadorCollection.Coordinador[0].CodigoProyectoCurricular)"
				outputError = errorctrl.Error(funcion+eval, err, "500")
				return
			}

			formatted := fmt.Sprintf("%02d", number)
			payload += "&codigoCarrera=" + formatted
		}
	}

	outputError = helpers.Get(new(models.RespuestaSolicitud), payload, &solicitudes)
	if outputError != nil {
		return
	}

	reqVinculacion := new(models.VinculacionTrabajoGrado)
	for _, s := range solicitudes {
		if s.EstadoSolicitud.Nombre == "Aprobado" { // ?
			continue
		}

		vinculacion := make([]models.VinculacionTrabajoGrado, 0)
		payload := "?limit=1&query=TrabajoGrado__Id:" + fmt.Sprint(s.SolicitudTrabajoGrado.TrabajoGrado.Id) + ",Usuario:" + dataUsuario.Documento
		helpers.Get(reqVinculacion, payload, &vinculacion)

		rol := ""
		if len(vinculacion) == 1 {
			rol = vinculacion[0].RolTrabajoGrado.CodigoAbreviacion
		} else if permisoCoordinador {
			rol = "COORDINADOR"
		}

		if rol != "" {
			s.Revisar = verificarSiPuedeAprobar(rol,
				s.SolicitudTrabajoGrado.ModalidadTipoSolicitud.Modalidad.Nombre,
				s.SolicitudTrabajoGrado.ModalidadTipoSolicitud.TipoSolicitud.Nombre,
				s.EstadoSolicitud.Nombre)
		}

	}

	return

}

func verificarSiPuedeAprobar(rol, modalidad, tipo, estado string) bool {
	permisosRol, ok := PermisosRevisionRol[rol]
	if !ok {
		return false
	}

	rolModalidad, ok := permisosRol[modalidad]
	if !ok {
		return false
	}

	modalidadTipo, ok := rolModalidad[tipo]
	if !ok {
		return false
	}

	for _, st := range modalidadTipo {
		if st == estado {
			return true
		}
	}

	return false

}

// Rol - Modalidad - TipoSolicitud - []Estados
var PermisosRevisionRol = map[string]map[string]map[string][]string{
	"DIRECTOR": {
		"Monografía": {
			"Solicitud inicial":          []string{},
			"Solicitud de socialización": []string{"Radicada"},
		},
	},
}
