package controllers

import (
	"bytes"
	"net/http"
	"testing"
)

func TestRegistrar(t *testing.T) {
	body := []byte(`{
		"Codigo": "20142020001",
		"Estado": "A",
		"Modalidad": 4,
		"Nivel": "PREGRADO",
		"Nombre": "HASTAMORIR HERNANDEZ EDUWIN YESID",
		"PorcentajeCursado": "95",
		"Promedio": "3.87",
		"Rendimiento": "0.7167",
		"Tipo": "POSGRADO",
		"TipoCarrera": "INGENIERIA"
	  }`)

	if response, err := http.Post("http://localhost:9001/v1/verificarRequisitos/Registrar", "application/json", bytes.NewBuffer(body)); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error Registrar Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("Registrar Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error Registrar:", err.Error())
		t.Fail()
	}
}

func TestCantidadModalidades(t *testing.T) {
	body := []byte(`{
		"Cantidad": "1",
		"Modalidad": "4"
	  }`)

	if response, err := http.Post("http://localhost:9001/v1/verificarRequisitos/CantidadModalidades", "application/json", bytes.NewBuffer(body)); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error CantidadModalidades Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("CantidadModalidades Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error CantidadModalidades:", err.Error())
		t.Fail()
	}
}
