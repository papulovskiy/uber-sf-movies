package main

import (
	"strconv"
)

var idSequencer int

func getNextId() int {
	i := idSequencer
	idSequencer++
	return i
}

type movie struct {
	id          int
	title       string
	year        int
	production  *company
	distributor *company
	places      []*place
	director    *person
	writer      *person
	actors      []*person
}

func (m *movie) addPlace(place *place) {
	m.places = append(m.places, place)
}

func (m *movie) addActor(actor *person) {
	m.actors = append(m.actors, actor)
}

type movies struct {
	list []*movie
}

func (m *movies) link(movie *movie) {
	m.list = append(m.list, movie)
}

type person struct {
	id     int
	name   string
	movies movies
}

type company struct {
	id     int
	name   string
	movies movies
}

type place struct {
	id        int
	name      string
	latitude  float64
	longitude float64
	movies    movies
}

type db struct {
	movies    map[string]*movie
	companies map[string]*company
	persons   map[string]*person
	places    map[string]*place
}

func (d *db) init() {
	d.movies = make(map[string]*movie)
	d.companies = make(map[string]*company)
	d.persons = make(map[string]*person)
	d.places = make(map[string]*place)
}

func (d *db) createPlace(name string, latitude float64, longitude float64) *place {
	_, ok := d.places[name]
	if !ok {
		d.places[name] = &place{id: getNextId(), name: name, latitude: latitude, longitude: longitude}
	}
	return d.places[name]
}

func (d *db) createCompany(name string) *company {
	_, ok := d.companies[name]
	if !ok {
		d.companies[name] = &company{id: getNextId(), name: name}
	}
	return d.companies[name]
}

func (d *db) createPerson(name string) *person {
	_, ok := d.persons[name]
	if !ok {
		d.persons[name] = &person{id: getNextId(), name: name}
	}
	return d.persons[name]
}

func (d *db) importArray(array [][]string) {
	d.init()
	headers := array[0]
	ref := make(map[string]int)

	// making reference between field name and array column
	for _, header := range []string{"Locations", "latitude", "longitude", "Production Company", "Distributor", "Director", "Writer", "Actor 1", "Actor 2", "Actor 3", "Title", "Year"} {
		ref[header] = -1 // set every field as not found
	}
	for i := 0; i < len(headers); i++ {
		ref[headers[i]] = i
	}

	for i := 1; i < len(array); i++ {
		row := array[i]

		var director, writer, actor1, actor2, actor3 *person
		var production, distributor *company
		var location *place

		// create place
		if ref["Locations"] >= 0 && ref["latitude"] >= 0 && ref["longitude"] >= 0 &&
			row[ref["Locations"]] != "" && row[ref["latitude"]] != "" && row[ref["longitude"]] != "" {
			lat, errLat := strconv.ParseFloat(row[ref["latitude"]], 64)
			lon, errLon := strconv.ParseFloat(row[ref["longitude"]], 64)
			if errLat == nil && errLon == nil {
				location = d.createPlace(row[ref["Locations"]], lat, lon)
			}
		}

		// create companies
		if ref["Production Company"] >= 0 && row[ref["Production Company"]] != "" {
			production = d.createCompany(row[ref["Production Company"]])
		}
		if ref["Distributor"] >= 0 && row[ref["Distributor"]] != "" {
			distributor = d.createCompany(row[ref["Distributor"]])
		}

		// create persons
		if ref["Director"] >= 0 && row[ref["Director"]] != "" {
			director = d.createPerson(row[ref["Director"]])
		}
		if ref["Writer"] >= 0 && row[ref["Writer"]] != "" {
			writer = d.createPerson(row[ref["Writer"]])
		}
		if ref["Actor 1"] >= 0 && row[ref["Actor 1"]] != "" {
			actor1 = d.createPerson(row[ref["Actor 1"]])
		}
		if ref["Actor 2"] >= 0 && row[ref["Actor 2"]] != "" {
			actor2 = d.createPerson(row[ref["Actor 2"]])
		}
		if ref["Actor 3"] >= 0 && row[ref["Actor 3"]] != "" {
			actor3 = d.createPerson(row[ref["Actor 3"]])
		}

		// create movie
		if ref["Title"] < 0 || row[ref["Title"]] == "" {
			continue
		}
		m, ok := d.movies[row[ref["Title"]]]
		if !ok {
			var year int64 = 0
			if ref["Year"] >= 0 {
				year, _ = strconv.ParseInt(row[ref["Year"]], 10, 32)
			}
			d.movies[row[ref["Title"]]] = &movie{
				id:          getNextId(),
				title:       row[ref["Title"]],
				year:        int(year),
				production:  production,
				distributor: distributor,
				director:    director,
				writer:      writer,
			}
			m = d.movies[row[ref["Title"]]]
		}
		m.addPlace(location)
		if actor1 != nil {
			m.addActor(actor1)
		}
		if actor2 != nil {
			m.addActor(actor2)
		}
		if actor3 != nil {
			m.addActor(actor3)
		}

	}
	// log.Printf("Places: %+v\n", d.companies)
}

func (d *db) createLinks() {
	for _, m := range d.movies {
		if m.director != nil {
			m.director.movies.link(m)
		}
		if m.writer != nil {
			m.writer.movies.link(m)
		}
		if m.production != nil {
			m.production.movies.link(m)
		}
		if m.distributor != nil {
			m.distributor.movies.link(m)
		}
		for _, p := range m.places {
			p.movies.link(m)
		}
		for _, a := range m.actors {
			a.movies.link(m)
		}
	}
}
