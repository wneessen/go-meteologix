// SPDX-FileCopyrightText: 2023 Winni Neessen <wn@neessen.dev>
//
// SPDX-License-Identifier: MIT

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

// SunsetByTime returns the date and time of the sunset on the give
// time as DateTime type.
// If the data point is not available in the AstronomicalInfo it will
// return DateTime in which the "not available" field will be true.
//
// Please keep in mind that the API only returns 14 days in the future.
// Any date given that exceeds that time, wil always return a
// "not available" value.
func (a *AstronomicalInfo) SunsetByTime(t time.Time) DateTime {
	if len(a.DailyData) < 1 {
		return DateTime{notAvailable: true}
	}
	var cdd APIAstronomicalDailyData
	for i := range a.DailyData {
		if a.DailyData[i].DateTime.Format(DateFormat) != t.Format(DateFormat) {
			continue
		}
		cdd = a.DailyData[i]
	}
	if cdd.DateTime.IsZero() {
		return DateTime{notAvailable: true}
	}
	return DateTime{
		dateTime:    cdd.DateTime.Time,
		name:        FieldSunset,
		source:      SourceForecast,
		dateTimeVal: *cdd.Sunset,
	}
}

// Sunset returns the date and time of the sunset on the current date
// as DateTime type.
// If the data point is not available in the AstronomicalInfo it will
// return DateTime in which the "not available" field will be true.
func (a *AstronomicalInfo) Sunset() DateTime {
	return a.SunsetByTime(time.Now())
}

// SunsetByDateString returns the date and time of the sunset at a
// given date string as DateTime type. Expected format is 2006-01-02.
// If the date wasn't able to be parsed or if the data point is not
// available in the AstronomicalInfo it will return DateTime in
// which the "not available" field will be true.
func (a *AstronomicalInfo) SunsetByDateString(ds string) DateTime {
	t, err := time.Parse(DateFormat, ds)
	if err != nil {
		return DateTime{notAvailable: true}
	}
	return a.SunsetByTime(t)
}

// SunsetAll returns a slice of all sunset data points in the given
// AstronomicalInfo instance as DateTime types. If no sunset data
// is available it will return an empty slice
func (a *AstronomicalInfo) SunsetAll() []DateTime {
	var sss []DateTime
	for _, cd := range a.DailyData {
		if cd.DateTime.IsZero() {
			continue
		}
		sss = append(sss, a.SunsetByTime(cd.DateTime.Time))
	}

	return sss
}

// SunriseByTime returns the date and time of the sunrise on the give
// time as DateTime type.
// If the data point is not available in the AstronomicalInfo it will
// return DateTime in which the "not available" field will be true.
//
// Please keep in mind that the API only returns 14 days in the future.
// Any date given that exceeds that time, wil always return a
// "not available" value.
func (a *AstronomicalInfo) SunriseByTime(t time.Time) DateTime {
	if len(a.DailyData) < 1 {
		return DateTime{notAvailable: true}
	}
	var cdd APIAstronomicalDailyData
	for i := range a.DailyData {
		if a.DailyData[i].DateTime.Format(DateFormat) != t.Format(DateFormat) {
			continue
		}
		cdd = a.DailyData[i]
	}
	if cdd.DateTime.IsZero() {
		return DateTime{notAvailable: true}
	}
	return DateTime{
		dateTime:    cdd.DateTime.Time,
		name:        FieldSunrise,
		source:      SourceForecast,
		dateTimeVal: *cdd.Sunrise,
	}
}

// Sunrise returns the date and time of the sunrise on the current date
// as DateTime type.
// If the data point is not available in the AstronomicalInfo it will
// return DateTime in which the "not available" field will be true.
func (a *AstronomicalInfo) Sunrise() DateTime {
	return a.SunriseByTime(time.Now())
}

// SunriseByDateString returns the date and time of the sunrise at a
// given date string as DateTime type. Expected format is 2006-01-02.
// If the date wasn't able to be parsed or if the data point is not
// available in the AstronomicalInfo it will return DateTime in
// which the "not available" field will be true.
func (a *AstronomicalInfo) SunriseByDateString(ds string) DateTime {
	t, err := time.Parse(DateFormat, ds)
	if err != nil {
		return DateTime{notAvailable: true}
	}
	return a.SunriseByTime(t)
}

// SunriseAll returns a slice of all sunrise data points in the given
// AstronomicalInfo instance as DateTime types. If no sunrise data
// is available it will return an empty slice
func (a *AstronomicalInfo) SunriseAll() []DateTime {
	var sss []DateTime
	for _, cd := range a.DailyData {
		if cd.DateTime.IsZero() {
			continue
		}
		sss = append(sss, a.SunriseByTime(cd.DateTime.Time))
	}

	return sss
}
