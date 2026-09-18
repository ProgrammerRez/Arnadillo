package dataregistry

import (
	utils "dmapi/csv_utils"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// Creates New Session Registry
func CreateNewRegistry(session_name string) *DataRegistry{

	// Creating New Session with random ID (will be stored in session_info map)
	session_info := make(map[string]string)
	session_id := uuid.New().String()

	session_info[session_name] = session_id

	// Creating SessionData Map with key value being session_id

	session_data := make(map[string]*DatasetRegistry)
	session_data[session_id] = &DatasetRegistry{}

	return &DataRegistry{
		SessionInfo: session_info,
		SessionData: session_data,
	}
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
// This function gets all the available csv object names
func (r *DataRegistry) GetAll(name string) (GenericSliceOutput, error){
	r.Mu.Lock()
	defer r.Mu.Unlock()

	csv_list := []string{}

	id, exists := r.SessionInfo[name]
	if !exists{
		return GenericSliceOutput{
			Items: csv_list,
			Count: len(csv_list),
			}, errors.New("Session Info not Found")
		}

		
	registry, exists := r.SessionData[id]

	if !exists{
		return GenericSliceOutput{
			Items: csv_list,
			Count: len(csv_list),
			}, errors.New("Session Data not Found")
		}

	for _, csv := range(registry.Registry){
		csv_list = append(csv_list, csv.Data.FileName)
	}
	return GenericSliceOutput{
			Items: csv_list,
			Count: len(csv_list),
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