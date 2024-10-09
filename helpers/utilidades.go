package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/xray"
)

type Origin interface {
	BasePath() string
	Endpoint() string
}

const (
	AppJson string = "application/json"
)

// Envia una petición al endpoint indicado y extrae la respuesta del campo Data para retornarla
func GetRequestNew(endpoint string, route string, target interface{}) error {
	url := beego.AppConfig.String(endpoint) + route
	fmt.Println("url GET", url)
	var response map[string]interface{}
	var err error
	err = request.GetJson(url, &response)
	err = ExtractData(response, &target, err)
	return err
}

// Envia una petición con datos al endpoint indicado y extrae la respuesta del campo Data para retornarla
// func SendRequestNew(endpoint string, route string, trequest string, target interface{}, datajson interface{}) (err error) {
// 	url := beego.AppConfig.String(endpoint) + route
// 	var response map[string]interface{}
// 	err = request.SendJson(url, trequest, &response, &datajson)
// 	err = ExtractData(response, target, err)
// 	return err
// }

func SendRequestNew(endpoint string, route string, trequest string, target interface{}, datajson interface{}) (status string, err error) {
	var response map[string]interface{}
	url := beego.AppConfig.String(endpoint) + route
	fmt.Println("url SEND", url)
	err = request.SendJson(url, trequest, &response, &datajson)
	//fmt.Println("RESPoNSE ", response)
	status, err = ExtractData2(response, target, err)
	//fmt.Println("status en send request ", status)
	return status, err
}

func GetJson(url string, target interface{}) error {
	req, _ := http.NewRequest("GET", url, nil)
	seg := xray.BeginSegmentSec(req)
	r, err := http.Get(url)
	xray.UpdateSegment(r, err, seg)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			beego.Error(err)
		}
	}()

	return json.NewDecoder(r.Body).Decode(target)
}

// Esta función extrae la información cuando se recibe encapsulada en una estructura
// y da manejo a las respuestas que contienen arreglos de objetos vacíos
func ExtractData(respuesta map[string]interface{}, v interface{}, err2 error) error {
	var err error
	if err2 != nil {
		return err2
	}
	if respuesta["Success"] == false {
		//err = errors.New(fmt.Sprint(respuesta["Data"], respuesta["Message"]))
		err := map[string]interface{}{"err": respuesta["Data"], "Message": respuesta["Message"], "Status": respuesta["Status"]}
		//panic(err)
		fmt.Println("Respuesta ExtractData 2", err)
		panic(err)
	}
	fmt.Println("Respuesta ExtractData Datatype 2", respuesta["Data"])
	datatype := fmt.Sprintf("%v", respuesta["Data"])
	switch datatype {
	case "map[]", "[map[]]": // response vacio
		break
	default:
		err = formatdata.FillStruct(respuesta["Data"], &v)
		respuesta = nil
	}
	fmt.Println("RETORNO EXTRACT DATA", err)
	return err
}
func ExtractData2(respuesta map[string]interface{}, v interface{}, err2 error) (status string, err error) {
	//fmt.Println("RESpUESTA ", respuesta)
	var statusAux = ""
	if err2 != nil {
		return "400", err2
	}
	if respuesta["Success"] == false {
		err := map[string]interface{}{"err": respuesta["Data"], "Message": respuesta["Message"], "Status": respuesta["Status"]}
		fmt.Println("Respuesta ExtractData 2", err)
	}
	datatype := fmt.Sprintf("%v", respuesta["Data"])
	fmt.Println("DATATYPE2 ", datatype)
	switch datatype {
	case "map[]", "[map[]]": // response vacio
		break
	default:
		err = formatdata.FillStruct(respuesta["Data"], &v)
		fmt.Println("STATUS2 ", respuesta["Status"].(string))
		statusAux = respuesta["Status"].(string)
		respuesta = nil
	}
	return statusAux, err
}

func Post(o Origin, data, response interface{}) (outputError map[string]interface{}) {
	funcion := "Post - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlCRUD := o.BasePath() + o.Endpoint()

	err := request.SendJson(urlCRUD, "POST", &response, &data)
	if err != nil {
		logs.Error(urlCRUD, err)
		outputError = errorctrl.Error("funcion+eval", err, "502")
		eval := `request.SendJson(urlCRUD, "POST", &response, &data)`
		outputError = errorctrl.Error(funcion+eval, err, "502")
	}
	return
}

// Manejo único de errores para controladores sin repetir código
func ErrorController(c beego.Controller, controller string) {
	var statusRes string
	var msgError string
	if err := recover(); err != nil {
		//fmt.Println("Lleva la secuencia pasando por el conrolador")
		logs.Error(err)
		localError := err.(map[string]interface{})
		c.Data["message"] = (beego.AppConfig.String("appname") + "/" + controller + "/" + (localError["funcion"]).(string))
		c.Data["data"] = (localError["err"])
		xray.EndSegmentErr(http.StatusBadRequest, localError["err"])
		if status, ok := localError["Status"]; ok {
			statusRes = status.(string)
			statusCode, _ := strconv.Atoi(status.(string))
			if msg, ok := localError["err"].(string); ok {
				msgError = msg
			} else if msg, ok := localError["err"].(error); ok {
				msgError = fmt.Sprint(msg)
			}
			c.Ctx.Output.SetStatus(statusCode)
		} else {
			statusRes = "500"
			c.Ctx.Output.SetStatus(500)
		}
		c.Data["json"] = map[string]interface{}{
			"Data":    "Error en " + (beego.AppConfig.String("appname") + "/" + controller + "/" + (localError["funcion"]).(string)),
			"Message": msgError,
			"Status":  statusRes,
			"Success": false,
		}
		c.ServeJSON()
	}
}

func DeferHelpers(funcion string, err interface{}) (outputError map[string]interface{}) {
	//fmt.Println("err ", err)
	var localError map[string]interface{}
	if errTemp, ok := err.(map[string]interface{}); ok {
		localError = errTemp
		//fmt.Println("STATUS ", localError["Status"])
	}
	if status, ok := localError["Status"]; ok {
		//fmt.Println("STATUS ", localError["Status"])
		outputError = map[string]interface{}{"funcion": funcion, "err": localError["Message"], "Status": status.(string)}
	} else {
		//fmt.Println("STATUS ", localError["Status"])
		outputError = map[string]interface{}{"funcion": funcion, "err": err, "status": "500"}
	}
	//fmt.Println("output ", outputError)

	return outputError
}
