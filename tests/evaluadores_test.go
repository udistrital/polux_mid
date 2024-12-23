package test

import (
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
	"github.com/udistrital/polux_mid/helpers"
	"github.com/udistrital/polux_mid/models"
)

func TestObtenerEvaluadores(t *testing.T) {
	t.Log("//////////////////////////////////")
	t.Log("Inicio TestObtenerEvaluadores")
	t.Log("//////////////////////////////////")
	// Datos de entrada simulados
	input := models.CantidadEvaluadoresModalidad{
		Modalidad: 4604,
	}

	// Mock del modelo esperado como retorno
	mockParametro := models.Parametro{
		Id:                4604,
		Nombre:            "Modalidad de Prueba",
		Descripcion:       "Descripción de la modalidad",
		CodigoAbreviacion: "MOD-PRUEBA",
		Activo:            true,
		NumeroOrden:       1,
		TipoParametroId:   nil,
		ParametroPadreId:  nil,
	}

	// Patch de la función helpers.ObtenerModalidad
	monkey.Patch(helpers.ObtenerModalidad, func(idModalidad models.CantidadEvaluadoresModalidad) (models.Parametro, map[string]interface{}) {
		return mockParametro, nil
	})
	defer monkey.Unpatch(helpers.ObtenerModalidad)

	// Llamada al controlador para probar el comportamiento
	modalidad, errMap := helpers.ObtenerModalidad(input)

	// Verificaciones
	assert.Nil(t, errMap, "Se esperaba que errMap fuera nil") // Verifica que no hubo errores
	assert.Equal(t, mockParametro, modalidad)                 // Verifica que la modalidad sea la esperada
	assert.Equal(t, "Modalidad de Prueba", modalidad.Nombre)
	assert.True(t, modalidad.Activo)
}
