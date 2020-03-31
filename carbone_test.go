package carbone

import (
	"errors"
	"log"
	"os"
	"testing"
)

var csdk CSDK

/**
	TODO:
	- render a report
	- full test: upload template, render report, get the report, delete the template
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
	resp, err := csdk.AddTemplate("./template.odt", "")
	if err != nil {
		t.Error(err)
	}
	if resp.Success == false {
		t.Error(resp.Error)
	}
	if len(resp.Data.TemplateID) <= 0 {
		t.Error(errors.New("templateId not returned from the api"))
	}
	// fmt.Println(resp.Data.TemplateID)
}

func TestGetTemplate(t *testing.T) {
	templateData, err := csdk.GetTemplate("f90e67221d7d5ee11058a000bdb997fb41bf149b1f88b45cb1aba9edcab8f868")
	if err != nil || len(templateData) <= 0 {
		t.Error(err)
	}
}

func TestAddDeleteTwiceTemplate(t *testing.T) {
	resp, err := csdk.AddTemplate("./template.odt", "")
	if err != nil {
		t.Error(err)
	}
	if resp.Success == false {
		t.Error(resp.Error)
	}
	if len(resp.Data.TemplateID) <= 0 {
		t.Error(errors.New("templateId not returned from the api"))
	}
	resp, err = csdk.DeleteTemplate(resp.Data.TemplateID)
	if err != nil {
		t.Error(err)
	}
	if resp.Success == false {
		t.Error(resp.Error)
	}
	resp, err = csdk.DeleteTemplate(resp.Data.TemplateID)
	if err != nil {
		t.Error(err)
	}
	if resp.Success == true {
		t.Error(errors.New("Error: the template should not be able to delete the template twice"))
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
	template, er := csdk.GetReport(cresp.Data.RenderID)
	if er != nil {
		t.Error(template)
	}
	if len(template) <= 0 {
		t.Error(errors.New("Rendered template empty"))
	}
}

// templateID := "f90e67221d7d5ee11058a000bdb997fb41bf149b1f88b45cb1aba9edcab8f868"
// template, err := csdk.GetTemplate(templateID)
// checkError(err)
// ioutil.WriteFile(templateID+"-template.odt", template, 0644)

// cresp, err := csdk.DeleteTemplate(templateID)
// checkError(err)
// fmt.Printf("%+v", cresp)
// if cresp.Success == false {
// 	log.Fatal(cresp.Error)
// }

// cresp, err := csdk.RenderReport(templateID, `{"data":{"firstname":"Felix","lastname":"Arvid Ulf Kjellberg","color":"#00FF00"},"convertTo":"pdf"}`)
// if err != nil {
// 	log.Fatal(err)
// }

// fmt.Println("Success:", cresp.Success)
// if cresp.Success == false {
// 	log.Fatal(cresp.Error)
// }
// fmt.Printf("%+v", cresp.Data)

// // renderID := "MTAuMjAuMTEuMTEgICAg01E4NAFFCFXM0SE3KVVT8GAK1C.pdf"
// file, err := csdk.GetReport(cresp.Data.RenderID)
// if err != nil {
// 	log.Fatal(err)
// }
// fmt.Println("Final file:\n", file)
// ioutil.WriteFile(cresp.Data.RenderID, file, 0644)
// }

// ======= Create a file to debug
// by, e := ioutil.ReadAll(req.Body)
// if e != nil {
// 	log.Fatal(e)
// }
// err = ioutil.WriteFile("http.log", by, 0644)
// if err != nil {
// 	log.Fatal(err)
// }
