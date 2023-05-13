// SPDX-FileCopyrightText: 2023 Winni Neessen <winni@neessen.dev>
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

// DefaultRadius is the default radius value that is used in the
// station search
const DefaultRadius int = 10

const (
	PrecisionHigh Precision = iota
	PrecisionMedium
	PrecisionLow
	PrecisionUnknown
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
	Precision StationPrecision `json:"precision"`
	// RecentlyActive represents if the station was recently active
	RecentlyActive bool `json:"recentlyActive"`
	// Type is the type of weather station
	Type string `json:"type"`
}

// Precision is a type wrapper for an int type
type Precision int

// StationPrecision is a type wrapper for an int type
type StationPrecision struct {
	Precision
}

// StationSearch returns a list of available weather stations based on the
// given Latitude, Longitude
//
// Results will be sorted by distance to the requested coordinates
// given Latitude, Longitude
//
// Depending on your subscription you may have access to one, two or
// unlimited locations for station observations.
// Finding a station with his endpoint does not automatically mean
// that you are allowed to get all data from this station.
//
// See: https://api.kachelmannwetter.com/v02/_doc.html#/operations/get_station_search
func (c *Client) StationSearch(la, lo float64) ([]Station, error) {
	return c.StationSearchWithRadius(la, lo, DefaultRadius)
}

// StationSearchByCity returns a list of available weather stations based
// on the given City name
//
// # Results will be sorted by distance to the requested coordinates given City
//
// Depending on your subscription you may have access to one, two or
// unlimited locations for station observations.
// Finding a station with his endpoint does not automatically mean
// that you are allowed to get all data from this station.
//
// See: https://api.kachelmannwetter.com/v02/_doc.html#/operations/get_station_search
func (c *Client) StationSearchByCity(ci string) ([]Station, error) {
	l, err := c.GetGeoLocationByCityName(ci)
	if err != nil {
		return nil, fmt.Errorf("failed too look up city details: %w", err)
	}
	return c.StationSearchWithRadius(l.Latitude, l.Longitude, DefaultRadius)
}

// StationSearchWithRadius returns a list of available weather stations based on the
// given Latitude, Longitude and Radius values
//
// # Results will be sorted by distance to the requested coordinates
//
// Depending on your subscription you may have access to one, two or
// unlimited locations for station observations.
// Finding a station with his endpoint does not automatically mean
// that you are allowed to get all data from this station.
//
// See: https://api.kachelmannwetter.com/v02/_doc.html#/operations/get_station_search
func (c *Client) StationSearchWithRadius(la, lo float64, ra int) ([]Station, error) {
	if ra < 1 {
		return nil, ErrRadiusTooSmall
	}

	u, err := url.Parse(fmt.Sprintf("%s/station/search/%f/%f",
		APIBaseURL, la, lo))
	if err != nil {
		return nil, fmt.Errorf("failed to parse station search URL: %w", err)
	}
	uq := u.Query()
	uq.Add("radius", fmt.Sprintf("%d", ra))
	u.RawQuery = uq.Encode()

	r, err := c.httpClient.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	var sl []Station
	if err := json.Unmarshal(r, &sl); err != nil {
		return nil, fmt.Errorf("failed to unmarshal API response JSON: %w", err)
	}
	if len(sl) < 1 {
		return nil, ErrNoStationFound
	}
	sort.SliceStable(sl, func(i, j int) bool { return sl[i].Distance < sl[j].Distance })

	return sl, nil
}

// UnmarshalJSON method for converting API precision responses into
// StationPrecision types
func (p *StationPrecision) UnmarshalJSON(s []byte) error {
	v := string(s)
	v = strings.ReplaceAll(v, `"`, ``)
	switch strings.ToLower(v) {
	case "high":
		p.Precision = PrecisionHigh
	case "medium":
		p.Precision = PrecisionMedium
	case "low":
		p.Precision = PrecisionLow
	default:
		p.Precision = PrecisionUnknown
	}
	return nil
}

// String satisfies the fmt.Stringer interface for the Precision type
func (p Precision) String() string {
	switch p {
	case PrecisionHigh:
		return "HIGH"
	case PrecisionMedium:
		return "MEDIUM"
	case PrecisionLow:
		return "LOW"
	case PrecisionUnknown:
		return "UNKNOWN"
	default:
		return "UNKNOWN"
	}
}
