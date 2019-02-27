package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"net/http"
	"bytes"
	"log"
	"io/ioutil"
)

type SavedObjectHeader struct {
        Id string `json:"id,omitempty"`
        ObjectType string `json:"type,omitempty"`
        UpdatedAt string `json:"updated_at,omitempty"`
        Version int `json:"version,omitempty"`
	Attributes map[string]interface{} `json:"attributes"`
}




func getKibClient(d *schema.ResourceData, meta interface{}) *http.Client {
	return &http.Client{}
}

func genericKibRequest(requestType string, d *schema.ResourceData, meta interface{}, url string, username string, password string, buffer string) (*[]byte, error ) {

	client := getKibClient(d, meta)

	var req *http.Request
	var err error
	if buffer != "" {
		log.Printf("Reauest %s, %s => %s", requestType, url, buffer)
		req, err = http.NewRequest(requestType, url, bytes.NewBufferString(buffer))
	} else {
		log.Printf("Reauest %s, %s", requestType, url)
		req, err = http.NewRequest(requestType, url, bytes.NewBufferString(buffer))
		req, err = http.NewRequest(requestType, url, nil)
	}
	if err != nil {
		return nil, err
	}

	if username != "" {
		req.SetBasicAuth(username, password)
	}

	req.Header.Add("kbn-xsrf", "true")
        req.Header.Add("content-type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	
	respBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("Respose: %s", respBody)
	if err != nil {
		return nil, err
	}
	
	return &respBody, nil
}

func postKibRequest(d *schema.ResourceData, meta interface{}, url string, username string, password string, buffer string) (*[]byte, error) {
	return genericKibRequest("POST", d, meta, url, username, password, buffer)
}

func getKibRequest(d *schema.ResourceData, meta interface{}, url string, username string, password string) (*[]byte, error) {
	return genericKibRequest("GET", d, meta, url, username, password, "")
}

func putKibRequest(d *schema.ResourceData, meta interface{}, url string, username string, password string, buffer string) (*[]byte, error) {
	return genericKibRequest("PUT", d, meta, url, username, password, buffer)
}

func deleteKibRequest(d *schema.ResourceData, meta interface{}, url string, username string, password string) (*[]byte, error) {
	return genericKibRequest("DELETE", d, meta, url, username, password, "")
}
