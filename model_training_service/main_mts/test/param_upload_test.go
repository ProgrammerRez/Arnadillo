package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mts/jobs"
	log "mts/logging"
	"reflect"
	"testing"

	// "mts/jobs"
	"net/http"
	"net/url"
	"strings"
)

// Struct for Testing
type TestingStuct struct{
	name 			string
	endpoint 		string
	requestType		string
	expectedOutput	any
	input 			map[string]any
	expectedError	bool
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


// This function is for testing the Add Job Handler
func TestGetEndpoint(t *testing.T){

	// Start the logger
	log.Init()
	log.SetLevel(log.InfoLevel)
	log.Info("Testing the Add Job Handler")


	// Generating Valid Tests First

	// Creating Working Struct
	working_struct := map[string]any{
		"name": "session1",
		"urgency": 1,
		"model_type": "traditional",
		"model_name": "LogisticRegression",
		"model_params": map[string]any{
			"limit": 1000,
		},
	}

	// Working Job Result
	working_result_slice := []jobs.Job{
		{
			Status: "queued",
			ID: 1,
			DataParams: nil,
			TrainingParams: jobs.TrainingParameters{
				ModelType: "traditional",
				ModelName: "LogisticRegression",
				ModelParams: map[string]any{
					"limit":1000.0, 
				},
			},
			DoOverParams: jobs.DoOverParameters{},
			EvalParams: jobs.EvaluationParameters{},
		},
	}

	get_jobs_filter := map[string]any{
		"filter": 1,
	}

	// Creating Test Scenarios

	tests := []TestingStuct{
		{
			name: "Add Job Struct",
			endpoint: "/data_params",
			requestType: "POST",
			input: working_struct,
			expectedOutput: "Worked :>",
			expectedError: false,
		},

		{
			name: "Get Job Struct",
			endpoint: "/get_jobs",
			requestType: "POST",
			input: get_jobs_filter,
			expectedOutput: working_result_slice,
			expectedError: false,
		},
	}

	runTests(t, tests)

}


// This function runs all the provided tests
func runTests(t *testing.T, tests []TestingStuct){
	for _, test :=  range(tests){
		t.Run(test.name, func(t *testing.T) {
	
			// Run the command
			data, err := GetEndpoint(test.endpoint, test.requestType, test.input)

			if test.expectedError{
				if err == nil{
					log.Error("Error was expected but didn't occur")
					t.Fatalf("Error was expected but didn't occur")
				}
				return
			}

			if err != nil{
				log.Error(err.Error())
				t.Fatalf("Unexpected Error Occurred: %v", err)
			}

			if test.expectedOutput != nil{
				json_bytes, err := json.Marshal(data)
				if err != nil{
					log.Error(err.Error())
					t.Fatalf("%s", fmt.Sprintf("%v",err.Error()))
					return
				}
				target_type := reflect.TypeOf(test.expectedOutput)
				val := reflect.New(target_type)

				if json.Unmarshal(json_bytes, val.Interface()) == nil{
					data = val.Elem().Interface()
				}
			}

			if !reflect.DeepEqual(data, test.expectedOutput){
				log.Error("Result Mismatch")
				t.Fatalf("Result Mismatch: Wanted: %v, Got: %v", test.expectedOutput, data)
			}
		})
	}
}