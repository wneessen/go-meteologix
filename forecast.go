// SPDX-FileCopyrightText: 2023 Winni Neessen <wn@neessen.dev>
//
// SPDX-License-Identifier: MIT

package meteologix

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

const (
	// ForecastSteps1h represents a weather forecast in an 1 hour time interval. It will return
	// up to 24 hours of forecast data.
	ForecastSteps1h ForecastTimeSteps = "1h"
	// ForecastSteps3h represents a weather forecast in a 3 hour time interval. It will return
	// up to 120 hours of forecast data.
	ForecastSteps3h ForecastTimeSteps = "3h"
	// ForecastSteps6h represents a weather forecast in a 6 hour time interval. It will return
	// up to 240 hours of forecast data.
	ForecastSteps6h ForecastTimeSteps = "6h"
)

const (
	// ForecastDetailStandard represents a standard level of detail for weather forecasts retrieved from the API.
	ForecastDetailStandard ForecastDetails = "standard"
	// ForecastDetailAdvanced represents an advanced level of detail for weather forecasts retrieved from the API.
	ForecastDetailAdvanced ForecastDetails = "advanced"
)

// WeatherForecast represents the weather forecast API response
type WeatherForecast struct {
	// Altitude represents the altitude of the location that has been queried
	Altitude int `json:"alt"`
	// Data holds the different APICurrentWeatherData points
	Data []APIWeatherForecastData `json:"data"`
	// Latitude represents the GeoLocation latitude coordinates for the weather data
	Latitude float64 `json:"lat"`
	// Longitude represents the GeoLocation longitude coordinates for the weather data
	Longitude float64 `json:"lon"`
	// Precision represents the weather models resolution
	Precision Precision `json:"resolution"`
	// Run represents the time when the weather forecast was generated.
	Run time.Time `json:"run"`
	// Timezone represents the timezone at the location
	Timezone string `json:"timeZone"`
	// UnitSystem is the unit system that is used for the results (we default to metric)
	UnitSystem string `json:"systemOfUnits"`
}

// ForecastTimeSteps represents a time step used in a weather forecast. It is an alias type for a string type
type ForecastTimeSteps string

// ForecastDetails represents a type of detail for weather forecasts retrieved from the API
type ForecastDetails string

// APIWeatherForecastData holds the different data points of the WeatherForecast as returned by the
// weather forecast API endpoints.
type APIWeatherForecastData struct {
	// DateTime represents the date and time for the forecast values
	DateTime time.Time `json:"dateTime"`
	// IsDay is true when it is date and time of forecast is at daytime
	IsDay bool `json:"isDay"`
	// Dewpoint represents the predicted dewpoint (at current timestamp)
	Dewpoint NilFloat64 `json:"dewpoint,omitempty"`
	// PressureMSL represents barometric air pressure at mean sea level (at current timestamp)
	PressureMSL NilFloat64 `json:"pressureMsl,omitempty"`
	// Temperature represents the predicted temperature at 2m height (at current timestamp)
	Temperature float64 `json:"temp"`
}

type WeatherForecastDatapoint struct {
	dateTime    time.Time
	isDay       bool
	dewpoint    NilFloat64
	pressureMSL NilFloat64
	temperature float64
}

// ForecastByCoordinates returns the WeatherForecast values for the given coordinates
func (c *Client) ForecastByCoordinates(latitude, longitude float64, timesteps ForecastTimeSteps,
	details ForecastDetails) (WeatherForecast, error) {
	var forecast WeatherForecast
	latitudeFormat := strconv.FormatFloat(latitude, 'f', -1, 64)
	longitudeFormat := strconv.FormatFloat(longitude, 'f', -1, 64)
	apiURL, err := url.Parse(fmt.Sprintf("%s/forecast/%s/%s/%s/%s", c.config.apiURL, latitudeFormat,
		longitudeFormat, details, timesteps))
	if err != nil {
		return forecast, fmt.Errorf("failed to parse weather forecast URL: %w", err)
	}
	queryString := apiURL.Query()
	queryString.Add("units", "metric")
	apiURL.RawQuery = queryString.Encode()

	response, err := c.httpClient.Get(apiURL.String())
	if err != nil {
		return forecast, fmt.Errorf("API request failed: %w", err)
	}

	if err = json.Unmarshal(response, &forecast); err != nil {
		return forecast, fmt.Errorf("failed to unmarshal API response JSON: %w", err)
	}

	return forecast, nil
}

// ForecastByLocation returns the WeatherForecast values for the given location
func (c *Client) ForecastByLocation(location string, timesteps ForecastTimeSteps,
	details ForecastDetails) (WeatherForecast, error) {
	geoLocation, err := c.GetGeoLocationByName(location)
	if err != nil {
		return WeatherForecast{}, fmt.Errorf("failed too look up geolocation: %w", err)
	}
	return c.ForecastByCoordinates(geoLocation.Latitude, geoLocation.Longitude, timesteps, details)
}

func (wf WeatherForecast) At(timestamp time.Time) WeatherForecastDatapoint {
	datapoint := findClosestForecast(wf.Data, timestamp)
	if datapoint == nil {
		return WeatherForecastDatapoint{}
	}
	return newWeatherForecastDataPoint(*datapoint)
}

func (wf WeatherForecast) All() []WeatherForecastDatapoint {
	datapoints := make([]WeatherForecastDatapoint, 0)
	for _, data := range wf.Data {
		datapoint := newWeatherForecastDataPoint(data)
		datapoints = append(datapoints, datapoint)
	}
	return datapoints
}

func newWeatherForecastDataPoint(data APIWeatherForecastData) WeatherForecastDatapoint {
	return WeatherForecastDatapoint{
		dateTime:    data.DateTime,
		isDay:       data.IsDay,
		dewpoint:    data.Dewpoint,
		pressureMSL: data.PressureMSL,
		temperature: data.Temperature,
	}
}

func (dp WeatherForecastDatapoint) DateTime() time.Time {
	return dp.dateTime
}

func (dp WeatherForecastDatapoint) Dewpoint() Temperature {
	if dp.dewpoint.IsNil() {
		return Temperature{notAvailable: true}
	}
	temperature := Temperature{
		dateTime: dp.dateTime,
		name:     FieldDewpoint,
		source:   SourceForecast,
		floatVal: dp.dewpoint.Get(),
	}
	return temperature
}

// PressureMSL returns the pressure at mean sea level data point as Pressure.
//
// If the data point is not available in the WeatherForecast it will return Pressure in which the
// "not available" field will be true.
func (dp WeatherForecastDatapoint) PressureMSL() Pressure {
	if dp.pressureMSL.IsNil() {
		return Pressure{notAvailable: true}
	}
	pressure := Pressure{
		dateTime: dp.dateTime,
		name:     FieldPressureMSL,
		source:   SourceForecast,
		floatVal: dp.pressureMSL.Get(),
	}
	return pressure
}

// Temperature returns the temperature data point as Temperature.
func (dp WeatherForecastDatapoint) Temperature() Temperature {
	return Temperature{
		dateTime: dp.DateTime(),
		name:     FieldTemperature,
		source:   SourceForecast,
		floatVal: dp.temperature,
	}
}

func findClosestForecast(items []APIWeatherForecastData, target time.Time) *APIWeatherForecastData {
	if len(items) <= 0 {
		return nil
	}

	closest := items[0]
	minDiff := target.Sub(closest.DateTime).Abs()

	for _, item := range items[1:] {
		diff := target.Sub(item.DateTime).Abs()
		if diff < minDiff {
			minDiff = diff
			closest = item
		}
	}

	return &closest
}
