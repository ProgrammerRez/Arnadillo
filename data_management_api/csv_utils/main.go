package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// Here is the new strcuture:
// Create an object for CSV files with all the necessary stuff like columns, shape, null and dupes, etc.
// Each object will have the following schema: 1. file path, data, file size, file shape, column list.
// The rest of the stats will have class methods.
// Goal is to make it as portable as possible

// Main CSV Object
type CSV struct{
	Data 		[][]string
	ColumnList []string
	FilePath  	string
	FileShape	[2]int64
	FileSize	int32
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
			FilePath: file_path,
			Data: data,
			ColumnList: column_list,
			FileShape: [2]int64{int64(len(data)), int64(len(column_list))}, 
		}

	}

	return new_csv_object, nil
}


// Now Object Methods

// This function will provide with stats such as data type, unique values and null and dupe values for pandas
func(csv CSV) GetStats() DataFrameStats{
	
	col_data := make(map[string]ColumnStats)

	dtypesMatrix := csv.getDTypes()
	uniqueValueMatrix, UniqueValuesCount := csv.getUniqueValues()
	nulls, NPC := csv.getNulls()
	dupeCount := csv.getDupes()
	

	for _, col := range(csv.ColumnList){
		col_data[col] = ColumnStats{
			UniqueValues: uniqueValueMatrix[col],
			NullsInCol: NPC[col],
			DataType: dtypesMatrix[col],
		}
	}
	

	
	return DataFrameStats{
		DtypesMatrix: dtypesMatrix,
		Dupes: int64(dupeCount),
		Nulls: int64(nulls),
		NullsByCol: NPC,
		UniqueValueMatrix: UniqueValuesCount,
		ColStats: col_data,
	}
}

// This function allows to dynamically adjust data types when loading the CSV file
func (csv CSV) getDTypes() (dTypeMatrix map[string]string){

	dTypeMatrix = make(map[string]string)

	for _, col := range(csv.ColumnList){
		dTypeMatrix[col] = ""
	}

	for i, column := range(csv.Data[0]){
		
		col := csv.ColumnList[i]

		trimmed := strings.TrimSpace(column)
		
		if _, err := strconv.ParseBool(trimmed); err == nil{
			dTypeMatrix[col] = "bool"
		}else if _, err := strconv.ParseInt(trimmed, 10, 64); err == nil{
			dTypeMatrix[col] = "int64"
		}else if _ , err := strconv.ParseFloat(trimmed, 64); err == nil{
			dTypeMatrix[col] = "float64"
		}else{
			dTypeMatrix[col] = "string"
		}
	}
	return dTypeMatrix
}


// This function will output the unqiue values for each column
func (csv CSV) getUniqueValues() (uniqueValues map[string][]string, uniqueValuesCount map[string]int64){
	
	uniqueValues = make(map[string][]string)
	uniqueValuesCount = make(map[string]int64)
	seenValues := make(map[string]map[string]bool)

	for _, col := range(csv.ColumnList){
		uniqueValuesCount[col] = 0
		uniqueValues[col] = []string{}
		seenValues[col] = make(map[string]bool)
	}

	for _, row := range(csv.Data){
		for i, cell := range(row){
			col_name := csv.ColumnList[i]

			if i > len(csv.ColumnList){
				break
			}
			
			if seenValues[col_name][cell]{
				continue
			}
			
			uniqueValues[col_name] = append(uniqueValues[col_name], cell)
			seenValues[col_name][cell] = true
			uniqueValuesCount[col_name]++
		
		}
	}

	return uniqueValues, uniqueValuesCount
}


// This function will provide null values with respect to the columns
func (csv CSV) getNulls() (nullCount int32, nullsPerCol map[string]int64) {

	nullsPerCol = make(map[string]int64)

	for _, col := range(csv.ColumnList){
		nullsPerCol[col] = 0
	}

	for _, row := range(csv.Data){
		for i, cell := range(row){
			if strings.TrimSpace(cell) == ""{
				nullCount++
				if i < len(csv.ColumnList){
					nullsPerCol[csv.ColumnList[i]]++
				}
			}

		}
	}
	return nullCount, nullsPerCol
}

// This function will provide dupe values and their population in the dataset
func (csv CSV) getDupes() (dupeCount int64) {
	seenRows := make(map[string]bool)

	for _, row := range(csv.Data){
		rowKey := strings.Join(row, "|")

		if seenRows[rowKey]{
			dupeCount++
		}else{
			seenRows[rowKey] = true
		}
	}
	return dupeCount
}






func main(){

	csv, err := loadCSV("/run/media/programmerrez/Field Testing/Side-Projects/Arnadillo/Data/customer_master.csv")

	if err != nil{
		fmt.Println(err)
	}

	// fmt.Println(csv)

	stats := csv.GetStats()

	fmt.Println(stats.DtypesMatrix)
	fmt.Println(stats.Nulls)
	fmt.Println(stats.NullsByCol)
	fmt.Println(stats.Dupes)
	fmt.Println(stats.UniqueValueMatrix)

	for _, stat := range(stats.ColStats){	
		fmt.Println(stat)
	}
}	