# Carbone Render Go SDK
![Version](https://img.shields.io/badge/version-1.0.1-blue.svg?cacheSeconds=2592000)
[![Documentation](https://img.shields.io/badge/documentation-yes-brightgreen.svg)](https://carbone.io/api-reference.html#carbone-sdk-go)

> The golang SDK to use Carbone Render easily.

Carbone is a report generator (PDF, DOCX, XLSX, ODT, PPTX, ODS, XML, CSV...) using templates and JSON data.
[Learn more about the Carbone ecosystem](https://carbone.io/documentation.html).

### 🔖 [API REFERENCE](./API-REFERENCE.md)

## Install

```sh
go get github.com/Ideolys/carbone-sdk-go
```

## Usage

You can copy and run the code bellow to try.
```go
package main

import (
	"io/ioutil"
	"log"

	"github.com/Ideolys/carbone-sdk-go/carbone"
)

func main() {
	// SDK constructor
	// The access token can be passed as an argument to NewCarboneSDK
	// Or by the environment variable "CARBONE_TOKEN", use the command "export CARBONE_TOKEN=secret-token"
	csdk, err := carbone.NewCarboneSDK("secret-token")
	if err != nil {
		log.Fatal(err)
	}

	// The template ID
	templateID := "template"
	// Data injected into the template to generate the report with Carbone
	jsonData := `{"data":{"id":42,"date":1492012745,"company":{"name":"myCompany","address":"here","city":"Notfar","postalCode":123456},"customer":{"name":"myCustomer","address":"there","city":"Faraway","postalCode":654321},"products":[{"name":"product 1","priceUnit":0.1,"quantity":10,"priceTotal":1}],"total":140},"convertTo":"pdf"}`

	// Render and return the report as []byte
	reportBuffer, err := csdk.Render(templateID, jsonData)
	if err != nil {
		log.Fatal(err)
	}
	// Create the file
	err = ioutil.WriteFile("Invoice.pdf", reportBuffer, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
```
## Documentation
- [API REFERENCE](./API-REFERENCE.md)

## Run tests
First, Go to the `carbone` package directory.
```bash
$ cd carbone
```
Then, create an environment variable `CARBONE_TOKEN` with the Carbone access token as value:
```bash
$ export CARBONE_TOKEN="YOUR_ACCESS_TOKEN"
```
Check if it is set by running:
```bash
$ printenv | grep "CARBONE_TOKEN"
```
To run all the tests (-v for verbose output):
```bash
$ go test -v
```
To run only one test:
```bash
$ go test -v -run NameOfTheTest
```
If you need to test the generation of templateId, you can use the nodejs `main.js` to test the sha256 generation.
```bash
$ node ./tests/main.js
```

## 👤 Author

- [**@steevepay**](https://github.com/steevepay)

## 🤝 Contributing

Contributions, issues and feature requests are welcome!

Feel free to check [issues page](https://github.com/Ideolys/carbone-sdk-go/issues).

## Show your support

Give a ⭐️ if this project helped you!