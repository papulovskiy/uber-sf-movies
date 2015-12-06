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

type Movie struct {
	Id          int
	Title       string
	Year        int
	Production  *Company
	Distributor *Company
	Places      []*Place
	Director    *Person
	Writer      *Person
	Actors      []*Person

	actorsMap map[int]bool
	placesMap map[int]bool
}

func (m *Movie) addPlace(place *Place) {
	if place != nil && m.placesMap[place.Id] == false {
		m.placesMap[place.Id] = true
		m.Places = append(m.Places, place)
	}
}

func (m *Movie) addActor(actor *Person) {
	if actor != nil && m.actorsMap[actor.Id] == false {
		m.actorsMap[actor.Id] = true
		m.Actors = append(m.Actors, actor)
	}
}

type Movies struct {
	List []*Movie
}

func (m *Movies) link(movie *Movie) {
	if m.List == nil {
		m.List = make([]*Movie, 0)
	}
	m.List = append(m.List, movie)
}

type Person struct {
	Id     int
	Name   string
	Movies Movies `json:"-"`
}

type Company struct {
	Id     int
	Name   string
	Movies Movies `json:"-"`
}

type Place struct {
	Id        int     `json:"id"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
	Movies    Movies  `json:"-"`
}

type db struct {
	movies    map[string]*Movie
	companies map[string]*Company
	persons   map[string]*Person
	places    map[string]*Place

	allPlaces    []*Place
	allMovies    []*Movie
	allCompanies []*Company
	allPersons   []*Person

	byIdPlaces    map[string]*Place
	byIdMovies    map[string]*Movie
	byIdCompanies map[string]*Company
	byIdPersons   map[string]*Person
}

func (d *db) init() {
	d.movies = make(map[string]*Movie)
	d.companies = make(map[string]*Company)
	d.persons = make(map[string]*Person)
	d.places = make(map[string]*Place)

	d.byIdMovies = make(map[string]*Movie)
	d.byIdCompanies = make(map[string]*Company)
	d.byIdPersons = make(map[string]*Person)
	d.byIdPlaces = make(map[string]*Place)
}

// Fill various lookups
func (d *db) updateLists() {
	d.allPlaces = make([]*Place, 0)
	for _, p := range d.places {
		d.allPlaces = append(d.allPlaces, p)
		d.byIdPlaces[strconv.Itoa(p.Id)] = p
	}

	d.allMovies = make([]*Movie, 0)
	for _, m := range d.movies {
		d.allMovies = append(d.allMovies, m)
		d.byIdMovies[strconv.Itoa(m.Id)] = m
	}

	d.allCompanies = make([]*Company, 0)
	for _, c := range d.companies {
		d.allCompanies = append(d.allCompanies, c)
		d.byIdCompanies[strconv.Itoa(c.Id)] = c
	}

	d.allPersons = make([]*Person, 0)
	for _, p := range d.persons {
		d.allPersons = append(d.allPersons, p)
		d.byIdPersons[strconv.Itoa(p.Id)] = p
	}
}

func (d *db) createPlace(name string, latitude float64, longitude float64) *Place {
	_, ok := d.places[name]
	if !ok {
		d.places[name] = &Place{Id: getNextId(), Name: name, Latitude: latitude, Longitude: longitude}
	}
	return d.places[name]
}

func (d *db) createCompany(name string) *Company {
	_, ok := d.companies[name]
	if !ok {
		d.companies[name] = &Company{Id: getNextId(), Name: name}
	}
	return d.companies[name]
}

func (d *db) createPerson(name string) *Person {
	_, ok := d.persons[name]
	if !ok {
		d.persons[name] = &Person{Id: getNextId(), Name: name}
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

		var director, writer, actor1, actor2, actor3 *Person
		var production, distributor *Company
		var location *Place

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
			d.movies[row[ref["Title"]]] = &Movie{
				Id:          getNextId(),
				Title:       row[ref["Title"]],
				Year:        int(year),
				Production:  production,
				Distributor: distributor,
				Director:    director,
				Writer:      writer,

				actorsMap: make(map[int]bool),
				placesMap: make(map[int]bool),
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

	d.createLinks()
	d.updateLists()
	// log.Printf("Places: %+v\n", d.companies)
}

func (d *db) createLinks() {
	for _, m := range d.movies {
		if m.Director != nil {
			m.Director.Movies.link(m)
		}
		if m.Writer != nil {
			m.Writer.Movies.link(m)
		}
		if m.Production != nil {
			m.Production.Movies.link(m)
		}
		if m.Distributor != nil {
			m.Distributor.Movies.link(m)
		}
		for _, p := range m.Places {
			// TODO: find out why the empty place is possible
			if p != nil {
				p.Movies.link(m)
			}
		}
		for _, a := range m.Actors {
			a.Movies.link(m)
		}
	}
}
