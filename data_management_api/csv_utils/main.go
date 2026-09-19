package utils
// package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"sort"

	// "fmt"
	log "dmapi/logging"
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
	ColumnList 	[]string
	FileName 	string
	FilePath  	string
	FileShape	[2]int64
	FileSize	int64
}


// Now the loading mechanism

func LoadCSV(file_path string) (CSV, error){

	// log.SetLevel(log.InfoLevel)

	new_csv_object := CSV{}
	// First open the file and check for any errors
	file, err := os.Open(file_path)
	if err != nil || file == nil{
		log.Error("Error Occurred during file opening: " + err.Error())
		return new_csv_object, errors.New("Error Occurred during file opening: " + err.Error())
	}

	log.Info("Opened File at path: " + file_path)
	defer file.Close()

	// Then getting the file size
	file_info, err := os.Stat(file_path)
	if err != nil{
		log.Error("Error Occurred during file info extraction: " + err.Error())
		return new_csv_object, errors.New("Error Occurred during file info extraction: " + err.Error())
	}

	log.Info("Extracted File Info")

	// Now reading the data
	csv_reader := csv.NewReader(file)
	
	
	// Read the header row first
	column_list, err := csv_reader.Read()

	// Handling Read Errors
	if err == io.EOF{
		log.Error("File is Empty" + err.Error())
		return new_csv_object, errors.New("File is Empty" + err.Error())
	}
	if err != nil{
		log.Error(fmt.Sprintf("unable to read header row: %w", err))
		return new_csv_object, fmt.Errorf("unable to read header row: %w", err)
	}
	
	log.Info("Extracted Column List")
	
	var records [][]string
	
	for{
		record, err := csv_reader.Read()
		
		if err == io.EOF{
			log.Error("File is Empty: EOF error")
			break
		}

		if err != nil{
			log.Error(fmt.Sprintf("error reading row %d: %w", len(record)+2, err))
			return new_csv_object, fmt.Errorf("error reading row %d: %w", len(record)+2, err)
		}

		records = append(records, record)
	}

	log.Info("Extracted all Records")
	// Creating the New Object
	if records != nil{	

		data := records

		new_csv_object = CSV{
			FileName: file_info.Name(),
			FilePath: file_path,
			Data: data,
			ColumnList: column_list,
			FileShape: [2]int64{int64(len(data)), int64(len(column_list))}, 
			FileSize: file_info.Size(),
		}

	}
	log.Info("Loaded CSV File: " + new_csv_object.FileName)
	return new_csv_object, nil
}


// Now Object Methods

// This function will provide with stats such as data type, unique values and null and dupe values for pandas
func(csv CSV) GetStats() DataFrameStats{

	// log.SetLevel(log.InfoLevel)
	
	col_data := make(map[string]ColumnStats)
	
	dtypesMatrix, err := csv.getDTypes()
	if err != nil{
		log.Error("Error Extracting Data Types: " + err.Error())
	}

	uniqueValueMatrix, UniqueValuesCount := csv.getUniqueValues()
	nulls, NPC := csv.getNulls()
	dupeCount := csv.getDupes()
	numStats := csv.getNumStats(dtypesMatrix)


	for _, col := range(csv.ColumnList){
		col_data[col] = ColumnStats{
			UniqueValues: uniqueValueMatrix[col],
			NullsInCol: NPC[col],
			DataType: dtypesMatrix[col],
			NumStats: *numStats[col],
		}
	}

	log.Info("Calculations and Assignment Working")
	log.Info("OKAY NP :)")
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
func (csv CSV) getDTypes() (dTypeMatrix map[string]string, err error){

	log.Info("Extracting Data Type Info")
	dTypeMatrix = make(map[string]string)
	err = nil
	
	if len(csv.ColumnList) == 0{
		log.Error("Column List is Empty")
		return dTypeMatrix, errors.New("Column List is Empty")
	}
	for _, col := range(csv.ColumnList){
		dTypeMatrix[col] = ""
	}

	for i, column := range(csv.Data[0]){
		
		col := csv.ColumnList[i]

		trimmed := strings.TrimSpace(strings.ToLower(column))
		
		if _, err := strconv.ParseBool(trimmed); err == nil{
			dTypeMatrix[col] = "bool"
		}else if _, err := strconv.ParseInt(trimmed, 10, 64); err == nil{
			dTypeMatrix[col] = "int64"
		}else if _ , err := strconv.ParseFloat(trimmed, 64); err == nil{
			dTypeMatrix[col] = "float64"
		}else{
			dTypeMatrix[col] = "string"
		}

		if dTypeMatrix[col] == ""{
			log.Error("No Data Type Assigned")
			return dTypeMatrix, errors.New("No Data Type Assigned")
		}
	}

	log.Info("Completing Data Type Extraction")
	return dTypeMatrix, err
}


// This function will output the unqiue values for each column
func (csv CSV) getUniqueValues() (uniqueValues map[string][]string, uniqueValuesCount map[string]int64){
	
	log.Info("Extracting Unique Values")
	
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
	log.Info("Completed Unique Value Extraction")
	return uniqueValues, uniqueValuesCount
}


// This function will provide null values with respect to the columns
func (csv CSV) getNulls() (nullCount int32, nullsPerCol map[string]int64) {

	log.Info("Extracting Null Value Info")
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
	log.Info("Completed Extracting Null Value Info")
	return nullCount, nullsPerCol
}

// This function will provide dupe values and their population in the dataset
func (csv CSV) getDupes() (dupeCount int64) {

	log.Info("Started Dupe Info Extraction")
	seenRows := make(map[string]bool)

	for _, row := range(csv.Data){
		rowKey := strings.Join(row, "|")

		if seenRows[rowKey]{
			dupeCount++
		}else{
			seenRows[rowKey] = true
		}
	}

	log.Info("Completed Dupe Info Extraction")
	return dupeCount
}

// This function fills_na of respective columns
func (csv CSV) FillNA(column_name string, replacement string) error{

	log.Info("Reached csv_utils")
	// First Let's get the stats
	stats := csv.GetStats()

	if stats.Nulls == 0{
		log.Error("No Null Values found in the dataframe")
		return errors.New("No Null Values found in the dataframe")
	}

	if stats.ColStats[column_name].NullsInCol == 0 && stats.NullsByCol[column_name] == 0{
		log.Error("No Null Values found in the column. Please Double Check Column Name" + string(stats.NullsByCol[column_name]))
		return errors.New("No Null Values found in the column. Please Double Check Column Name")
	}

	fmt.Println(stats.ColStats[column_name].NullsInCol)
	fmt.Println(stats.NullsByCol[column_name])

	// Get the replacement Value
	dtype := stats.DtypesMatrix[column_name]
	trimmed := strings.TrimSpace(replacement)

	// Numerical Substitution
	if dtype == "int64" || dtype == "float64"{
		switch trimmed {
		case "mode":
			csv.fillNa(column_name, stats.ColStats[column_name].NumStats.Mode)
		case "median":
			csv.fillNa(column_name, stats.ColStats[column_name].NumStats.Median)
		case "mean":
			csv.fillNa(column_name, stats.ColStats[column_name].NumStats.Mean)
		default:
			log.Error("No Valid Option Provided for Numerical Substitution")
			return errors.New("No Valid Option Provided for Numerical Substitution")
		}
	}
	
	// Object Substitution
	if dtype == "string" || dtype == "bool"{
		if trimmed == "mode"{
			csv.fillNa(column_name, stats.ColStats[column_name].NumStats.Mode)
		}else{
			csv.fillNa(column_name, trimmed)
		}
	}

	return nil
}


// This function fills the value in the specified column
func (csv *CSV) fillNa(column_name string, value any) error{

	log.Info("Filling NA: Utils")
	fmt.Println(value)
	id := -1

	trimmed := strings.TrimSpace(column_name)
	
	for i, col := range(csv.ColumnList){
		if col == trimmed{
			id = i
		}
	}

	if id == -1{
		log.Error("Column Does Not Exist: " + trimmed)
		return errors.New("Column Does Not Exist: " + trimmed)
	}

	var str_value string

	switch v := value.(type){
	case float32, float64:
		str_value = fmt.Sprintf("%.2f", v)
	default:
		str_value = fmt.Sprintf("%v", v)
	}

	log.Info("Successfully Parsed Replacement value")

	for row := range(csv.Data){
		if id < len(csv.Data[row]){
			if strings.TrimSpace(csv.Data[row][id]) == ""{
				csv.Data[row][id] = str_value
			}
		}
	}

	log.Info("Filled NA value")

	return nil
}


// This function gets the numerical stats for each column
func (csv CSV) getNumStats(dtype_matrix map[string]string) (num_stats map[string]*NumericalStats){

	// Starting Logging
	log.Info("Starting Extraction for Numerical Stats")

	var err error

	// Creating num_stats to avoid nil map assignment error
	num_stats = map[string]*NumericalStats{}

	// Creating empty map values
	for _, col := range(csv.ColumnList){
		num_stats[col] = &NumericalStats{}
	}
	for i, col:= range(csv.ColumnList){
		var data []string
		dtype := dtype_matrix[col]

		log.Info(fmt.Sprintf("Getting data for column: %s", col))
		// Getting Each Row's Data
		for _, val := range(csv.Data){
			data = append(data, val[i])
		}

		log.Info("Column Name: " + col + " is passing through the float loop")
		
		if dtype == "int64" || dtype == "float64"{
		
			// Converting to Float Data
			float_data := ParsingFloats(data)

			// Calculating Mean
			num_stats[col].Mean, err = CalculateMean(float_data)	
			if err != nil{
				log.Error("Mean Error" + err.Error()) 
			}

			// Calculating Median
			num_stats[col].Median, err = CalculateMedian(float_data)
			if err != nil{
				log.Error("Median Error" + err.Error()) 
			}

			// Calculating Mode
			num_stats[col].Mode, err = CalculateModeNum(float_data)
			if err != nil{
				log.Error("Numerical Mode Error" + err.Error()) 
			}
		}

		// String Pathway
		if dtype == "string" || dtype == "bool"{
			log.Info("Column Name: " + col + " is passing through the string loop")
			num_stats[col].Mode, _ = CalculateModeString(data)
			if err != nil{
				log.Error("String Mode Error" + err.Error()) 
			}
		}
		log.Info("Completed Numerical Stat Extraction for column: " + col)
	}

	log.Info("Completed Numerical Stat Extraction for all columns")

	return num_stats
}

// This function converts `[]string` into `[]float64`
func ParsingFloats(data []string) []float64{
	log.Info("Parsing Floats")
	var nums []float64
	for i, val := range(data){
		num, err := strconv.ParseFloat(val, 64) 
		if err != nil{
			
			log.Warn(fmt.Sprintf("Parsing Error Caused at line: %d. Changing to Int Conversion",i))
			fmt.Println("Parsing Error Caused at line: ",i)
			
			num_int, err := strconv.ParseInt(val, 10, 64)
			if err != nil{
				log.Error(fmt.Sprintf("Integer Conversion didn't work for: %w", val))
				}
			num = float64(num_int)
			
			log.Info(fmt.Sprintf("Integer Finally Converted to Float64: %.2f", num))
		}
		nums = append(nums, num)
	}

	log.Info("Completed Float Parsing")
	return nums
}

// This function calculates the mean value of a `[]float64`
func CalculateMean(nums []float64) (float64, error){
	log.Info("Calculating Mean")
	if len(nums) == 0{
		log.Error("Provided Float Slice is empty")
		return 0, errors.New("Provided Float Slice is empty")
	}
	var sum float64
	for _, num := range(nums){
		sum += num
	}
	return sum / float64(len(nums)), nil
}

// This function calculates the median value of a `[]float64`
func CalculateMedian(nums []float64) (float64, error){
	log.Info("Calculating Median")
	n := len(nums)
	if  n == 0{
		log.Error("Provided Float Slice is empty")
		return 0, errors.New("Provided Float Slice is empty")
	}

	sorted := make([]float64, n)
	copy(sorted, nums)
	sort.Float64s(sorted)

	if n%2 == 0{
		return (sorted[n/2-1] + sorted[n/2]) / 2.0, nil
	}

	return	sorted[n/2], nil
	
}

// This function calculates the mode value of a `[]float64`
func CalculateModeNum(nums []float64) (float64, error){
	
	log.Info("Calculating Numerical Mode")
	if len(nums) == 0{
		log.Error("Provided Float Slice is empty")
		return 0, errors.New("Provided Float Slice is empty")
	}

	counts := make(map[float64]int)
	maxFreq := 0

	var mode float64
	for _, num := range(nums){
		counts[num]++
		if counts[num] > maxFreq{
			mode = num
		}
	}
	return mode, nil
}

// This function calculates mode for a `[]string`
func CalculateModeString(strs []string) (string, error){
	log.Info("Calculating String Mode")
	if len(strs) == 0{
		log.Error("Provided String Slice is empty")
		return "", errors.New("Provided String Slice is empty")
	}

	counts := make(map[string]int)
	maxFreq := 0

	var mode string
	for _, str := range(strs){
		counts[str]++
		if counts[str] > maxFreq{
			mode = str
		}
	}
	return mode, nil
}



// func main(){

// 	csv, err := LoadCSV("/run/media/programmerrez/Field Testing/Side-Projects/Arnadillo/Data/customer_master.csv")

// 	if err != nil{
// 		fmt.Println(err)
// 	}

	// fmt.Println(csv)

	// fmt.Println(csv.getNulls())
	
	// fmt.Println(csv.Data[0])

	// fmt.Print(csv.FileShape)

	// stats := csv.GetStats()

	// fmt.Println(stats.DtypesMatrix)
	// fmt.Println(stats.Nulls)
	// fmt.Println("Nulls by Col: ", stats.NullsByCol)
	// fmt.Println(stats.Dupes)
	// fmt.Println(stats.UniqueValueMatrix)

	// csv.FillNA("customer_state", "mode")
	// fmt.Println(csv.GetStats().ColStats["customer_state"].NumStats.Mode)

	// fmt.Println(csv.getNulls())

	// for _, stat := range(stats.ColStats){	
	// 	fmt.Println(stat)
	// }


// }	