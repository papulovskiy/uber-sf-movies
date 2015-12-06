# uber-sf-movies
This is just a test implementation of SF movies story for [Uber coding challenge](https://github.com/uber/coding-challenge-tools/blob/master/coding_challenge.md).

### Preparation
[Given list](https://data.sfgov.org/Culture-and-Recreation/Film-Locations-in-San-Francisco/yitu-d5am) of places was saved as CSV and to show it on a map I'd prefer to geocode it to avoid real-time Google Maps API usage.
I made a simple perl script to geocode locations via GMaps and store another CSV file.

### First step
I decided not to use any database for this exersize, because whole dataset perfectly fits in memory, far that I've built data structure and implemented a method to import from two-dimensional array (CSV).

### Autocomplete
For autocomplete I chose radix-tree data structure and implementation from [Armon Dadgar](github.com/armon/go-radix)

### List of points
I'd prefer to stick to standards, that's why I chose GeoJSON as a format for geodata that will be shown on the map. I used nice [implementation of GeoJSON](https://github.com/kpawlik/geojson) fo Go. 