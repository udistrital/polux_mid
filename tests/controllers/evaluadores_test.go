package controllers

import (
	"bytes"
	"net/http"
	"testing"
)

func TestObtenerEvaluadores(t *testing.T) {
	body := []byte(`{
		"Modalidad": 2
	}`)

	if response, err := http.Post("http://localhost:9001/v1/evaluadores/ObtenerEvaluadores", "application/json", bytes.NewBuffer(body)); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error ObtenerEvaluadores Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("ObtenerEvaluadores Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error ObtenerEvaluadores:", err.Error())
		t.Fail()
	}
}
