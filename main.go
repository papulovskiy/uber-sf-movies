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
	// log.Printf("%+v", list)
	r := csv.NewReader(strings.NewReader(string(list)))
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// log.Printf("%+v\n", len(records))
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
		// rest.Get("/geo/by_movie/:id", getPlacesByMovie),
		// rest.Get("/geo/by_company/:id", getPlacesByCompany),
		// rest.Get("/geo/by_person/:id", getPlacesByPerson),
		// rest.Put("/users/:id", users.PutUser),
		// rest.Delete("/users/:id", users.DeleteUser),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	http.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))
	http.Handle("/", http.FileServer(http.Dir("./static/")))
	// http.Handle("/", http.StripPrefix("/static", http.FileServer(http.Dir("."))))
	// http.FileServer(http.Dir("/usr/share/doc"))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getAllPlaces(w rest.ResponseWriter, r *rest.Request) {
	places := make([]*Place, 0)
	for _, p := range database.places {
		places = append(places, p)
		log.Printf("Place: %+v\n", *p)
		// p[i] = interface{label: p.name, lat: p.latitude, lon: p.longitude}
	}
	log.Printf("Places: \n%+v\n%+v\n", database.places, places)
	w.WriteJson(&places)
}
