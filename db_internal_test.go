package main

import (
	"testing"
)

func testEmptyness(database *db, t *testing.T) {
	if len(database.movies) != 0 {
		t.Error("Database (movies) is not empty")
	}
	if len(database.places) != 0 {
		t.Error("Database (places) is not empty")
	}
	if len(database.persons) != 0 {
		t.Error("Database (persons) is not empty")
	}
	if len(database.companies) != 0 {
		t.Error("Database (companies) is not empty")
	}
}

func TestEmptyDb(t *testing.T) {
	database := new(db)
	testEmptyness(database, t)
}

func TestCompanies(t *testing.T) {
	database := new(db)
	database.init()

	if len(database.companies) != 0 {
		t.Error("Database (companies) is not empty")
	}

	database.createCompany("TheCompany")
	if len(database.companies) != 1 {
		t.Error("Database (companies) does not have a company")
	}

	database.createCompany("TheCompany")
	if len(database.companies) != 1 {
		t.Error("Database (companies) has incorrect number of companies")
	}

	database.createCompany("TheCompany2")
	if len(database.companies) != 2 {
		t.Error("Database (companies) has incorrect number of companies")
	}
}

func TestPlaces(t *testing.T) {
	database := new(db)
	database.init()

	if len(database.places) != 0 {
		t.Error("Database (places) is not empty")
	}

	database.createPlace("ThePlace", 10.2, 30.4)
	if len(database.places) != 1 {
		t.Error("Database does not have a place")
	}

	database.createPlace("ThePlace", 10.3, 31.5)
	if len(database.places) != 1 {
		t.Error("Database has incorrect number of places")
	}

	database.createPlace("ThePlace2", 22.0, -18.4)
	if len(database.places) != 2 {
		t.Error("Database has incorrect number of places")
	}
}

func TestPersons(t *testing.T) {
	database := new(db)
	database.init()

	if len(database.persons) != 0 {
		t.Error("Database (persons) is not empty")
	}

	database.createPerson("ThePerson")
	if len(database.persons) != 1 {
		t.Error("Database does not have a person")
	}

	database.createPerson("ThePerson")
	if len(database.persons) != 1 {
		t.Error("Database has incorrect number of persons")
	}

	database.createPerson("ThePerson2")
	if len(database.persons) != 2 {
		t.Error("Database has incorrect number of persons")
	}
}

func TestImportEmptyDb(t *testing.T) {
	database := new(db)

	database.importArray([][]string{{"one", "two", "three"}})
	testEmptyness(database, t)
}

func TestImportPlace(t *testing.T) {
	database := new(db)
	database.importArray([][]string{{"Locations", "latitude", "longitude"}, {"Somewhere", "1.54", "-10.152"}})
	if len(database.places) != 1 {
		t.Error("Database contains incorrect number of places")
	}
	v, ok := database.places["Somewhere"]
	if !ok {
		t.Error("Database does not contain required place")
	}
	if v.Latitude != 1.54 || v.Longitude != -10.152 {
		t.Error("Place has incorrect coordinates")
	}
}

func TestImportCompanies(t *testing.T) {
	database := new(db)
	database.importArray([][]string{{"Production Company", "Distributor"}, {"TheFirst", "TheSecond"}})
	if len(database.companies) != 2 {
		t.Error("Database contains incorrect number of companies")
	}
	_, ok := database.companies["TheFirst"]
	if !ok {
		t.Error("Database does not contain required company")
	}
	_, ok = database.companies["TheSecond"]
	if !ok {
		t.Error("Database does not contain required company")
	}
}

func TestImportPersons(t *testing.T) {
	database := new(db)
	database.importArray([][]string{{"Director", "Writer", "Actor 1", "Actor 2", "Actor 3"}, {"First", "Second", "First", "Third", "Fourth"}})
	if len(database.persons) != 4 {
		t.Error("Database contains incorrect number of persons")
	}
	_, ok := database.persons["First"]
	if !ok {
		t.Error("Database does not contain required person")
	}
	_, ok = database.persons["Fourth"]
	if !ok {
		t.Error("Database does not contain required person")
	}
}

func TestImportMovieWithSingleLocation(t *testing.T) {
	database := new(db)
	database.importArray([][]string{{"Title", "Year", "Locations", "latitude", "longitude"}, {"TheMovie", "2000", "Somewhere", "20.5", "-12.99"}})
	if len(database.movies) != 1 {
		t.Error("Database contains incorrect number of movies")
	}
	if len(database.places) != 1 {
		t.Error("Database contains incorrect number of places")
	}
	m, ok := database.movies["TheMovie"]
	if !ok {
		t.Error("Database does not contain required movie")
	}
	if m.Year != 2000 {
		t.Error("Year transformed incorrectly")
	}
	if len(m.Places) != 1 {
		t.Error("Database does not contain required place")
	}
	if m.Places[0].Name != "Somewhere" {
		t.Error("Location name imported incorrectly")
	}
}

func TestImportMovieWithMultipleLocations(t *testing.T) {
	database := new(db)
	database.importArray([][]string{{"Title", "Year", "Locations", "latitude", "longitude"}, {"TheMovie", "2000", "Somewhere", "20.5", "-12.99"}, {"TheMovie", "2000", "Somewhere2", "33.22", "66.55"}})
	if len(database.movies) != 1 {
		t.Error("Database contains incorrect number of movies")
	}
	if len(database.places) != 2 {
		t.Error("Database contains incorrect number of places")
	}
	m, ok := database.movies["TheMovie"]
	if !ok {
		t.Error("Database does not contain required movie")
	}
	if m.Year != 2000 {
		t.Error("Year transformed incorrectly")
	}
	if len(m.Places) != 2 {
		t.Error("Database does not contain required place")
	}
	if m.Places[0].Name != "Somewhere" {
		t.Error("Location name imported incorrectly")
	}
	if m.Places[1].Name != "Somewhere2" {
		t.Error("Location name imported incorrectly")
	}
}

func TestImportMoviesWithMultipleLocations(t *testing.T) {
	database := new(db)
	database.importArray([][]string{{"Title", "Year", "Locations", "latitude", "longitude"}, {"TheMovie", "2000", "Somewhere", "20.5", "-12.99"}, {"TheMovie", "2000", "Somewhere2", "33.22", "66.55"}, {"TheMovie2", "2002", "Somewhere3", "10.0", "-2"}})
	if len(database.movies) != 2 {
		t.Error("Database contains incorrect number of movies")
	}
	if len(database.places) != 3 {
		t.Error("Database contains incorrect number of places")
	}
	m, _ := database.movies["TheMovie"]
	if len(m.Places) != 2 {
		t.Error("Database does not contain required place")
	}
	m, _ = database.movies["TheMovie2"]
	if len(m.Places) != 1 {
		t.Error("Database does not contain required place")
	}
}

func TestImportMoviesWithSingleLocation(t *testing.T) {
	database := new(db)
	database.importArray([][]string{{"Title", "Year", "Locations", "latitude", "longitude"}, {"TheMovie", "2000", "Somewhere", "20.5", "-12.99"}, {"TheMovie2", "2002", "Somewhere", "20.5", "-12.99"}})
	if len(database.movies) != 2 {
		t.Error("Database contains incorrect number of movies")
	}
	if len(database.places) != 1 {
		t.Error("Database contains incorrect number of places")
	}
	m, _ := database.movies["TheMovie"]
	if len(m.Places) != 1 {
		t.Error("Database does not contain required place")
	}
	m, _ = database.movies["TheMovie2"]
	if len(m.Places) != 1 {
		t.Error("Database does not contain required place")
	}
}

func TestImportedLinks(t *testing.T) {
	database := new(db)
	database.importArray([][]string{
		{"Title", "Year", "Locations", "latitude", "longitude", "Director", "Actor 3", "Distributor"},
		{"TheMovie", "2000", "Somewhere", "20.5", "-12.99", "TheGuy", "TheDude", "TheCompany"},
		{"TheMovie2", "2002", "Somewhere", "20.5", "-12.99", "OtherGuy", "TheDude", "TheDistributor"},
	})

	if len(database.movies) != 2 {
		t.Error("Database contains incorrect number of movies")
	}
	if len(database.places) != 1 {
		t.Error("Database contains incorrect number of places")
	}
	if len(database.persons) != 3 {
		t.Error("Database contains incorrect number of persons")
	}
	if len(database.companies) != 2 {
		t.Error("Database contains incorrect number of companies")
	}

	database.createLinks()

	if len(database.persons["TheGuy"].Movies.List) != 1 {
		t.Error("Incorrect number of movies for person")
	}
	if len(database.persons["TheDude"].Movies.List) != 2 {
		t.Error("Incorrect number of movies for person")
	}

	if len(database.companies["TheDistributor"].Movies.List) != 1 {
		t.Error("Incorrect number of movies for company")
	}

	if len(database.places["Somewhere"].Movies.List) != 2 {
		t.Error("Incorrect number of movies for place")
	}

}

// TODO: test for empty place in movie
// TODO: test for multiple actors in movie
