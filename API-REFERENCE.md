# GO Carbone SDK

The GO Carbone SDK provide an simple interface to communicate with Carbone Render easily.

## Install the GO SDK

```sh
go get github.com/Ideolys/carbone-sdk-go
```

## Quickstart with the GO SDK

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

## GO SDK API

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
Example
```go
// Carbone access token passed as parameter
csdk, err := carbone.NewCarboneSDK("YOUR-ACCESS-TOKEN")
// Carbone access token passed as environment variable "Carbone TOKEN"
csdk, err := carbone.NewCarboneSDK()
```
### Render
```go
func (csdk *CSDK) Render(pathOrTemplateID string, jsonData string, payload string) ([]byte, error)
```
The render function takes `pathOrTemplateID` the path of your local file OR a templateID, `jsonData` a stringify json and `payload` (optional). It returns the report as a []byte. Carbone engine deleted files which has not been used since a while. By using this method, if your file has been deleted, the SDK will automatically upload it again and return you the result.

When a **template file path** is passed as argument, the function verifies if the template has been uploaded to render the report. If not, it call `AddTemplate` to upload the template to the server and generate a new template ID. Then it call `RenderReport` and `GetReport` to generate the report. If the path does not exist, an error is return.

When a **templateID** is passed as argument, the function render with `RenderReport` then call `GetReport` to return the report. If the templateID does not exist, an error is return.

**Example**
```go
reportBuffer, err := csdk.Render("./templates/invoice.docx", `{"data":{"nane":"eric"},"convertTo":"pdf"}`, "")
if err != nil {
	log.Fatal(err)
}
// create the file
err = ioutil.WriteFile("Report.pdf", reportBuffer, 0644)
if err != nil {
	log.Fatal(err)
}
```


### AddTemplate
```go
func (csdk *CSDK) AddTemplate(templateFileName string, payload string) (APIResponse, error)
```
Add the template to the API and returns an `APIResponse` struct with a `TemplateID`.
You can add multiple times the same template and get different templateId thanks to the payload.

**Example**
```go
resp, err := csdk.AddTemplate("./tests/template.test.odt", "")
if err != nil {
	t.Error(err)
}
if resp.Success == false {
	t.Error(resp.Error)
}
if len(resp.Data.TemplateID) <= 0 {
	t.Error(errors.New("templateId not returned from the api"))
}
fmt.Println("templateID:", resp.Data.TemplateID)
```

### GetTemplate
```go
func (csdk *CSDK) GetTemplate(templateID string) ([]byte, error)
```

Pass a `templateID` to the function and it returns the template as `[]byte`. The templateID must exist otherwise an error is returned by the server.

### DeleteTemplate
```go
func (csdk *CSDK) DeleteTemplate(templateID string) (APIResponse, error)
```

### RenderReport
```go
func (csdk *CSDK) RenderReport(templateID string, jsonData string) (APIResponse, error)
```

### GetReport
```go
func (csdk *CSDK) GetReport(renderID string) ([]byte, error)
```

### GenerateTemplateID
```go
func (csdk *CSDK) GenerateTemplateID(filepath string, payload string) (string, error)
```
The Template ID is predictable and indempotent, pass the template path and it will return the "templateID".

### SetAccessToken
```go
func (csdk *CSDK) SetAccessToken(newToken string)
```