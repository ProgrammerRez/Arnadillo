package main

// This file would go like this
// 1. This would be an API where all the endpoints are available
// 2. The routers are coming from other dir structures but everything meets here
// 3. Make it scalable with data

import (
	"fmt"
	// "log"
	"net/http"
	// utils "dmapi/server_utils"
)

// This one is the main Object for handling all operations
type DMService struct{
	Data 	CSVStorage
}

// Creates new DM Service
func newDMS() DMService{
	return DMService{}
}

func main(){
	dms := newDMS()
	mux := http.NewServeMux()
	mux.HandleFunc("POST /upload", dms.storeUpload)
	mux.HandleFunc("GET /stats", dms.outputStats)
	mux.HandleFunc("GET /files", dms.outputFilesDetails)
	mux.HandleFunc("POST /file", dms.specificFileDetails)
	// mux.HandleFunc("GET /out", dms.showEM)

	serv := http.Server{
		Addr: ":8000",
		Handler: mux,
	}

	fmt.Println("Server Started at Port 8000")
	serv.ListenAndServe()
}
