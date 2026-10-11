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
	raw_int := r.FormValue("urgency")

	urgency, err := strconv.ParseInt(raw_int, 10, 32)
	if err != nil{
		log.Error("Please provide a proper Urgency Integer")
		http.Error(w, "Please provide a proper Urgency Integer", http.StatusBadRequest)
		return
	}

	// Now Adding the Training Parameters Support

	model_type := r.FormValue("model_type")
	model_name := r.FormValue("model_name")
	model_params := r.FormValue("model_params")

	// Unmarshalling the model_params map
	var model_params_map map[string]any
	if err := json.Unmarshal([]byte(model_params), &model_params_map); err != nil{
		log.Error("Model Params is Invalid Type")
		http.Error(w, "Model Params is Invalid Type", http.StatusBadRequest)
		return
	}

	train_data := map[string]any{
		"model_type": model_type,
		"model_name": model_name,
		"model_params": model_params_map,
	}

	// Now Adding the Evaluation Parameters Support

	problem_type := r.FormValue("problem_type")
	tests_required := r.FormValue("tests_required")
	

	var tests_required_map map[string]float64
	if err := json.Unmarshal([]byte(tests_required), &tests_required_map); err != nil{
		log.Error("Invalid Type for Required Tests")
		http.Error(w, "Invalid Type for Required Tests", http.StatusBadRequest)
		return
	}
	

	eval_data := map[string]any{
		"problem_type": problem_type,
		"tests_required": tests_required_map,
	}

	// Now Adding Do Over Parameters Support

	fallback_params := r.FormValue("fallback_params")
	max_retries := r.FormValue("max_retries")
	max_deviance := r.FormValue("max_deviance")

	// Converting to valid numerical types
	max_retries_int, err := strconv.ParseInt(max_retries, 10, 64)
	if err != nil{
		log.Error("Unable to Parse Max Retries")
		http.Error(w, "Unable to Parse Max Retries", http.StatusBadRequest)
		return
	}

	max_deviance_float, err := strconv.ParseFloat(max_deviance, 64)
	if err != nil{
		log.Error("Unable to Parse Max Deviance")
		http.Error(w, "Unable to Parse Max Deviance", http.StatusBadRequest)
		return
	}


	var fallback_params_map	map[string]any
	if err := json.Unmarshal([]byte(fallback_params), &fallback_params_map); err != nil{
		log.Error("Invalid Type for Fallback Params")
		http.Error(w, "Invalid Type for Fallback Params", http.StatusBadRequest)
		return
	}

	do_over_data := map[string]any{
		"fallback_params": fallback_params_map,
		"max_retries": max_retries_int,
		"max_deviance": max_deviance_float,
	}

	// Then finally gathering all the Parameters Data
	train_params, err := jobs.CreateTrainParamObject(train_data)
	if err != nil{
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return 
	}
	log.Info("Training Parameters Extracted")

	eval_params, err := jobs.CreateEvalParamObjects(eval_data)
	if err != nil{
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	do_over_params, err := jobs.CreateDoOverParamObjects(do_over_data)
	if err != nil{
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Getting the Data Parameters
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
		data_params,
		*train_params,
		*eval_params,
		*do_over_params,
		int(urgency),
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