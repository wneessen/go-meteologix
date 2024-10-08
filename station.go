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
	"strings"
)

// DefaultRadius is the default radius value that is used in the station search
const DefaultRadius int = 10

const (
	// PrecisionSuperHigh represents the precision level of data corresponding
	// to a resolution of less than or approximately equal to 4 kilometers.
	// This is the highest level of precision, usually associated with highly
	// detailed measurements or observations.
	PrecisionSuperHigh Precision = iota
	// PrecisionHigh represents the precision level of data corresponding to a
	// resolution between 4 kilometers and 10 kilometers. This is a high precision
	// level, suitable for most operational needs that require a balance between
	// detail and processing requirements.
	PrecisionHigh
	// PrecisionStandard represents the precision level of data corresponding to
	// a resolution of 10 kilometers or more. This is the standard level of
	// precision, generally used for large-scale analysis and modeling.
	PrecisionStandard
	// PrecisionUnknown is used when the precision level of a weather station
	// is unknown. This constant can be used as a placeholder when the resolution
	// data is not available.
	PrecisionUnknown
)

// Precision levels defined as strings to allow for clear, consistent use throughout the
// application.
const (
	// PrecisionStringSuperHigh represents the super high precision level string.
	PrecisionStringSuperHigh = "SUPER_HIGH"

	// PrecisionStringHigh represents the high precision level string.
	PrecisionStringHigh = "HIGH"

	// PrecisionStringStandard represents the standard precision level string.
	PrecisionStringStandard = "STANDARD"

	// PrecisionStringUnknown represents an unknown precision level string.
	PrecisionStringUnknown = "UNKNOWN"
)

var (
	// ErrRadiusTooSmall is returned if a given radius value is too small
	ErrRadiusTooSmall = errors.New("given radius is too small")
	// ErrNoStationFound is returned if a station search did not return any results
	ErrNoStationFound = errors.New("no station found in requested location")
)

// Station is a weather station as returned by the Meteologix API
type Station struct {
	// Altitude is the altitude of the station
	Altitude int `json:"alt"`
	// Distance is the distatnce of the station to the provided coordinates
	Distance float64 `json:"distance"`
	// ID is the station ID
	ID string `json:"id"`
	// Latitude is the latitude of the station
	Latitude float64 `json:"lat"`
	// Longitude is the latitude of the station
	Longitude float64 `json:"lon"`
	// Name is the name or location of the station
	Name string `json:"name"`
	// Precision is the precision string returned by the API
	Precision *Precision `json:"precision,omitempty"`
	// RecentlyActive represents if the station was recently active
	RecentlyActive bool `json:"recentlyActive"`
	// Type is the type of weather station
	Type *string `json:"type,omitempty"`
}

// Precision is a type wrapper for an int type
type Precision int

// StationSearchByCoordinates returns a list of available weather stations
// based on the given latitude, longitude coordinates within the default
// radius
//
// Results will be sorted by distance to the requested coordinates.
//
// Depending on your subscription you may have access to one, two or
// unlimited locations for station observations.
// Finding a station with his endpoint does not automatically mean
// that you are allowed to get all data from this station.
//
// See: https://api.kachelmannwetter.com/v02/_doc.html#/operations/get_station_search
func (c *Client) StationSearchByCoordinates(latitude, longitude float64) ([]Station, error) {
	return c.StationSearchByCoordinatesWithinRadius(latitude, longitude, DefaultRadius)
}

// StationSearchByLocation returns a list of available weather stations
// based on the given location string within the default radius
//
// # Results will be sorted by distance to the requested location
//
// Depending on your subscription you may have access to one, two or
// unlimited locations for station observations.
// Finding a station with his endpoint does not automatically mean
// that you are allowed to get all data from this station.
//
// See: https://api.kachelmannwetter.com/v02/_doc.html#/operations/get_station_search
func (c *Client) StationSearchByLocation(location string) ([]Station, error) {
	return c.StationSearchByLocationWithinRadius(location, DefaultRadius)
}

// StationSearchByLocationWithinRadius returns a list of available weather
// stations based on the given location string and radius.
//
// Results will be sorted by distance to the requested location.
//
// Depending on your subscription you may have access to one, two or
// unlimited locations for station observations.
// Finding a station with his endpoint does not automatically mean
// that you are allowed to get all data from this station.
//
// See: https://api.kachelmannwetter.com/v02/_doc.html#/operations/get_station_search
func (c *Client) StationSearchByLocationWithinRadius(location string, radius int) ([]Station, error) {
	geoLocation, err := c.GetGeoLocationByName(location)
	if err != nil {
		return nil, fmt.Errorf("failed too look up location details: %w", err)
	}
	return c.StationSearchByCoordinatesWithinRadius(geoLocation.Latitude, geoLocation.Longitude, radius)
}

// StationSearchByCoordinatesWithinRadius returns a list of available weather stations
// based on the given latitude, longitude coordinates and radius.
//
// Results will be sorted by distance to the requested coordinates.
//
// Depending on your subscription you may have access to one, two or
// unlimited locations for station observations.
// Finding a station with his endpoint does not automatically mean
// that you are allowed to get all data from this station.
//
// See: https://api.kachelmannwetter.com/v02/_doc.html#/operations/get_station_search
func (c *Client) StationSearchByCoordinatesWithinRadius(latitude, longitude float64, radius int) ([]Station, error) {
	if radius < 1 {
		return nil, ErrRadiusTooSmall
	}

	apiURL, err := url.Parse(fmt.Sprintf("%s/station/search/%f/%f",
		c.config.apiURL, latitude, longitude))
	if err != nil {
		return nil, fmt.Errorf("failed to parse station search URL: %w", err)
	}
	query := apiURL.Query()
	query.Add("radius", fmt.Sprintf("%d", radius))
	apiURL.RawQuery = query.Encode()

	response, err := c.httpClient.Get(apiURL.String())
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	var stations []Station
	if err = json.Unmarshal(response, &stations); err != nil {
		return nil, fmt.Errorf("failed to unmarshal API response JSON: %w", err)
	}
	if len(stations) < 1 {
		return nil, ErrNoStationFound
	}
	sort.SliceStable(stations, func(i, j int) bool { return stations[i].Distance < stations[j].Distance })

	return stations, nil
}

// UnmarshalJSON method for converting API precision responses into
// StationPrecision types
func (p *Precision) UnmarshalJSON(data []byte) error {
	value := string(data)
	value = strings.ReplaceAll(value, `"`, ``)
	switch strings.ToUpper(value) {
	case PrecisionStringSuperHigh:
		*p = PrecisionSuperHigh
	case PrecisionStringHigh:
		*p = PrecisionHigh
	case PrecisionStringStandard:
		*p = PrecisionStandard
	default:
		*p = PrecisionUnknown
	}
	return nil
}

// String satisfies the fmt.Stringer interface for the Precision type
func (p *Precision) String() string {
	switch *p {
	case PrecisionSuperHigh:
		return PrecisionStringSuperHigh
	case PrecisionHigh:
		return PrecisionStringHigh
	case PrecisionStandard:
		return PrecisionStringStandard
	case PrecisionUnknown:
		return PrecisionStringUnknown
	default:
		return PrecisionStringUnknown
	}
}
