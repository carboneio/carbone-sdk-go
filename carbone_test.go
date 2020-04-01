package carbone

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var csdk *CSDK

/**
	TODO:
	- Générer le templateID avec la lib crypto
**/

func TestMain(m *testing.M) {
	var e error
	csdk, e = NewCarboneSDK("eyJhbGciOiJFUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiIxNjY3IiwiYXVkIjoiY2FyYm9uZSIsImV4cCI6MjIwNzQwNjQ0NywiZGF0YSI6eyJpZEFjY291bnQiOjE2Njd9fQ.AH2NiPdd8dRC_FNsd4aJ1DHy2wNNhXFmRvyh6PM-jkksfPn7hIIgiUfZ-L7Ng9Jou3eCeLrymjcPuABFVcaGiGvCATAICKX_j7WKBdMO_iPzD1LvL5j35FX1_i513OLqSvqTY_3KvBZO2RXMh4tLWlMn-dhNFLn-aE6IcS3lpce_A2PB")
	if e != nil {
		log.Fatal(e)
	}
	os.Exit(m.Run())
}

func TestAddTemplate(t *testing.T) {
	resp, err := csdk.AddTemplate("./tests/template.odt", "")
	if err != nil {
		t.Error(err)
	}
	if resp.Success == false {
		t.Error(resp.Error)
	}
	if len(resp.Data.TemplateID) <= 0 {
		t.Error(errors.New("templateId not returned from the api"))
	}
}

func TestAddTemplateWithEmptyFilePath(t *testing.T) {
	resp, err := csdk.AddTemplate("", "")
	if err == nil || resp.Success == true {
		t.Error(errors.New("Test failled: the file path argument is empty and the method should have thrown an error"))
	}
}

func TestGetTemplate(t *testing.T) {
	templateData, err := csdk.GetTemplate("f90e67221d7d5ee11058a000bdb997fb41bf149b1f88b45cb1aba9edcab8f868")
	if err != nil || len(templateData) <= 0 {
		t.Error(err)
	}
}

func TestGetTemplateAndCreateFile(t *testing.T) {
	os.Remove("template.test.odt")
	templateData, err := csdk.GetTemplate("f90e67221d7d5ee11058a000bdb997fb41bf149b1f88b45cb1aba9edcab8f868")
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
}

func TestGetTemplateWithEmptyTemplateID(t *testing.T) {
	templateData, err := csdk.GetTemplate("")
	if err == nil || len(templateData) > 0 {
		t.Error(errors.New("Test failled: the templateID argument is empty and the method should have thrown an error"))
	}
}

func TestAddDeleteTwiceTemplate(t *testing.T) {
	res, err := csdk.AddTemplate("./tests/template.odt", "")
	if err != nil {
		t.Error(err)
	}
	if res.Success == false {
		t.Error(res.Error)
	}
	if len(res.Data.TemplateID) <= 0 {
		t.Error(errors.New("templateId not returned from the api"))
	}
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
}

func TestDeleteTemplateWithEmptyTemplateID(t *testing.T) {
	resp, err := csdk.DeleteTemplate("")
	if err == nil || resp.Success == true {
		t.Error(errors.New("Test failled: the templateID argument is empty and the method should have thrown an error"))
	}
}

func TestRenderTemplate(t *testing.T) {
	templateID := "f90e67221d7d5ee11058a000bdb997fb41bf149b1f88b45cb1aba9edcab8f868"
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
}

func TestRenderAndGetReport(t *testing.T) {
	templateID := "f90e67221d7d5ee11058a000bdb997fb41bf149b1f88b45cb1aba9edcab8f868"
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
}

func TestRenderAndGetReportAndCreateFile(t *testing.T) {
	os.Remove("./report.test.pdf")
	templateID := "f90e67221d7d5ee11058a000bdb997fb41bf149b1f88b45cb1aba9edcab8f868"
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
}

func TestRenderReportWithEmptyTemplateId(t *testing.T) {
	cresp, err := csdk.RenderReport("", ``)
	if err == nil || cresp.Success == true {
		t.Error(errors.New("Test failled: the templateID argument is empty and the method should have thrown an error"))
	}
}

func TestRenderReportWithEmptyJsonData(t *testing.T) {
	cresp, err := csdk.RenderReport("fewfwefwe", ``)
	if err == nil || cresp.Success == true {
		t.Error(errors.New("Test failled: the jsonData argument is empty and the method should have thrown an error"))
	}
}

func TestGetReportWithEmptyRenderId(t *testing.T) {
	file, err := csdk.GetReport("")
	if err == nil || len(file) > 0 {
		t.Error(errors.New("Test failled: the renderID argument is empty and the method should have thrown an error"))
	}
}
