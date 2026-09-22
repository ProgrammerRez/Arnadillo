package main

import (
	log "mts/logging"
	"net/http"
)

func main(){

	// Starting the Logging
	log.Init()
	log.SetLevel(log.InfoLevel)
	log.Info("Started Logging in the model training service")

	mux := http.NewServeMux()

	serve := http.Server{
		Addr: ":8001",
		Handler: mux,
	}

	log.Info("Model Training Service Started at Port 8001")
	serve.ListenAndServe()
}
