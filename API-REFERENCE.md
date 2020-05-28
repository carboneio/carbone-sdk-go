# GO Carbone SDK

The GO Carbone SDK provide an simple interface to communicate with Carbone Render easily.

## Install the GO SDK

```sh
go get github.com/Ideolys/carbone-sdk-go
```

## Quickstart with the GO SDK

Try the following code to render a report in 10 seconds. Just replace your API key, the template you want to render, and the data as a stringified JSON.

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
	reportBuffer, err := csdk.Render(templateID, jsonData)
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
func NewCarboneSDK(SecretAccessToken ...string) (*CSDK, error)
```
Function to create a new instance of CSDK (CarboneSDK).
The access token can be pass as an argument to NewCarboneSDK (`args[0]`) or by the environment variable "CARBONE_TOKEN".
To set a new environment variable, use the command:
```bash
$ export CARBONE_TOKEN=your-secret-token
```
Check if it is set by running:
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
func (csdk *CSDK) Render(pathOrTemplateID string, jsonData string, payload ...string) ([]byte, error)
```
The render function takes `pathOrTemplateID` the path of your local file OR a templateID, `jsonData` a stringified JSON, and an optional `payload`.

It returns the report as a `[]byte`. Carbone engine deleted files that have not been used for a while. By using this method, if your file has been deleted, the SDK will automatically upload it again and return you the result.

When a **template file path** is passed as an argument, the function verifies if the template has been uploaded to render the report. If not, it calls [AddTemplate](#AddTemplate) to upload the template to the server and generate a new template ID. Then it calls [RenderReport](#RenderReport) and [GetReport](#GetReport) to generate the report. If the path does not exist, an error is returned.

When a **templateID** is passed as an argument, the function renders with [RenderReport](#RenderReport) then call [GetReport](#GetReport) to return the report. If the templateID does not exist, an error is returned.

**Example**
```go
reportBuffer, err := csdk.Render("./templates/invoice.docx", `{"data":{"nane":"eric"},"convertTo":"pdf"}`, "OptionalPayload1234")
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
func (csdk *CSDK) AddTemplate(templateFileName string, payload ...string) (APIResponse, error)
```
Add the template to the API and returns an `APIResponse` struct (that contains a `TemplateID`).
You can add multiple times the same template and get different templateId thanks to the optional `payload`.

**Example**
```go
resp, err := csdk.AddTemplate("./tests/template.test.odt")
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

```go
	templateData, err := csdk.GetTemplate("TemplateId")
	if err != nil || len(templateData) <= 0 {
		t.Error(err)
	}
	err = ioutil.WriteFile(filename, templateData, 0644)
	if err != nil {
		t.Error(err)
	}
```

### DeleteTemplate
```go
func (csdk *CSDK) DeleteTemplate(templateID string) (APIResponse, error)
```
**Example**
```go
resp, err := csdk.DeleteTemplate(templateID)
if err != nil {
	t.Error(err)
}
if resp.Success == false {
	t.Error(resp.Error)
}
```

### RenderReport
```go
func (csdk *CSDK) RenderReport(templateID string, jsonData string) (APIResponse, error)
```
Function to render the report from a templateID and a stringified JSON Object with [datas and options](https://carbone.io/api-reference.html#rendering-a-report). It returns a APIResponse struct. The generated report and link are destroyed one hour after rendering.


**Example**
```go
cresp, err := csdk.RenderReport(templateID, `{"data":{},"convertTo":"pdf"}`)
if err != nil {
	t.Error(err)
}
if cresp.Success == false {
	t.Error(cresp.Error)
}
if len(cresp.Data.RenderID) <= 0 {
	t.Error(errors.New("renderId has not been returned"))
}
// `cresp.Data.RenderID` can be used to get the report from the `GetReport` function
```

### GetReport
```go
func (csdk *CSDK) GetReport(renderID string) ([]byte, error)
```
Return the Report from a renderID.

**Example**
```go
reportBuffer, err := csdk.GetReport(cresp.Data.RenderID)
if err != nil {
	t.Error(report)
}
if len(report) <= 0 {
	t.Error(errors.New("Report empty"))
}
err = ioutil.WriteFile("ReportName.pdf", reportBuffer, 0644)
if err != nil {
	log.Fatal(err)
}

```
### GenerateTemplateID
```go
func (csdk *CSDK) GenerateTemplateID(filepath string, payload ...string) (string, error)
```
The Template ID is predictable and idempotent, pass the template path and it will return the `templateID`.
You can get a different templateId thanks to the optional `payload`.


### SetAccessToken
```go
func (csdk *CSDK) SetAccessToken(newToken string)
```
It sets the Carbone access token.

### SetAPIVersion
```go
func (csdk *CSDK) SetAPIVersion(version int)
```
It sets the the Carbone version requested. By default, it is calling the version `2` of Carbone.

*Note:* You can only set a major version of carbone.

### GetAPIVersion
```go
func (csdk *CSDK) GetAPIVersion() (int, error)
```
It returns the Carbone version.