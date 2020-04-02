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
	csdk, e = NewCarboneSDK("eyJhbGciOiJFUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiIxNjY3IiwiYXVkIjoiY2FyYm9uZSIsImV4cCI6MjIwNzQwNjQ0NywiZGF0YSI6eyJpZEFjY291bnQiOjE2Njd9fQ.AH2NiPdd8dRC_FNsd4aJ1DHy2wNNhXFmRvyh6PM-jkksfPn7hIIgiUfZ-L7Ng9Jou3eCeLrymjcPuABFVcaGiGvCATAICKX_j7WKBdMO_iPzD1LvL5j35FX1_i513OLqSvqTY_3KvBZO2RXMh4tLWlMn-dhNFLn-aE6IcS3lpce_A2PB")
	if e != nil {
		log.Fatal(e)
	}
	os.Exit(m.Run())
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
		template := "./tests/index.test.html"
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
		template := "./tests/index.test.html"
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
		filename := "./tests/index.test.html"
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
		filename := "./tests/index.test.html"
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
		os.Remove("template.test.odt")
		templateData, err := csdk.GetTemplate(templateID)
		if err != nil || len(templateData) <= 0 {
			t.Error(err)
		}
		err = ioutil.WriteFile("template.test.odt", templateData, 0644)
		if err != nil {
			t.Error(err)
		}
		err = os.Remove("template.test.odt")
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
		os.Remove("./report.test.pdf")
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
		er = ioutil.WriteFile("report.test.pdf", report, 0644)
		if er != nil {
			t.Error(er)
		}
		er = os.Remove("./report.test.pdf")
		if er != nil {
			t.Error(er)
		}
	})
}

/**
	- [x] Render à partir d'un templateID existant
	- [ ] Render à partir d'un templateID qui n'existe pas
	- [ ] Rendre un template qui n'existe pas dans le server
	- [ ] Rendre un template qui existe déjà
	- [ ] Rendre un template qui n'existe pas du tout dans le pc
**/

func TestRender(t *testing.T) {
	t.Run("Render a report from a templateID and create the file", func(t *testing.T) {
		os.Remove("./tests/report.test.pdf")
		templateID := "f90e67221d7d5ee11058a000bdb997fb41bf149b1f88b45cb1aba9edcab8f868"
		jsonData := `{"data":{"firstname":"Felix","lastname":"Arvid Ulf Kjellberg","color":"#00FF00"},"convertTo":"pdf"}`
		report, err := csdk.Render(templateID, jsonData, "")
		if err != nil {
			t.Fatal(err)
		}
		if len(report) <= 0 {
			t.Fatal(errors.New("The report is empty"))
		}
		err = ioutil.WriteFile("./tests/report.test.pdf", report, 0644)
		if err != nil {
			t.Fatal(err)
		}
		err = os.Remove("./tests/report.test.pdf")
		if err != nil {
			t.Fatal(err)
		}
	})
}
