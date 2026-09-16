package main

import(
	"fmt"
	"encoding/csv"
	"os"
	"log"
)





func main(){
	file, err := os.Open("Data/customer_master.csv")

	if err != nil {
		log.Fatal("Can't load. Smth happened: ", err)
	}

	csv_reader := csv.NewReader(file)

	records, err := csv_reader.ReadAll()
	    // Checks for the error
    if err != nil{
        fmt.Println("Error reading records")
    }

	// for _, rec := range(records){
	// 	fmt.Println(rec)
	// }

	fmt.Print("%T", records[1])
}

