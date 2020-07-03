// Package carbone provide an SDK to communicate with Carbone Render
// Carbone is the most efficient report generator
// It render from a JSON and template into PDF, DOCX, XLSX, PPTX, ODS and many more reports
package carbone

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
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
	apiVersion     string
	apiAccessToken string
	apiURL         string
	apiTimeOut     time.Duration
	apiHTTPClient  *http.Client
}

// NewCarboneSDK is a constructor and return a new instance of CSDK
func NewCarboneSDK(args ...string) (*CSDK, error) {
	apiAccessToken := os.Getenv("CARBONE_TOKEN")
	if len(args) > 0 && args[0] != "" {
		apiAccessToken = args[0]
	}
	if apiAccessToken == "" {
		return nil, errors.New(`NewCarboneSDK error: "apiAccessToken" argument OR "CARBONE_TOKEN" env variable is missing`)
	}
	csdk := &CSDK{
		apiVersion:     "2",
		apiAccessToken: apiAccessToken,
		apiURL:         "https://render.carbone.io",
		apiTimeOut:     time.Second * 10,
		apiHTTPClient:  &http.Client{Timeout: time.Second * 10},
	}
	return csdk, nil
}

// AddTemplate upload your template to Carbone Render. The first parameter is the template file path, the second is an optional payload.
func (csdk *CSDK) AddTemplate(templateFileName string, args ...string) (APIResponse, error) {
	payload := ""
	if len(args) > 0 && args[0] != "" {
		payload = args[0]
	}
	cResp := APIResponse{}
	if templateFileName == "" {
		return cResp, errors.New("Carbone SDK AddTemplate error: argument is missing: templateFileName")
	}
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
	headerRequest := map[string]string{
		"Content-Type": w.FormDataContentType(),
	}
	resp, err := csdk.doHTTPRequest("POST", csdk.apiURL+"/template", headerRequest, buf)
	if err != nil {
		return cResp, err
	}
	// Read the stream
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return cResp, errors.New("Carbone SDK request error: failled to read the body: " + err.Error())
	}
	// Close the connection https://stackoverflow.com/questions/33238518/what-could-happen-if-i-dont-close-response-body
	defer resp.Body.Close()
	// Parse JSON body and store into the APIResponse Struct
	err = json.Unmarshal(body, &cResp)
	if err != nil {
		return cResp, errors.New("Carbone SDK request error: failled to parse the JSON response from the body: " + err.Error())
	}
	return cResp, nil
}

// GetTemplate returns the original template from the templateId (Unique identifier of the template)
func (csdk *CSDK) GetTemplate(templateID string) ([]byte, error) {
	if templateID == "" {
		return []byte{}, errors.New("Carbone SDK GetTemplate error: argument is missing: templateID")
	}
	// Create the request
	resp, err := csdk.doHTTPRequest("GET", csdk.apiURL+"/template/"+templateID, nil, nil)
	if err != nil {
		return []byte{}, err
	}
	// Read the response data and return a []byte. The http package automatically decodes chunking when reading response body.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, errors.New("Carbone SDK GetTemplate request error: failled to read the body: " + err.Error())
	}
	// Close the connection
	defer resp.Body.Close()
	if len(body) == 0 {
		return body, errors.New("Carbone SDK GetTemplate request error: The response body is empty")
	}
	return body, nil
}

// DeleteTemplate Delete an uploaded template from a templateID.
func (csdk *CSDK) DeleteTemplate(templateID string) (APIResponse, error) {
	cResp := APIResponse{}
	if templateID == "" {
		return cResp, errors.New("Carbone SDK DeleteTemplate error: argument is missing: templateID")
	}
	// HTTP Request
	resp, err := csdk.doHTTPRequest("DELETE", csdk.apiURL+"/template/"+templateID, nil, nil)
	if err != nil {
		return cResp, err
	}
	// Read body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return cResp, errors.New("Carbone SDK DeleteTemplate request error: failled to read the body: " + err.Error())
	}
	// Close the connection
	defer resp.Body.Close()
	// Parse JSON body and store into the APIResponse Struct
	err = json.Unmarshal(body, &cResp)
	if err != nil {
		return cResp, errors.New("Carbone SDK DeleteTemplate request error: failled to parse the JSON response from the body: " + err.Error())
	}
	return cResp, nil
}

// RenderReport a report from a templateID and a json data
func (csdk *CSDK) RenderReport(templateID string, jsonData string) (APIResponse, error) {
	cResp := APIResponse{}
	if templateID == "" {
		return cResp, errors.New("Carbone SDK RenderReport error: argument is missing: templateID")
	}
	if jsonData == "" {
		return cResp, errors.New("Carbone SDK RenderReport error: argument is missing: jsonData")
	}
	headerRequest := map[string]string{
		"Content-Type": "application/json",
	}
	resp, err := csdk.doHTTPRequest("POST", csdk.apiURL+"/render/"+templateID, headerRequest, bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		return cResp, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return cResp, errors.New("Carbone SDK RenderReport request error: failled to read the body: " + err.Error())
	}
	// Close the connection
	defer resp.Body.Close()
	err = json.Unmarshal(body, &cResp)
	if err != nil {
		return cResp, errors.New("Carbone SDK RenderReport request error: failled to parse the JSON response from the body: " + err.Error())
	}
	return cResp, nil
}

// GetReport Request Carbone Render and return a generated report
func (csdk *CSDK) GetReport(renderID string) ([]byte, error) {
	if renderID == "" {
		return []byte{}, errors.New("Carbone SDK GetReport error: argument is missing: renderID")
	}
	// http request
	resp, err := csdk.doHTTPRequest("GET", csdk.apiURL+"/render/"+renderID, nil, nil)
	if err != nil {
		return []byte{}, err
	}
	// Read the response data and return a []byte. The http package automatically decodes chunking when reading response body.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, errors.New("Carbone SDK GetReport request error: failled to read the body: " + err.Error())
	}
	// Close the connection
	defer resp.Body.Close()
	if len(body) == 0 {
		return body, errors.New("Carbone SDK GetReport request error: The response body is empty: Render again and generate a new renderId")
	}
	return body, nil
}

// Render render a report from a templateID OR a template path.
// pathOrTemplateID {string}: Accept a file path OR a template ID returned by AddTemplate
// jsonData {string}: stringify json, all options here: https://carbone.io/api-reference.html#rendering-a-report
// args {...string}: You can pass an optinal payload used during the template upload (AddTemplate) to create a different templateID.
// It returns a []byte of the file.
func (csdk *CSDK) Render(pathOrTemplateID string, jsonData string, args ...string) ([]byte, error) {
	var cresp APIResponse
	var er error
	payload := ""
	if len(args) > 0 && args[0] != "" {
		payload = args[0]
	}
	info, err := os.Stat(pathOrTemplateID)
	if os.IsNotExist(err) == true {
		// The first argument `pathOrTemplateID` is a templateID
		cresp, er = csdk.RenderReport(pathOrTemplateID, jsonData)
		if er != nil {
			return []byte{}, er
		}
	} else if info.IsDir() == true {
		return []byte{}, errors.New("Carbone SDK Render error: the path passed as argument is a directory")
	} else {
		// The first argument `pathOrTemplateID` is maybe a file
		templateID, e := csdk.GenerateTemplateID(pathOrTemplateID, payload)
		if e != nil {
			return []byte{}, errors.New("Carbone SDK Render error: failled to generate the templateID hash:" + e.Error())
		}
		cresp, er = csdk.RenderReport(templateID, jsonData)
		if er != nil {
			return []byte{}, er
		} else if cresp.Success == false {
			// if RenderReport return one of the following error, it means the template does not exist
			// - Error while rendering template Error: ENOENT:File not found
			// - Error while rendering template Error: 404 Not Found
			// Then call add template and render again
			cres, e := csdk.AddTemplate(pathOrTemplateID, payload)
			if e != nil {
				return []byte{}, errors.New("Carbone SDK Render error:" + e.Error())
			}
			cresp, er = csdk.RenderReport(cres.Data.TemplateID, jsonData)
			if er != nil {
				return []byte{}, errors.New("Carbone SDK Render error:" + er.Error())
			}
		}
	}
	if cresp.Success == false {
		// If an error is returned, it means something went wrong.
		// if the error is "Error while rendering template Error: 404 Not Found" or "ENOENT:File not found" it means TemplateID does not exist
		return []byte{}, errors.New(cresp.Error)
	}
	if len(cresp.Data.RenderID) <= 0 {
		return []byte{}, errors.New("Carbone SDK Render error: renderID is empty")
	}
	// Return the report
	return csdk.GetReport(cresp.Data.RenderID)
}

// GenerateTemplateID Generate the templateID from a template
// filepath {string}: file path
// args {...string}: You can set a payload (args[0])
// returns the file TemplateId
func (csdk *CSDK) GenerateTemplateID(filepath string, args ...string) (string, error) {
	payload := ""
	if len(args) > 0 && args[0] != "" {
		payload = args[0]
	}
	// Open the file
	f, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	// New HASH
	h := sha256.New()
	// Write payload
	h.Write([]byte(payload))
	// Write file buffer
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	// Return the sha256 has as hexadecimal
	return hex.EncodeToString(h.Sum(nil)), nil
}

// SetAccessToken set the Carbone Render access token
func (csdk *CSDK) SetAccessToken(newToken string) {
	csdk.apiAccessToken = newToken
}

// SetAPIVersion set the Carbone Render version
func (csdk *CSDK) SetAPIVersion(version int) {
	csdk.apiVersion = strconv.Itoa(version)
}

// GetAPIVersion get the Carbone Render version
func (csdk *CSDK) GetAPIVersion() (int, error) {
	return strconv.Atoi(csdk.apiVersion)
}

// ------------------ private function
func (csdk *CSDK) doHTTPRequest(method string, url string, headers map[string]string,
	body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, errors.New("Carbone SDK request: failled to create a new request: " + err.Error())
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// User Api Token
	req.Header.Set("Authorization", "Bearer "+csdk.apiAccessToken)
	req.Header.Set("carbone-version", csdk.apiVersion)

	// Send request
	resp, err := csdk.apiHTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Carbone SDK request error: %v", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Carbone SDK request error: status code %d", resp.StatusCode)
	}
	return resp, nil
}
