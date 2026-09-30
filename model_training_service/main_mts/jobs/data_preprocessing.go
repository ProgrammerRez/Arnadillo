package jobs
// package main

import (
	"errors"
	"fmt"
	log "mts/logging"
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

		log.Info(fmt.Sprintf("Processed %d object(s)", i))

		// Finally appending the data
		collected_data = append(collected_data, DataParameters{
			Data: row_data,
			ColumnList: column_list,
			TargetLabel: target_col,
		})
	}

	return collected_data, nil
}

// Helper Function for checking value existence quickly
func getValue(object map[string]any, check_value string) (any, error){

	// Checking if value exists
	value, exists := object[check_value]
	if !exists{
		log.Error(fmt.Sprintf("Value: %s doesn't exist in object %s", check_value, object))
		return nil, fmt.Errorf("Value: %s doesn't exist in object %s", check_value, object)
	}
	log.Info(fmt.Sprintf("Value %s exists", check_value))
	return value, nil
}

// This function handles value existence and type checkking for String Arrays
func getStringArrays(value any) ([]string, error){

	// Making sure it's a slice
	raw_data, ok := value.([]any)
	if !ok {
		return nil, errors.New("Data should be an array")
	}

	// Creating Transfer Object
	result := make([]string, len(raw_data))

	// Using a for loop to deal with it
	for i, raw_value := range(raw_data){
		value, ok := raw_value.(string)
		if !ok{
			return nil, errors.New("Values of array are not string")
		}
		result[i] = value
	}

	return result, nil
}

// This function handles value existence and type checking for String Row Data (looks like this `[][]string`)
func getStringRows(value any) ([][]string, error){

	// First Verifying the Data as a slice of `any` data type
	raw_data, ok := value.([]any)
	if !ok{
		return nil, errors.New("Data needs to be Row String Array")
	}	

	// Creating Transfer Object
	results := make([][]string, len(raw_data))

	// Using a for loop to transfer the data
	for i, row_data := range(raw_data){
		row, err := getStringArrays(row_data)
		if err != nil{
			return nil , errors.New(err.Error())
		}
		results[i] = row
	}

	return results, nil
}