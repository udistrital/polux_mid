package helpers

import (
	"fmt"

	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/request"
)

type Origin interface {
	BasePath() string
	Endpoint() string
}

func Get(o Origin, payload string, response interface{}) (outputError map[string]interface{}) {

	funcion := "Get - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlCRUD := "http://" + o.BasePath() + o.Endpoint() + payload

	APIresponse, err := request.GetJsonTest(urlCRUD, &response)
	if err != nil || APIresponse.StatusCode != 200 {
		logs.Error(urlCRUD, err)
		eval := "request.GetJsonTest(urlCRUD, &response)"
		outputError = errorctrl.Error(funcion+eval, err, "502")
	}

	return

}

func GetXML(o Origin, payload string, response interface{}) (outputError map[string]interface{}) {

	funcion := "GetXML - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlCRUD := "http://" + o.BasePath() + o.Endpoint() + payload

	err := request.GetJsonWSO2(urlCRUD, &response)
	if err != nil {
		logs.Error(urlCRUD, err)
		eval := "request.GetJsonWSO2(urlCRUD, &response)"
		outputError = errorctrl.Error(funcion+eval, err, "502")
	}

	return

}

func Put(o Origin, id int, data, response interface{}) (outputError map[string]interface{}) {

	funcion := "Put - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlCRUD := "http://" + o.BasePath() + o.Endpoint() + "/" + fmt.Sprint(id)

	err := request.SendJson(urlCRUD, "PUT", &response, &data)
	if err != nil {
		logs.Error(urlCRUD, err)
		eval := `request.SendJson(urlCRUD, "PUT", &response, &data)`
		outputError = errorctrl.Error(funcion+eval, err, "502")
	}

	return

}

func Post(o Origin, data, response interface{}) (outputError map[string]interface{}) {

	funcion := "Post - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlCRUD := "http://" + o.BasePath() + o.Endpoint()

	err := request.SendJson(urlCRUD, "POST", &response, &data)
	if err != nil {
		logs.Error(urlCRUD, err)
		outputError = errorctrl.Error("funcion+eval", err, "502")
		eval := `request.SendJson(urlCRUD, "POST", &response, &data)`
		outputError = errorctrl.Error(funcion+eval, err, "502")
	}

	return

}
