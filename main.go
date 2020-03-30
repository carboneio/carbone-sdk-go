package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"
)

// CarboneResponseData object created during Carbone Render response.
type CarboneResponseData struct {
	TemplateID            string `json:"templateId,omitempty"`
	RenderID              string `json:"renderId,omitempty"`
	TemplateFileExtension string `json:"inputFileExtension,omitempty"`
}

// CarboneResponse object created during Carbone Render response.
type CarboneResponse struct {
	Success bool                `json:"success"`
	Error   string              `json:"error,omitempty"`
	Data    CarboneResponseData `json:"data"`
}

// CarboneSDK to use Carbone render API easily.
type CarboneSDK struct {
	apiAccessToken string
	apiURL         string
	apiTimeOut     time.Duration
}

// NewCarboneSDK is a constructor and return a new instance of carboneSDK
func NewCarboneSDK(apiAccessToken string) (CarboneSDK, error) {
	csdk := CarboneSDK{}
	if apiAccessToken == "" {
		return csdk, errors.New("Carbone SDK constructor error: argument is missing: apiAccessToken")
	}
	csdk.apiAccessToken = apiAccessToken
	csdk.apiURL = "https://render.carbone.io"
	csdk.apiTimeOut = time.Second * 10
	return csdk, nil
}

// AddTemplate upload your template to Carbone Render.
func (csdk CarboneSDK) AddTemplate(templateFileName string, payload string) (CarboneResponse, error) {
	cResp := CarboneResponse{}
	// Create buffer
	buf := new(bytes.Buffer)
	// create a tmpfile and assemble your multipart from there
	w := multipart.NewWriter(buf)
	// Create the data object to send
	// { "payload":"", "template": readstream(file...) }
	label, err := w.CreateFormField("payload")
	if err != nil {
		return cResp, err
	}
	// Write payload content (empty for now)
	label.Write([]byte(payload))
	// Create the FormData
	fw, err := w.CreateFormFile("template", templateFileName)
	if err != nil {
		return cResp, err
	}
	// Open Template
	fd, err := os.Open(templateFileName)
	if err != nil {
		return cResp, err
	}
	defer fd.Close()
	// Write file field from file to upload
	_, err = io.Copy(fw, fd)
	if err != nil {
		return cResp, err
	}
	// Important if you do not close the multipart writer you will not have a terminating boundry
	w.Close()
	// Create the request
	req, err := http.NewRequest("POST", csdk.apiURL+"/template", buf)
	if err != nil {
		return cResp, errors.New("Carbone SDK request: failled to create a new request: " + err.Error())
	}
	// Set headers
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", csdk.apiAccessToken)
	// Set client timeout
	client := &http.Client{Timeout: csdk.apiTimeOut}
	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return cResp, errors.New("Carbone SDK request error: " + strconv.Itoa(resp.StatusCode))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return cResp, fmt.Errorf("Carbone SDK request error: status code %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return cResp, errors.New("Carbone SDK request error: failled to read the body: " + err.Error())
	}

	err = json.Unmarshal(body, &cResp)
	if err != nil {
		return cResp, errors.New("Carbone SDK request error: failled to parse the JSON response from the body: " + err.Error())
	}
	return cResp, nil
}

// Render a report from a templateID and a json data
func (csdk CarboneSDK) Render(templateID string, jsonData string) (CarboneResponse, error) {
	cResp := CarboneResponse{}
	req, err := http.NewRequest("POST", csdk.apiURL+"/render/"+templateID, bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		return cResp, errors.New("Carbone SDK request: failled to create a new request: " + err.Error())
	}
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", csdk.apiAccessToken)
	// Set client timeout
	client := &http.Client{Timeout: csdk.apiTimeOut}
	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return cResp, fmt.Errorf("Carbone SDK request error: status code %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return cResp, errors.New("Carbone SDK request error: failled to read the body: " + err.Error())
	}
	err = json.Unmarshal(body, &cResp)
	if err != nil {
		return cResp, errors.New("Carbone SDK request error: failled to parse the JSON response from the body: " + err.Error())
	}
	return cResp, nil
}

// GetReport Request Carbone Render and return a generated report
func (csdk CarboneSDK) GetReport(renderID string) ([]byte, error) {
	req, err := http.NewRequest("GET", csdk.apiURL+"/render/"+renderID, nil)
	if err != nil {
		return []byte{}, errors.New("Carbone SDK request: failled to create a new request: " + err.Error())
	}
	req.Header.Set("Authorization", csdk.apiAccessToken)
	// Set client timeout
	client := &http.Client{Timeout: csdk.apiTimeOut}
	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, fmt.Errorf("Carbone SDK request error: status code %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	// Read the response data and return a []byte. The http package automatically decodes chunking when reading response body.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, errors.New("Carbone SDK request error: failled to read the body: " + err.Error())
	}
	if len(body) == 0 {
		return body, errors.New("Carbone SDK request error: The response body is empty: Render again and generate a new renderId")
	}
	return body, nil
}

func main() {
	csdk, err := NewCarboneSDK("eyJhbGciOiJFUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiIxNjY3IiwiYXVkIjoiY2FyYm9uZSIsImV4cCI6MjIwNzQwNjQ0NywiZGF0YSI6eyJpZEFjY291bnQiOjE2Njd9fQ.AH2NiPdd8dRC_FNsd4aJ1DHy2wNNhXFmRvyh6PM-jkksfPn7hIIgiUfZ-L7Ng9Jou3eCeLrymjcPuABFVcaGiGvCATAICKX_j7WKBdMO_iPzD1LvL5j35FX1_i513OLqSvqTY_3KvBZO2RXMh4tLWlMn-dhNFLn-aE6IcS3lpce_A2PB")
	if err != nil {
		log.Fatal(err)
	}

	// cresp, err := csdk.AddTemplate("./template.odt", "")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// templateID := "f90e67221d7d5ee11058a000bdb997fb41bf149b1f88b45cb1aba9edcab8f868"
	// cresp, err := csdk.Render(templateID, `{"data":{"firstname":"Steeve","lastname":"Payraudeau"},"convertTo":"pdf"}`)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("Success:", cresp.Success)
	// if cresp.Success == false {
	// 	log.Fatal(cresp.Error)
	// }
	// fmt.Printf("%+v", cresp.Data.RenderID)

	renderID := "MTAuMjAuMTEuMTEgICAg01E4NAFFCFXM0SE3KVVT8GAK1C.pdf"
	file, err := csdk.GetReport(renderID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Final file:\n", file)
	ioutil.WriteFile("2"+renderID, file, 0644)
}

// ======= Create a file to debug
// by, e := ioutil.ReadAll(req.Body)
// if e != nil {
// 	log.Fatal(e)
// }
// err = ioutil.WriteFile("http.log", by, 0644)
// if err != nil {
// 	log.Fatal(err)
// }
