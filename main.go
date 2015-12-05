package main

import (
	"encoding/csv"
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	log.Println("Hello")

	list, err := ioutil.ReadFile("./data/list.csv")
	if err != nil {
		log.Fatal("Cannot open data file: %+v\n", err)
	}
	// log.Printf("%+v", list)
	r := csv.NewReader(strings.NewReader(string(list)))
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%+v\n", len(records))
	database := new(db)
	database.importArray(records)
}
