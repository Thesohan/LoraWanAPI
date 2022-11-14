package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type LorawanService struct{}

// RegisterDevice calls the Lorawan API to register a new device.
// It returns true/false based upon if the device is successfully registered or not.
func (s *LorawanService) RegisterDevice(requestPayload LorawanAPIRequestBody, count int) bool {
	//Encode the data
	postBody, _ := json.Marshal(requestPayload)
	responseBody := bytes.NewBuffer(postBody)
	//Leverage Go's HTTP Post function to make request
	url := LorawanBaseURL + RegisterDeviceEndPoint
	resp, err := http.Post(url, "application/json", responseBody)
	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	fmt.Println(" ", resp.StatusCode, " ", count, " ", requestPayload.Deveui)
	return resp.StatusCode == http.StatusOK
}
