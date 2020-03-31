package carbone

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"
)

// APIResponseData object created during Carbone Render response.
type APIResponseData struct {
	TemplateID            string `json:"templateId,omitempty"`
	RenderID              string `json:"renderId,omitempty"`
	TemplateFileExtension string `json:"inputFileExtension,omitempty"`
}

// APIResponse object created during Carbone Render response.
type APIResponse struct {
	Success bool            `json:"success"`
	Error   string          `json:"error,omitempty"`
	Data    APIResponseData `json:"data"`
}

// CSDK (CarboneSDK) to use Carbone render API easily.
type CSDK struct {
	apiAccessToken string
	apiURL         string
	apiTimeOut     time.Duration
	apiHTTPClient  *http.Client
}

// NewCarboneSDK is a constructor and return a new instance of CSDK
func NewCarboneSDK(apiAccessToken string) (*CSDK, error) {
	if apiAccessToken == "" {
		return nil, errors.New("Carbone SDK constructor error: argument is missing: apiAccessToken")
	}
	csdk := &CSDK{
		apiAccessToken: apiAccessToken,
		apiURL:         "https://render.carbone.io",
		apiTimeOut:     time.Second * 10,
		apiHTTPClient:  &http.Client{Timeout: time.Second * 10},
	}
	return csdk, nil
}

// AddTemplate upload your template to Carbone Render.
func (csdk *CSDK) AddTemplate(templateFileName string, payload string) (APIResponse, error) {
	cResp := APIResponse{}
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

// GetTemplate returns the original template from the templateId (Unique identifier of the template)
func (csdk *CSDK) GetTemplate(templateID string) ([]byte, error) {
	req, err := http.NewRequest("GET", csdk.apiURL+"/template/"+templateID, nil)
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
		return body, errors.New("Carbone SDK request error: The response body is empty")
	}
	return body, nil
}

// DeleteTemplate Delete an uploaded template from a templateID.
func (csdk *CSDK) DeleteTemplate(templateID string) (APIResponse, error) {
	cResp := APIResponse{}
	// Create client and set timeout
	client := &http.Client{Timeout: csdk.apiTimeOut}
	// Create Request
	req, err := http.NewRequest("DELETE", csdk.apiURL+"/template/"+templateID, nil)
	if err != nil {
		return cResp, errors.New("Carbone SDK request: failled to create a new request: " + err.Error())
	}
	req.Header.Set("Authorization", csdk.apiAccessToken)
	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return cResp, fmt.Errorf("Carbone SDK request error: status code %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	// Read body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return cResp, errors.New("Carbone SDK request error: failled to read the body: " + err.Error())
	}
	// Parse JSON body and store into the APIResponse Struct
	err = json.Unmarshal(body, &cResp)
	if err != nil {
		return cResp, errors.New("Carbone SDK request error: failled to parse the JSON response from the body: " + err.Error())
	}
	return cResp, nil

}

// RenderReport a report from a templateID and a json data
func (csdk *CSDK) RenderReport(templateID string, jsonData string) (APIResponse, error) {
	cResp := APIResponse{}
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
func (csdk *CSDK) GetReport(renderID string) ([]byte, error) {
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

// ------------------ private function

func (csdk *CSDK) doHTTPRequest(method, url string, headers map[string]string,
	body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// User Api Token
	req.Header.Set("Authorization", csdk.apiAccessToken)

	// https://code.google.com/p/go/issues/detail?id=6738
	// if method == "PUT" || method == "POST" {
	// 	length := req.Header.Get("Content-Length")
	// 	if length != "" {
	// 		req.ContentLength, _ = strconv.ParseInt(length, 10, 64)
	// 	}
	// }

	// Send request
	resp, err := csdk.apiHTTPClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("Carbone SDK request error: status code %d: %v", resp.StatusCode, err.Error())
	}
	defer resp.Body.Close()

	return resp, nil
}
