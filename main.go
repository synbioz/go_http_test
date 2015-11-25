package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Car struct {
	Manufacturer string `json:"manufacturer"`
	Design       string `json:"design"`
	Doors        uint8  `json:"doors"`
}

type Cars []Car

func main() {
	client := &http.Client{}
	car := &Car{Manufacturer: "renault", Design: "r5", Doors: 3}

	returnedCar, err := AddCar(car, client, "http://localhost:8080")

	if err != nil {
		panic(err)
	}

	fmt.Println(returnedCar)

	cars, err := GetAllCars(client)

	if err != nil {
		panic(err)
	}

	fmt.Println(cars)
}

func GetAllCars(client *http.Client) (*Cars, error) {
	resp, err := client.Get("http://localhost:8080/cars")

	if err != nil {
		return nil, err
	}

	cars, err := ParseGetAllCarsResponse(resp)

	return cars, err
}

func AddCar(car *Car, client *http.Client, url string) (*Car, error) {
	req, err := AddCarRequest(car, url)

	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	returnedCar, err := ParseAddCarResponse(resp)

	return returnedCar, err
}

func AddCarRequest(car *Car, url string) (*http.Request, error) {
	encodedJson, err := json.Marshal(car)

	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest("POST", url+"/cars", bytes.NewBuffer(encodedJson))
	req.Header.Set("Content-Type", "application/json")

	return req, err
}

func ParseAddCarResponse(resp *http.Response) (*Car, error) {
	var returnedCar Car

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &returnedCar)

	return &returnedCar, err
}

func ParseGetAllCarsResponse(resp *http.Response) (*Cars, error) {
	var cars *Cars

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &cars)

	return cars, err
}
