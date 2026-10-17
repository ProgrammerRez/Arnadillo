package jobs

// package main

import (
	"errors"
	"fmt"
	log "mts/logging"
	"slices"
)

// This file will handle the preprocessing of Data provided

// This function will parse the retrieved data to transform into a DataParamsObject
func CreateDataParamObject(data_object map[string]any) ([]DataParameters, error){

	log.Info("Converting into Data Params Object")

	// Checking for Regitstry Key
	val, err := getValue(data_object, "registry")

	if err != nil{
		return nil, errors.New(err.Error())
	}

	// Checking if registry is an slice / array 
	registry, ok := val.([]any)
	if !ok || len(registry) == 0 {
		log.Error("Regitry is not an array / slice / is empty")
		return nil, errors.New("Regitry is not an array / slice / is empty")
	}

	// Now pulling a for loop over the registry

	// Creating a variable to store collected data
	collected_data := make([]DataParameters, 0, len(registry))

	for i, reg_entry := range(registry){

		// Checking if `reg_entry` is a map
		reg_entry, ok := reg_entry.(map[string]any)
		if !ok{
			log.Error("Regisrty Data isn't a map")
			return nil, errors.New("Regisrty Data isn't a map")
		}

		log.Info("Processed Registry Data")

		// Checking `data` key exists and is a map
		raw, err := getValue(reg_entry, "data")
		if err != nil{
			return nil, errors.New(err.Error())
		}

		// Checking if `data` key is a map value
		data, ok := raw.(map[string]any)
		if !ok{
			log.Error("Data isn't a map")
			return nil, errors.New("Data isn't a map")
		}

		log.Info("Processed Data Map")

		// Checking if `data` map has a `target_col` key
		target, err := getValue(reg_entry, "target_col")
		if err != nil{
			return nil, errors.New(err.Error())
		}

		// Checking if `target_col` is a string
		target_col, ok := target.(string)
		if !ok{
			log.Error("Target Col is not a string")
			return nil, errors.New("Target Col is not a string")
		}

		log.Info("Traget Processing Complete")

		// Checking if `data` key exists in `data` map
		raw_data, err := getValue(data, "data")
		if err != nil{
			return nil, errors.New(err.Error())
		}

		// Checking if `raw_data` is a slice of `string`(s)
		row_data, err := getStringRows(raw_data)
		if err != nil{
			log.Error(err.Error())
			return nil, errors.New(err.Error())
		}

		log.Info("Proccessed Row Data")

		// Checking if `column_list` key exists in `data` map
		raw_list, err := getValue(data, "column_list")
		if err != nil{
			return nil, errors.New(err.Error())
		}

		// Checking if `column_list` is a slice of `string`(s)
		column_list, err := getStringArrays(raw_list)
		if err != nil{
			log.Error(err.Error())
			return nil, errors.New(err.Error())
		}

		log.Info("Proccessed Column Data")

		log.Info(fmt.Sprintf("Processed %d object(s)", i+1))

		// Finally appending the data
		collected_data = append(collected_data, DataParameters{
			Data: row_data,
			ColumnList: column_list,
			TargetLabel: target_col,
		})
	}

	return collected_data, nil
}

// This function will set the operations for numerical columns and categorical columns
func (dp *DataParameters) SetOps(num_ops []string, cat_ops map[string][]string) error{
	log.Info("Setting the Num Ops and Categorical Ops")

	if len(num_ops) == 0 || len(cat_ops) == 0{
		log.Error("Proper Populated Operations Not Provided")
		return errors.New("Proper Populated Operations Not Provided")
	}

	// Simply appending the numerical operations
	dp.NumericalOperations = append(dp.NumericalOperations, num_ops...)
	
	// Appending to existing Operations
	for col, ops := range cat_ops{
		if !slices.Contains(dp.ColumnList, col){
			log.Error(fmt.Sprintf("Column %s doesn't exist in the column list", col))
			return fmt.Errorf("Column %s doesn't exist in the column list", col)
			
		}
		log.Info(fmt.Sprintf("Added Object Operations for columns: %s", col))
		dp.ObjectOperations[col] = append(dp.ObjectOperations[col], ops...)
	}
	
	log.Info("Successfully Set the Operations")
	return nil
}