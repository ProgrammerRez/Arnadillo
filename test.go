package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func main() {
	f, err := os.Open("Data/dataset_statistics.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil || len(data) == 0 {
		fmt.Println("Error reading CSV or file is empty")
		return
	}
	fmt.Println("Original Column Row:", data[0])
	fmt.Println("Original Data Row:", data[1])

	target := "Total Transactions"

	column_list := data[0]

	// 1. Find the target column index in the header row (data[0])
	targetColIdx := -1
	for colIdx, colName := range data[0] {
		if colName == target {
			targetColIdx = colIdx
			break
		}
	}

	if targetColIdx == -1 {
		fmt.Printf("Column %q not found\n", target)
		return
	}

	// 2. Remove targetColIdx from EVERY row (including header)
	for i := range data {
		// Guard against rows that might have fewer columns
		if targetColIdx < len(data[i]) {
			data[i] = append(data[i][:targetColIdx], data[i][targetColIdx+1:]...)
			column_list = append(column_list[:targetColIdx], column_list[targetColIdx:]...)
		}
	}

	// Output verification
	fmt.Printf("Deleted column %q at index %d\n", target, targetColIdx)
	fmt.Println("New header row:", column_list)
	fmt.Println("New Data Row:", data[1])
}