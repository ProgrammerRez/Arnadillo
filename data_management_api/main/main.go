package main

import (
	"fmt"
	// "log"
	"net/http"
	// utils "dmapi/server_utils"
	registry "dmapi/data_registry"
)

// This one is the main Object for handling all operations
type DMService struct{
	Data 	registry.DataRegistry
}

// Creates new DM Service
func newDMS(name string) DMService{
	new_registry := registry.CreateNewRegistry(name)
	return DMService{
		Data: *new_registry,
	}
}

func main(){
	dms := newDMS("session1")
	mux := http.NewServeMux()
	mux.HandleFunc("POST /upload", dms.storeUpload)
	// mux.HandleFunc("GET /stats", dms.outputStats)
	mux.HandleFunc("GET /files", dms.outputFilesDetails)
	mux.HandleFunc("GET /session", dms.getSession)
	mux.HandleFunc("POST /file", dms.specificFileDetails)
	// mux.HandleFunc("GET /out", dms.showEM)

	serv := http.Server{
		Addr: ":8000",
		Handler: mux,
	}

	fmt.Println("Server Started at Port 8000")
	serv.ListenAndServe()
}
