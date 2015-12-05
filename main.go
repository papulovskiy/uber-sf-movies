package main

import (
	"encoding/csv"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/ant0ine/go-json-rest/rest"
)

var database *db

func csv2database() *db {
	list, err := ioutil.ReadFile("./data/list.csv")
	if err != nil {
		log.Fatal("Cannot open data file: %+v\n", err)
	}

	r := csv.NewReader(strings.NewReader(string(list)))
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	database := new(db)
	database.importArray(records)
	return database
}

func main() {
	log.Println("Hello")
	database = csv2database()

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		// rest.Get("/", index),
		rest.Get("/places", getAllPlaces),
		// rest.Get("/places/by_movie/:id", getPlacesByMovie),
		// rest.Get("/places/by_company/:id", getPlacesByCompany),
		// rest.Get("/places/by_person/:id", getPlacesByPerson),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	http.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))
	http.Handle("/", http.FileServer(http.Dir("./static/")))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getAllPlaces(w rest.ResponseWriter, r *rest.Request) {
	places := make([]*Place, 0)
	// TODO: fix to avoid copying on every request
	for _, p := range database.places {
		places = append(places, p)
	}
	w.WriteJson(&places)
}
