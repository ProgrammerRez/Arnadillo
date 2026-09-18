package dataregistry

import (
	utils "dmapi/csv_utils"
	"sync"
)

// First High Level Object associated with DMS

type DataRegistry struct{
	SessionInfo 	map[string]string //Key would include the recongizable names while Value would have session ID
	SessionData		map[string]*DatasetRegistry
	Mu				sync.RWMutex
}

// Second Level Object controlling data for each individual session
type DatasetRegistry struct{
	Registry 	[]ManagedDataObjects
}

// Third Level Object communicating with csv_utils 
type ManagedDataObjects struct{
	TargetCol	string		`json:"target_col"`
	Data 		utils.CSV	`json:"-"`
	ID 			int			`json:"id"`
}

// This type is to provide easy API Output for generic slices
type GenericSliceOutput struct{
	Items 	[]string	`json:"items"`
	Count	int			`json:"count"`
}