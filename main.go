package main

import (
	"io/ioutil"
	"log"

	"github.com/Ideolys/carbone-sdk-go/carbone"
)

func main() {
	csdk, err := carbone.NewCarboneSDK("secret-token")
	if err != nil {
		log.Fatal(err)
	}
	templateID := "template"
	jsonData := `{"data":{"id":42,"date":1492012745,"company":{"name":"myCompany","address":"here","city":"Notfar","postalCode":123456},"customer":{"name":"myCustomer","address":"there","city":"Faraway","postalCode":654321},"products":[{"name":"product 1","priceUnit":0.1,"quantity":10,"priceTotal":1}],"total":140},"convertTo":"pdf"}`
	reportBuffer, err := csdk.Render(templateID, jsonData, "")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("Invoice.pdf", reportBuffer, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
