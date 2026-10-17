package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	jobs "mts/jobs"
	log "mts/logging"
	"net/http"
	"net/url"
	"strconv"
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

	urgency, err := strconv.ParseInt(raw_int, 10, 64)
	if err != nil{
		log.Error("Please provide a proper Urgency Integer")
		http.Error(w, "Please provide a proper Urgency Integer", http.StatusBadRequest)
		return
	}

	// Creating a Helper Function to Avoid a Messy Function
	parseJsonForm := func(form_key string, target any) error{
		val := r.FormValue(form_key)
		if strings.TrimSpace(val) == ""{
			log.Error(fmt.Sprintf("Form Value: %s can't be found in form", form_key))
			return fmt.Errorf("Form Value: %s can't be found in form", form_key)
		}
		return json.Unmarshal([]byte(val), target)
	}

	// Creating all the necessary variables 
	var dataParamsMap, trainParamsMap, evalParamsMap, doOverParamsMap map[string]any


	// Then Unpacking all the parameters
	if err := parseJsonForm("data_params", &dataParamsMap); err != nil{
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if err := parseJsonForm("train_params", &trainParamsMap); err != nil{
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := parseJsonForm("eval_params", &evalParamsMap); err != nil{
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := parseJsonForm("do_over_params", &doOverParamsMap); err != nil{
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Then finally gathering all the Parameters Data
	
	
	// 1. Resolve Data Parameters Issue First
	var final_data_params []jobs.DataParameters
	
	data_params_value := r.FormValue("data_params")

	if (
		len(dataParamsMap) == 0 || 
		strings.TrimSpace(data_params_value) == "" || 
		strings.TrimSpace(data_params_value) == "{}" ||
		strings.TrimSpace(data_params_value) == "[]"){

		log.Warn("Data Params is unavailable in form data. Routing to DMS API")
		// First the dms condition
		dmsData, err := mts.getData(name)
		if err != nil{
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Registry Existence Condition
		registry, _ := dmsData["registry"]
		if registry != nil{	
			final_data_params, err = jobs.CreateDataParamObject(dmsData)
			if err != nil{
				log.Error(err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}else{
			log.Warn("API Output is Nill, so returning nil")
			final_data_params = nil
		}
	}else{
		// Then the normal condition
		final_data_params, err = jobs.CreateDataParamObject(dataParamsMap)
		if err != nil{
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	log.Info("Data Parameters Extracted")

	
	final_train_params, err := jobs.CreateTrainParamObject(trainParamsMap)
	if err != nil{
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return 
	}
	log.Info("Training Parameters Extracted")


	final_eval_params, err := jobs.CreateEvalParamObjects(evalParamsMap)
	if err != nil{
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Info("Eval Params Extracted")

	final_do_over_params, err := jobs.CreateDoOverParamObjects(doOverParamsMap)
	if err != nil{
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Info("Do Over Params Extracted")

	mts.JobsQueue.CreateNewJob(
		final_data_params,
		*final_train_params,
		*final_eval_params,
		*final_do_over_params,
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

	priority_filter, err := strconv.ParseInt(value, 10, 64)
	if err != nil{
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	log.Info("Parsed Filter Input")

	list, err := mts.JobsQueue.GetJobs(int(priority_filter))
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

// This function will set the data ops for the object
func (mts *ModelTrainingService) setDataOps(w http.ResponseWriter, r *http.Request){

	log.Info("Setting / Adding Data Ops")

	// Gathering the Arguments 

	id := r.FormValue("id")
	file_id := r.FormValue("file_id")
	filter := r.FormValue("filter")

	num_ops_values := r.FormValue("num_ops")
	cat_ops_values := r.FormValue("cat_ops")

	// Checking for Null Values
	if strings.TrimSpace(num_ops_values) == ""{
		log.Error("Numerical Ops is Empty")
		http.Error(w, "Numerical Ops is Empty", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(cat_ops_values) == ""{
		log.Error("Categorical Ops is Empty")
		http.Error(w, "Categorical Ops is Empty", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(id) == ""{
		log.Error("ID is Empty")
		http.Error(w, "ID is Empty", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(filter) == ""{
		log.Error("Filter is Empty")
		http.Error(w, "Filter is Empty", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(file_id) == ""{
		log.Error("File ID is Empty")
		http.Error(w, "File ID is Empty", http.StatusBadRequest)
		return
	}

	// Now Unmarshalling and converting the data into suitable types

	id_int, err := strconv.ParseInt(id, 10, 64)
	if err != nil{
		log.Error("Unable to parse ID")
		http.Error(w, "Unable to parse ID", http.StatusBadRequest)
	}

	filter_int, err := strconv.ParseInt(filter, 10, 64)
	if err != nil{
		log.Error("Unable to parse filter")
		http.Error(w, "Unable to parse filter", http.StatusBadRequest)
	}

	file_id_int, err := strconv.ParseInt(file_id, 10, 64)
	if err != nil{
		log.Error("Unable to parse file_id")
		http.Error(w, "Unable to parse file_id", http.StatusBadRequest)
	}

	var num_ops []string
	var cat_ops map[string][]string

	err = json.Unmarshal([]byte(num_ops_values), &num_ops)
	if err != nil{
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal([]byte(cat_ops_values), &cat_ops)
	if err != nil{
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Now Finding the Job
	job_pointer, err := mts.JobsQueue.FindJob(int(id_int), int(filter_int))
	if err != nil{
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Now Setting the Ops
	job_pointer.DataParams[int(file_id_int)].SetOps(num_ops, cat_ops)

	log.Error("Hopefully Job Data Ops are Set as Operation was completed Successfully")

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