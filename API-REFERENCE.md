# GO Carbone SDK

The GO Carbone SDK provide an simple interface to communicate with Carbone Render easily

## Install sdk go

```sh
go get github.com/Ideolys/carbone-sdk-go
```

## Quickstart sdk go

Try the following code to render a report in 10 seconds. Just replace your API key, the template you want to render and the data as stringify JSON.

```go
package main

import (
	"io/ioutil"
	"log"

	"github.com/Ideolys/carbone-sdk-go/carbone"
)

func main() {
	csdk, err := carbone.NewCarboneSDK("YOUR-ACCESS-TOKEN")
	if err != nil {
		log.Fatal(err)
  }
  // Path to your template
  templateID := "./folder/template.odt"
  // Add your data here
	jsonData := `{"data":{},"convertTo":"pdf"}`
	reportBuffer, err := csdk.Render(templateID, jsonData, "")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("Report.pdf", reportBuffer, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

```

## API sdk go

### NewCarboneSDK
```go
func NewCarboneSDK(args ...string) (*CSDK, error)
```
Function to create a new instance of CSDK (CarboneSDK).
The access token can be pass as argument to NewCarboneSDK or by the environment variable "CARBONE_TOKEN".
To set a new environment variable, use the command:
```bash
$ export CARBONE_TOKEN=your-secret-token
```
You can check this by running:
```bash
$ printenv | grep "CARBONE_TOKEN"
```
#### Example
```go
// Carbone access token passed as parameter
csdk, err := carbone.NewCarboneSDK("YOUR-ACCESS-TOKEN")
// Carbone access token passed as environment variable "Carbone TOKEN"
csdk, err := carbone.NewCarboneSDK()
```
### Render
```go
func (csdk *CSDK) GetReport(renderID string) ([]byte, error)
```

### AddTemplate
```go
func (csdk *CSDK) AddTemplate(templateFileName string, payload string) (APIResponse, error)
```

### GetTemplate
```go
func (csdk *CSDK) GetTemplate(templateID string) ([]byte, error)
```

### DeleteTemplate
```go
func (csdk *CSDK) DeleteTemplate(templateID string) (APIResponse, error)
```

### RenderReport
```go
func (csdk *CSDK) RenderReport(templateID string, jsonData string) (APIResponse, error)
```

### GenerateTemplateID
```go
func (csdk *CSDK) GenerateTemplateID(filepath string, payload string) (string, error)
```

### SetAccessToken
```go
func (csdk *CSDK) SetAccessToken(newToken string)
```