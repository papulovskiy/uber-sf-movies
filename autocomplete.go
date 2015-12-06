package main

import (
	"regexp"
	"strings"

	"github.com/armon/go-radix"
)

const (
	TYPE_MOVIE   = iota
	TYPE_COMPANY = iota
	TYPE_PLACE   = iota
	TYPE_PERSON  = iota
)

type entity struct {
	label  string
	object int
}

type Result struct {
	Type   string      `json:"type"`
	Object interface{} `json:"object"`
}

type autocomplete struct {
	d *db
	r *radix.Tree
}

func (auto *autocomplete) init() {
	auto.r = radix.New()
}

func filterString(s string) string {
	re := regexp.MustCompile("(\\W+)")
	return re.ReplaceAllString(strings.ToLower(s), "")
}

func (auto *autocomplete) generateTree(d *db) {
	auto.init()

	auto.d = d

	// Add movies
	for _, movie := range auto.d.movies {
		auto.r.Insert(filterString(movie.Title), &entity{label: movie.Title, object: TYPE_MOVIE})
	}

	// Add companies
	for _, company := range auto.d.companies {
		auto.r.Insert(filterString(company.Name), &entity{label: company.Name, object: TYPE_COMPANY})
	}

	// Add places
	for _, place := range auto.d.places {
		auto.r.Insert(filterString(place.Name), &entity{label: place.Name, object: TYPE_PLACE})
	}

	// Add persons
	for _, person := range auto.d.persons {
		auto.r.Insert(filterString(person.Name), &entity{label: person.Name, object: TYPE_PERSON})
	}
}

func (auto *autocomplete) entity2data(e entity) interface{} {
	switch e.object {
	case TYPE_MOVIE:
		return &Result{Type: "movie", Object: auto.d.movies[e.label]}
	case TYPE_COMPANY:
		return &Result{Type: "company", Object: auto.d.companies[e.label]}
	case TYPE_PLACE:
		return &Result{Type: "place", Object: auto.d.places[e.label]}
	case TYPE_PERSON:
		return &Result{Type: "person", Object: auto.d.persons[e.label]}
	}
	return nil
}

func (auto *autocomplete) search(query string) []entity {
	result := make([]entity, 0)

	fn := func(s string, v interface{}) bool {
		e := entity(*v.(*entity))
		result = append(result, e)
		return false
	}

	auto.r.WalkPrefix(filterString(query), fn)
	return result
}

func (auto *autocomplete) searchObjects(query string) []interface{} {
	result := make([]interface{}, 0)
	entities := auto.search(query)
	for _, e := range entities {
		result = append(result, auto.entity2data(e))
	}
	return result
}
