package main

import (
	"io/ioutil"
	"log"

	"github.com/Ideolys/carbone-sdk-go/carbone"
)

func main() {
	// secret-token
	csdk, err := carbone.NewCarboneSDK("")
	if err != nil {
		log.Fatal(err)
	}

	templateID := "f90e67221d7d5ee11058a000bdb997fb41bf149b1f88b45cb1aba9edcab8f868"
	templateBuffer, err := csdk.GetTemplate(templateID)
	if err != nil || len(templateBuffer) <= 0 {
		log.Fatal(err)
	}
	ioutil.WriteFile("template.odt", templateBuffer, 0644)
	err = ioutil.WriteFile("report.pdf", templateBuffer, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// templateID := "template"
	// jsonData := `{"data":{"id":42,"date":1492012745,"company":{"name":"myCompany","address":"here","city":"Notfar","postalCode":123456},"customer":{"name":"myCustomer","address":"there","city":"Faraway","postalCode":654321},"products":[{"name":"product 1","priceUnit":0.1,"quantity":10,"priceTotal":1}],"total":140},"convertTo":"pdf"}`
	// reportBuffer, err := csdk.Render(templateID, jsonData, "")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// err = ioutil.WriteFile("report.pdf", templateBuffer, 0644)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
