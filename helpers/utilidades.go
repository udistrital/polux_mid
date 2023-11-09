package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/formatdata"
)

// Envia una petición al endpoint indicado y extrae la respuesta del campo Data para retornarla
func GetRequestNew(endpoint string, route string, target interface{}) error {
	url := beego.AppConfig.String(endpoint) + route
	fmt.Println("url ", url)
	var response map[string]interface{}
	var err error
	err = GetJson(url, &response)
	err = ExtractData(response, &target)
	return err
}

func GetJson(url string, target interface{}) error {
	r, err := http.Get(url)
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
