// SPDX-FileCopyrightText: 2023 Winni Neessen <wn@neessen.dev>
//
// SPDX-License-Identifier: MIT

// Package meteologix provides bindings to the Meteologix/Kachelmann-Wetter weather API
package meteologix

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// AstronomicalInfo provides astronomical data for the next 14 days.
// This includes moon and sun information.
type AstronomicalInfo struct {
	// DailyData holds the different APIAstronomicalDailyData data
	// points for the next 14 days
	DailyData []APIAstronomicalDailyData `json:"dailyData"`
	// Latitude represents the GeoLocation latitude coordinates for the weather data
	Latitude float64 `json:"lat"`
	// Longitude represents the GeoLocation longitude coordinates for the weather data
	Longitude float64 `json:"lon"`
	// NextFullMoon represent the date and time of the next full moon
	NextFullMoon time.Time `json:"nextFullMoon"`
	// NextNewMoon represent the date and time of the next new moon
	NextNewMoon time.Time `json:"nextNewMoon"`
	// Run represents when astronomical values have been calculated
	Run time.Time `json:"run"`
	// TimeZone is the timezone at the queried location
	TimeZone string `json:"timeZone"`
}

// APIAstronomicalDailyData holds the API response date for the daily
// details in the AstronomicalInfo.
type APIAstronomicalDailyData struct {
	// AstronomicalDawn represents the date and time when civil dawn begins
	AstronomicalDawn *time.Time `json:"astronomicalDawn,omitempty"`
	// AstronomicalDusk represents the date and time when civil dusk ends
	AstronomicalDusk *time.Time `json:"astronomicalDusk,omitempty"`
	// CivilDawn represents the date and time when civil dawn begins
	CivilDawn *time.Time `json:"civilDawn,omitempty"`
	// CivilDusk represents the date and time when civil dusk ends
	CivilDusk *time.Time `json:"civilDusk,omitempty"`
	// DateTime represents the date for the forecast values
	DateTime APIDate `json:"dateTime"`
	// MoonIllumination represents how much of the moon is illuminated in %
	MoonIllumination float64 `json:"moonIllumination"`
	// MoonPhase represents the moon phase in %
	MoonPhase int `json:"moonPhase"`
	// MoonRise represents the date and time when the moon rises
	MoonRise *time.Time `json:"moonRise,omitempty"`
	// MoonSet represents the date and time when the moon sets
	MoonSet *time.Time `json:"moonSet,omitempty"`
	// NauticalDawn represents the date and time when nautical dawn begins
	NauticalDawn *time.Time `json:"nauticalDawn,omitempty"`
	// NauticalDusk represents the date and time when nautical dusk ends
	NauticalDusk *time.Time `json:"nauticalDusk,omitempty"`
	// Sunrise represents the date and time of the sunrise
	Sunrise *time.Time `json:"sunrise,omitempty"`
	// Sunset represents the date and time of the sunset
	Sunset *time.Time `json:"sunset,omitempty"`
	// Transit represents the date and time when the sun is at
	// its zenith
	Transit *time.Time `json:"transit,omitempty"`
}

// AstronomicalInfoByCoordinates returns the AstronomicalInfo values for
// the given coordinates
func (c *Client) AstronomicalInfoByCoordinates(la, lo float64) (AstronomicalInfo, error) {
	var ai AstronomicalInfo
	lat := strconv.FormatFloat(la, 'f', -1, 64)
	lon := strconv.FormatFloat(lo, 'f', -1, 64)
	u := fmt.Sprintf("%s/tools/astronomy/%s/%s", c.config.apiURL, lat, lon)

	r, err := c.httpClient.Get(u)
	if err != nil {
		return ai, fmt.Errorf("API request failed: %w", err)
	}

	if err := json.Unmarshal(r, &ai); err != nil {
		return ai, fmt.Errorf("failed to unmarshal API response JSON: %w", err)
	}

	return ai, nil
}

// AstronomicalInfoByLocation returns the AstronomicalInfo values for
// the given location
func (c *Client) AstronomicalInfoByLocation(lo string) (AstronomicalInfo, error) {
	gl, err := c.GetGeoLocationByName(lo)
	if err != nil {
		return AstronomicalInfo{}, fmt.Errorf("failed too look up geolocation: %w", err)
	}
	return c.AstronomicalInfoByCoordinates(gl.Latitude, gl.Longitude)
}

// Sunset returns the date and time of the sunset on the current day
// as DateTime type
// If the data point is not available in the AstronomicalInfo it will
// return DateTime in which the "not available" field will be true.
func (a *AstronomicalInfo) Sunset() DateTime {
	n := time.Now()
	if len(a.DailyData) < 1 {
		return DateTime{na: true}
	}
	cdd := a.DailyData[0]
	if cdd.DateTime.Format("2006-01-02") != n.Format("2006-01-02") {
		return DateTime{na: true}
	}
	return DateTime{
		dt: a.Run,
		n:  FieldSunset,
		s:  SourceForecast,
		dv: *cdd.Sunset,
	}
}
