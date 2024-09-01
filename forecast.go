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
	// CloudCoverage represents the effective cloud coverage within the preceding timespan
	// in % (e.g. low clouds have more priority than high clouds)
	CloudCoverage NilFloat64 `json:"cloudCoverage,omitempty"`
	// DateTime represents the date and time for the forecast values
	DateTime time.Time `json:"dateTime"`
	// Humidity represents the relative humidity value of a weather forecast
	Humidity NilFloat64 `json:"humidityRelative"`
	// IsDay is true when it is date and time of forecast is at daytime
	IsDay bool `json:"isDay"`
	// Dewpoint represents the predicted dewpoint (at current timestamp)
	Dewpoint NilFloat64 `json:"dewpoint,omitempty"`
	// PressureMSL represents barometric air pressure at mean sea level (at current timestamp)
	PressureMSL NilFloat64 `json:"pressureMsl,omitempty"`
	// SunHours represents the most probable amount of hours the sun will be visible
	SunHours NilFloat64 `json:"sunHours,omitempty"`
	// Temperature represents the predicted temperature at 2m height (at current timestamp)
	Temperature float64 `json:"temp"`
	// WeatherSymbol is a text representation of the current weather conditions
	WeatherSymbol NilString `json:"weatherSymbol,omitempty"`
	// WindDirection represents the average direction from which the wind originates in degree
	WindDirection NilFloat64 `json:"windDirection,omitempty"`
	// WindGust represents the wind gust speed in m/s (for a timespan)
	WindGust NilFloat64 `json:"windGust,omitempty"`
	// WindGust3h represents the wind gust speed in m/s over the last 3 hours
	WindGust3h NilFloat64 `json:"windGust3h,omitempty"`
	// WindSpeed represents the average wind speed (for a timespan) in m/s
	WindSpeed NilFloat64 `json:"windspeed,omitempty"`
}

// WeatherForecastDatapoint represents a single data point in a weather forecast.
type WeatherForecastDatapoint struct {
	cloudCoverage NilFloat64
	dateTime      time.Time
	dewpoint      NilFloat64
	humidity      NilFloat64
	isDay         bool
	pressureMSL   NilFloat64
	sunhours      NilFloat64
	temperature   float64
	weatherSymbol NilString
	winddirection NilFloat64
	windgust      NilFloat64
	windgust3h    NilFloat64
	windspeed     NilFloat64
}

// ForecastByCoordinates returns the WeatherForecast values for the given coordinates
func (c *Client) ForecastByCoordinates(latitude, longitude float64, timespan Timespan,
	details ForecastDetails,
) (WeatherForecast, error) {
	var forecast WeatherForecast
	var steps string
	switch timespan {
	case Timespan1Hour, Timespan3Hours, Timespan6Hours:
		steps = timespan.String()
	default:
		return forecast, fmt.Errorf("unsupported timespan for weather forecasts: %s", timespan)
	}

	latitudeFormat := strconv.FormatFloat(latitude, 'f', -1, 64)
	longitudeFormat := strconv.FormatFloat(longitude, 'f', -1, 64)
	apiURL, err := url.Parse(fmt.Sprintf("%s/forecast/%s/%s/%s/%s", c.config.apiURL, latitudeFormat,
		longitudeFormat, details, steps))
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
func (c *Client) ForecastByLocation(location string, timesteps Timespan,
	details ForecastDetails,
) (WeatherForecast, error) {
	geoLocation, err := c.GetGeoLocationByName(location)
	if err != nil {
		return WeatherForecast{}, fmt.Errorf("failed too look up geolocation: %w", err)
	}
	return c.ForecastByCoordinates(geoLocation.Latitude, geoLocation.Longitude, timesteps, details)
}

// At returns the WeatherForecastDatapoint for the specified timestamp. It will try to find the closest datapoint
// in the forecast that matches the given timestamp. If no matching datapoint is found, an empty
// WeatherForecastDatapoint is returned.
func (wf WeatherForecast) At(timestamp time.Time) WeatherForecastDatapoint {
	datapoint := findClosestForecast(wf.Data, timestamp)
	if datapoint == nil {
		return WeatherForecastDatapoint{}
	}
	return newWeatherForecastDataPoint(*datapoint)
}

// All returns a slice of WeatherForecastDatapoint representing all forecasted data points.
func (wf WeatherForecast) All() []WeatherForecastDatapoint {
	datapoints := make([]WeatherForecastDatapoint, 0)
	for _, data := range wf.Data {
		datapoint := newWeatherForecastDataPoint(data)
		datapoints = append(datapoints, datapoint)
	}
	return datapoints
}

// CloudCoverage returns the cloud coverage data point as Coverage.
//
// If the data point is not available in the WeatherForecast it will return Coverage in which
// the "not available" field will be true.
func (dp WeatherForecastDatapoint) CloudCoverage() Coverage {
	if dp.cloudCoverage.IsNil() {
		return Coverage{notAvailable: true}
	}
	coverage := Coverage{
		dateTime: dp.dateTime,
		name:     FieldCloudCoverage,
		source:   SourceForecast,
		floatVal: dp.cloudCoverage.value,
	}
	return coverage
}

// DateTime returns the date and time of the WeatherForecastDatapoint.
func (dp WeatherForecastDatapoint) DateTime() time.Time {
	return dp.dateTime
}

// Dewpoint returns the dewpoint data point as Temperature.
//
// If the data point is not available in the WeatherForecast it will return Temperature in which the
// "not available" field will be true.
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

// HumidityRelative returns the relative humidity data point as Humidity.
//
// If the data point is not available in the WeatherForecast it will return Humidity in which the
// "not available" field will be true.
func (dp WeatherForecastDatapoint) HumidityRelative() Humidity {
	if dp.humidity.IsNil() {
		return Humidity{notAvailable: true}
	}
	humidity := Humidity{
		dateTime: dp.dateTime,
		name:     FieldHumidityRelative,
		source:   SourceForecast,
		floatVal: dp.humidity.Get(),
	}
	return humidity
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

// SunHours returns the sun hours data point as Duration.
//
// If the data point is not available in the WeatherForecast it will return Duration in which the
// "not available" field will be true.
func (dp WeatherForecastDatapoint) SunHours() Duration {
	if dp.winddirection.IsNil() {
		return Duration{notAvailable: true}
	}
	duration := Duration{
		dateTime: dp.dateTime,
		name:     FieldSunhours,
		source:   SourceForecast,
		floatVal: dp.sunhours.Get(),
	}
	return duration
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

// WeatherSymbol returns a text representation of the weather forecast as Condition.
//
// If the data point is not available in the WeatherForecast, it will return Condition in which
// the "not available" field will be true.
func (dp WeatherForecastDatapoint) WeatherSymbol() Condition {
	if dp.weatherSymbol.IsNil() {
		return Condition{notAvailable: true}
	}
	condition := Condition{
		dateTime:  dp.dateTime,
		name:      FieldWeatherSymbol,
		source:    SourceForecast,
		stringVal: dp.weatherSymbol.value,
	}
	return condition
}

// WindDirection returns the wind direction data point as Direction.
//
// If the data point is not available in the WeatherForecast it will return Direction in which the
// "not available" field will be true.
func (dp WeatherForecastDatapoint) WindDirection() Direction {
	if dp.winddirection.IsNil() {
		return Direction{notAvailable: true}
	}
	direction := Direction{
		dateTime: dp.dateTime,
		name:     FieldWindDirection,
		source:   SourceForecast,
		floatVal: dp.winddirection.Get(),
	}
	return direction
}

// WindGust returns the wind gust data point as Speed.
//
// If the data point is not available in the WeatherForecast it will return Speed in which the
// "not available" field will be true.
func (dp WeatherForecastDatapoint) WindGust() Speed {
	if dp.windgust.IsNil() {
		return Speed{notAvailable: true}
	}
	speed := Speed{
		dateTime: dp.dateTime,
		name:     FieldWindGust,
		source:   SourceForecast,
		floatVal: dp.windgust.Get(),
	}
	return speed
}

// WindGust3h returns the wind gust over the last 3 hours data point as Speed.
//
// If the data point is not available in the WeatherForecast it will return Speed in which the
// "not available" field will be true.
func (dp WeatherForecastDatapoint) WindGust3h() Speed {
	if dp.windgust3h.IsNil() {
		return Speed{notAvailable: true}
	}
	speed := Speed{
		dateTime: dp.dateTime,
		name:     FieldWindGust3h,
		source:   SourceForecast,
		floatVal: dp.windgust3h.Get(),
	}
	return speed
}

// WindSpeed returns the average wind speed data point as Speed.
//
// If the data point is not available in the WeatherForecast it will return Speed in which the
// "not available" field will be true.
func (dp WeatherForecastDatapoint) WindSpeed() Speed {
	if dp.windspeed.IsNil() {
		return Speed{notAvailable: true}
	}
	speed := Speed{
		dateTime: dp.dateTime,
		name:     FieldWindSpeed,
		source:   SourceForecast,
		floatVal: dp.windspeed.Get(),
	}
	return speed
}

// findClosestForecast finds the APIWeatherForecastData item in the given items slice
// that has the closest DateTime value to the target time. It returns a pointer to
// the closest item. If the items slice is empty, it returns nil.
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

// newWeatherForecastDataPoint creates a new WeatherForecastDatapoint from the provided APIWeatherForecastData.
// It extracts the necessary data from the APIWeatherForecastData and sets them in the WeatherForecastDatapoint
// structure. The new WeatherForecastDatapoint is then returned.
func newWeatherForecastDataPoint(data APIWeatherForecastData) WeatherForecastDatapoint {
	return WeatherForecastDatapoint{
		cloudCoverage: data.CloudCoverage,
		dateTime:      data.DateTime,
		dewpoint:      data.Dewpoint,
		humidity:      data.Humidity,
		isDay:         data.IsDay,
		pressureMSL:   data.PressureMSL,
		sunhours:      data.SunHours,
		temperature:   data.Temperature,
		weatherSymbol: data.WeatherSymbol,
		winddirection: data.WindDirection,
		windgust:      data.WindGust,
		windgust3h:    data.WindGust3h,
		windspeed:     data.WindSpeed,
	}
}
