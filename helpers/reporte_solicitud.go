package helpers

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/udistrital/polux_mid/models"
	"github.com/xuri/excelize/v2"
)

func BuildReporteSolicitud(filtros *models.FiltrosReporte) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Error: ", err)
			panic(DeferHelpers("AddTransaccionSolicitud", err))
		}
	}()

	var reporteSolicitud []models.ReporteSolicitud

	//Se consultan los Estado de Solicitud "Aprobada por consejo de carrera" y "Rechazada por consejo de carrera" para aplicar el filtro de Estado
	var parametroEstadoTg []models.Parametro
	var url = "parametro?query=CodigoAbreviacion:RCC_PLX"
	if err := GetRequestNew("UrlCrudParametros", url, &parametroEstadoTg); err != nil {
		//logs.Error(err.Error())
		panic(err.Error())
	}

	filtros.IdEstCancelado = parametroEstadoTg[0].Id

	url = "parametro?query=CodigoAbreviacion:ACC_PLX"
	if err := GetRequestNew("UrlCrudParametros", url, &parametroEstadoTg); err != nil {
		//logs.Error(err.Error())
		panic(err.Error())
	}

	filtros.IdEstFinalizado = parametroEstadoTg[0].Id

	//Se traen todos los datos de reporte solicitud del CRUD
	url = "/v1/reporte_solicitud"
	if status, err := SendRequestNew("PoluxCrudUrl", url, "POST", &reporteSolicitud, &filtros); err != nil && status == "201" {
		logs.Error("Error al obtener ReporteSolicitud")
		panic(err.Error())
	}

	var parametros []models.Parametro

	//Se trae los Estados, la Modalidades, los Tipo Solicitud y los Estados de Solicitud de Trabajo de Grado de Parametros
	url = "parametro?query=TipoParametroId__in:73|76|77|78&limit=0"
	if err := GetRequestNew("UrlCrudParametros", url, &parametros); err != nil {
		logs.Error("Error al obtener Parametros")
		panic(err.Error())
	}

	//Crear un mapa de parámetros para facilitar la búsqueda
	parametroMap := make(map[int]string)
	for _, parametro := range parametros {
		parametroMap[parametro.Id] = parametro.Nombre
	}

	//Mapa para almacenar los nombres y carreras de estudiantes ya consultados
	nombresCache := make(map[string]models.DatosBasicosEstudiante)

	//Iterar sobre reporteSolicitud y modificar los campos necesarios
	for i, rs := range reporteSolicitud {
		if modalidadID, err := strconv.Atoi(rs.Modalidad); err == nil {
			if nombre, ok := parametroMap[modalidadID]; ok {
				reporteSolicitud[i].Modalidad = nombre
			}
		}

		if estadoID, err := strconv.Atoi(rs.EstadoTrabajoGrado); err == nil {
			if nombre, ok := parametroMap[estadoID]; ok {
				reporteSolicitud[i].EstadoTrabajoGrado = nombre
			}
		}

		if tipoSolicitudID, err := strconv.Atoi(rs.Solicitud); err == nil {
			if nombre, ok := parametroMap[tipoSolicitudID]; ok {
				reporteSolicitud[i].Solicitud = nombre
			}
		}

		if estadoSolicitudID, err := strconv.Atoi(rs.Respuesta); err == nil {
			if nombre, ok := parametroMap[estadoSolicitudID]; ok {
				reporteSolicitud[i].Respuesta = nombre
			}
		}

		//Procesar IdEstudiante
		if rs.IdEstudiante != "" {
			if datos, exists := nombresCache[rs.IdEstudiante]; exists {
				//Si el nombre y carrera ya están en el cache, usarlos directamente
				reporteSolicitud[i].NombreEstudiante = datos.Nombre
				reporteSolicitud[i].ProgramaAcademico = datos.Carrera
			} else {
				//Si no están en el cache, obtenerlos y guardarlos
				datos, err := obtenerDatosEstudiante(rs.IdEstudiante)
				if err != nil {
					logs.Error("Error al obtener datos del estudiante")
					panic(err.Error())
				} else {
					reporteSolicitud[i].NombreEstudiante = datos.Nombre
					reporteSolicitud[i].ProgramaAcademico = datos.Carrera
					nombresCache[rs.IdEstudiante] = datos // Guardar en el cache
				}
			}
		}

		//Procesar IdCoestudiante (sin modificar ProgramaAcademico)
		if rs.IdCoestudiante != "" {
			if datos, exists := nombresCache[rs.IdCoestudiante]; exists {
				//Si el nombre ya está en el cache, usarlo directamente
				reporteSolicitud[i].NombreCoestudiante = datos.Nombre
			} else {
				//Si no están en el cache, obtenerlos y guardarlos
				datos, err := obtenerDatosEstudiante(rs.IdCoestudiante)
				if err != nil {
					logs.Error("Error al obtener datos del coestudiante")
					panic(err.Error())
				} else {
					reporteSolicitud[i].NombreCoestudiante = datos.Nombre
					nombresCache[rs.IdCoestudiante] = datos // Guardar en el cache
				}
			}
		}
	}

	var reporteSolicitudFiltered []models.ReporteSolicitud

	//Se aplica el filtro de Proyecto Curricular creando un nuevo slice
	for _, rs := range reporteSolicitud {
		if rs.ProgramaAcademico == filtros.ProyectoCurricular {
			reporteSolicitudFiltered = append(reporteSolicitudFiltered, rs)
		}
	}

	//si el reporte traido desde el crud o el reporte filtrado por carrera no tiene registros, se ejecuta un error 400
	if(len(reporteSolicitud) <= 0 || len(reporteSolicitudFiltered) <= 0){
		return "", errors.New("no se encontraron registros")
	}

	//Traer docentes
	docenteMap, err := obtenerDocentes()
	if err != nil {
		logs.Error("Error al obtener docentes")
		panic(err.Error())
	}

	//Mapa para almacenar los nombres de carreras ya consultadas
	coordinadorCache := make(map[string]string)

	//Mapa para almacenar los nombres de carreras ya consultadas
	carreraCache := make(map[string]string)

	//Hubo la necesidad de iterar nuevamente sobre reporteSolicitud, ya que se necesitaba que se añadiera primero el id de la carrera a ProgramaAcademico a través del anterior for, para luego obtener el nombre de la carrera y está fue la única manera (aunque a mi parecer no tan óptima)
	for i, rs := range reporteSolicitudFiltered {
		//Obtener nombre del coordinador a partir del ID almacenado en ProgramaAcademico
		if rs.ProgramaAcademico != "" {
			if nombreCoordinador, exists := coordinadorCache[rs.ProgramaAcademico]; exists {
				//Si el nombre de la carrera ya está en el cache, usarlo directamente
				reporteSolicitudFiltered[i].NombreCoordinador = nombreCoordinador
			} else {
				//Si no está en el cache, obtenerlo y guardarlo
				nombreCoordinador, err := obtenerNombreCoordinador(rs.ProgramaAcademico)
				if err != nil {
					logs.Error("Error al obtener el nombre de la carrera")
					panic(err.Error())
				} else {
					reporteSolicitudFiltered[i].NombreCoordinador = nombreCoordinador
					coordinadorCache[rs.ProgramaAcademico] = nombreCoordinador // Guardar en el cache
				}
			}
		}

		//Obtener nombre de la carrera a partir del ID almacenado en ProgramaAcademico
		if rs.ProgramaAcademico != "" {
			if nombreCarrera, exists := carreraCache[rs.ProgramaAcademico]; exists {
				//Si el nombre de la carrera ya está en el cache, usarlo directamente
				reporteSolicitudFiltered[i].ProgramaAcademico = nombreCarrera
			} else {
				//Si no está en el cache, obtenerlo y guardarlo
				nombreCarrera, err := obtenerNombreCarrera(rs.ProgramaAcademico)
				if err != nil {
					logs.Error("Error al obtener el nombre de la carrera")
					panic(err.Error())
				} else {
					reporteSolicitudFiltered[i].ProgramaAcademico = nombreCarrera
					carreraCache[rs.ProgramaAcademico] = nombreCarrera // Guardar en el cache
				}
			}
		}

		//Asignar nombres de docentes
		if nombre, exists := docenteMap[rs.DocenteDirector]; exists {
			reporteSolicitudFiltered[i].NombreDocenteDirector = nombre
		}
		if nombre, exists := docenteMap[rs.DocenteCodirector]; exists {
			reporteSolicitudFiltered[i].NombreDocenteCodirector = nombre
		}
		if nombre, exists := docenteMap[rs.Evaluador]; exists {
			reporteSolicitudFiltered[i].NombreEvaluador = nombre
		}
	}

	//Título de las Columnas del Excel
	headers := map[string]string{
		"A1": "ID Solicitud",
		"B1": "Trabajo Grado",
		"C1": "Título",
		"D1": "Modalidad",
		"E1": "Estado Trabajo Grado",
		"F1": "ID Estudiante",
		"G1": "Nombre Estudiante",
		"H1": "ID Estudiante",
		"I1": "Nombre Estudiante",
		"J1": "Programa Academico",
		"K1": "Nombre Coordinador",
		"L1": "ID Docente Director",
		"M1": "Nombre Docente Director",
		"N1": "ID Docente Codirector",
		"O1": "Nombre Docente Codirector",
		"P1": "ID Evaluador",
		"Q1": "Nombre Evaluador",
		"R1": "Fecha Solicitud",
		"S1": "Fecha Revision",
		"T1": "Concepto de Revision",
		"U1": "Observaciones",
		"V1": "Respuesta",
	}

	//Creación e Inicialización del Excel
	file := excelize.NewFile()

	//Definir el estilo para los Encabezados
	style := &excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#FFFF00"},
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	}

	//Precarsar los estilos a la hoja de calculo
	styleID, err := file.NewStyle(style)
	if err != nil {
		logs.Error("Error al cargar los estilos a la hoja de calculo")
		panic(err.Error())
	}

	//Recorrer los headers y añadir a la hoja de cálculo del Excel
	for k, v := range headers {
		file.SetCellValue("Sheet1", k, v)
		file.SetCellStyle("Sheet1", k, k, styleID) // Aplicar el estilo de fondo amarillo
	}

	//Recorrer cada elemento del Slice de ReporteSolicitud y escribir en cada fila del Excel su respectiva información
	for i := 0; i < len(reporteSolicitudFiltered); i++ {
		rowCount := i + 2

		file.SetCellValue("Sheet1", fmt.Sprintf("A%v", rowCount), reporteSolicitudFiltered[i].Id)
		file.SetCellValue("Sheet1", fmt.Sprintf("B%v", rowCount), reporteSolicitudFiltered[i].TrabajoGrado)
		file.SetCellValue("Sheet1", fmt.Sprintf("C%v", rowCount), reporteSolicitudFiltered[i].Titulo)
		file.SetCellValue("Sheet1", fmt.Sprintf("D%v", rowCount), reporteSolicitudFiltered[i].Modalidad)
		file.SetCellValue("Sheet1", fmt.Sprintf("E%v", rowCount), reporteSolicitudFiltered[i].EstadoTrabajoGrado)
		file.SetCellValue("Sheet1", fmt.Sprintf("F%v", rowCount), reporteSolicitudFiltered[i].IdEstudiante)
		file.SetCellValue("Sheet1", fmt.Sprintf("G%v", rowCount), reporteSolicitudFiltered[i].NombreEstudiante)
		file.SetCellValue("Sheet1", fmt.Sprintf("H%v", rowCount), reporteSolicitudFiltered[i].IdCoestudiante)
		file.SetCellValue("Sheet1", fmt.Sprintf("I%v", rowCount), reporteSolicitudFiltered[i].NombreCoestudiante)
		file.SetCellValue("Sheet1", fmt.Sprintf("J%v", rowCount), reporteSolicitudFiltered[i].ProgramaAcademico)
		file.SetCellValue("Sheet1", fmt.Sprintf("K%v", rowCount), reporteSolicitudFiltered[i].NombreCoordinador)
		file.SetCellValue("Sheet1", fmt.Sprintf("L%v", rowCount), reporteSolicitudFiltered[i].DocenteDirector)
		file.SetCellValue("Sheet1", fmt.Sprintf("M%v", rowCount), reporteSolicitudFiltered[i].NombreDocenteDirector)
		file.SetCellValue("Sheet1", fmt.Sprintf("N%v", rowCount), reporteSolicitudFiltered[i].DocenteCodirector)
		file.SetCellValue("Sheet1", fmt.Sprintf("O%v", rowCount), reporteSolicitudFiltered[i].NombreDocenteCodirector)
		file.SetCellValue("Sheet1", fmt.Sprintf("P%v", rowCount), reporteSolicitudFiltered[i].Evaluador)
		file.SetCellValue("Sheet1", fmt.Sprintf("Q%v", rowCount), reporteSolicitudFiltered[i].NombreEvaluador)
		file.SetCellValue("Sheet1", fmt.Sprintf("R%v", rowCount), reporteSolicitudFiltered[i].FechaSolicitud.Format("2006-01-02"))
		file.SetCellValue("Sheet1", fmt.Sprintf("S%v", rowCount), reporteSolicitudFiltered[i].FechaRevision.Format("2006-01-02"))
		file.SetCellValue("Sheet1", fmt.Sprintf("T%v", rowCount), reporteSolicitudFiltered[i].Solicitud)
		file.SetCellValue("Sheet1", fmt.Sprintf("U%v", rowCount), reporteSolicitudFiltered[i].Observacion)
		file.SetCellValue("Sheet1", fmt.Sprintf("V%v", rowCount), reporteSolicitudFiltered[i].Respuesta)
	}

	//Guardar el archivo Excel en este caso en la raíz del proyecto
	/*if err := file.SaveAs("ReporteSolicitud.xlsx"); err != nil {
		fmt.Println(err)
	}*/

	//Guardar el archivo en memoria
	var buffer bytes.Buffer
	if err := file.Write(&buffer); err != nil {
		logs.Error("Error al escribir archivo en buffer")
		panic(err.Error())
	}

	// Codificar el archivo en Base64
	encodedFile := base64.StdEncoding.EncodeToString(buffer.Bytes())

	return encodedFile, nil
}

func obtenerNombreCoordinador(idCarrera string) (string, error) {
	url := fmt.Sprintf("http://busservicios.intranetoas.udistrital.edu.co:8282/wso2eiserver/services/servicios_academicos/coordinador_proyecto/%s", idCarrera)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Error al realizar la solicitud: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var coordinadorCollection struct {
		Coordinador struct {
			Nombre string `xml:"nombre_coordinador"`
		} `xml:"coordinador"`
	}

	if err := xml.Unmarshal(body, &coordinadorCollection); err != nil {
		return "", err
	}

	return coordinadorCollection.Coordinador.Nombre, nil
}
