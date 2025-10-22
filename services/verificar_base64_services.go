package services

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
)

type LambdaResponse struct {
	Status    string `json:"status"`
	RawOutput string `json:"raw_output"`
}

type LambdaRawResponse struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}

func VerificarBase64(pdfBase64 string) requestresponse.APIResponse {
	if pdfBase64 == "" {
		return requestresponse.APIResponseDTO(false, 400, nil, "El contenido PDF base64 está vacío")
	}

	payload := map[string]string{"pdf_base64": pdfBase64}
	var rawResponse LambdaRawResponse

	urlEscanear := "https://" + beego.AppConfig.String("EscanearArchivo") + "verificar"
	err := request.SendJson(
		urlEscanear,
		"POST", &rawResponse, payload,
	)
	if err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, "Error al llamar a Lambda: "+err.Error())
	}
	var lambdaResult LambdaResponse
	if err := json.Unmarshal([]byte(rawResponse.Body), &lambdaResult); err != nil {
		return requestresponse.APIResponseDTO(false, 500, nil, "Error al decodificar cuerpo de respuesta del antivirus: "+err.Error())
	}

	if lambdaResult.Status != "clean" && lambdaResult.Status != "infected" {
		return requestresponse.APIResponseDTO(false, 500, nil, "Respuesta inválida del antivirus")
	}

	return requestresponse.APIResponseDTO(true, 200, map[string]interface{}{
		"Virus": map[string]interface{}{
			"message":    "Verificación de virus completada correctamente.",
			"archive":    lambdaResult.Status,
			"statusCode": 200,
		},
	}, nil)
}
