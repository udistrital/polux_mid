package controllers

import (
	"net/http"
	"testing"
)

func TestObtenerMinimo(t *testing.T) {
	if response, err := http.Get("http://localhost:9001/v1/creditos_materias/ObtenerCreditos"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error ObtenerMinimo Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("ObtenerMinimo Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error ObtenerMinimo:", err.Error())
		t.Fail()
	}
}
