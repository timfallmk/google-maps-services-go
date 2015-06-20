// Copyright 2015 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// More information about Google Directions API is available on
// https://developers.google.com/maps/documentation/directions/

package maps // import "google.golang.org/maps"

import (
	"net/http"
	"reflect"
	"testing"
)

func TestElevationDenver(t *testing.T) {

	// Elevation of Denver, the mile high city
	response := `{
   "results" : [
      {
         "elevation" : 1608.637939453125,
         "location" : {
            "lat" : 39.73915360,
            "lng" : -104.98470340
         },
         "resolution" : 4.771975994110107
      }
   ],
   "status" : "OK"
}`

	server := mockServer(200, response)
	defer server.Close()
	client := &http.Client{}
	ctx := newContextWithBaseURL(apiKey, client, server.URL)
	r := &ElevationRequest{
		Locations: []LatLng{
			LatLng{
				Lat: 39.73915360,
				Lng: -104.9847034,
			},
		},
	}

	resp, err := r.Get(ctx)

	if len(resp) != 1 {
		t.Errorf("Expected 1 result, got %v", len(resp))
	}
	if err != nil {
		t.Errorf("r.Get returned non nil error: %v", err)
	}

	correctResponse := ElevationResult{
		Location: &LatLng{
			Lat: 39.73915360,
			Lng: -104.98470340,
		},
		Elevation:  1608.637939453125,
		Resolution: 4.771975994110107,
	}

	if !reflect.DeepEqual(resp[0], correctResponse) {
		t.Errorf("Actual response != expected")
	}
}

func TestElevationSampledPath(t *testing.T) {

	// Elevation of Denver, the mile high city
	response := `{
  "results" : [
        {
           "elevation" : 4411.941894531250,
           "location" : {
              "lat" : 36.5785810,
              "lng" : -118.2919940
           },
           "resolution" : 19.08790397644043
        },
        {
           "elevation" : 1381.861694335938,
           "location" : {
              "lat" : 36.41150289067028,
              "lng" : -117.5602607523847
           },
           "resolution" : 19.08790397644043
        },
        {
           "elevation" : -84.61699676513672,
           "location" : {
              "lat" : 36.239980,
              "lng" : -116.831710
           },
           "resolution" : 19.08790397644043
        }
     ],
     "status" : "OK"
  }`

	server := mockServer(200, response)
	defer server.Close()
	client := &http.Client{}
	ctx := newContextWithBaseURL(apiKey, client, server.URL)
	r := &ElevationRequest{
		Path: []LatLng{
			LatLng{Lat: 36.578581, Lng: -118.291994},
			LatLng{Lat: 36.23998, Lng: -116.83171},
		},
		Samples: 3,
	}

	resp, err := r.Get(ctx)

	if len(resp) != 3 {
		t.Errorf("Expected 3 results, got %v", len(resp))
	}
	if err != nil {
		t.Errorf("r.Get returned non nil error: %v", err)
	}

	correctResponse := ElevationResult{
		Location: &LatLng{
			Lat: 36.5785810,
			Lng: -118.2919940,
		},
		Elevation:  4411.941894531250,
		Resolution: 19.08790397644043,
	}

	if !reflect.DeepEqual(resp[0], correctResponse) {
		t.Errorf("Actual response != expected")
	}
}

func TestElevationNoPathOrLocations(t *testing.T) {
	client := &http.Client{}
	ctx := NewContext(apiKey, client)
	r := &ElevationRequest{}

	if _, err := r.Get(ctx); err == nil {
		t.Errorf("Missing both Path and Locations should return error")
	}
}

func TestElevationPathWithNoSamples(t *testing.T) {
	client := &http.Client{}
	ctx := NewContext(apiKey, client)
	r := &ElevationRequest{
		Path: []LatLng{
			LatLng{Lat: 36.578581, Lng: -118.291994},
			LatLng{Lat: 36.23998, Lng: -116.83171},
		},
	}

	if _, err := r.Get(ctx); err == nil {
		t.Errorf("Missing both Path and Locations should return error")
	}
}

func TestElevationFailingServer(t *testing.T) {
	server := mockServer(500, `{"status" : "ERROR"}`)
	defer server.Close()
	client := &http.Client{}
	ctx := newContextWithBaseURL(apiKey, client, server.URL)
	r := &ElevationRequest{
		Path: []LatLng{
			LatLng{Lat: 36.578581, Lng: -118.291994},
			LatLng{Lat: 36.23998, Lng: -116.83171},
		},
		Samples: 3,
	}

	if _, err := r.Get(ctx); err == nil {
		t.Errorf("Failing server should return error")
	}
}
