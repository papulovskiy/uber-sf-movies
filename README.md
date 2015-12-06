# uber-sf-movies
This is just a test implementation of SF movies story for [Uber coding challenge](https://github.com/uber/coding-challenge-tools/blob/master/coding_challenge.md).

### Preparation
[Given list](https://data.sfgov.org/Culture-and-Recreation/Film-Locations-in-San-Francisco/yitu-d5am) of places was saved as CSV and to show it on a map I'd prefer to geocode it to avoid real-time Google Maps API usage.
I made a simple Perl script to geocode locations via GMaps and store another CSV file.

### First step
I decided not to use any database for this exercise, because the whole dataset perfectly fits in memory, far that I've built data structure and implemented a method to import from a two-dimensional array (CSV).

### Autocomplete
For autocomplete I chose radix-tree data structure and implementation from [Armon Dadgar](github.com/armon/go-radix)

### List of points
I'd prefer to stick to standards, that's why I chose GeoJSON as a format for geodata that will be shown on the map. I used nice [implementation of GeoJSON](https://github.com/kpawlik/geojson) fo Go. 

### Map
The map is built with Leaflet and some jQuery, pretty simple and straightforward.

## Credits
### Backend part
[Go-Json-Rest](https://github.com/ant0ine/go-json-rest/rest) - Go library for making JSON REST API.
[GEOJSON](https://github.com/kpawlik/geojson) -  Go package to manipulate GeoJSON data, just to avoid manual JSON marshalling.
[go-radix](https://github.com/armon/go-radix) - Radix tree implementation in Go, for autocomplete search.

### Frontend part
[Leaflet](http://leafletjs.com/) - lightweight web map library.
[leaflet-ajax](https://github.com/calvinmetcalf/leaflet-ajax) - GeoJSON Leaflet layer with a remote data source.
[jQuery](https://jquery.com/) for AJAX data load.
[jQuery UI](https://jqueryui.com/) for autocomplete.

## Live demo
https://uber.play.c17.nl/
The project is running on ARM server by [Scaleway](https://www.scaleway.com/) with SSL certificate from [Let's Encrypt](https://letsencrypt.org/).

Special thanks to supervisord.
