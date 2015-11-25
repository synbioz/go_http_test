package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestAddCarRequest(t *testing.T) {
	var carTests = []struct {
		Car  *Car
		Json string
	}{
		{&Car{}, `{"manufacturer":"","design":"","doors":0}`},
		{&Car{Manufacturer: "renault", Design: "r5", Doors: 3}, `{"manufacturer":"renault","design":"r5","doors":3}`},
	}

	for _, tt := range carTests {
		request, err := AddCarRequest(tt.Car, "http://localhost")

		defer request.Body.Close()
		json, _ := ioutil.ReadAll(request.Body)

		if tt.Json != string(json) {
			t.Errorf("AddCarRequest: expected: %s, actual: %s", tt.Json, string(json))
		}
		if err != nil {
			t.Errorf("AddCarRequest: error: %s", err)
		}
	}
}

func TestParseAddCarResponse(t *testing.T) {
	var carTests = []struct {
		Car  *Car
		Json string
	}{
		{&Car{}, `{"manufacturer":"","design":"","doors":0}`},
		{&Car{Manufacturer: "renault", Design: "r5", Doors: 3}, `{"manufacturer":"renault","design":"r5","doors":3}`},
	}

	for _, tt := range carTests {
		resp := &http.Response{Body: ioutil.NopCloser(strings.NewReader(tt.Json))}

		car, err := ParseAddCarResponse(resp)

		if !reflect.DeepEqual(tt.Car, car) {
			t.Errorf("ParseAddCarResponse: expected: %v, actual: %v", tt.Car, car)
		}

		if err != nil {
			t.Errorf("ParseAddCarResponse: error: %s", err)
		}
	}
}

func TestParseGetAllCarsResponse(t *testing.T) {
	var carTests = []struct {
		Cars *Cars
		Json string
	}{
		{new(Cars), `[]`},
		{&Cars{Car{Manufacturer: "renault", Design: "r5", Doors: 3}}, `[{"manufacturer":"renault","design":"r5","doors":3}]`},
	}

	for _, tt := range carTests {
		resp := &http.Response{Body: ioutil.NopCloser(strings.NewReader(tt.Json))}

		cars, err := ParseGetAllCarsResponse(resp)

		if len(*tt.Cars) == 0 && len(*tt.Cars) != len(*cars) {
			t.Errorf("ParseGetAllCarsResponse: expected len: %d, actual len: %d", len(*tt.Cars), len(*cars))
		} else if !reflect.DeepEqual(tt.Cars, cars) && len(*tt.Cars) != 0 {
			t.Errorf("ParseGetAllCarsResponse: expected: %v, actual: %v", tt.Cars, cars)
		}

		if err != nil {
			t.Errorf("ParseGetAllCarsResponse: error: %s", err)
		}
	}
}

func TestAddCar(t *testing.T) {
	var returnJsonCarHandler = func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"manufacturer":"citroen", "design":"ds3", "door":0}`)
	}

	server := httptest.NewServer(http.HandlerFunc(returnJsonCarHandler))
	defer server.Close()

	httpClient := &http.Client{}

	car, _ := AddCar(&Car{}, httpClient, server.URL)

	if car.Manufacturer != "citroen" || car.Design != "ds3" {
		t.Errorf("AddCar: %v && %s", car, server.URL)
	}
}
