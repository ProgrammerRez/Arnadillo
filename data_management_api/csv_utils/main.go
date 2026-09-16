package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
)

// Here is the new strcuture:
// Create an object for CSV files with all the necessary stuff like columns, shape, null and dupes, etc.
// Each object will have the following schema: 1. file path, data, file size, file shape, column list.
// The rest of the stats will have class methods.
// Goal is to make it as portable as possible

// Main CSV Object
type CSV struct{
	data 		[][]string
	column_list []string
	file_path  	string
	file_shape	[2]int64
	file_size	int32
}


// Now the loading mechanism

func loadCSV(file_path string) (CSV, error){
	
	new_csv_object := CSV{}

	// First open the file and check for any errors
	file, err := os.Open(file_path)

	if err != nil || file == nil{
		log.Fatalf("Error Occurred during file opening: ", err)
		return new_csv_object, errors.New(err.Error())
	}

	defer file.Close()

	// Now reading the data
	csv_reader := csv.NewReader(file)

	records, err := csv_reader.ReadAll()

	if err != nil{
		log.Fatalf("Error occurred while reading the data: ", err)
		return new_csv_object, errors.New(err.Error())
	}

	// Creating the New Object
	if records != nil{	

		data := records[1:]
		column_list := records[0]

		new_csv_object = CSV{
			file_path: file_path,
			data: data,
			column_list: column_list,
			file_shape: [2]int64{int64(len(data)), int64(len(column_list))}, 
		}

	}

	return new_csv_object, nil
}


// Now Object Methods

// This function will provide with stats such as .describe() and .info() for pandas
func(csv CSV) getStats(){
	return 
}

// This function will provide null values with respect to the columns
func (csv CSV) getNulls(){
	return
}

// This function will provide dupe values and their population in the dataset
func (csv CSV) getDupes(){
	return
}






func main(){

	csv, err := loadCSV("Data/customer_master.csv")

	if err != nil{
		fmt.Println(err)
	}

	fmt.Println(csv.file_shape)
}	