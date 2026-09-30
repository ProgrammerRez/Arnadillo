package main

import (
	"encoding/json"
	"errors"
	"strconv"
	"io"
	jobs "mts/jobs"
	log "mts/logging"
	"net/http"
	"net/url"
	"strings"
)

const (
	DATA_ENDPOINT = "http://localhost:8000/the_good_stuff"
)

// This is a function to Initialize MTS Object
func InitMTS() *ModelTrainingService{
	return &ModelTrainingService{
		JobsQueue: jobs.InitializeQueue(),
	}
}

// This function will recieve data and create a job for the request
func (mts *ModelTrainingService) addJob(w http.ResponseWriter, r *http.Request){

	log.Info("Adding New Job")

	name := r.FormValue("name")
	raw_data, err := mts.getData(name)

	if err != nil{
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Info("Converted the Data Successfully")
	queue := mts.JobsQueue
	data_params, _ := jobs.CreateDataParamObject(raw_data)
	queue.CreateNewJob(
		data_params[0],
		jobs.TrainingParameters{},
		jobs.EvaluationParameters{},
		jobs.DoOverParameters{},
		1,
	)

	w.Write([]byte("Worked :>"))
}

// This function will retrieve all the jobs created
func (mts *ModelTrainingService) getJobs(w http.ResponseWriter, r *http.Request){

	value := r.FormValue("filter")

	if strings.TrimSpace(value) == ""{
		log.Error("Filter not Provided, Defaulting to All")
		return
	}

	priority_filter, err := strconv.ParseInt(value, 10, 32)
	if err != nil{
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	log.Info("Parsed Filter Input")

	list, err := mts.JobsQueue.GetJobs(int32(priority_filter))
	if err != nil{
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	
	w.Header().Set("Content-Type","application/json")
	bytes, err := json.Marshal(list)
	if err != nil{
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}

	w.Write(bytes) 
}


// ----------------------------------------------------------------------------------------------
// These functions are for Data Exchange between both the services
// ----------------------------------------------------------------------------------------------


// This function will get the necessary data for the job creation
func (mts *ModelTrainingService) getData(name string) (map[string]any, error){

	log.Info("Retrieving Data")

	// Check for Session Name
	if strings.TrimSpace(name) == ""{
		log.Error("Session Name was not Provided")
		return nil, errors.New("Session Name was not Provided")
	}

	// Creating the form for the arguments
	form := url.Values{
		"name":{name},
	}

	// Sending the PostForm Request for it
	resp, err := http.PostForm(DATA_ENDPOINT, form)
	if err != nil{
		log.Error(err.Error())
		return nil, errors.New(err.Error())
	}

	defer resp.Body.Close()

	// Waiting for Data Transmission
	var incomingData map[string]any
	
	body, err := io.ReadAll(resp.Body)
	if err != nil{
		log.Error(err.Error())
		return nil, errors.New(err.Error())
	}

	err = json.Unmarshal(body, &incomingData)
	if err != nil{
		log.Error(err.Error())
		return nil, errors.New(err.Error())
	}

	return incomingData, nil 
}