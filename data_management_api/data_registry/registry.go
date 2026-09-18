package dataregistry

import (
	utils "dmapi/csv_utils"
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

	id, exists := r.SessionInfo[name]
	if exists == false{
		return errors.New("Session Info not Found")
	}

	registry, exists := r.SessionData[id]

	if !exists || registry == nil{
		return errors.New("Registry Not Found for session")
	}

	for _, csv := range(data){
		new_managed_csv_object := ManagedDataObjects{
		ID: len(registry.Registry) + 1,
		Data: csv,
		TargetCol: "",
	}

	registry.Registry = append(registry.Registry, new_managed_csv_object)
	}	
	
	// for _, csv := range(registry.Registry){
	// 	fmt.Println(csv.ID)
	// 	fmt.Println(csv.Data.FileName)
	// 	fmt.Println(csv.Data.FilePath)
	// }
	return nil
}

func (r *DataRegistry) Delete(name string, csv_id int) error{
	r.Mu.Lock()
	defer r.Mu.Unlock()

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
	return nil
}

// This function gets all the available csv object names
func (r *DataRegistry) GetAll(name string) (GenericMapOutput, error){
	r.Mu.Lock()
	defer r.Mu.Unlock()

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
	return GenericMapOutput{
			Map: csv_map,
			Count: len(csv_map),
		}, nil
}


// This function will pick one csv and then analyze the stats for it
func (r *DataRegistry) GetStats(name string, csv_name string) (utils.DataFrameStats, error){
	r.Mu.Lock()
	defer r.Mu.Unlock()

	var selected_csv ManagedDataObjects
	
	found := false

	id, exists := r.SessionInfo[name]

	if !exists{
		return selected_csv.Data.GetStats(), errors.New("Session Info not Found")
	}

	registry, exists := r.SessionData[id]

	if !exists{
		return utils.DataFrameStats{}, errors.New("Session Data not Found")
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
	
	return selected_csv.Data.GetStats(), nil

}
// Session Functions

// DEBUG Method: Checking session data storage

func (r *DataRegistry) GetSession(name string){
	r.Mu.Lock()
	defer r.Mu.Unlock()

	fmt.Println(name)
	id, exists := r.SessionInfo[name]
	if exists == false{
		return
	}

	fmt.Println(r.SessionData[id])
}

// This function will list all the active sessions 
func (r *DataRegistry) ListSession() (GenericSliceOutput, error){

	r.Mu.Lock()
	defer r.Mu.Unlock()

	var session_list GenericSliceOutput

	if len(r.SessionInfo) == 0{
		return session_list, errors.New("No Active Sessions Found")
	}

	session_list.Items =slices.Collect(maps.Keys(r.SessionInfo))

	session_list.Count = len(session_list.Items)

	return session_list, nil
}

