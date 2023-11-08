package controllers

import (
	"net/http"
	"testing"
)

func TestObtenerCupos(t *testing.T) {
	if response, err := http.Get("http://localhost:9001/v1/cupos/Obtener"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error ObtenerCupos Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("ObtenerCupos Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error ObtenerCupos:", err.Error())
		t.Fail()
	}
}
