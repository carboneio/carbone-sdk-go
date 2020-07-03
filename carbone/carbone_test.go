package carbone

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
)

var csdk *CSDK

func TestMain(m *testing.M) {
	var e error
	csdk, e = NewCarboneSDK()
	if e != nil {
		log.Fatal(e)
	}
	status := m.Run()
	deleteTmpFiles("tests")
	os.Exit(status)
}

func TestGenerateTemplateID(t *testing.T) {
	t.Run("(node test 1) generate a templateID from a file without payload", func(t *testing.T) {
		template := "./tests/template.test.odt"
		expectedHash := "f90e67221d7d5ee11058a000bdb997fb41bf149b1f88b45cb1aba9edcab8f868"
		resultHash, err := csdk.GenerateTemplateID(template)
		if err != nil {
			t.Error(err)
			return
		}
		if expectedHash != resultHash {
			t.Error(errors.New("Generated templateID not equal"))
		}
	})

	t.Run("(node test 2) generate a templateID from a file with payload", func(t *testing.T) {
		template := "./tests/template.test.odt"
		payload := "ThisIsAPayload"
		expectedHash := "2903587f87c80bd3a9b25d17f5db344fa6d276db17363403ac0fcf351a70d5a5"
		resultHash, err := csdk.GenerateTemplateID(template, payload)
		if err != nil {
			t.Error(err)
			return
		}
		if expectedHash != resultHash {
			t.Error(errors.New("Generated templateID not equal"))
		}
	})
	t.Run("(node test 3) generate a templateID from a file with payload", func(t *testing.T) {
		template := "./tests/template.test.odt"
		payload := "8B5PmafbjdRqHuksjHNw83mvPiGj7WTE"
		expectedHash := "12bcc644ff8479a09f80c01bbd05614599dd102478bc2aa40881a11cf20af21a"
		resultHash, err := csdk.GenerateTemplateID(template, payload)
		if err != nil {
			t.Error(err)
			return
		}
		if expectedHash != resultHash {
			t.Error(errors.New("Generated templateID not equal"))
		}
	})
	t.Run("(node test 4) generate a templateID from an HTML file without payload", func(t *testing.T) {
		template := "./tests/template.test.html"
		expectedHash := "75256dd5c260cdf039ae807d3a007e78791e2d8963ea1aa6aff87ba03074df7f"
		resultHash, err := csdk.GenerateTemplateID(template)
		if err != nil {
			t.Error(err)
			return
		}
		if expectedHash != resultHash {
			t.Error(errors.New("Generated templateID not equal"))
		}
	})
	t.Run("(node test 5) generate a templateID from an HTML file without payload", func(t *testing.T) {
		template := "./tests/template.test.html"
		payload := "This is a long payload with different characters 1 *5 &*9 %$ 3%&@9 @(( 3992288282 29299 9299929"
		expectedHash := "70799b421cc9cf75d9112273a8e054c141d484eb8d5988bd006fac83e3990707"
		resultHash, err := csdk.GenerateTemplateID(template, payload)
		if err != nil {
			t.Error(err)
			return
		}
		if expectedHash != resultHash {
			t.Error(errors.New("Generated templateID not equal"))
		}
	})
}

func TestAddTemplate(t *testing.T) {

	t.Run("Basic Add and test access token", func(t *testing.T) {

		templateID := "f90e67221d7d5ee11058a000bdb997fb41bf149b1f88b45cb1aba9edcab8f868"
		fakeToken := "ThisIsAToken1234"

		csdk.SetAccessToken(fakeToken)

		// ---- httpmock
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("POST", "https://render.carbone.io/template", func(req *http.Request) (*http.Response, error) {
			if strings.Contains(req.Header.Get("Authorization"), fakeToken) == false {
				return httpmock.NewStringResponse(500, "Access token not provided"), nil
			}
			return httpmock.NewStringResponse(200, `{ "success" : true, "data": {"templateId" : "`+templateID+`" }}`), nil
		})
		// ----

		resp, err := csdk.AddTemplate("./tests/template.test.odt")
		if err != nil {
			t.Error(err)
			return
		}
		if resp.Success == false {
			t.Error(resp.Error)
			return
		}
		if len(resp.Data.TemplateID) <= 0 {
			t.Error(errors.New("templateId not returned from the api"))
			return
		}
		if resp.Data.TemplateID != templateID {
			t.Error(errors.New("The template id is different"))
			return
		}
		// -- test httpmock
		if httpmock.GetTotalCallCount() != 1 {
			t.Fatal(errors.New("HTTPMOCH error - the number of requests is invalid"))
		}
	})

	t.Run("Basic Add with payload", func(t *testing.T) {

		templateID := "12345"

		// ---- httpmock
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		mockResp := httpmock.NewStringResponder(200, `{ "success" : true, "data": {"templateId" : "`+templateID+`" }}`)
		httpmock.RegisterResponder("POST", "https://render.carbone.io/template", mockResp)
		// ----

		resp, err := csdk.AddTemplate("./tests/template.test.odt", "This is an optional payload")
		if err != nil {
			t.Error(err)
		}
		if resp.Success == false {
			t.Error(resp.Error)
		}
		if len(resp.Data.TemplateID) <= 0 {
			t.Error(errors.New("templateId not returned from the api"))
		}
		if resp.Data.TemplateID != templateID {
			t.Error(errors.New("The template id is different"))
		}
		// -- test httpmock
		if httpmock.GetTotalCallCount() != 1 {
			t.Fatal(errors.New("HTTPMOCH error - the number of requests is invalid"))
		}
	})

	t.Run("WithEmptyFilePath", func(t *testing.T) {
		resp, err := csdk.AddTemplate("")
		if err == nil || resp.Success == true {
			t.Error(errors.New("Test failled: the file path argument is empty and the method should have thrown an error"))
		}
	})

	t.Run("WithWrongFilePath", func(t *testing.T) {
		resp, err := csdk.AddTemplate("./fewijwoeij.odt")
		if err == nil || resp.Success == true {
			t.Error(errors.New("Test failled: the file path argument is empty and the method should have thrown an error"))
		}
	})
}

func TestGetTemplate(t *testing.T) {
	templateID := "f90e67221d7d5ee11058a000bdb997fb41bf149b1f88b45cb1aba9edcab8f868"

	t.Run("Should Get the template", func(t *testing.T) {
		htmlBytes, err := ioutil.ReadFile("./tests/template.test.html")
		if err != nil {
			t.Error(err)
			return
		}
		// ---- httpmock
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		mockResp := httpmock.NewBytesResponder(200, htmlBytes)
		httpmock.RegisterResponder("GET", "https://render.carbone.io/template/"+templateID, mockResp)
		// ----

		templateData, err := csdk.GetTemplate(templateID)
		if err != nil || len(string(templateData)) <= 0 {
			t.Error(err)
		}

		if string(templateData) != string(htmlBytes) {
			t.Error(errors.New("The content is wrong"))
		}

		// -- test httpmock
		if httpmock.GetTotalCallCount() != 1 {
			t.Fatal(errors.New("HTTPMOCH error - the number of requests is invalid"))
		}
	})

	t.Run("Should Get the template with an missing template ID and Throw and error", func(t *testing.T) {
		templateData, err := csdk.GetTemplate("")
		if err == nil || len(templateData) > 0 {
			t.Error(errors.New("Test failled: the templateID argument is empty and the method should have thrown an error"))
		}
	})
}

func TestDeleteTemplate(t *testing.T) {
	templateID := "f90e67221d7d5ee11058a000bdb997fb41bf149b1f88b45cb1aba9edcab8f868"

	t.Run("Should delete only one time (delete called twice)", func(t *testing.T) {
		// ---- httpmock
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("DELETE", "https://render.carbone.io/template/"+templateID, httpmock.NewStringResponder(200, `{"success" : true, "error"   : null}`))
		// ----

		resp, err := csdk.DeleteTemplate(templateID)
		if err != nil {
			t.Error(err)
		}
		if resp.Success == false {
			t.Error(resp.Error)
		}

		// ---- httpmock
		httpmock.RegisterResponder("DELETE", "https://render.carbone.io/template/"+templateID, httpmock.NewStringResponder(200, `{"success" : false, "error"   : "The template doesn't exist"}`))
		// ---
		resp, err = csdk.DeleteTemplate(templateID)
		if err != nil {
			t.Error(err)
		}
		if resp.Success == true {
			t.Error(errors.New("Error: the template should not be able to delete the template twice"))
		}

		// -- test httpmock
		if httpmock.GetTotalCallCount() != 2 {
			t.Fatal(errors.New("HTTPMOCH error - the number of requests is invalid"))
		}
	})

	t.Run("Should throw an error because of a missing templateID as argument", func(t *testing.T) {
		resp, err := csdk.DeleteTemplate("")
		if err == nil || resp.Success == true {
			t.Error(errors.New("Test failled: the templateID argument is empty and the method should have thrown an error"))
		}
	})
}

func TestRenderReport(t *testing.T) {
	templateID := "f90e67221d7d5ee11058a000bdb997fb41bf149b1f88b45cb1aba9edcab8f868"
	renderID := "r3209jf903j2f90j2309fj3209fj"
	t.Run("Should Render basic a report", func(t *testing.T) {

		// ---- httpmock
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("POST", "https://render.carbone.io/render/"+templateID, httpmock.NewStringResponder(200, `{"success" : true,"error"   : null,"data": {"renderId": "`+renderID+`"}}`))
		// ----

		cresp, err := csdk.RenderReport(templateID, `{"data":{"firstname":"Felix","lastname":"Arvid Ulf Kjellberg","color":"#00FF00"},"convertTo":"pdf"}`)
		if err != nil {
			t.Error(err)
		}
		if cresp.Success == false {
			t.Error(cresp.Error)
		}
		if len(cresp.Data.RenderID) <= 0 {
			t.Error(errors.New("renderId has not been returned"))
		}
		if cresp.Data.RenderID != renderID {
			t.Error(errors.New("renderId has not been returned"))
		}

		// -- test httpmock
		if httpmock.GetTotalCallCount() != 1 {
			t.Fatal(errors.New("HTTPMOCH error - the number of requests is invalid"))
		}
	})

	t.Run("Should throw an error because the templateID arg is missing", func(t *testing.T) {
		cresp, err := csdk.RenderReport("", ``)
		if err == nil || cresp.Success == true {
			t.Error(errors.New("Test failled: the templateID argument is empty and the method should have thrown an error"))
		}
	})

	t.Run("Should throw an error because the jsonData arg is missing", func(t *testing.T) {
		cresp, err := csdk.RenderReport("fewfwefwe", ``)
		if err == nil || cresp.Success == true {
			t.Error(errors.New("Test failled: the jsonData argument is empty and the method should have thrown an error"))
		}
	})
}

func TestGetReport(t *testing.T) {
	// Setup
	// templateID := "f90e67221d7d5ee11058a000bdb997fb41bf149b1f88b45cb1aba9edcab8f868"
	renderID := "r3209jf903j2f90j2309fj3209fj"
	t.Run("Should throw an error because the renderID arg is missing", func(t *testing.T) {
		file, err := csdk.GetReport("")
		if err == nil || len(file) > 0 {
			t.Error(errors.New("Test failled: the renderID argument is empty and the method should have thrown an error"))
		}
	})

	t.Run("Should Get a report and create the file", func(t *testing.T) {
		htmlBytes, err := ioutil.ReadFile("./tests/template.test.html")
		if err != nil {
			t.Error(err)
			return
		}

		// ---- httpmock
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("GET", "https://render.carbone.io/render/"+renderID, httpmock.NewBytesResponder(200, htmlBytes))
		// ----

		report, er := csdk.GetReport(renderID)
		if er != nil {
			t.Error(report)
			return
		}
		if len(report) <= 0 {
			t.Error(errors.New("Rendered report empty"))
			return
		}

		// --- test compare content
		if string(report) != string(htmlBytes) {
			t.Error(errors.New("The content is not equal"))
			return
		}
		// -- test httpmock
		if httpmock.GetTotalCallCount() != 1 {
			t.Fatal(errors.New("HTTPMOCH error - the number of requests is invalid"))
		}
	})
}

func TestRender(t *testing.T) {

	t.Run("Render a report from an existing templateID and create the file", func(t *testing.T) {
		templateID := "f90e67221d7d5ee11058a000bdb997fb41bf149b1f88b45cb1aba9edcab8f868"
		renderID := "r32rsqdwq2fg2f32r90e67221d7d5ee11058a000bdb997fb41bf149b1f"
		content := "<xml>File Content</xml>"

		// ---- httpmock
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("POST", "https://render.carbone.io/render/"+templateID, httpmock.NewStringResponder(200, `{"success" : true,"error":null,"data": {"renderId": "`+renderID+`"}}`))
		httpmock.RegisterResponder("GET", "https://render.carbone.io/render/"+renderID, httpmock.NewBytesResponder(200, []byte(content)))
		// ----

		jsonData := `{"data":{"firstname":"Felix","lastname":"Arvid Ulf Kjellberg","color":"#00FF00"},"convertTo":"pdf"}`
		report, err := csdk.Render(templateID, jsonData)
		if err != nil {
			t.Fatal(err)
			return
		}
		if len(report) <= 0 {
			t.Fatal(errors.New("The report is empty"))
			return
		}
		// --- test compare content
		if string(report) != content {
			t.Error(errors.New("The content is not equal"))
		}
		// -- test httpmock
		if httpmock.GetTotalCallCount() != 2 {
			t.Fatal(errors.New("HTTPMOCH error - the number of requests is invalid"))
		}
	})

	t.Run("Render a report from an fake templateID/path (the file does not exist and hasn't been uploaded)", func(t *testing.T) {

		templateID := "ItsnotGonnaWork"
		errorMessage := "The templateID does not exist"

		// ---- httpmock
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("POST", "https://render.carbone.io/render/"+templateID, httpmock.NewStringResponder(200, `{"success" : false,"error":"`+errorMessage+`"}`))
		// ----

		jsonData := `{"data":{"firstname":"Felix","lastname":"Arvid Ulf Kjellberg","color":"#00FF00"},"convertTo":"pdf"}`
		report, err := csdk.Render(templateID, jsonData, "")
		if err == nil && err.Error() != errorMessage {
			t.Fatal(errors.New("Should have thrown an error"))
		}
		if len(report) > 0 {
			t.Fatal(errors.New("Should have been empty"))
		}
		if httpmock.GetTotalCallCount() != 1 {
			t.Fatal(errors.New("HTTPMOCH error - the number of requests is invalid"))
		}
	})

	t.Run("Render a report from a fresh new template (path as argument + payload)", func(t *testing.T) {
		payload := getUnixTime()
		templateID, err := csdk.GenerateTemplateID("./tests/template.test.html", payload)
		renderID := "feowjf329jf9823jf9823j9f238jf9832jf"
		fileContent := "<xml>File Content</xml>"
		nbrPOSTrequestCall := 0
		if err != nil {
			t.Fatal(errors.New("Can't generate the templateId"))
		}

		// ---- httpmock
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("POST", "https://render.carbone.io/render/"+templateID,
			func(req *http.Request) (*http.Response, error) {
				if nbrPOSTrequestCall >= 1 {
					// 3 - render the report and return a reportID // REQUEST NO CALLED
					return httpmock.NewStringResponse(200, `{"success" : true, "error":null, "data": {"renderId": "`+renderID+`"}}`), nil
				}
				nbrPOSTrequestCall++
				// 1 - test render from generated templateID but it does not exist
				return httpmock.NewStringResponse(200, `{"success" : false, "error": "Error while rendering template Error: 404 Not Found"}`), nil
			})
		// 2 - upload the template and return the template ID
		httpmock.RegisterResponder("POST", "https://render.carbone.io/template", httpmock.NewStringResponder(200, `{ "success" : false, "data": {"templateId" : "`+templateID+`" }}`))
		// 4 - Get the template
		httpmock.RegisterResponder("GET", "https://render.carbone.io/render/"+renderID, httpmock.NewBytesResponder(200, []byte(fileContent)))
		// ----

		jsonData := `{"data":{"name":"Arvid Ulf Kjellberg"},"convertTo":"html"}`
		reportBuffer, err := csdk.Render("./tests/template.test.html", jsonData, payload)
		if err != nil {
			t.Fatal(err)
		}
		if len(reportBuffer) <= 0 {
			t.Fatal(errors.New("The report is empty"))
		}
		if string(reportBuffer) != fileContent {
			t.Fatal(errors.New("The content is different"))
		}
		// --- test httpmock
		if httpmock.GetTotalCallCount() != 4 {
			t.Fatal(errors.New("HTTPMOCH error - the number of requests is invalid"))
		}
	})

	t.Run("Render a report from an existing template, the template has already been uploaded (templateID)", func(t *testing.T) {
		templatePath := "./tests/template.test.html"
		templateID, err := csdk.GenerateTemplateID(templatePath)
		renderID := "feowjf329jf9823jf9823j9f238jf9832jf"
		fileContent := "<xml>File Content</xml>"
		if err != nil {
			t.Fatal(errors.New("Can't generate the templateId"))
		}

		// ---- httpmock
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		// 1 - generate templateID and try to render, it already exist, it returns the renderID
		httpmock.RegisterResponder("POST", "https://render.carbone.io/render/"+templateID, httpmock.NewStringResponder(200, `{"success" : true,"error":null,"data": {"renderId": "`+renderID+`"}}`))
		// 2 - Get the template
		httpmock.RegisterResponder("GET", "https://render.carbone.io/render/"+renderID, httpmock.NewBytesResponder(200, []byte(fileContent)))
		// ----

		jsonData := `{"data":{"firstname":"Felix","lastname":"Arvid Ulf Kjellberg","color":"#00FF00"},"convertTo":"html"}`
		reportBuffer, err := csdk.Render(templatePath, jsonData, "")
		if err != nil {
			t.Fatal(err)
		}
		if len(reportBuffer) <= 0 {
			t.Fatal(errors.New("The report is empty"))
		}
		if string(reportBuffer) != fileContent {
			t.Fatal(errors.New("The content is different"))
		}
		if httpmock.GetTotalCallCount() != 2 {
			t.Fatal(errors.New("HTTPMOCH error - the number of requests is invalid"))
		}
	})

	t.Run("Render a report from a new template but the API return an error twice", func(t *testing.T) {
		payload := getUnixTime()
		templateID, err := csdk.GenerateTemplateID("./tests/template.test.html", payload)
		renderID := "feowjf329jf9823jf9823j9f238jf9832jf"
		fileContent := "<xml>File Content</xml>"
		errorMessage := "Error while rendering template Error: 404 Not Found"
		if err != nil {
			t.Fatal(errors.New("Can't generate the templateId"))
		}

		// ---- httpmock
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		// 1 - {CALLED TWICE} test render from generated templateID but it does not exist
		httpmock.RegisterResponder("POST", "https://render.carbone.io/render/"+templateID, httpmock.NewStringResponder(200, `{"success" : false, "error": "`+errorMessage+`"}`))
		// 2 - upload the template and return the template ID
		httpmock.RegisterResponder("POST", "https://render.carbone.io/template", httpmock.NewStringResponder(200, `{ "success" : false, "data": {"templateId" : "`+templateID+`" }}`))
		// 4 - Get the template
		httpmock.RegisterResponder("GET", "https://render.carbone.io/render/"+renderID, httpmock.NewBytesResponder(200, []byte(fileContent)))
		// ----

		jsonData := `{"data":{"name":"Arvid Ulf Kjellberg"},"convertTo":"html"}`
		reportBuffer, err := csdk.Render("./tests/template.test.html", jsonData, payload)
		if err.Error() != errorMessage {
			t.Fatal(errors.New("Should have returned an error"))
		}
		if string(reportBuffer) != "" {
			t.Fatal(errors.New("The buffer should have been empty"))
		}
		// --- test httpmock
		if httpmock.GetTotalCallCount() != 3 {
			t.Fatal(errors.New("HTTPMOCH error - the number of requests is invalid"))
		}
	})

	t.Run("Render and pass a directory path", func(t *testing.T) {
		templateID := "./tests/"
		jsonData := `{"data":{"firstname":"Felix","lastname":"Arvid Ulf Kjellberg","color":"#00FF00"},"convertTo":"pdf"}`
		report, err := csdk.Render(templateID, jsonData, "")
		if err == nil {
			t.Fatal(errors.New("Should have thrown an error"))
		}
		if len(report) > 0 {
			t.Fatal(errors.New("Should have been empty"))
		}
	})

	t.Run("Render the Carbone public report and change the access token.", func(t *testing.T) {
		reportName := "./tests/homereport.tmp.test.pdf"
		csdk.SetAccessToken("secret-token")
		templateID := "template"
		jsonData := `{"data":{"id":42,"date":1492012745,"company":{"name":"myCompany","address":"here","city":"Notfar","postalCode":123456},"customer":{"name":"myCustomer","address":"there","city":"Faraway","postalCode":654321},"products":[{"name":"product 1","priceUnit":0.1,"quantity":10,"priceTotal":1}],"total":140},"convertTo":"pdf"}`
		reportBuffer, err := csdk.Render(templateID, jsonData, "")
		if err != nil {
			t.Fatal(err)
		}
		if len(string(reportBuffer)) <= 0 {
			t.Fatal(errors.New("The report is empty"))
		}
		err = ioutil.WriteFile(reportName, reportBuffer, 0644)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Should update the default carbone render version", func(t *testing.T) {
		v, err := csdk.GetAPIVersion()
		if err != nil {
			t.Fatal(err)
		}
		if v != 2 {
			t.Fatal(errors.New("The API version is unvalid"))
		}
		csdk.SetAPIVersion(3)
		v, err = csdk.GetAPIVersion()
		if err != nil {
			t.Fatal(err)
		}
		if v != 3 {
			t.Fatal(errors.New("The API version is unvalid"))
		}
		csdk.SetAPIVersion(2)
	})
}
