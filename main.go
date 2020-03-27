package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

type carboneResponseData struct {
	TemplateID string `json:"templateId,omitempty"`
	RenderID   string `json:"renderId,omitempty"`
}

type carboneResponse struct {
	Success bool                `json:"success"`
	Error   string              `json:"error,omitempty"`
	Data    carboneResponseData `json:"data"`
}

type carboneSDK struct {
	apiAccessToken string
	apiURL         string
}

// addTemplate upload your template to Carbone Render.
func (csdk carboneSDK) addTemplate(templateFileName string, payload string) (carboneResponse, error) {
	fmt.Println("Payload:", payload)
	// Create buffer
	buf := new(bytes.Buffer)
	// create a tmpfile and assemble your multipart from there
	w := multipart.NewWriter(buf)

	// Create the data to send
	// {
	//		"payload":"",
	// 		"template": readstream(file...)
	// }
	label, err := w.CreateFormField("payload")
	if err != nil {
		log.Fatal("Error:", err)
	}
	// Write payload content (empty for now)
	label.Write([]byte(payload))

	fw, err := w.CreateFormFile("template", templateFileName)
	if err != nil {
		log.Fatal(err)
	}

	// Open Template
	fd, err := os.Open(templateFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()
	// Write file field from file to upload
	_, err = io.Copy(fw, fd)
	if err != nil {
		log.Fatal(err)
	}
	// Important if you do not close the multipart writer you will not have a terminating boundry
	w.Close()

	req, err := http.NewRequest("POST", csdk.apiURL, buf)
	if err != nil {
		log.Fatal("Error POST request:", err)
	}

	// Set headers
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", csdk.apiAccessToken)

	// Set client timeout
	client := &http.Client{Timeout: time.Second * 10}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()

	// fmt.Println("response Status:", resp.Status)

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("bad status: %s", resp.Status)
		fmt.Println("response Headers:", resp.Header)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body. ", err)
	}
	fmt.Printf("%s\n", body)
	cResp := carboneResponse{}
	err = json.Unmarshal(body, &cResp)
	if err != nil {
		return cResp, err
	}
	return cResp, nil
}

// func

func main() {
	csdk := carboneSDK{
		apiURL:         "https://render.carbone.io/template",
		apiAccessToken: "eyJhbGciOiJFUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiIxNjY3IiwiYXVkIjoiY2FyYm9uZSIsImV4cCI6MjIwNzQwNjQ0NywiZGF0YSI6eyJpZEFjY291bnQiOjE2Njd9fQ.AH2NiPdd8dRC_FNsd4aJ1DHy2wNNhXFmRvyh6PM-jkksfPn7hIIgiUfZ-L7Ng9Jou3eCeLrymjcPuABFVcaGiGvCATAICKX_j7WKBdMO_iPzD1LvL5j35FX1_i513OLqSvqTY_3KvBZO2RXMh4tLWlMn-dhNFLn-aE6IcS3lpce_A2PB",
	}
	templateFileName := "./template.odt"
	cresp, err := csdk.addTemplate(templateFileName, "")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Success:", cresp.Success)
	if cresp.Success == false {
		log.Fatal(cresp.Error)
	} else {
		fmt.Printf("%+v", cresp.Data)
	}
}

// templateId: f90e67221d7d5ee11058a000bdb997fb41bf149b1f88b45cb1aba9edcab8f868
// templateId avec payload 1234: dd226478563e4e1f8a2c38b97c71005b68cdb0f45ce9f9c2155aae4b4fd341d2

// ======= Log cookie and headers are attached
// fmt.Println(req.Cookies())
// fmt.Println(req.Header)
// fmt.Println(req.Body)
// ======= Create a file to debug
// by, e := ioutil.ReadAll(req.Body)
// if e != nil {
// 	log.Fatal(e)
// }
// err = ioutil.WriteFile("http.log", by, 0644)
// if err != nil {
// 	log.Fatal(err)
// }
