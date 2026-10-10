package jobs

import(
	log "mts/logging"
	"fmt"
	"errors"
)


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