package main

import (
	"encoding/csv"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/kpawlik/geojson"
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

	// Building database
	database = csv2database()
	// Building autocomplete lookup
	auto := new(autocomplete)
	auto.generateTree(database)

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		// All places
		rest.Get("/places", func(w rest.ResponseWriter, r *rest.Request) {
			encodePlaces(w, r, database.allPlaces)
		}),

		// Places by specific object
		rest.Get("/places/by_movie/:id", func(w rest.ResponseWriter, r *rest.Request) {
			m, ok := database.byIdMovies[r.PathParam("id")]
			if ok {
				encodePlaces(w, r, m.Places)
			} else {
				rest.Error(w, "Not Found", http.StatusNotFound)
			}
		}),
		rest.Get("/places/by_company/:id", func(w rest.ResponseWriter, r *rest.Request) {
			c, ok := database.byIdCompanies[r.PathParam("id")]
			if ok {
				places := make([]*Place, 0)
				for _, m := range c.Movies.List {
					for _, p := range m.Places {
						places = append(places, p)
					}
				}
				encodePlaces(w, r, places)
			} else {
				rest.Error(w, "Not Found", http.StatusNotFound)
			}
		}),
		rest.Get("/places/by_person/:id", func(w rest.ResponseWriter, r *rest.Request) {
			p, ok := database.byIdPersons[r.PathParam("id")]
			if ok {
				places := make([]*Place, 0)
				for _, m := range p.Movies.List {
					for _, place := range m.Places {
						places = append(places, place)
					}
				}
				encodePlaces(w, r, places)
			} else {
				rest.Error(w, "Not Found", http.StatusNotFound)
			}
		}),

		// Just lists of everything, mostly for debugging purpose
		rest.Get("/movies", func(w rest.ResponseWriter, r *rest.Request) {
			encodeList(w, r, database.allMovies)
		}),
		rest.Get("/persons", func(w rest.ResponseWriter, r *rest.Request) {
			encodeList(w, r, database.allPersons)
		}),
		rest.Get("/companies", func(w rest.ResponseWriter, r *rest.Request) {
			encodeList(w, r, database.allCompanies)
		}),

		// Autocomplete
		rest.Get("/autocomplete", func(w rest.ResponseWriter, r *rest.Request) {
			q := strings.ToLower(r.URL.Query().Get("q"))
			// Do we want to handle UTF here?
			if len(q) > 1 {
				encodeList(w, r, auto.searchObjects(q))
			} else {
				rest.Error(w, "Minimum query length is 2 characters", http.StatusBadRequest)
			}

		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)

	http.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))
	http.Handle("/", http.FileServer(http.Dir("./static/")))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Places to GeoJSON conversion
func encodePlaces(w rest.ResponseWriter, r *rest.Request, places []*Place) {
	fc := geojson.NewFeatureCollection(make([]*geojson.Feature, 0))
	// TODO: fix to avoid copying on every request
	for _, p := range places {
		properties := make(map[string]interface{})
		properties["name"] = p.Name
		fc.AddFeatures(geojson.NewFeature(geojson.NewPoint(geojson.Coordinate{geojson.CoordType(p.Longitude), geojson.CoordType(p.Latitude)}), properties, p.Id))
	}
	w.WriteJson(&fc)
}

// Generic output for a list
func encodeList(w rest.ResponseWriter, r *rest.Request, list interface{}) {
	w.WriteJson(&list)
}
