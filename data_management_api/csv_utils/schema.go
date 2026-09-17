package utils



// This will be the schema for representing data for each column
type ColumnStats struct{
	UniqueValues  	[]string	`json:"unique_values"`
	DataType 		string		`json:"dtypes"`
	NullsInCol		int64		`json:"nulls_in_col"`
}


// This will be the final output of the getStats method
type DataFrameStats struct{
	DtypesMatrix		map[string]string		`json:"dtype_matrix"`
	ColStats			map[string]ColumnStats	`json:"col_stats"`
	Dupes 				int64					`json:"dupes"`
	Nulls				int64					`json:"nulls"`
	NullsByCol			map[string]int64		`json:"null_matrix"`
	UniqueValueMatrix 	map[string]int64		`json:"unique_value_matrix"`
}


