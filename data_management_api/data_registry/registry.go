package dataregistry

import (
	utils "dmapi/csv_utils"
	log "dmapi/logging"
	"errors"
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"

	"github.com/google/uuid"
)

// Creates New Session Registry
func CreateNewRegistry() *DataRegistry{

	log.SetLevel(log.InfoLevel)
	log.Info("Creating New Registry")
	// Creating empty maps to avoid assignment to nil errors
	session_info := make(map[string]string)
	session_data := make(map[string]*DatasetRegistry)

	return &DataRegistry{	
		SessionInfo: session_info,
		SessionData: session_data,
	}
}

// This function creates new session info 
func (r *DataRegistry) NewSession(session_name string){
	r.Mu.Lock()
	defer r.Mu.Unlock()

	log.Info("Creating New Session: " + session_name)
	// Creating New Session with random ID (will be stored in session_info map)
	session_id := uuid.New().String()
	// Finally setting as session data and info in the final object
	r.SessionInfo[session_name] = session_id
	r.SessionData[session_id] = &DatasetRegistry{}
}

// Deletes Session Data 
func (r *DataRegistry) DeleteSession(session_name string){
	r.Mu.Lock()
	defer r.Mu.Unlock()

	log.Info("Deleting Existing Session: " + session_name)
	// Get the Id first
	id := r.SessionInfo[session_name]
	// Delete the Data
	delete(r.SessionData, id)
	delete(r.SessionInfo, session_name)
}

// Adds different documents to a particular session
func (r *DataRegistry) Add(name string, data []utils.CSV) error{
	r.Mu.Lock()
	defer r.Mu.Unlock()

	log.Info("Adding New Files")

	id, exists := r.SessionInfo[name]
	if !exists{
		return errors.New("Session Info not Found")
	}

	registry, exists := r.SessionData[id]
	if !exists || registry == nil{
		return errors.New("Registry Not Found for session")
	}

	log.Info("Creating New Registry")

	for _, csv := range(data){
		new_managed_csv_object := ManagedDataObjects{
			ID: len(registry.Registry) + 1,
			Data: csv,
			TargetCol: "",
		}
		
		registry.Registry = append(registry.Registry, new_managed_csv_object)
		}	

	log.Info("Added Files to the Registry")

	return nil
}

func (r *DataRegistry) Delete(name string, csv_id int) error{
	r.Mu.Lock()
	defer r.Mu.Unlock()

	log.Info("Deleting Existing Files")

	id, exists := r.SessionInfo[name]
	if exists == false{
		return errors.New("Session Info not Found")
	}

	registry, exists := r.SessionData[id]
	if !exists || registry == nil{
		return errors.New("Registry Not Found for session")
	}

	for i, csv := range(registry.Registry){
		if csv.ID == csv_id{
			err := os.Remove(csv.Data.FilePath)
			if err != nil{
				if os.IsNotExist(err){
					return errors.New("File Doesn't Exist")
				}
				return errors.New("Failed to Delete File: " + err.Error())
			}
			registry.Registry = slices.Delete(registry.Registry,i,i+1)
		}
	}

	log.Info("Deleted Existing Files")

	return nil
}


// This function fills the Null Values in the specified Columns
func (r *DataRegistry) FillNulls(csv_id int, name, column_name, replacement string) error{

	r.Mu.Lock()
	defer r.Mu.Unlock()

	log.Info("Filling in Null Values: Registry")

	id, exists := r.SessionInfo[name]
	if exists == false{
		return errors.New("Session Info not Found")
	}

	registry, exists := r.SessionData[id]
	if !exists || registry == nil{
		return errors.New("Registry Not Found for session")
	}

	for _, csv := range(registry.Registry){
		if csv.ID == csv_id{
			err := csv.Data.FillNA(column_name, replacement)
			if err !=  nil{
				return errors.New(err.Error())
			}
		}
	}

	log.Info("Filled in Null Values")

	return nil


}
// This function gets all the available csv object names
func (r *DataRegistry) GetAll(name string) (GenericMapOutput, error){
	r.Mu.Lock()
	defer r.Mu.Unlock()

	log.Info("Getting Existing Files in a Session")

	csv_map := make(map[string]int)

	id, exists := r.SessionInfo[name]
	if !exists{
		return GenericMapOutput{
			Map: csv_map,
			Count: len(csv_map),
			}, errors.New("Session Info not Found")
		}

		
	registry, exists := r.SessionData[id]

	if !exists{
		return GenericMapOutput{
			Map: csv_map,
			Count: len(csv_map),
			}, errors.New("Session Data not Found")
		}

	for _, csv := range(registry.Registry){
		csv_map[csv.Data.FileName] = int(csv.ID)
	}

	log.Info("Completed Getting all Existing Files in a Session")

	return GenericMapOutput{
			Map: csv_map,
			Count: len(csv_map),
		}, nil
}


// This function will pick one csv and then analyze the stats for it
func (r *DataRegistry) GetStats(name string, csv_name string) (utils.DataFrameStats, error){
	r.Mu.Lock()
	defer r.Mu.Unlock()

	log.Info("Getting Stats for a Specific File in the Session")

	var selected_csv ManagedDataObjects
	
	found := false

	id, exists := r.SessionInfo[name]

	if !exists{
		return selected_csv.Data.GetStats(), errors.New("Session Info not Found")
	}

	registry, exists := r.SessionData[id]

	if !exists{
		return selected_csv.Data.GetStats(), errors.New("Session Data not Found")
		}

	for _, csv := range(registry.Registry){
		if csv.Data.FileName == strings.TrimSpace(strings.ToLower(csv_name)){
			selected_csv = csv
			found =  true
		}
	}

	if !found{
		return selected_csv.Data.GetStats(), errors.New("File Not Found")
	}

	log.Info("Got Specfic Details about the Specific File in the Session")

	return selected_csv.Data.GetStats(), nil

}


// Session Functions
// DEBUG Method: Checking session data storage

func (r *DataRegistry) GetSession(name string) string{
	r.Mu.Lock()
	defer r.Mu.Unlock()

	log.Info("HELPER FUNCTION: Getting Session ID")

	fmt.Println(name)
	id, exists := r.SessionInfo[name]
	if exists == false{
		return ""
	}

	return id
}

// This function will list all the active sessions 
func (r *DataRegistry) ListSession() (GenericSliceOutput, error){

	r.Mu.Lock()
	defer r.Mu.Unlock()

	log.Info("Listing Active Sessions")

	var session_list GenericSliceOutput

	if len(r.SessionInfo) == 0{
		return session_list, errors.New("No Active Sessions Found")
	}

	session_list.Items =slices.Collect(maps.Keys(r.SessionInfo))

	session_list.Count = len(session_list.Items)

	log.Info("Got all the list details of the sessions")

	return session_list, nil
}