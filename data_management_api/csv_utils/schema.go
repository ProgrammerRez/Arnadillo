package utils
// package main


// This will be the schema for representing data for each column
type ColumnStats struct{
	UniqueValues  	[]string		`json:"unique_values"`
	DataType 		string			`json:"dtype"`
	NumStats		NumericalStats	`json:"num_stats"`
	NullsInCol		int64			`json:"nulls_in_col"`
}

// This will be the schema for Numerical Stats for a Column
type NumericalStats struct{
	Mode	any			`json:"mode"`
	Mean 	float64 	`json:"mean"`	
	Median	float64		`json:"median"`
}


// This will be the final output of the getStats method
type DataFrameStats struct{
	DtypesMatrix		map[string]string		`json:"dtype_matrix"`
	ColStats			map[string]*ColumnStats	`json:"col_stats"`
	NullsByCol			map[string]int64		`json:"null_matrix"`
	UniqueValueMatrix 	map[string]int64		`json:"unique_value_matrix"`
	Dupes 				int64					`json:"dupes"`
	Nulls				int64					`json:"nulls"`
}


