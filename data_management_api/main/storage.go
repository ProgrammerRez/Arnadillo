package main

import (
	utils "dmapi/csv_utils"
	"strings"
	"sync"
)

// This struct stores a slice of CSV Objects
type CSVStorage struct{
	CSVS  	[]utils.CSV
	mu		sync.RWMutex
}

// Adds CSV objects to the storage
func (s *CSVStorage) Add(csv utils.CSV){
	s.mu.Lock()
	defer s.mu.Unlock()
	s.CSVS = append(s.CSVS, csv)
}


// This function lists all the available CSV files
func (s *CSVStorage) GetAll() GenericSliceOutput {
	s.mu.Lock()
	defer s.mu.Unlock()
	available_csvs := []string{}
	
	for _, csv := range(s.CSVS){
		available_csvs = append(available_csvs, csv.FileName)
	}
	return GenericSliceOutput{
		Items: available_csvs,
		Count: len(available_csvs),
	}
}

// This function gets the specific csv file user is aiming for
func (s *CSVStorage) GetCSV(name string) utils.CSV{
	s.mu.Lock()
	defer s.mu.Unlock()

	final_file := utils.CSV{}
	for _, csv := range(s.CSVS){
		if strings.TrimSpace(strings.ToLower(csv.FileName)) == strings.TrimSpace(strings.ToLower(name)){
			final_file = csv
		}
	}

	return final_file
}
