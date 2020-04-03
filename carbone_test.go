package carbone

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var csdk *CSDK

func TestMain(m *testing.M) {
	var e error
	apiAccessToken := os.Getenv("CARBONE_API_TOKEN")
	if apiAccessToken == "" {
		log.Fatal(errors.New("Carbone Access token missing: in your terminal set the carbone access token: CARBONE_API_TOKEN"))
	}
	csdk, e = NewCarboneSDK(apiAccessToken)
	if e != nil {
		log.Fatal(e)
	}
	status := m.Run()
	deleteTmpFiles("tests")
	deleteTmpFiles("")
	os.Exit(status)
}

func TestGenerateTemplateID(t *testing.T) {
	t.Run("(node test 1) generate a templateID from a file without payload", func(t *testing.T) {
		template := "./tests/template.test.odt"
		payload := ""
		expectedHash := "f90e67221d7d5ee11058a000bdb997fb41bf149b1f88b45cb1aba9edcab8f868"
		resultHash, err := csdk.GenerateTemplateID(template, payload)
		if err != nil {
			t.Error(err)
			return
		}
		// fmt.Printf("Expected:%v\nResult:%v\n", expectedHash, resultHash)
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
		// fmt.Printf("Expected:%v\nResult:%v\n", expectedHash, resultHash)
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
		// fmt.Printf("Expected:%v\nResult:%v\n", expectedHash, resultHash)
		if expectedHash != resultHash {
			t.Error(errors.New("Generated templateID not equal"))
		}
	})
	t.Run("(node test 4) generate a templateID from an HTML file without payload", func(t *testing.T) {
		template := "./tests/template.test.html"
		payload := ""
		expectedHash := "75256dd5c260cdf039ae807d3a007e78791e2d8963ea1aa6aff87ba03074df7f"
		resultHash, err := csdk.GenerateTemplateID(template, payload)
		if err != nil {
			t.Error(err)
			return
		}
		// fmt.Printf("Expected:%v\nResult:%v\n", expectedHash, resultHash)
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
		// fmt.Printf("Expected:%v\nResult:%v\n", expectedHash, resultHash)
		if expectedHash != resultHash {
			t.Error(errors.New("Generated templateID not equal"))
		}
	})
	t.Run("Upload a template without payload and compare the templateID with the generated templateID", func(t *testing.T) {
		filename := "./tests/template.test.html"
		payload := ""
		resp, err := csdk.AddTemplate(filename, payload)
		if err != nil {
			t.Error(err)
			return
		}
		if resp.Success == false {
			t.Error(resp.Error)
			return
		}
		expectedTemplateID, err := csdk.GenerateTemplateID(filename, payload)
		if err != nil {
			t.Error(err)
		}
		if expectedTemplateID != resp.Data.TemplateID {
			t.Error(errors.New("Generated templateID not equal"))
		}
	})
	t.Run("Upload a template and compare the templateID with the generated templateID", func(t *testing.T) {
		filename := "./tests/template.test.html"
		payload := "7uE5G24ad2Vgnj2zFyiaqfN4dHzm4Xrq"
		resp, err := csdk.AddTemplate(filename, payload)
		if err != nil {
			t.Error(err)
			return
		}
		if resp.Success == false {
			t.Error(resp.Error)
			return
		}
		expectedTemplateID, err := csdk.GenerateTemplateID(filename, payload)
		if err != nil {
			t.Error(err)
		}
		if expectedTemplateID != resp.Data.TemplateID {
			t.Error(errors.New("Generated templateID not equal"))
		}
	})
}

func TestAddTemplate(t *testing.T) {
	t.Run("Basic Add", func(t *testing.T) {
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
	})

	t.Run("WithEmptyFilePath", func(t *testing.T) {
		resp, err := csdk.AddTemplate("", "")
		if err == nil || resp.Success == true {
			t.Error(errors.New("Test failled: the file path argument is empty and the method should have thrown an error"))
		}
	})

	t.Run("WithWrongFilePath", func(t *testing.T) {
		resp, err := csdk.AddTemplate("./fewijwoeij.odt", "")
		if err == nil || resp.Success == true {
			t.Error(errors.New("Test failled: the file path argument is empty and the method should have thrown an error"))
		}
	})
}

func TestGetTemplate(t *testing.T) {
	templateID := "f90e67221d7d5ee11058a000bdb997fb41bf149b1f88b45cb1aba9edcab8f868"
	t.Run("Should Get the template", func(t *testing.T) {
		templateData, err := csdk.GetTemplate(templateID)
		if err != nil || len(templateData) <= 0 {
			t.Error(err)
		}
	})

	t.Run("Should Get the template and create a file", func(t *testing.T) {
		filename := "./tests/template.tmp." + getUnixTime() + ".test.odt"
		templateData, err := csdk.GetTemplate(templateID)
		if err != nil || len(templateData) <= 0 {
			t.Error(err)
		}
		err = ioutil.WriteFile(filename, templateData, 0644)
		if err != nil {
			t.Error(err)
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
	// Setup before deleting
	res, err := csdk.AddTemplate("./tests/template.test.odt", "")
	if err != nil {
		t.Error(err)
	}
	if res.Success == false {
		t.Error(res.Error)
	}
	if len(res.Data.TemplateID) <= 0 {
		t.Error(errors.New("templateId not returned from the api"))
	}

	t.Run("Should delete only one time (delete called twice)", func(t *testing.T) {
		resp, err := csdk.DeleteTemplate(res.Data.TemplateID)
		if err != nil {
			t.Error(err)
		}
		if resp.Success == false {
			t.Error(resp.Error)
		}
		resp, err = csdk.DeleteTemplate(res.Data.TemplateID)
		if err != nil {
			t.Error(err)
		}
		if resp.Success == true {
			t.Error(errors.New("Error: the template should not be able to delete the template twice"))
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
	t.Run("Should Render basic a report", func(t *testing.T) {
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
	templateID := "f90e67221d7d5ee11058a000bdb997fb41bf149b1f88b45cb1aba9edcab8f868"

	t.Run("Should throw an error because the renderID arg is missing", func(t *testing.T) {
		file, err := csdk.GetReport("")
		if err == nil || len(file) > 0 {
			t.Error(errors.New("Test failled: the renderID argument is empty and the method should have thrown an error"))
		}
	})

	t.Run("Should Get a report", func(t *testing.T) {
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
		report, er := csdk.GetReport(cresp.Data.RenderID)
		if er != nil {
			t.Error(report)
		}
		if len(report) <= 0 {
			t.Error(errors.New("Rendered report empty"))
		}
	})

	t.Run("Should Get a report and create a file", func(t *testing.T) {
		reportname := "./tests/report.tmp." + getUnixTime() + ".test.pdf"
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
		report, er := csdk.GetReport(cresp.Data.RenderID)
		if er != nil {
			t.Error(report)
		}
		if len(report) <= 0 {
			t.Error(errors.New("Rendered report empty"))
		}
		er = ioutil.WriteFile(reportname, report, 0644)
		if er != nil {
			t.Error(er)
		}
	})
}

func TestRender(t *testing.T) {

	t.Run("Render a report from an existing templateID and create the file", func(t *testing.T) {
		templateID := "f90e67221d7d5ee11058a000bdb997fb41bf149b1f88b45cb1aba9edcab8f868"
		jsonData := `{"data":{"firstname":"Felix","lastname":"Arvid Ulf Kjellberg","color":"#00FF00"},"convertTo":"pdf"}`
		report, err := csdk.Render(templateID, jsonData, "")
		if err != nil {
			t.Fatal(err)
		}
		if len(report) <= 0 {
			t.Fatal(errors.New("The report is empty"))
		}
		err = ioutil.WriteFile("./tests/report."+templateID+".tmp.test.pdf", report, 0644)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Render a report from an fake templateID/path (the file does not exist and hasn't been uploaded)", func(t *testing.T) {
		templateID := "ItsnotGonnaWorkSoSad"
		jsonData := `{"data":{"firstname":"Felix","lastname":"Arvid Ulf Kjellberg","color":"#00FF00"},"convertTo":"pdf"}`
		report, err := csdk.Render(templateID, jsonData, "")
		if err == nil {
			t.Fatal(errors.New("Should have thrown an error"))
		}
		if len(report) > 0 {
			t.Fatal(errors.New("Should have been empty"))
		}
	})

	t.Run("Render a report from a fresh new template (path as argument) (create/delete template and create/delete report)", func(t *testing.T) {
		tmpTemplateNamePath, templateName := createTmpFile("", "")
		reportNamePath := "./tests/report." + templateName
		jsonData := `{"data":{"name":"Arvid Ulf Kjellberg"},"convertTo":"html"}`
		reportBuffer, err := csdk.Render(tmpTemplateNamePath, jsonData, "")
		if err != nil {
			t.Fatal(err)
		}
		if len(reportBuffer) <= 0 {
			t.Fatal(errors.New("The report is empty"))
		}
		err = ioutil.WriteFile(reportNamePath, reportBuffer, 0644)
		if err != nil {
			t.Fatal(err)
		}
		templateID, err := csdk.GenerateTemplateID(tmpTemplateNamePath, "")
		if err != nil {
			t.Fatal(err)
		}
		cresp, err := csdk.DeleteTemplate(templateID)
		if err != nil {
			t.Fatal(err)
		}
		if cresp.Success == false {
			t.Fatal(cresp.Error)
		}
	})

	t.Run("Render a report from an existing template and create the file (template has already been uploaded)", func(t *testing.T) {
		templatePath := "./tests/template.test.odt"
		jsonData := `{"data":{"firstname":"Felix","lastname":"Arvid Ulf Kjellberg","color":"#00FF00"},"convertTo":"pdf"}`
		reportBuffer, err := csdk.Render(templatePath, jsonData, "")
		if err != nil {
			t.Fatal(err)
		}
		if len(reportBuffer) <= 0 {
			t.Fatal(errors.New("The report is empty"))
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
}
