package main

import (
	"fmt"
	// "log"
	"net/http"
	// utils "dmapi/server_utils"
	registry "dmapi/data_registry"
	log "dmapi/logging"
)

// This one is the main Object for handling all operations
type DMService struct{
	Data 	registry.DataRegistry
}



// Creates new DM Service
func newDMS(name string) DMService{
	new_registry := registry.CreateNewRegistry()
	new_registry.NewSession(name)
	log.Info("New Data Management Service Created")
	return DMService{
		Data: *new_registry,
	}
}

func main(){
	// Initializing Logger
	log.Init()
	log.SetLevel(log.InfoLevel)
	log.Info("Starting the API")


	dms := newDMS("session1")
	
	mux := http.NewServeMux()
	mux.HandleFunc("POST /upload", dms.storeUpload)
	mux.HandleFunc("GET /files", dms.outputFilesDetails)
	mux.HandleFunc("GET /session", dms.getSession)
	mux.HandleFunc("POST /file", dms.specificFileDetails)
	mux.HandleFunc("GET /sessions", dms.listSessions)
	mux.HandleFunc("POST /new_session", dms.createNewSession)
	mux.HandleFunc("POST /delete_session", dms.deleteSpecificSession)
	mux.HandleFunc("POST /delete_file", dms.deleteFile)
	mux.HandleFunc("POST /fill_na", dms.fillNulls)
	mux.HandleFunc("POST /delete_column", dms.deleteColumn)
	mux.HandleFunc("POST /set_target", dms.setTarget)
	mux.HandleFunc("GET /the_good_stuff", dms.getDataObjects)
	mux.HandleFunc("GET /AAJA", dms.exportToFile)

	serv := http.Server{
		Addr: ":8000",
		Handler: mux,
	}

	fmt.Println("Server Started at Port 8000")
	serv.ListenAndServe()
}
