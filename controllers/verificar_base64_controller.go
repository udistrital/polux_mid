package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/udistrital/polux_mid/models"
	"github.com/udistrital/polux_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
)

// VerificarBase64Controller operations for VerficarBase64
type VerificarBase64Controller struct {
	beego.Controller
}

// URLMapping ...
func (c *VerificarBase64Controller) URLMapping() {
	c.Mapping("PostVerificarBase64", c.PostVerificarBase64)
}

// PostVerificarBase64 ...
// @Title PostVerificarBase64
// @Description Verifica si un PDF base64 está limpio usando ClamAV
// @Param	body		body 	models.EmailAttachment	true		"Base64 del PDF"
// @Success 200 {object} map[string]interface{}
// @Failure 404 body is empty
// @router / [post]
func (c *VerificarBase64Controller) PostVerificarBase64() {
	defer errorhandler.HandlePanic(&c.Controller)

	var archivos []models.EmailAttachment
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &archivos); err != nil {
		beego.Error(err)
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, "Error al decodificar el cuerpo de la solicitud: "+err.Error())
		c.ServeJSON()
		return
	}

	if len(archivos) == 0 {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, "El array de archivos está vacío")
		c.ServeJSON()
		return
	}

	archivo := archivos[0]

	if archivo.PdfBase64 == "" || archivo.UrlFileUp == "" {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, "Los campos pdf_base64 y urlFileUp son obligatorios")
		c.ServeJSON()
		return
	}

	resultadoClamAV := services.VerificarBase64(archivo.PdfBase64)
	if !resultadoClamAV.Success {
		c.Ctx.Output.SetStatus(resultadoClamAV.Status)
		c.Data["json"] = resultadoClamAV
		c.ServeJSON()
		return
	}

	var virusResult map[string]interface{}
	if clamAVData, ok := resultadoClamAV.Data.(map[string]interface{}); ok {
		if virusData, ok := clamAVData["Virus"].(map[string]interface{}); ok {
			virusResult = virusData
		} else {
			virusResult = map[string]interface{}{
				"message":    "Verificación de virus completada correctamente.",
				"archive":    clamAVData["status"],
				"statusCode": 200,
			}
		}
	} else {
		virusResult = map[string]interface{}{
			"message":    "Error al procesar el archivo",
			"archive":    "unknown",
			"statusCode": 500,
		}
	}

	mensajeFinal := "El archivo PDF está limpio."
	if virusResult["archive"] == "infected" {
		mensajeFinal = "El archivo contiene virus."
	}

	dataFinal := map[string]interface{}{
		"Virus": virusResult,
	}

	respuestaFinal := requestresponse.APIResponseDTO(
		true,
		200,
		dataFinal,
		mensajeFinal,
	)

	c.Ctx.Output.SetStatus(200)
	c.Data["json"] = respuestaFinal
	c.ServeJSON()

}
