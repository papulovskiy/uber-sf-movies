(function() {
    'use strict';

    var map,
        geojsonLayer;

    function makeUrl(type, id) {
        if(type === undefined || id === undefined) {
            return '/api/places';
        }
        return '/api/places/by_' + type + '/' + id;
    }

    function updateGeoJsonLayer(url) {
        geojsonLayer.refresh(url);
    }

    function renderPopup(popup, place_id) {
        $.get('/api/place/' + place_id).done(function(data) {
            var html = '';
            if(data.name) {
                html += '<h3>' + data.name + '</h3>';
            }
            if(data.movies.List) {
                data.movies.List.forEach(function(i) {
                    html += '<div class="movie">Movie: ';
                    html += '<a href="#" class="movie" data-type="movie" data-id="' + i.id + '" data-title="' + i.title + '">';
                    html += i.title;
                    if(i.year > 0) {
                        html += ' (' + i.year + ')';
                    }
                    html += '</a></div>';
                });
            }
            popup.setContent(html);
            popup.update();
        });

    }

    function initMap() {
        map = L.map('map').setView([37.77, -122.42], 13);

        L.tileLayer('https://api.tiles.mapbox.com/v4/{id}/{z}/{x}/{y}.png?access_token=pk.eyJ1IjoibWFwYm94IiwiYSI6IjZjNmRjNzk3ZmE2MTcwOTEwMGY0MzU3YjUzOWFmNWZhIn0.Y8bhBaUMqFiPrDRW9hieoQ', {
            maxZoom: 18,
            attribution: 'Map data &copy; <a href="http://openstreetmap.org">OpenStreetMap</a> contributors, ' +
                '<a href="http://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>, ' +
                'Imagery © <a href="http://mapbox.com">Mapbox</a>, ' +
                '<a href="https://www.scaleway.com/">Hosted on Scaleway</a>, ' +
                'SSL by <a href="https://letsencrypt.org/">Let\'s Encrypt</a>, ' +
                '<a href="https://github.com/papulovskiy/uber-sf-movies">Source Code</a>',
            id: 'mapbox.streets'
        }).addTo(map);

        var geojsonMarkerOptions = {
            radius: 6,
            fillColor: '#ff7800',
            color: '#ff9900',
            weight: 1,
            opacity: 0.9,
            fillOpacity: 0.7
        };

        // Initial map shows all known places
        geojsonLayer = new L.GeoJSON.AJAX('/api/places', {
            pointToLayer: function (feature, latlng) {
                var m = L.circleMarker(latlng, geojsonMarkerOptions);
                m.bindPopup('<div>Loading...</div>');
                m.on('click', function(e) {
                    var p = e.target._popup;
                    renderPopup(p, feature.id);
                });
                return m;
            }
        });

        geojsonLayer.addTo(map);

    }

    initMap();


    $('input[name=search]').autocomplete({
        minLength: 2,
        source: function(request, response) {
            $.ajax({
                type: 'GET',
                url: '/api/autocomplete',
                data: {
                    q: request.term
                },
                dataType: 'json',
                success: function(msg){
                    var result = [];
                   for (var i in msg) {
                        var item = msg[i];
                        result.push({
                            label: item.type + ': ' + (item.object.title ? item.object.title : item.object.name),
                            value: (item.object.title ? item.object.title : item.object.name),
                            type: item.type,
                            id: item.object.id
                        });
                   }
                    response(result);
                }
            });
        },
        select: function(event, ui) {
            updateGeoJsonLayer(makeUrl(ui.item.type, ui.item.id));
        }
    });

    // Switch map to places by movie
    $("html").on('click', 'a.movie', function(e) {
        e.stopPropagation();
        var a = e.target;
        $('input[name=search]').val($(a).data('title'))
        updateGeoJsonLayer(makeUrl('movie', $(a).data('id')));
    });

    // Cleanup search on header click
    $('h1 a').on('click', function(e) {
        e.stopPropagation();
        $('input[name=search]').val('')
        updateGeoJsonLayer(makeUrl());
    });
})();