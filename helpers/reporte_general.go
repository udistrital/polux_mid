package helpers

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/udistrital/polux_mid/models"
	"github.com/xuri/excelize/v2"
)

func BuildReporteGeneral() (string, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Error: ", err)
			panic(DeferHelpers("AddTransaccionSolicitud", err))
		}
	}()

	var reporteGeneral []models.ReporteGeneral

	//Se traen todos los datos de reporte general
	url := "/v1/reporte_general"
	if err := GetRequestNew("PoluxCrudUrl", url, &reporteGeneral); err != nil {
		logs.Error("Error al obtener ReporteGeneral")
		panic(err.Error())
	}

	var parametros []models.Parametro

	//Se trae los Estados, la Modalidad del Trabajo de Grado y las Areas de Conocimiento
	url = "parametro?query=TipoParametroId__in:73|76|3|4&limit=0"
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

	//Iterar sobre reporteGeneral y modificar los campos necesarios
	for i, rg := range reporteGeneral {
		if modalidadID, err := strconv.Atoi(rg.Modalidad); err == nil {
			if nombre, ok := parametroMap[modalidadID]; ok {
				reporteGeneral[i].Modalidad = nombre
			}
		}

		if estadoID, err := strconv.Atoi(rg.EstadoTrabajoGrado); err == nil {
			if nombre, ok := parametroMap[estadoID]; ok {
				reporteGeneral[i].EstadoTrabajoGrado = nombre
			}
		}

		if areaID, err := strconv.Atoi(rg.AreaConocimiento); err == nil {
			if nombre, ok := parametroMap[areaID]; ok {
				reporteGeneral[i].AreaConocimiento = nombre
			}
		}

		//Procesar IdEstudiante
		if rg.IdEstudiante != "" {
			if datos, exists := nombresCache[rg.IdEstudiante]; exists {
				//Si el nombre y carrera ya están en el cache, usarlos directamente
				reporteGeneral[i].NombreEstudiante = datos.Nombre
				reporteGeneral[i].ProgramaAcademico = datos.Carrera
			} else {
				//Si no están en el cache, obtenerlos y guardarlos
				datos, err := obtenerDatosEstudiante(rg.IdEstudiante)
				if err != nil {
					logs.Error("Error al obtener datos del estudiante")
					panic(err.Error())
				} else {
					reporteGeneral[i].NombreEstudiante = datos.Nombre
					reporteGeneral[i].ProgramaAcademico = datos.Carrera
					nombresCache[rg.IdEstudiante] = datos // Guardar en el cache
				}
			}
		}

		//Procesar IdCoestudiante (sin modificar ProgramaAcademico)
		if rg.IdCoestudiante != "" {
			if datos, exists := nombresCache[rg.IdCoestudiante]; exists {
				//Si el nombre ya está en el cache, usarlo directamente
				reporteGeneral[i].NombreCoestudiante = datos.Nombre
			} else {
				//Si no están en el cache, obtenerlos y guardarlos
				datos, err := obtenerDatosEstudiante(rg.IdCoestudiante)
				if err != nil {
					logs.Error("Error al obtener datos del coestudiante")
					panic(err.Error())
				} else {
					reporteGeneral[i].NombreCoestudiante = datos.Nombre
					nombresCache[rg.IdCoestudiante] = datos // Guardar en el cache
				}
			}
		}
	}

	//Traer docentes
	docenteMap, err := obtenerDocentes()
	if err != nil {
		logs.Error("Error al obtener docentes")
		panic(err.Error())
	}

	//Mapa para almacenar los nombres de carreras ya consultadas
	carreraCache := make(map[string]string)

	//Hubo la necesidad de iterar nuevamente sobre reporteGeneral, ya que se necesitaba que se añadiera primero el id de la carrera a ProgramaAcademico a través del anterior for, para luego obtener el nombre de la carrera y está fue la única manera (aunque a mi parecer no tan óptima)
	for i, rg := range reporteGeneral {
		//Obtener nombre de la carrera a partir del ID almacenado en ProgramaAcademico
		if rg.ProgramaAcademico != "" {
			if nombreCarrera, exists := carreraCache[rg.ProgramaAcademico]; exists {
				//Si el nombre de la carrera ya está en el cache, usarlo directamente
				reporteGeneral[i].ProgramaAcademico = nombreCarrera
			} else {
				//Si no está en el cache, obtenerlo y guardarlo
				nombreCarrera, err := obtenerNombreCarrera(rg.ProgramaAcademico)
				if err != nil {
					logs.Error("Error al obtener el nombre de la carrera")
					panic(err.Error())
				} else {
					reporteGeneral[i].ProgramaAcademico = nombreCarrera
					carreraCache[rg.ProgramaAcademico] = nombreCarrera // Guardar en el cache
				}
			}
		}

		//Asignar nombres de docentes
		if nombre, exists := docenteMap[rg.DocenteDirector]; exists {
			reporteGeneral[i].NombreDocenteDirector = nombre
		}
		if nombre, exists := docenteMap[rg.DocenteCodirector]; exists {
			reporteGeneral[i].NombreDocenteCodirector = nombre
		}
		if nombre, exists := docenteMap[rg.Evaluador]; exists {
			reporteGeneral[i].NombreEvaluador = nombre
		}
	}

	//Título de las Columnas del Excel
	headers := map[string]string{
		"A1": "Trabajo Grado",
		"B1": "Título",
		"C1": "Modalidad",
		"D1": "Estado",
		"E1": "ID Estudiante",
		"F1": "Nombre Estudiante",
		"G1": "ID Estudiante",
		"H1": "Nombre Estudiante",
		"I1": "Programa Academico",
		"J1": "Area Conocimiento",
		"K1": "ID Docente Director",
		"L1": "Nombre Docente Director",
		"M1": "ID Docente Codirector",
		"N1": "Nombre Docente Codirector",
		"O1": "ID Evaluador",
		"P1": "Nombre Evaluador",
		"Q1": "Fecha Inicio",
		"R1": "Fecha Fin",
		"S1": "Calificación 1",
		"T1": "Calificación 2",
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

	//Precargar los estilos a la hoja de calculo
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

	//Recorrer cada elemento del Slice de ReporteGeneral y escribir en cada fila del Excel su respectiva información
	for i := 0; i < len(reporteGeneral); i++ {
		rowCount := i + 2

		file.SetCellValue("Sheet1", fmt.Sprintf("A%v", rowCount), reporteGeneral[i].TrabajoGrado)
		file.SetCellValue("Sheet1", fmt.Sprintf("B%v", rowCount), reporteGeneral[i].Titulo)
		file.SetCellValue("Sheet1", fmt.Sprintf("C%v", rowCount), reporteGeneral[i].Modalidad)
		file.SetCellValue("Sheet1", fmt.Sprintf("D%v", rowCount), reporteGeneral[i].EstadoTrabajoGrado)
		file.SetCellValue("Sheet1", fmt.Sprintf("E%v", rowCount), reporteGeneral[i].IdEstudiante)
		file.SetCellValue("Sheet1", fmt.Sprintf("F%v", rowCount), reporteGeneral[i].NombreEstudiante)
		file.SetCellValue("Sheet1", fmt.Sprintf("G%v", rowCount), reporteGeneral[i].IdCoestudiante)
		file.SetCellValue("Sheet1", fmt.Sprintf("H%v", rowCount), reporteGeneral[i].NombreCoestudiante)
		file.SetCellValue("Sheet1", fmt.Sprintf("I%v", rowCount), reporteGeneral[i].ProgramaAcademico)
		file.SetCellValue("Sheet1", fmt.Sprintf("J%v", rowCount), reporteGeneral[i].AreaConocimiento)
		file.SetCellValue("Sheet1", fmt.Sprintf("K%v", rowCount), reporteGeneral[i].DocenteDirector)
		file.SetCellValue("Sheet1", fmt.Sprintf("L%v", rowCount), reporteGeneral[i].NombreDocenteDirector)
		file.SetCellValue("Sheet1", fmt.Sprintf("M%v", rowCount), reporteGeneral[i].DocenteCodirector)
		file.SetCellValue("Sheet1", fmt.Sprintf("N%v", rowCount), reporteGeneral[i].NombreDocenteCodirector)
		file.SetCellValue("Sheet1", fmt.Sprintf("O%v", rowCount), reporteGeneral[i].Evaluador)
		file.SetCellValue("Sheet1", fmt.Sprintf("P%v", rowCount), reporteGeneral[i].NombreEvaluador)
		file.SetCellValue("Sheet1", fmt.Sprintf("Q%v", rowCount), reporteGeneral[i].FechaInicio.Format("2006-01-02"))
		file.SetCellValue("Sheet1", fmt.Sprintf("R%v", rowCount), reporteGeneral[i].FechaFin.Format("2006-01-02"))
		file.SetCellValue("Sheet1", fmt.Sprintf("S%v", rowCount), reporteGeneral[i].CalificacionUno)
		file.SetCellValue("Sheet1", fmt.Sprintf("T%v", rowCount), reporteGeneral[i].CalificacionDos)
	}

	//Guardar el archivo Excel en este caso en la raíz del proyecto
	/*if err := file.SaveAs("ReporteGeneral.xlsx"); err != nil {
		fmt.Println(err)
	}*/

	//Guardar el archivo en memoria
	var buffer bytes.Buffer
	if err := file.Write(&buffer); err != nil {
		logs.Error("Error al escribir archivo en buffer")
		panic(err.Error())
	}

	//Codificar el archivo en Base64
	encodedFile := base64.StdEncoding.EncodeToString(buffer.Bytes())

	return encodedFile, nil
}

func obtenerDatosEstudiante(idEstudiante string) (models.DatosBasicosEstudiante, error) {
	url := fmt.Sprintf("http://busservicios.intranetoas.udistrital.edu.co:8282/wso2eiserver/services/servicios_academicos/datos_basicos_estudiante/%s", idEstudiante)

	resp, err := http.Get(url)
	if err != nil {
		return models.DatosBasicosEstudiante{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.DatosBasicosEstudiante{}, fmt.Errorf("Error al realizar la solicitud: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.DatosBasicosEstudiante{}, err
	}

	var datos models.DatosEstudianteCollection
	if err := xml.Unmarshal(body, &datos); err != nil {
		return models.DatosBasicosEstudiante{}, err
	}

	if len(datos.DatosBasicosEstudiante) > 0 {
		return datos.DatosBasicosEstudiante[0], nil
	}

	return models.DatosBasicosEstudiante{}, fmt.Errorf("No se encontraron datos para el estudiante %s", idEstudiante)
}

func obtenerNombreCarrera(idCarrera string) (string, error) {
	url := fmt.Sprintf("http://busservicios.intranetoas.udistrital.edu.co:8282/wso2eiserver/services/servicios_academicos/carrera/%s", idCarrera)

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

	var carreraCollection struct {
		Carrera struct {
			Nombre string `xml:"nombre"`
		} `xml:"carrera"`
	}

	if err := xml.Unmarshal(body, &carreraCollection); err != nil {
		return "", err
	}

	return carreraCollection.Carrera.Nombre, nil
}

func obtenerDocentes() (map[int]string, error) {
	url := "http://busservicios.intranetoas.udistrital.edu.co:8282/wso2eiserver/services/servicios_academicos/get_docentes_tg"

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error al realizar la solicitud: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var docentes struct {
		Docentes []struct {
			ID     int    `xml:"id"`
			Nombre string `xml:"NOMBRE"`
		} `xml:"docente"`
	}

	if err := xml.Unmarshal(body, &docentes); err != nil {
		return nil, err
	}

	docenteMap := make(map[int]string)
	for _, docente := range docentes.Docentes {
		docenteMap[docente.ID] = docente.Nombre
	}

	return docenteMap, nil
}
