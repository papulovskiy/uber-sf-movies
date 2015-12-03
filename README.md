# uber-sf-movies
This is just a test implementation of SF movies story for [Uber coding challenge](https://github.com/uber/coding-challenge-tools/blob/master/coding_challenge.md).

### Preparation
[Given list](https://data.sfgov.org/Culture-and-Recreation/Film-Locations-in-San-Francisco/yitu-d5am) of places was saved as CSV and to show it on a map I'd prefer to geocode it to avoid real-time Google Maps API usage.
I made a simple perl script to geocode locations via GMaps and store another CSV file.