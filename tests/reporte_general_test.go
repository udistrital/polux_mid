package test

import (
	"encoding/base64"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
	"github.com/udistrital/polux_mid/helpers"
	"github.com/udistrital/polux_mid/models"
)

func TestBuildReporteGeneral(t *testing.T) {
	t.Log("//////////////////////////////////")
	t.Log("Inicio TestBuildReporteGeneral")
	t.Log("//////////////////////////////////")

	t.Run("Caso 1: Generación exitosa del reporte", func(t *testing.T) {
		// Mock de dependencias
		monkey.Patch(helpers.GetRequestNew, func(baseURL, url string, target interface{}) error {
			switch url {
			case "parametro?query=CodigoAbreviacion:CNC_PLX":
				*target.(*[]models.Parametro) = []models.Parametro{
					{Id: 1, Nombre: "Cancelado"},
				}
			case "parametro?query=CodigoAbreviacion:NTF_PLX":
				*target.(*[]models.Parametro) = []models.Parametro{
					{Id: 2, Nombre: "Notificado"},
				}
			case "parametro?query=TipoParametroId__in:73|76|3|4&limit=0":
				*target.(*[]models.Parametro) = []models.Parametro{
					{Id: 3, Nombre: "Modalidad"},
				}
			}
			return nil
		})
		defer monkey.Unpatch(helpers.GetRequestNew)

		monkey.Patch(helpers.SendRequestNew, func(baseURL, url, method string, target interface{}, body interface{}) (string, error) {
			*target.(*[]models.ReporteGeneral) = []models.ReporteGeneral{
				{
					TrabajoGrado:       1,
					Titulo:             "Trabajo de Grado 1",
					Modalidad:          "3",
					EstadoTrabajoGrado: "1",
					AreaConocimiento:   "2",
					IdEstudiante:       "123",
					FechaInicio:        time.Now(),
					FechaFin:           time.Now(),
				},
			}
			return "201", nil
		})
		defer monkey.Unpatch(helpers.SendRequestNew)

		monkey.Patch(helpers.ObtenerDatosEstudiante, func(idEstudiante string) (models.DatosBasicosEstudiante, error) {
			if idEstudiante == "123" {
				return models.DatosBasicosEstudiante{
					Nombre:  "Juan Pérez",
					Carrera: "Ingeniería de Sistemas",
				}, nil
			}
			return models.DatosBasicosEstudiante{}, assert.AnError
		})
		defer monkey.Unpatch(helpers.ObtenerDatosEstudiante)

		// Actualización del mock de obtenerNombreCarrera
		monkey.Patch(helpers.ObtenerNombreCarrera, func(idCarrera string) (string, error) {
			// Simular una respuesta válida para cualquier ID
			switch idCarrera {
			case "2", "3":
				return "Ingeniería de Sistemas", nil
			default:
				return "Carrera Desconocida", nil
			}
		})
		defer monkey.Unpatch(helpers.ObtenerNombreCarrera)

		// Input de prueba
		filtros := &models.FiltrosReporte{
			ProyectoCurricular: "Ingeniería de Sistemas",
		}

		// Llamar a la función
		result, err := helpers.BuildReporteGeneral(filtros)

		// Validar resultados
		assert.NoError(t, err)
		assert.NotEmpty(t, result)

		// Validar Base64
		decoded, err := base64.StdEncoding.DecodeString(result)
		assert.NoError(t, err)
		assert.NotEmpty(t, decoded)
	})

	t.Run("Caso 2: Estudiante no encontrado", func(t *testing.T) {
		monkey.Patch(helpers.ObtenerDatosEstudiante, func(idEstudiante string) (models.DatosBasicosEstudiante, error) {
			return models.DatosBasicosEstudiante{}, assert.AnError
		})
		defer monkey.Unpatch(helpers.ObtenerDatosEstudiante)

		datos, err := helpers.ObtenerDatosEstudiante("999")
		assert.Error(t, err)
		assert.Empty(t, datos.Nombre)
	})
}

func TestObtenerNombreCarrera(t *testing.T) {
	t.Log("//////////////////////////////////")
	t.Log("Inicio TestObtenerNombreCarrera")
	t.Log("//////////////////////////////////")

	t.Run("Caso 1: Obtención exitosa", func(t *testing.T) {
		monkey.Patch(helpers.ObtenerNombreCarrera, func(idCarrera string) (string, error) {
			return "Ingeniería de Sistemas", nil
		})
		defer monkey.Unpatch(helpers.ObtenerNombreCarrera)

		nombre, err := helpers.ObtenerNombreCarrera("3")
		assert.NoError(t, err)
		assert.Equal(t, "Ingeniería de Sistemas", nombre)
	})

	t.Run("Caso 2: Carrera no encontrada", func(t *testing.T) {
		monkey.Patch(helpers.ObtenerNombreCarrera, func(idCarrera string) (string, error) {
			return "", assert.AnError
		})
		defer monkey.Unpatch(helpers.ObtenerNombreCarrera)

		nombre, err := helpers.ObtenerNombreCarrera("999")
		assert.Error(t, err)
		assert.Empty(t, nombre)
	})
}

func TestObtenerDocentes(t *testing.T) {
	t.Log("//////////////////////////////////")
	t.Log("Inicio TestObtenerDocentes")
	t.Log("//////////////////////////////////")

	t.Run("Caso 1: Obtención exitosa", func(t *testing.T) {
		// Parchear correctamente la función ObtenerDocentes
		monkey.Patch(helpers.ObtenerDocentes, func() (map[int]string, error) {
			return map[int]string{
				1: "Docente Director",
				2: "Codirector",
			}, nil
		})
		defer monkey.Unpatch(helpers.ObtenerDocentes) // Deshacer el parche al final del caso

		// Llamar a la función y validar los resultados
		docentes, err := helpers.ObtenerDocentes()
		assert.NoError(t, err)
		assert.Equal(t, "Docente Director", docentes[1])
		assert.Equal(t, "Codirector", docentes[2])
	})

	t.Run("Caso 2: Error en la obtención de docentes", func(t *testing.T) {
		// Simular error en la función ObtenerDocentes
		monkey.Patch(helpers.ObtenerDocentes, func() (map[int]string, error) {
			return nil, assert.AnError
		})
		defer monkey.Unpatch(helpers.ObtenerDocentes) // Deshacer el parche al final del caso

		// Llamar a la función y validar el error
		docentes, err := helpers.ObtenerDocentes()
		assert.Error(t, err)
		assert.Nil(t, docentes)
	})
}
