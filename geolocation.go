// SPDX-FileCopyrightText: 2023 Winni Neessen <wn@neessen.dev>
//
// SPDX-License-Identifier: MIT

package meteologix

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strconv"
)

// OSMNominatimURL is the API endpoint URL for the OpenStreetMaps Nominatim API
const OSMNominatimURL = "https://nominatim.openstreetmap.org/search"

// ErrCityNotFound is returned if a requested city was not found in the OSM API
var ErrCityNotFound = errors.New("requested city not found in OSM Nominatim API")

// GeoLocation represent the GPS GeoLocation coordinates of a City
type GeoLocation struct {
	// Importance is the OSM computed importance rank
	Importance float64 `json:"importance"`
	// Latitude represents the GPS Latitude coordinates of the requested City as Float
	Latitude float64
	// LatitudeString represents the GPS Latitude coordinates of the requested City as string
	LatitudeString string `json:"lat"`
	// Longitude represents the GPS Longitude coordinates of the requested City as Float
	Longitude float64
	// LongitudeString represents the GPS Longitude coordinates of the requested City as String
	LongitudeString string `json:"lon"`
	// Name represents the requested City
	Name string `json:"display_name"`
	// PlaceID is the OSM Nominatim internal database ID
	PlaceID int64 `json:"place_id"`
}

// GetGeoLocationByName returns the GeoLocation with the highest importance based on
// the given City name
//
// This method makes use of the OSM Nominatim API
func (c *Client) GetGeoLocationByName(ci string) (GeoLocation, error) {
	ga, err := c.GetGeoLocationsByName(ci)
	if err != nil || len(ga) < 1 {
		return GeoLocation{}, err
	}
	return ga[0], nil
}

// GetGeoLocationsByName returns a slice of GeoLocation based on the requested City name
// The returned slice will be sorted by Importance of the results with the highest
// importance as first entry
//
// This method makes use of the OSM Nominatim API
func (c *Client) GetGeoLocationsByName(city string) ([]GeoLocation, error) {
	locations := make([]GeoLocation, 0)

	apiURL, err := url.Parse(OSMNominatimURL)
	if err != nil {
		return locations, fmt.Errorf("failed to parse OSM Nominatim URL: %w", err)
	}
	query := apiURL.Query()
	query.Add("format", "json")
	query.Add("q", city)
	apiURL.RawQuery = query.Encode()

	response, err := c.httpClient.Get(apiURL.String())
	if err != nil {
		return locations, fmt.Errorf("OSM Nominatim API request failed: %w", err)
	}
	var jsonLocations []GeoLocation
	if err = json.Unmarshal(response, &jsonLocations); err != nil {
		return locations, fmt.Errorf("failed to unmarshal API response JSON: %w", err)
	}
	if len(jsonLocations) < 1 {
		return locations, ErrCityNotFound
	}

	for _, location := range jsonLocations {
		latitude, err := strconv.ParseFloat(location.LatitudeString, 64)
		if err != nil {
			return locations, fmt.Errorf("failed to convert latitude string to float value: %w", err)
		}
		longitude, err := strconv.ParseFloat(location.LongitudeString, 64)
		if err != nil {
			return locations, fmt.Errorf("failed to convert longitude string to float value: %w", err)
		}
		location.Latitude = latitude
		location.Longitude = longitude
		locations = append(locations, location)
	}
	sort.SliceStable(locations, func(i, j int) bool { return locations[i].Importance > locations[j].Importance })

	return locations, nil
}
