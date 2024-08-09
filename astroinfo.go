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

// AstronomicalInfoByCoordinates returns the AstronomicalInfo values for the given coordinates
func (c *Client) AstronomicalInfoByCoordinates(latitude, longitude float64) (AstronomicalInfo, error) {
	var astroInfo AstronomicalInfo
	latitudeFormat := strconv.FormatFloat(latitude, 'f', -1, 64)
	longitudeFormat := strconv.FormatFloat(longitude, 'f', -1, 64)
	apiURL := fmt.Sprintf("%s/tools/astronomy/%s/%s", c.config.apiURL, latitudeFormat, longitudeFormat)

	response, err := c.httpClient.Get(apiURL)
	if err != nil {
		return astroInfo, fmt.Errorf("API request failed: %w", err)
	}

	if err = json.Unmarshal(response, &astroInfo); err != nil {
		return astroInfo, fmt.Errorf("failed to unmarshal API response JSON: %w", err)
	}

	return astroInfo, nil
}

// AstronomicalInfoByLocation returns the AstronomicalInfo values for the given location
func (c *Client) AstronomicalInfoByLocation(location string) (AstronomicalInfo, error) {
	geoLocation, err := c.GetGeoLocationByName(location)
	if err != nil {
		return AstronomicalInfo{}, fmt.Errorf("failed too look up geolocation: %w", err)
	}
	return c.AstronomicalInfoByCoordinates(geoLocation.Latitude, geoLocation.Longitude)
}

// SunsetByTime returns the date and time of the sunset on the given time as DateTime type.
//
// If the data point is not available in the AstronomicalInfo it will return DateTime in
// which the "not available" field will be true.
//
// Please keep in mind that the API only returns 14 days in the future. Any date given
// that exceeds that time, wil always return a "not available" value.
func (a *AstronomicalInfo) SunsetByTime(timeVal time.Time) DateTime {
	if len(a.DailyData) < 1 {
		return DateTime{notAvailable: true}
	}
	var currentDayData APIAstronomicalDailyData
	for i := range a.DailyData {
		if a.DailyData[i].DateTime.Format(DateFormat) != timeVal.Format(DateFormat) {
			continue
		}
		currentDayData = a.DailyData[i]
	}
	if currentDayData.DateTime.IsZero() {
		return DateTime{notAvailable: true}
	}
	return DateTime{
		dateTime:    currentDayData.DateTime.Time,
		name:        FieldSunset,
		source:      SourceForecast,
		dateTimeVal: *currentDayData.Sunset,
	}
}

// Sunset returns the date and time of the sunset on the current date as DateTime type.
//
// If the data point is not available in the AstronomicalInfo it will return DateTime
// in which the "not available" field will be true.
func (a *AstronomicalInfo) Sunset() DateTime {
	return a.SunsetByTime(time.Now())
}

// SunsetByDateString returns the date and time of the sunset at a given date string as
// DateTime type. Expected Go format template is: 2006-01-02.
//
// If the date wasn't able to be parsed or if the data point is not available in the
// AstronomicalInfo it will return DateTime in which the "not available" field will be
// true.
func (a *AstronomicalInfo) SunsetByDateString(date string) DateTime {
	timeVal, err := time.Parse(DateFormat, date)
	if err != nil {
		return DateTime{notAvailable: true}
	}
	return a.SunsetByTime(timeVal)
}

// SunsetAll returns a slice of all sunset data points in the given AstronomicalInfo instance
// as DateTime types.
//
// If no sunset data is available an empty slice is returned
func (a *AstronomicalInfo) SunsetAll() []DateTime {
	var sunsets []DateTime
	for _, dayData := range a.DailyData {
		if dayData.DateTime.IsZero() {
			continue
		}
		sunsets = append(sunsets, a.SunsetByTime(dayData.DateTime.Time))
	}

	return sunsets
}

// SunriseByTime returns the date and time of the sunrise on the give  time as DateTime type.
//
// If the data point is not available in the AstronomicalInfo it will return DateTime in
// which the "not available" field will be true.
//
// Please keep in mind that the API only returns 14 days in the future. Any date given that
// exceeds that time, wil always return a "not available" value.
func (a *AstronomicalInfo) SunriseByTime(timeVal time.Time) DateTime {
	if len(a.DailyData) < 1 {
		return DateTime{notAvailable: true}
	}
	var currentDayData APIAstronomicalDailyData
	for i := range a.DailyData {
		if a.DailyData[i].DateTime.Format(DateFormat) != timeVal.Format(DateFormat) {
			continue
		}
		currentDayData = a.DailyData[i]
	}
	if currentDayData.DateTime.IsZero() {
		return DateTime{notAvailable: true}
	}
	return DateTime{
		dateTime:    currentDayData.DateTime.Time,
		name:        FieldSunrise,
		source:      SourceForecast,
		dateTimeVal: *currentDayData.Sunrise,
	}
}

// Sunrise returns the date and time of the sunrise on the current date as DateTime type.
//
// If the data point is not available in the AstronomicalInfo it will return DateTime in
// which the "not available" field will be true.
func (a *AstronomicalInfo) Sunrise() DateTime {
	return a.SunriseByTime(time.Now())
}

// SunriseByDateString returns the date and time of the sunrise at a given date string as
// DateTime type. Expected Go format template is: 2006-01-02.
//
// If the date wasn't able to be parsed or if the data point is not available in the
// AstronomicalInfo it will return DateTime in which the "not available" field will be
// true.
func (a *AstronomicalInfo) SunriseByDateString(date string) DateTime {
	timeVal, err := time.Parse(DateFormat, date)
	if err != nil {
		return DateTime{notAvailable: true}
	}
	return a.SunriseByTime(timeVal)
}

// SunriseAll returns a slice of all sunrise data points in the given AstronomicalInfo instance
// as DateTime types.
//
// If no sunrise data is available it will return an empty slice
func (a *AstronomicalInfo) SunriseAll() []DateTime {
	var sunrises []DateTime
	for _, dayData := range a.DailyData {
		if dayData.DateTime.IsZero() {
			continue
		}
		sunrises = append(sunrises, a.SunriseByTime(dayData.DateTime.Time))
	}

	return sunrises
}
