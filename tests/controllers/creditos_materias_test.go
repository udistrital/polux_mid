package controllers

import (
	"net/http"
	"testing"
)

func TestObtenerCreditos(t *testing.T) {
	if response, err := http.Get("http://localhost:9001/v1/creditos_materias/ObtenerCreditos"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error ObtenerCreditos Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("ObtenerCreditos Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error ObtenerCreditos:", err.Error())
		t.Fail()
	}
}
