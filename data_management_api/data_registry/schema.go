package dataregistry

import (
	utils "dmapi/csv_utils"
	"sync"
)

// First High Level Object associated with DMS

type DataRegistry struct{
	SessionInfo 	map[string]string 				`json:"session_info"`//Key would include the recongizable names while Value would have session ID
	SessionData		map[string]*DatasetRegistry		`json:"session_data"`
	Mu				sync.RWMutex					`json:"mutex (dont touch)"`
}

// Second Level Object controlling data for each individual session
type DatasetRegistry struct{
	Registry 	[]ManagedDataObjects	`json:"registry"`
}

// Third Level Object communicating with csv_utils 
type ManagedDataObjects struct{
	TargetCol	string		`json:"target_col"`
	Data 		utils.CSV	`json:"data"`
	ID 			int			`json:"id"`
}

// This type is to provide easy API Output for generic slices
type GenericSliceOutput struct{
	Items 	[]string	`json:"items"`
	Count	int			`json:"count"`
}

// This type will provide easy API output for generic maps
type GenericMapOutput struct{
	Map		map[string]int	`json:"map"`
	Count	int				`json:"count"`
}
