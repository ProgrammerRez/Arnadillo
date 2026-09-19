package main

import (
	utils "dmapi/csv_utils"
	log "dmapi/logging"
	"encoding/json"
	// "fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// File Handling Section

// This function handles file uploads and loading to the API
func (dms *DMService) storeUpload(w http.ResponseWriter, r *http.Request){

	// Setting logger level to 1
	log.SetLevel(log.InfoLevel)

	// Handling Uploading File Logic with LIMITS
	r.Body = http.MaxBytesReader(w, r.Body, 10 << 20)
	if err := r.ParseMultipartForm(5 << 20); err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error(err.Error())
		return
	}
	// Removing all temp files
	defer r.MultipartForm.RemoveAll()

	// Getting the file data from media key in the POST request
	files := r.MultipartForm.File["media"]

	if len(files) == 0{
		http.Error(w, "No files provided", http.StatusBadRequest)
		log.Error("No files provided")
		return
	}

	// Creating path for file storage
	file_path := "/run/media/programmerrez/Field Testing/Side-Projects/Arnadillo/data_management_api/.uploads/"	
	if err := os.MkdirAll(filepath.Dir(file_path), os.ModePerm); err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err.Error())
		return
	}

	var loadedCSVs []utils.CSV

	for _, file_header := range(files){

		// Checks if file is CSV
		if strings.HasSuffix(file_header.Filename, ".csv") == false{
			http.Error(w, "File format not supported", http.StatusBadRequest)
			log.Error("File format not supported")
			return 
		}

		file, err := file_header.Open()

		if err != nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error(err.Error())
			return
		}
		
		path := filepath.Join(file_path, file_header.Filename)
		f, err := os.Create(path)
		
		if err != nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error(err.Error())
			file.Close()
			return 
		}

		// Copying File Data Over 
		if _, err := io.Copy(f, file); err != nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error(err.Error())
			return 
		}
		file.Close()
		f.Close()
		
		
		// Appending to CSV Storage Object
		csv_file, err := utils.LoadCSV(string(path))
		
		if err != nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error(err.Error())
			return
		}

		loadedCSVs = append(loadedCSVs, csv_file)
	}

	// Get the name value from form body
	name := r.FormValue("name")

	if name == ""{
		http.Error(w, "name for the session was not provided", http.StatusBadRequest)
		log.Error("name for the session was not provided")
		return
	}
	err := dms.Data.Add(name, loadedCSVs)

	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error(err.Error())
		return
	}
	w.Write([]byte("uploaded"))
}

// Function to delete CSV file from specified session
func (dms *DMService) deleteFile(w http.ResponseWriter, r *http.Request){

	// Get the name and csv_name value from form body
	name := r.FormValue("name")
	csv_name := r.FormValue("csv_name")
	if name == "" || csv_name == ""{
		http.Error(w, "Session or Filer name not Provided", http.StatusBadRequest)
		log.Error("Session or Filer name not Provided")
		return
	}

	csv_id, err := strconv.ParseInt(csv_name, 10, 64)

	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error(err.Error())
		return
	}

	dms.Data.Delete(name, int(csv_id))
}

// This function fills nulls in the stored data; Persists only with export option
func (dms *DMService) fillNulls(w http.ResponseWriter, r *http.Request){
	// Get the name and csv_name value from form body
	log.Info("Filling NA Values: Handler")
	name := r.FormValue("name")
	csv_name := r.FormValue("csv_name")
	column_name := r.FormValue("column_name")
	replacement := r.FormValue("replacement")
	if name == "" || csv_name == "" || column_name == "" || replacement == ""{
		http.Error(w, "Session or File name not Provided", http.StatusBadRequest)
		log.Error("Session or File name not Provided")
		return
	}

	csv_id, err := strconv.ParseInt(csv_name, 10, 64)

	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error(err.Error())
		return
	}

	if err := dms.Data.FillNulls(int(csv_id), name, column_name, replacement); err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err.Error())
		return 
	}
}

// Outputs all the statistics provided by the getStats class method
// func (dms *DMService) outputStats(w http.ResponseWriter, r *http.Request){

// 	// Getting the name for the session
// 	name := r.FormValue("name")

// 	if name == ""{
// 		http.Error(w, "name for the session was not provided", http.StatusBadRequest)
// 	}

// 	bytes, err := json.Marshal(dms.Data.CSVS[0].GetStats())
// 	if err != nil{
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// 	w.Write(bytes)
// } 


// // This function is a utility for getting a list of available CSV files 
func (dms *DMService) outputFilesDetails(w http.ResponseWriter, r *http.Request){

	// Setting the log level to 1
	log.SetLevel(log.InfoLevel)

	name := r.FormValue("name")

	if name == ""{
		http.Error(w, "no session name was provided", http.StatusBadRequest)
		log.Error("no session name was provided")
		return
	}
	
	csv_list, err := dms.Data.GetAll(name)

	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err.Error())
		return
	}

	bytes, err := json.Marshal(csv_list)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err.Error())
	}
	w.Write(bytes)
}





// // This function outputs the specific file you're targetting
func (dms *DMService) specificFileDetails(w http.ResponseWriter, r *http.Request){

	name := r.FormValue("name")
	csv_name := r.FormValue("csv_name")

	if strings.TrimSpace(strings.ToLower(name)) == "" || strings.TrimSpace(strings.ToLower(csv_name)) == ""{
		http.Error(w, "Session or File does not exist", http.StatusBadRequest)
		log.Error("Session or File does not exist")
		return
	}

	stats , err:= dms.Data.GetStats(name, csv_name)

	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err.Error())
		return
	}

	bytes, err := json.Marshal(stats)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err.Error())
		return
	}
	w.Write(bytes)

}

// ----

// Session Handlers

// DEBUG Handler: Checks on session Data structure
func (dms *DMService) getSession(w http.ResponseWriter, r *http.Request){

	name := r.FormValue("name")

	if name == ""{
		http.Error(w, "no session name was provided", http.StatusBadRequest)
		log.Error("no session name was provided")
		return
	}
	
	dms.Data.GetSession(name)
}

// This function will provide a list of all created sessions in the API
func (dms *DMService) listSessions(w http.ResponseWriter, r *http.Request){
	session_list, err := dms.Data.ListSession()

	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err.Error())
		return
	}

	bytes, err := json.Marshal(session_list)

	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err.Error())
		return
	}

	w.Write(bytes)
}

// This function creates a new session
func (dms *DMService) createNewSession(w http.ResponseWriter, r *http.Request){

	name := r.FormValue("name")

	if name == ""{
		http.Error(w, "no session name was provided", http.StatusBadRequest)
		log.Error("no session name was provided")
		return
	}

	dms.Data.NewSession(name)
}

// This function deletes the specified session
func (dms *DMService) deleteSpecificSession(w http.ResponseWriter, r *http.Request){

	name := r.FormValue("name")

	if name == ""{
		http.Error(w, "no session name was provided", http.StatusBadRequest)
		log.Error("no session name was provided")
		return
	}

	dms.Data.DeleteSession(name)
}