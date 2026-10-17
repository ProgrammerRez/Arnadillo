package utils

import(
	"fmt"
	"strings"
	"errors"
	log "mts/logging"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

// Struct for Testing
type TestingStuct struct{
	Name 			string
	Endpoint 		string
	RequestType		string
	ExpectedOutput	any
	Input 			map[string]any
	ExpectedError	bool
}

const ENDPOINT_TEMPLATE string = "http://localhost:8001"

// This function is for a POST request to any endpoint pointed
func GetEndpoint(endpoint, request_type string, data map[string]any) (any, error){
	
	// Starting logging
	log.SetLevel(log.InfoLevel)
	log.Info("Retrieving Test Data")

	// Checking Endpoint 
	if strings.TrimSpace(endpoint) == ""{
		log.Error("No Endpoint was Provided for Testing")
		return nil, errors.New("No Endpoint was Provided for Testing")
	}

	final_url := ENDPOINT_TEMPLATE + endpoint

	
	// Checking and Routing for Request Type
	if strings.TrimSpace(request_type) == ""{
		log.Error("No Request Type Provided")
		return  nil, errors.New("No Request Type Provided")
	}
	
	switch request_type{
	case "GET":
		log.Info(fmt.Sprintf("Starting GET Request on Endpoint %s", endpoint))
		resp, err := http.Get(final_url)
		if err != nil{
			log.Error(err.Error())
			return nil, errors.New(err.Error())
		}
		defer resp.Body.Close()
		
		var incomingData map[string]any
		body, err := io.ReadAll(resp.Body)
		
		if err != nil{
			log.Error(err.Error())
			return nil, errors.New(err.Error())
		}
		
		err = json.Unmarshal(body, &incomingData)
		if err != nil{
			log.Error(err.Error())
			return string(body), nil
		}
		return incomingData, nil
		
	case "POST":
		log.Info(fmt.Sprintf("Starting POST Request on Endpoint %s", endpoint))
		// Now Creating the form
		if len(data) == 0{
			log.Error("No Data Parameters Provided")
			return nil, errors.New("No Data Parameters Provided")
		}
		form := url.Values{}

		for key, val := range data{
			switch v := val.(type){
			case string:
				form.Add(key, v)
			default:
				// Then Add as JSON string
				data, err := json.Marshal(v)
				if err != nil{
					log.Error(fmt.Sprintf("Unable to Marshal data %v", v))
					return nil, fmt.Errorf("Unable to Marshal data %v", v)
				}
				form.Add(key, string(data))
			}
		}
		resp, err := http.PostForm(final_url, form)
		if err != nil{
			log.Error(err.Error())
			return nil, errors.New(err.Error())
		}
		defer resp.Body.Close()

		var incomingData any
		body, err := io.ReadAll(resp.Body)
		
		if err != nil{
			log.Error(err.Error())
			return nil, errors.New(err.Error())
		}

		err = json.Unmarshal(body, &incomingData)
		if err != nil{
			log.Error(err.Error())
			return string(body), nil
		}

		return incomingData, nil

	default: 
		log.Error("Method Not Supported")
		return nil, errors.New("Method Not Supported")
	}
}
