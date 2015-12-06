(function() {
    "use strict";

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

    function initMap() {
        map = L.map('map').setView([37.77, -122.42], 13);

        L.tileLayer('https://api.tiles.mapbox.com/v4/{id}/{z}/{x}/{y}.png?access_token=pk.eyJ1IjoibWFwYm94IiwiYSI6IjZjNmRjNzk3ZmE2MTcwOTEwMGY0MzU3YjUzOWFmNWZhIn0.Y8bhBaUMqFiPrDRW9hieoQ', {
            maxZoom: 18,
            attribution: 'Map data &copy; <a href="http://openstreetmap.org">OpenStreetMap</a> contributors, ' +
                '<a href="http://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>, ' +
                'Imagery © <a href="http://mapbox.com">Mapbox</a>',
            id: 'mapbox.streets'
        }).addTo(map);

        var geojsonMarkerOptions = {
            radius: 6,
            fillColor: "#ff7800",
            color: "#ff9900",
            weight: 1,
            opacity: 0.9,
            fillOpacity: 0.7
        };

        geojsonLayer = new L.GeoJSON.AJAX("/api/places", {
            pointToLayer: function (feature, latlng) {
                return L.circleMarker(latlng, geojsonMarkerOptions);
            }
        });

        geojsonLayer.addTo(map);

        // function popUp(f,l){
        //     var out = [];
        //     if (f.properties){
        //         for(key in f.properties){
        //             out.push(key+": "+f.properties[key]);
        //         }
        //         l.bindPopup(out.join("<br />"));
        //     }
        // }
        // var jsonTest = new L.GeoJSON.AJAX(["/api/places"],{onEachFeature:popUp}).addTo(map);


    }

    initMap();



    $("input[name=search]").autocomplete({
        minLength: 2,
        source: function(request, response) {
            $.ajax({
                type: "GET",
                url: "/api/autocomplete",
                data: {
                    q: request.term
                },
                dataType: "json",
                success: function(msg){
                    var result = [];
                   console.log(msg);
                   for (var i in msg) {
                        var item = msg[i];
                        result.push({
                            label: item.type + ": " + (item.object.title ? item.object.title : item.object.name),
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
            console.log(event, ui, makeUrl(ui.item.type, ui.item.id));
            updateGeoJsonLayer(makeUrl(ui.item.type, ui.item.id));
        }
    });

        // function onMapClick(e) {
        //     popup
        //         .setLatLng(e.latlng)
        //         .setContent("You clicked the map at " + e.latlng.toString())
        //         .openOn(map);
        // }

        // map.on('click', onMapClick);

})();