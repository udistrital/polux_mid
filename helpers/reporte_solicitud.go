package helpers

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/udistrital/polux_mid/models"
	"github.com/xuri/excelize/v2"
)

func BuildReporteSolicitud() error {
	var reporteSolicitud []models.ReporteSolicitud

	//Se traen todos los datos de reporte solicitud del CRUD
	url := "/v1/reporte_solicitud"
	if err := GetRequestNew("PoluxCrudUrl", url, &reporteSolicitud); err != nil {
		logs.Error("Error al obtener ReporteSolicitud: ", err.Error())
		return err
	}

	var parametros []models.Parametro

	//Se trae los Estados, la Modalidades, los Tipo Solicitud y los Estados de Solicitud de Trabajo de Grado de Parametros
	url = "parametro?query=TipoParametroId__in:73|76|77|78&limit=0"
	if err := GetRequestNew("UrlCrudParametros", url, &parametros); err != nil {
		logs.Error("Error al obtener Parametros: ", err.Error())
		return err
	}

	//Crear un mapa de parámetros para facilitar la búsqueda
	parametroMap := make(map[int]string)
	for _, parametro := range parametros {
		parametroMap[parametro.Id] = parametro.Nombre
	}

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
	}

	//Título de las Columnas del Excel
	headers := map[string]string{
		"A1": "ID Solicitud",
		"B1": "Trabajo Grado",
		"C1": "Título",
		"D1": "Modalidad",
		"E1": "Estado Trabajo Grado",
		"F1": "Estudiante 1",
		"G1": "Estudiante 2",
		"H1": "Programa Academico",
		"I1": "Coordinador",
		"J1": "Docente Director",
		"K1": "Docente Codirector",
		"L1": "Evaluador",
		"M1": "Fecha Solicitud",
		"N1": "Fecha Revision",
		"O1": "Concepto de Revision",
		"P1": "Observaciones",
		"Q1": "Respuesta",
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
		fmt.Println(err)
		return err
	}

	//Recorrer los headers y añadir a la hoja de cálculo del Excel
	for k, v := range headers {
		file.SetCellValue("Sheet1", k, v)
		file.SetCellStyle("Sheet1", k, k, styleID) // Aplicar el estilo de fondo amarillo
	}

	//Recorrer cada elemento del Slice de ReporteSolicitud y escribir en cada fila del Excel su respectiva información
	for i := 0; i < len(reporteSolicitud); i++ {
		rowCount := i + 2

		file.SetCellValue("Sheet1", fmt.Sprintf("A%v", rowCount), reporteSolicitud[i].Id)
		file.SetCellValue("Sheet1", fmt.Sprintf("B%v", rowCount), reporteSolicitud[i].TrabajoGrado)
		file.SetCellValue("Sheet1", fmt.Sprintf("C%v", rowCount), reporteSolicitud[i].Titulo)
		file.SetCellValue("Sheet1", fmt.Sprintf("D%v", rowCount), reporteSolicitud[i].Modalidad)
		file.SetCellValue("Sheet1", fmt.Sprintf("E%v", rowCount), reporteSolicitud[i].EstadoTrabajoGrado)
		file.SetCellValue("Sheet1", fmt.Sprintf("F%v", rowCount), reporteSolicitud[i].IdEstudiante)
		file.SetCellValue("Sheet1", fmt.Sprintf("G%v", rowCount), reporteSolicitud[i].IdCoestudiante)
		file.SetCellValue("Sheet1", fmt.Sprintf("H%v", rowCount), reporteSolicitud[i].ProgramaAcademico)
		file.SetCellValue("Sheet1", fmt.Sprintf("I%v", rowCount), reporteSolicitud[i].Coordinador)
		file.SetCellValue("Sheet1", fmt.Sprintf("J%v", rowCount), reporteSolicitud[i].DocenteDirector)
		file.SetCellValue("Sheet1", fmt.Sprintf("K%v", rowCount), reporteSolicitud[i].DocenteCodirector)
		file.SetCellValue("Sheet1", fmt.Sprintf("L%v", rowCount), reporteSolicitud[i].Evaluador)
		file.SetCellValue("Sheet1", fmt.Sprintf("M%v", rowCount), reporteSolicitud[i].FechaSolicitud.Format("2006-01-02"))
		file.SetCellValue("Sheet1", fmt.Sprintf("N%v", rowCount), reporteSolicitud[i].FechaRevision.Format("2006-01-02"))
		file.SetCellValue("Sheet1", fmt.Sprintf("O%v", rowCount), reporteSolicitud[i].Solicitud)
		file.SetCellValue("Sheet1", fmt.Sprintf("P%v", rowCount), reporteSolicitud[i].Observacion)
		file.SetCellValue("Sheet1", fmt.Sprintf("Q%v", rowCount), reporteSolicitud[i].Respuesta)
	}

	//Guardar el archivo Excel en este caso en la raíz del proyecto
	if err := file.SaveAs("ReporteSolicitud.xlsx"); err != nil {
		fmt.Println(err)
	}

	return nil
}
