package main

import (
	log "mts/logging"
	jobs "mts/jobs"
	"net/http"
)

// Main MTS Object
type ModelTrainingService struct{
	JobsQueue	*jobs.JobQeueu
}

func main(){
	// Starting the Logging
	log.Init()
	log.SetLevel(log.InfoLevel)
	log.Info("Started Logging in the model training service")

	// Creating a New Service 
	log.Info("Creating New Model Training Service")
	new_mts := InitMTS()

	mux := http.NewServeMux()

	serve := http.Server{
		Addr: ":8001",
		Handler: mux,
	}

	// Adding all the endpoints
	mux.HandleFunc("POST /data_params", new_mts.addJob)
	mux.HandleFunc("POST /get_jobs", new_mts.getJobs)
	mux.HandleFunc("POST /set_ops", new_mts.setDataOps)

	log.Info("Model Training Service Started at Port 8001")
	serve.ListenAndServe()
}
