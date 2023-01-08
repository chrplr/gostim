package main

import (
	"os"
	"encoding/csv"
	"log"
	"fmt"
)

func read_from_csv(fname string) (dataframe [][]string, err error) {
	f, err := os.Open(fname)
	if err != nil { return nil, err }
	defer f.Close()

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1

	dataframe, err = reader.ReadAll()
	if err != nil { return nil, err }
	return dataframe, nil
}


func main() {
	var df [][]string
	var err error
	
	df, err = read_from_csv("test.csv")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	for _, line := range df {
		for _, word := range line {
			fmt.Print(word, " ")
		}
		fmt.Println()
	}
}
