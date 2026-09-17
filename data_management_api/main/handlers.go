package main

import (
	utils "dmapi/csv_utils"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// This function handles file uploads and loading to the API
func (dms *DMService) storeUpload(w http.ResponseWriter, r *http.Request){

	// Handling Uploading File Logic with LIMITS
	r.Body = http.MaxBytesReader(w, r.Body, 10 << 20)
	if err := r.ParseMultipartForm(5 << 20); err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Removing all temp files
	defer r.MultipartForm.RemoveAll()

	// Getting the file data from media key in the POST request
	uf, ufh, err := r.FormFile("media")

	// Checks if file is CSV
	if strings.HasSuffix(ufh.Filename, ".csv") == false{
		http.Error(w, "File format not supported", http.StatusBadRequest)
		return 
	}

	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}

	defer uf.Close()
	
	// Creating path for file storage
	filepathflag := "/run/media/programmerrez/Field Testing/Side-Projects/Arnadillo/data_management_api/.uploads/"
	path := filepath.Join(filepathflag, ufh.Filename)
	
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	f, err := os.Create(path)

	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}

	defer f.Close()

	// Copying File Data Over 
	if _, err := io.Copy(f, uf); err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}
	
	w.Write([]byte("uploaded"))

	// Appending to CSV Storage Object
	csv_file, err := utils.LoadCSV(string(path))
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	
	dms.Data.CSVS = append(dms.Data.CSVS, csv_file)

	return 
}

// This function is a utility for loading and returing the CSV file



