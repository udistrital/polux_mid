package helpers

import (
	"bytes"
	"encoding/json"
	"errors"
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
	fmt.Println("url ", url)
	var response map[string]interface{}
	var err error
	err = request.GetJson(url, &response)
	err = ExtractData(response, &target)
	return err
}

// Envia una petición con datos al endpoint indicado y extrae la respuesta del campo Data para retornarla
func SendRequestNew(endpoint string, route string, trequest string, response interface{}, datajson interface{}) (err error) {
	url := beego.AppConfig.String(endpoint) + route
	//var response map[string]interface{}
	var statusCode int
	statusCode, err = SendJson(url, trequest, &response, &datajson)
	//err = ExtractData(response, target)
	if statusCode != 200 && statusCode != 201 {
		err = errors.New(fmt.Sprint("Error con status " + strconv.Itoa(statusCode)))
		fmt.Println("ERR ", err)
	}
	return err
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

func SendJson(url string, trequest string, target interface{}, datajson interface{}) (status int, err error) {
	b := new(bytes.Buffer)
	if datajson != nil {
		if err := json.NewEncoder(b).Encode(datajson); err != nil {
			beego.Error(err)
		}
	}
	client := &http.Client{}
	req, err := http.NewRequest(trequest, url, b)
	req.Header.Set("Accept", AppJson)
	req.Header.Add("Content-Type", AppJson)
	seg := xray.BeginSegmentSec(req)
	r, err := client.Do(req)
	xray.UpdateSegment(r, err, seg)
	if err != nil {
		beego.Error("error", err)
		return r.StatusCode, err
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			beego.Error(err)
		}
	}()

	return r.StatusCode, json.NewDecoder(r.Body).Decode(target)
}

// Esta función extrae la información cuando se recibe encapsulada en una estructura
// y da manejo a las respuestas que contienen arreglos de objetos vacíos
func ExtractData(respuesta map[string]interface{}, v interface{}) error {
	var err error
	if respuesta["Success"] == false {
		err = errors.New(fmt.Sprint(respuesta["Data"], respuesta["Message"]))
		panic(err)
	}
	datatype := fmt.Sprintf("%v", respuesta["Data"])
	switch datatype {
	case "map[]", "[map[]]": // response vacio
		break
	default:
		err = formatdata.FillStruct(respuesta["Data"], &v)
		respuesta = nil
	}
	return err
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
