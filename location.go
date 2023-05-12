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

// GetGeoLocationByCity returns the GeoLocation with the highest importance based on
// the given City name
//
// This method makes use of the OSM Nominatim API
func (c *Client) GetGeoLocationByCity(ci string) (GeoLocation, error) {
	ga, err := c.GetGeoLocationsByCity(ci)
	return ga[0], err
}

// GetGeoLocationsByCity returns a slice of GeoLocation based on the requested City name
// The returned slice will be sorted by Importance of the results with the highest
// importance as first entry
//
// This method makes use of the OSM Nominatim API
func (c *Client) GetGeoLocationsByCity(ci string) ([]GeoLocation, error) {
	ga := make([]GeoLocation, 0)

	u, err := url.Parse(OSMNominatimURL)
	if err != nil {
		return ga, fmt.Errorf("failed to parse OSM Nominatim URL: %w", err)
	}
	uq := u.Query()
	uq.Add("format", "json")
	uq.Add("q", ci)
	u.RawQuery = uq.Encode()

	r, err := c.hc.Get(u.String())
	if err != nil {
		return ga, fmt.Errorf("OSM Nominatim API request failed: %w", err)
	}
	var la []GeoLocation
	if err := json.Unmarshal(r, &la); err != nil {
		return ga, fmt.Errorf("failed to unmarshal API response JSON: %w", err)
	}
	if len(la) < 1 {
		return ga, ErrCityNotFound
	}

	for _, l := range la {
		lat, err := strconv.ParseFloat(l.LatitudeString, 64)
		if err != nil {
			return ga, fmt.Errorf("failed to convert latitude string to float value: %w", err)
		}
		lon, err := strconv.ParseFloat(l.LongitudeString, 64)
		if err != nil {
			return ga, fmt.Errorf("failed to convert longitude string to float value: %w", err)
		}
		l.Latitude = lat
		l.Longitude = lon
		ga = append(ga, l)
	}
	sort.SliceStable(ga, func(i, j int) bool { return ga[i].Importance > ga[j].Importance })

	return ga, nil
}
