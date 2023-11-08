package controllers

import (
	"net/http"
	"testing"
)

func TestObtenerFechas(t *testing.T) {
	if response, err := http.Get("http://localhost:9001/v1/fechas/ObtenerFechas"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error ObtenerFechas Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("ObtenerFechas Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error ObtenerFechas:", err.Error())
		t.Fail()
	}
}
