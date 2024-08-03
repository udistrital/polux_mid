package helpers

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/udistrital/polux_mid/models"
	"github.com/xuri/excelize/v2"
)

func BuildReporteGeneral() error {
	var reporteGeneral []models.ReporteGeneral

	//Se traen todos los datos de reporte general
	url := "/v1/reporte_general"
	if err := GetRequestNew("PoluxCrudUrl", url, &reporteGeneral); err != nil {
		logs.Error("Error al obtener ReporteGeneral: ", err.Error())
		return err
	}

	var parametros []models.Parametro

	//Se trae los Estados, la Modalidad del Trabajo de Grado y las Areas de Conocimiento
	url = "parametro?query=TipoParametroId__in:73|76|3|4&limit=0"
	if err := GetRequestNew("UrlCrudParametros", url, &parametros); err != nil {
		logs.Error("Error al obtener Parametros: ", err.Error())
		return err
	}

	//Crear un mapa de parámetros para facilitar la búsqueda
	parametroMap := make(map[int]string)
	for _, parametro := range parametros {
		parametroMap[parametro.Id] = parametro.Nombre
	}

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
	}

	//Título de las Columnas del Excel
	headers := map[string]string{
		"A1": "Trabajo Grado",
		"B1": "Título",
		"C1": "Modalidad",
		"D1": "Estado",
		"E1": "Estudiante 1",
		"F1": "Estudiante 2",
		"G1": "Area Conocimiento",
		"H1": "Docente Director",
		"I1": "Docente Codirector",
		"J1": "Evaluador",
		"K1": "Fecha Inicio",
		"L1": "Fecha Fin",
		"M1": "Calificación 1",
		"N1": "Calificación 2",
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
		fmt.Println(err)
		return err
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
		file.SetCellValue("Sheet1", fmt.Sprintf("F%v", rowCount), reporteGeneral[i].IdCoestudiante)
		file.SetCellValue("Sheet1", fmt.Sprintf("G%v", rowCount), reporteGeneral[i].AreaConocimiento)
		file.SetCellValue("Sheet1", fmt.Sprintf("H%v", rowCount), reporteGeneral[i].DocenteDirector)
		file.SetCellValue("Sheet1", fmt.Sprintf("I%v", rowCount), reporteGeneral[i].DocenteCodirector)
		file.SetCellValue("Sheet1", fmt.Sprintf("J%v", rowCount), reporteGeneral[i].Evaluador)
		file.SetCellValue("Sheet1", fmt.Sprintf("K%v", rowCount), reporteGeneral[i].FechaInicio.Format("2006-01-02"))
		file.SetCellValue("Sheet1", fmt.Sprintf("L%v", rowCount), reporteGeneral[i].FechaFin.Format("2006-01-02"))
		file.SetCellValue("Sheet1", fmt.Sprintf("M%v", rowCount), reporteGeneral[i].CalificacionUno)
		file.SetCellValue("Sheet1", fmt.Sprintf("N%v", rowCount), reporteGeneral[i].CalificacionDos)
	}

	//Guardar el archivo Excel en este caso en la raíz del proyecto
	if err := file.SaveAs("ReporteGeneral.xlsx"); err != nil {
		fmt.Println(err)
	}

	return nil
}
