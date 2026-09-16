package main



// This will be the schema for representing data for each column
type ColumnStats struct{
	UniqueValues  	[]string
	DataType 		string
	NullsInCol		int64
}


// This will be the final output of the getStats method
type DataFrameStats struct{
	DtypesMatrix		map[string]string
	ColStats			map[string]ColumnStats
	Dupes 				int64
	Nulls				int64
	NullsByCol			map[string]int64
	UniqueValueMatrix 	map[string]int64
}


