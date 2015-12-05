package main

import "github.com/armon/go-radix"

// import "fmt"

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

type autocomplete struct {
	d *db
	r *radix.Tree
}

func (auto *autocomplete) init() {
	auto.r = radix.New()
}

func (auto *autocomplete) generateTree(d *db) {
	auto.init()

	auto.d = d

	// Add movies
	for _, movie := range auto.d.movies {
		auto.r.Insert(movie.title, &entity{label: movie.title, object: TYPE_MOVIE})
	}

	// Add companies
	for _, company := range auto.d.companies {
		auto.r.Insert(company.name, &entity{label: company.name, object: TYPE_COMPANY})
	}

	// Add places
	for _, place := range auto.d.places {
		auto.r.Insert(place.name, &entity{label: place.name, object: TYPE_PLACE})
	}

	// Add persons
	for _, person := range auto.d.persons {
		auto.r.Insert(person.name, &entity{label: person.name, object: TYPE_PERSON})
	}
}

func (auto *autocomplete) entity2data(e entity) interface{} {
	switch e.object {
	case TYPE_MOVIE:
		return auto.d.movies[e.label]
	case TYPE_COMPANY:
		return auto.d.companies[e.label]
	case TYPE_PLACE:
		return auto.d.places[e.label]
	case TYPE_PERSON:
		return auto.d.persons[e.label]
	}
	return nil
}

func (auto *autocomplete) search(query string) []entity {
	result := make([]entity, 0)

	fn := func(s string, v interface{}) bool {
		e := entity(*v.(*entity))
		result = append(result, e)
		// } else {
		// 	fmt.Printf("Assertion error: %+v %+v\n", v, e)
		// }
		return false
	}

	auto.r.WalkPrefix(query, fn)
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
