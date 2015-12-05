package main

import (
	// "fmt"
	"testing"
)

func TestSearch(t *testing.T) {
	auto := new(autocomplete)
	auto.init()

	if len(auto.search("Hello")) != 0 {
		t.Error("Search tree should be empty")
	}

	auto.r.Insert("Hello", &entity{label: "Hello", object: TYPE_MOVIE})
	if len(auto.search("Hello")) != 1 {
		t.Error("One result expected")
	}
	if len(auto.search("He")) != 1 {
		t.Error("One result expected")
	}

	auto.r.Insert("Healthy Food", &entity{label: "Healthy Food", object: TYPE_PLACE})
	if len(auto.search("Hello")) != 1 {
		t.Error("One result expected")
	}
	if len(auto.search("Hea")) != 1 {
		t.Error("One result expected")
	}
	if len(auto.search("He")) != 2 {
		t.Error("Two results expected")
	}
}
