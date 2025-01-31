package test

import (
	"errors"
	"strconv"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
	"github.com/udistrital/polux_mid/helpers"
	"github.com/udistrital/polux_mid/models"
)

func TestAddTransaccionSubirArl(t *testing.T) {

	t.Log("//////////////////////////////////")
	t.Log("Inicio TestAddTransaccionSubirArl")
	t.Log("//////////////////////////////////")

	t.Run("Caso #1: Primer cargue de ARL sin errores", func(t *testing.T) {

		transaccion := &models.TrSubirArl{
			DocumentoEscrito: &models.DocumentoEscrito{
				Id:					  0,
				Titulo:               "Pasantía de pruebas unitarias",
				Enlace:               "abcdefg12345",
				Resumen:              "",
				TipoDocumentoEscrito: 12345,
			},
			DocumentoTrabajoGrado: &models.DocumentoTrabajoGrado{
				TrabajoGrado: &models.TrabajoGrado{
					Id:					12345,
				},      
				DocumentoEscrito:&models.DocumentoEscrito{
					Id:					  0,
				},
			},
			TrabajoGrado: &models.TrabajoGrado{
				Id:					12345,
				EstadoTrabajoGrado: 22,
			},
		}

		monkey.Patch(helpers.GetRequestNew, func (baseURL, url string, target interface{}) error {
			switch url {
			case "parametro?query=CodigoAbreviacion:ARC_PLX":
				*target.(*[]models.Parametro) = []models.Parametro{
					{Id: 1, Nombre: "ARL Rechazada"},
				}
			case "parametro?query=CodigoAbreviacion:ACEA_PLX":
				*target.(*[]models.Parametro) = []models.Parametro{
					{Id: 2, Nombre: "ARL Cargada, en espera de aprobación"},
				}
		}
			return nil
		})

		defer monkey.Unpatch(helpers.GetRequestNew)

		monkey.Patch(helpers.SendRequestNew, func(baseURL, url, method string, target interface{}, body interface{}) (string, error) {

			switch url {
			case "/v1/documento_escrito":
				*target.(*map[string]interface{}) = map[string]interface{}{"Id": float64(88)}
			case "/v1/documento_trabajo_grado":
				*target.(*map[string]interface{}) = map[string]interface{}{"Id": float64(99)}
			case "/v1/trabajo_grado/":
				*target.(*map[string]interface{}) = map[string]interface{}{"Id": float64(12345)}
			}

			return "201",nil
		})

		defer monkey.Unpatch(helpers.SendRequestNew)

		// Llamar a la función
		result, err := helpers.AddTransaccionSubirArl(transaccion)

		// Validar resultados
		assert.Empty(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("Caso #2: Segundo cargue de ARL sin errores", func(t *testing.T) {

		transaccion := &models.TrSubirArl{
			DocumentoEscrito: &models.DocumentoEscrito{
				Id:					  0,
				Titulo:               "Pasantía de pruebas unitarias",
				Enlace:               "abcdefg12345",
				Resumen:              "",
				TipoDocumentoEscrito: 12345,
			},
			DocumentoTrabajoGrado: &models.DocumentoTrabajoGrado{
				TrabajoGrado: &models.TrabajoGrado{
					Id:					12345,
				},      
				DocumentoEscrito:&models.DocumentoEscrito{
					Id:					  0,
				},
			},
			TrabajoGrado: &models.TrabajoGrado{
				Id:					12345,
				EstadoTrabajoGrado: 1,
			},
		}

		monkey.Patch(helpers.GetRequestNew, func (baseURL, url string, target interface{}) error {
			switch url {
			case "parametro?query=CodigoAbreviacion:ARC_PLX":
				*target.(*[]models.Parametro) = []models.Parametro{
					{Id: 1, Nombre: "ARL Rechazada"},
				}
			case "/v1/documento_trabajo_grado?query=trabajo_grado__Id:"+strconv.Itoa(transaccion.TrabajoGrado.Id)+",documento_escrito__tipo_documento_escrito:1":
				*target.(*[]models.DocumentoTrabajoGrado) = []models.DocumentoTrabajoGrado{
					{Id: 77},
				}
			case "parametro?query=CodigoAbreviacion:ACEA_PLX":
				*target.(*[]models.Parametro) = []models.Parametro{
					{Id: 2, Nombre: "ARL Cargada, en espera de aprobación"},
				}
		}
			return nil
		})

		defer monkey.Unpatch(helpers.GetRequestNew)

		monkey.Patch(helpers.GetJson, func (url string, target interface{}) error {
			*target.(*[]models.TipoDocumento) = []models.TipoDocumento{
				{Id: 1, Nombre: "Documentos Pasantía"},
			}
			return nil
		})

		defer monkey.Unpatch(helpers.GetJson)

		monkey.Patch(helpers.SendRequestNew, func(baseURL, url, method string, target interface{}, body interface{}) (string, error) {

			switch url {
			case "/v1/documento_escrito":
				*target.(*map[string]interface{}) = map[string]interface{}{"Id": float64(88)}
			case "/v1/documento_trabajo_grado/77":
				*target.(*map[string]interface{}) = map[string]interface{}{"Id": float64(99)}
				return "200",nil
			case "/v1/trabajo_grado":
				*target.(*map[string]interface{}) = map[string]interface{}{"Id": float64(12345)}
			}

			return "201",nil
		})

		defer monkey.Unpatch(helpers.SendRequestNew)

		// Llamar a la función
		result, err := helpers.AddTransaccionSubirArl(transaccion)

		// Validar resultados
		assert.Empty(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("Caso #3: Fallo al momento de cargar el documento escrito", func(t *testing.T) {

		transaccion := &models.TrSubirArl{
			DocumentoEscrito: &models.DocumentoEscrito{
				Id:					  0,
				Titulo:               "Pasantía de pruebas unitarias",
				Enlace:               "abcdefg12345",
				Resumen:              "",
				TipoDocumentoEscrito: 12345,
			},
			DocumentoTrabajoGrado: &models.DocumentoTrabajoGrado{
				TrabajoGrado: &models.TrabajoGrado{
					Id:					12345,
				},      
				DocumentoEscrito:&models.DocumentoEscrito{
					Id:					  0,
				},
			},
			TrabajoGrado: &models.TrabajoGrado{
				Id:					12345,
				EstadoTrabajoGrado: 1,
			},
		}

		monkey.Patch(helpers.SendRequestNew, func(baseURL, url, method string, target interface{}, body interface{}) (string, error) {

			*target.(*map[string]interface{}) = map[string]interface{}{"Id": float64(88)}

			return "404", errors.New("sintaxis incorrecta")
		})

		defer monkey.Unpatch(helpers.SendRequestNew)

		// Valida si se ejecuta el Panic
		defer func() {
			if err := recover(); err != nil {
				t.Logf("Panic en la transacción esperado: %v", err)
			} else {
				t.Errorf("No se ejecuta el Panic esperado")
			}
		}()

		// Llamar a la función
		_ : helpers.AddTransaccionSubirArl(transaccion)
	})

	t.Run("Caso #4: Fallo al momento de cargar Documento Trabajo Grado", func(t *testing.T) {

		transaccion := &models.TrSubirArl{
			DocumentoEscrito: &models.DocumentoEscrito{
				Id:					  0,
				Titulo:               "Pasantía de pruebas unitarias",
				Enlace:               "abcdefg12345",
				Resumen:              "",
				TipoDocumentoEscrito: 12345,
			},
			DocumentoTrabajoGrado: &models.DocumentoTrabajoGrado{
				TrabajoGrado: &models.TrabajoGrado{
					Id:					12345,
				},      
				DocumentoEscrito:&models.DocumentoEscrito{
					Id:					  0,
				},
			},
			TrabajoGrado: &models.TrabajoGrado{
				Id:					12345,
				EstadoTrabajoGrado: 22,
			},
		}

		monkey.Patch(helpers.GetRequestNew, func (baseURL, url string, target interface{}) error {
			switch url {
			case "parametro?query=CodigoAbreviacion:ARC_PLX":
				*target.(*[]models.Parametro) = []models.Parametro{
					{Id: 1, Nombre: "ARL Rechazada"},
				}
			case "parametro?query=CodigoAbreviacion:ACEA_PLX":
				*target.(*[]models.Parametro) = []models.Parametro{
					{Id: 2, Nombre: "ARL Cargada, en espera de aprobación"},
				}
		}
			return nil
		})

		defer monkey.Unpatch(helpers.GetRequestNew)

		monkey.Patch(helpers.SendRequestNew, func(baseURL, url, method string, target interface{}, body interface{}) (string, error) {

			switch url {
			case "/v1/documento_escrito":
				*target.(*map[string]interface{}) = map[string]interface{}{"Id": float64(88)}
			case "/v1/documento_trabajo_grado":
				*target.(*map[string]interface{}) = map[string]interface{}{"Id": float64(99)}
				return "401",assert.AnError
			case "/v1/documento_escrito/88":
				*target.(*map[string]interface{}) = map[string]interface{}{"Id": float64(88)}
				return "200",nil
			}

			return "201",nil
		})

		defer monkey.Unpatch(helpers.SendRequestNew)

		// Valida si se ejecuta el Panic
		defer func() {
			if err := recover(); err != nil {
				t.Logf("Panic en la transacción esperado: %v", err)
			} else {
				t.Errorf("No se ejecuta el Panic esperado")
			}
		}()

		// Llamar a la función
		_ : helpers.AddTransaccionSubirArl(transaccion)

	})

	t.Run("Caso #5: Fallo al momento de actualizar estado Trabajo Grado", func(t *testing.T) {

		transaccion := &models.TrSubirArl{
			DocumentoEscrito: &models.DocumentoEscrito{
				Id:					  0,
				Titulo:               "Pasantía de pruebas unitarias",
				Enlace:               "abcdefg12345",
				Resumen:              "",
				TipoDocumentoEscrito: 12345,
			},
			DocumentoTrabajoGrado: &models.DocumentoTrabajoGrado{
				TrabajoGrado: &models.TrabajoGrado{
					Id:					12345,
				},      
				DocumentoEscrito:&models.DocumentoEscrito{
					Id:					  0,
				},
			},
			TrabajoGrado: &models.TrabajoGrado{
				Id:					12345,
				EstadoTrabajoGrado: 22,
			},
		}

		monkey.Patch(helpers.GetRequestNew, func (baseURL, url string, target interface{}) error {
			switch url {
			case "parametro?query=CodigoAbreviacion:ARC_PLX":
				*target.(*[]models.Parametro) = []models.Parametro{
					{Id: 1, Nombre: "ARL Rechazada"},
				}
			case "parametro?query=CodigoAbreviacion:ACEA_PLX":
				*target.(*[]models.Parametro) = []models.Parametro{
					{Id: 2, Nombre: "ARL Cargada, en espera de aprobación"},
				}
		}
			return nil
		})

		defer monkey.Unpatch(helpers.GetRequestNew)

		monkey.Patch(helpers.SendRequestNew, func(baseURL, url, method string, target interface{}, body interface{}) (string, error) {

			switch url {
			case "/v1/documento_escrito":
				*target.(*map[string]interface{}) = map[string]interface{}{"Id": float64(88)}
			case "/v1/documento_trabajo_grado":
				*target.(*map[string]interface{}) = map[string]interface{}{"Id": float64(99)}
			case "/v1/trabajo_grado/12345":
				*target.(*map[string]interface{}) = map[string]interface{}{"Id": float64(12345)}
				return "404",assert.AnError
			case "/v1/documento_trabajo_grado/99":
				*target.(*map[string]interface{}) = map[string]interface{}{"Id": float64(88)}
				return "200",nil
			case "/v1/documento_escrito/88":
				*target.(*map[string]interface{}) = map[string]interface{}{"Id": float64(88)}
				return "200",nil
			}

			return "201",nil
		})

		defer monkey.Unpatch(helpers.SendRequestNew)

		// Valida si se ejecuta el Panic
		defer func() {
			if err := recover(); err != nil {
				t.Logf("Panic en la transacción esperado: %v", err)
			} else {
				t.Errorf("No se ejecuta el Panic esperado")
			}
		}()

		// Llamar a la función
		_ : helpers.AddTransaccionSubirArl(transaccion)
	})
}