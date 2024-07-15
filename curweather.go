// SPDX-FileCopyrightText: 2023 Winni Neessen <wn@neessen.dev>
//
// SPDX-License-Identifier: MIT

package meteologix

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// CurrentWeather represents the current weather API response
type CurrentWeather struct {
	// Data holds the different APICurrentWeatherData points
	Data APICurrentWeatherData `json:"data"`
	// Latitude represents the GeoLocation latitude coordinates for the weather data
	Latitude float64 `json:"lat"`
	// Longitude represents the GeoLocation longitude coordinates for the weather data
	Longitude float64 `json:"lon"`
	// UnitSystem is the unit system that is used for the results (we default to metric)
	UnitSystem string `json:"systemOfUnits"`
}

// APICurrentWeatherData holds the different data points of the CurrentWeather as returned by the
// current weather API endpoints.
//
// Please keep in mind that different Station types return different values, therefore all values
// are represented as pointer type returning nil if the data point in question is not returned for
// the requested Station.
type APICurrentWeatherData struct {
	// Dewpoint represents the dewpoint in °C
	Dewpoint *APIFloat `json:"dewpoint,omitempty"`
	// HumidityRelative represents the relative humidity in percent
	HumidityRelative *APIFloat `json:"humidityRelative,omitempty"`
	// IsDay is true when it is currently daytime
	IsDay *APIBool `json:"isDay"`
	// Precipitation represents the current amount of precipitation
	Precipitation *APIFloat `json:"prec,omitempty"`
	// Precipitation10m represents the amount of precipitation over the last 10 minutes
	Precipitation10m *APIFloat `json:"prec10m,omitempty"`
	// Precipitation1h represents the amount of precipitation over the last hour
	Precipitation1h *APIFloat `json:"prec1h,omitempty"`
	// Precipitation24h represents the amount of precipitation over the last 24 hours
	Precipitation24h *APIFloat `json:"prec24h,omitempty"`
	// PressureMSL represents the pressure at mean sea level (MSL) in hPa
	PressureMSL *APIFloat `json:"pressureMsl,omitempty"`
	// PressureQFE represents the pressure at station level (QFE) in hPa
	PressureQFE *APIFloat `json:"pressure,omitempty"`
	// SnowAmount represents the the amount of snow in kg/m3
	SnowAmount *APIFloat `json:"snowAmount,omitempty"`
	// SnowHeight represents the the height of snow in m
	SnowHeight *APIFloat `json:"snowHeight,omitempty"`
	// Temperature represents the temperature in °C
	Temperature *APIFloat `json:"temp,omitempty"`
	// WindDirection represents the direction from which the wind
	// originates in degree (0=N, 90=E, 180=S, 270=W)
	WindDirection *APIFloat `json:"windDirection,omitempty"`
	// WindGust represents the wind gust speed in m/s
	WindGust *APIFloat `json:"windGust,omitempty"`
	// WindSpeed represents the wind speed in m/s
	WindSpeed *APIFloat `json:"windSpeed,omitempty"`
	// WeatherSymbol is a text representation of the current weather
	// conditions
	WeatherSymbol *APIString `json:"weatherSymbol,omitempty"`
}

// CurrentWeatherByCoordinates returns the CurrentWeather values for the given coordinates
func (c *Client) CurrentWeatherByCoordinates(latitude, longitude float64) (CurrentWeather, error) {
	var currentWeather CurrentWeather
	latitudeFormat := strconv.FormatFloat(latitude, 'f', -1, 64)
	longitudeFormat := strconv.FormatFloat(longitude, 'f', -1, 64)
	apiURL, err := url.Parse(fmt.Sprintf("%s/current/%s/%s", c.config.apiURL, latitudeFormat, longitudeFormat))
	if err != nil {
		return currentWeather, fmt.Errorf("failed to parse current weather URL: %w", err)
	}
	queryString := apiURL.Query()
	queryString.Add("units", "metric")
	apiURL.RawQuery = queryString.Encode()

	response, err := c.httpClient.Get(apiURL.String())
	if err != nil {
		return currentWeather, fmt.Errorf("API request failed: %w", err)
	}

	if err := json.Unmarshal(response, &currentWeather); err != nil {
		return currentWeather, fmt.Errorf("failed to unmarshal API response JSON: %w", err)
	}

	return currentWeather, nil
}

// CurrentWeatherByLocation returns the CurrentWeather values for the given location
func (c *Client) CurrentWeatherByLocation(location string) (CurrentWeather, error) {
	geoLocation, err := c.GetGeoLocationByName(location)
	if err != nil {
		return CurrentWeather{}, fmt.Errorf("failed too look up geolocation: %w", err)
	}
	return c.CurrentWeatherByCoordinates(geoLocation.Latitude, geoLocation.Longitude)
}

// Dewpoint returns the dewpoint data point as Temperature.
//
// If the data point is not available in the CurrentWeather it will return Temperature in which
// the "not available" field will be true.
func (cw CurrentWeather) Dewpoint() Temperature {
	if cw.Data.Dewpoint == nil {
		return Temperature{notAvailable: true}
	}
	temperature := Temperature{
		dateTime: cw.Data.Dewpoint.DateTime,
		name:     FieldDewpoint,
		source:   SourceUnknown,
		floatVal: cw.Data.Dewpoint.Value,
	}
	if cw.Data.Dewpoint.Source != nil {
		temperature.source = StringToSource(*cw.Data.Dewpoint.Source)
	}
	return temperature
}

// HumidityRelative returns the relative humidity data point as Humidity.
//
// If the data point is not available in the CurrentWeather it will return Humidity in which
// the "not available" field will be true.
func (cw CurrentWeather) HumidityRelative() Humidity {
	if cw.Data.HumidityRelative == nil {
		return Humidity{notAvailable: true}
	}
	humidity := Humidity{
		dateTime: cw.Data.HumidityRelative.DateTime,
		name:     FieldHumidityRelative,
		source:   SourceUnknown,
		floatVal: cw.Data.HumidityRelative.Value,
	}
	if cw.Data.HumidityRelative.Source != nil {
		humidity.source = StringToSource(*cw.Data.HumidityRelative.Source)
	}
	return humidity
}

// IsDay returns true if it is day time at the current location.
func (cw CurrentWeather) IsDay() bool {
	if cw.Data.IsDay == nil {
		return false
	}
	return cw.Data.IsDay.Value
}

// Precipitation returns the current amount of precipitation (mm) as Precipitation
// If the data point is not available in the CurrentWeather it will return
// Precipitation in which the "not available" field will be true.
//
// At this point of development, it looks like currently only the 1 Hour value
// is returned by the endpoint, so expect non-availability for any other Timespan
// at this point.
func (cw CurrentWeather) Precipitation(timeSpan Timespan) Precipitation {
	var apiFloat *APIFloat
	var fieldName Fieldname
	switch timeSpan {
	case TimespanCurrent:
		apiFloat = cw.Data.Precipitation
		fieldName = FieldPrecipitation
	case Timespan10Min:
		apiFloat = cw.Data.Precipitation10m
		fieldName = FieldPrecipitation10m
	case Timespan1Hour:
		apiFloat = cw.Data.Precipitation1h
		fieldName = FieldPrecipitation1h
	case Timespan24Hours:
		apiFloat = cw.Data.Precipitation24h
		fieldName = FieldPrecipitation24h
	default:
		return Precipitation{notAvailable: true}
	}

	if apiFloat == nil {
		return Precipitation{notAvailable: true}
	}
	precipitation := Precipitation{
		dateTime: apiFloat.DateTime,
		name:     fieldName,
		source:   SourceUnknown,
		floatVal: apiFloat.Value,
	}
	if apiFloat.Source != nil {
		precipitation.source = StringToSource(*apiFloat.Source)
	}
	return precipitation
}

// PressureMSL returns the pressure at mean sea level data point as Pressure.
//
// If the data point is not available in the CurrentWeather it will return Pressure in which
// the "not available" field will be true.
func (cw CurrentWeather) PressureMSL() Pressure {
	if cw.Data.PressureMSL == nil {
		return Pressure{notAvailable: true}
	}
	pressure := Pressure{
		dateTime: cw.Data.PressureMSL.DateTime,
		name:     FieldPressureMSL,
		source:   SourceUnknown,
		floatVal: cw.Data.PressureMSL.Value,
	}
	if cw.Data.PressureMSL.Source != nil {
		pressure.source = StringToSource(*cw.Data.PressureMSL.Source)
	}
	return pressure
}

// PressureQFE returns the pressure at mean sea level data point as Pressure.
//
// If the data point is not available in the CurrentWeather it will return Pressure in which
// the "not available" field will be true.
func (cw CurrentWeather) PressureQFE() Pressure {
	if cw.Data.PressureQFE == nil {
		return Pressure{notAvailable: true}
	}
	pressure := Pressure{
		dateTime: cw.Data.PressureQFE.DateTime,
		name:     FieldPressureQFE,
		source:   SourceUnknown,
		floatVal: cw.Data.PressureQFE.Value,
	}
	if cw.Data.PressureQFE.Source != nil {
		pressure.source = StringToSource(*cw.Data.PressureQFE.Source)
	}
	return pressure
}

// SnowAmount returns the amount of snow data point as Density.
//
// If the data point is not available in the CurrentWeather it will return Density in which
// the "not available" field will be true.
func (cw CurrentWeather) SnowAmount() Density {
	if cw.Data.SnowAmount == nil {
		return Density{notAvailable: true}
	}
	density := Density{
		dateTime: cw.Data.SnowAmount.DateTime,
		name:     FieldSnowAmount,
		source:   SourceUnknown,
		floatVal: cw.Data.SnowAmount.Value,
	}
	if cw.Data.SnowAmount.Source != nil {
		density.source = StringToSource(*cw.Data.SnowAmount.Source)
	}
	return density
}

// SnowHeight returns the snow height data point as Height.
//
// If the data point is not available in the CurrentWeather it will return Height in which
// the "not available" field will be true.
func (cw CurrentWeather) SnowHeight() Height {
	if cw.Data.SnowHeight == nil {
		return Height{notAvailable: true}
	}
	height := Height{
		dateTime: cw.Data.SnowHeight.DateTime,
		name:     FieldSnowHeight,
		source:   SourceUnknown,
		floatVal: cw.Data.SnowHeight.Value,
	}
	if cw.Data.SnowHeight.Source != nil {
		height.source = StringToSource(*cw.Data.SnowHeight.Source)
	}
	return height
}

// Temperature returns the temperature data point as Temperature.
//
// If the data point is not available in the CurrentWeather it will return Temperature in which
// the "not available" field will be true.
func (cw CurrentWeather) Temperature() Temperature {
	if cw.Data.Temperature == nil {
		return Temperature{notAvailable: true}
	}
	temperature := Temperature{
		dateTime: cw.Data.Temperature.DateTime,
		name:     FieldTemperature,
		source:   SourceUnknown,
		floatVal: cw.Data.Temperature.Value,
	}
	if cw.Data.Temperature.Source != nil {
		temperature.source = StringToSource(*cw.Data.Temperature.Source)
	}
	return temperature
}

// WeatherSymbol returns a text representation of the current weather as Condition.
//
// If the data point is not available in the CurrentWeather it will return Condition in which
// the "not available" field will be true.
func (cw CurrentWeather) WeatherSymbol() Condition {
	if cw.Data.WeatherSymbol == nil {
		return Condition{notAvailable: true}
	}
	condition := Condition{
		dateTime:  cw.Data.WeatherSymbol.DateTime,
		name:      FieldWeatherSymbol,
		source:    SourceUnknown,
		stringVal: cw.Data.WeatherSymbol.Value,
	}
	if cw.Data.WeatherSymbol.Source != nil {
		condition.source = StringToSource(*cw.Data.WeatherSymbol.Source)
	}
	return condition
}

// WindDirection returns the wind direction data point as Direction.
//
// If the data point is not available in the CurrentWeather it will return Direction in which
// the "not available" field will be true.
func (cw CurrentWeather) WindDirection() Direction {
	if cw.Data.WindDirection == nil {
		return Direction{notAvailable: true}
	}
	direction := Direction{
		dateTime: cw.Data.WindDirection.DateTime,
		name:     FieldWindDirection,
		source:   SourceUnknown,
		floatVal: cw.Data.WindDirection.Value,
	}
	if cw.Data.WindDirection.Source != nil {
		direction.source = StringToSource(*cw.Data.WindDirection.Source)
	}
	return direction
}

// WindGust returns the wind gust data point as Speed.
//
// If the data point is not available in the CurrentWeather it will return Speed in which
// the "not available" field will be true.
func (cw CurrentWeather) WindGust() Speed {
	if cw.Data.WindGust == nil {
		return Speed{notAvailable: true}
	}
	speed := Speed{
		dateTime: cw.Data.WindGust.DateTime,
		name:     FieldWindGust,
		source:   SourceUnknown,
		floatVal: cw.Data.WindGust.Value,
	}
	if cw.Data.WindGust.Source != nil {
		speed.source = StringToSource(*cw.Data.WindGust.Source)
	}
	return speed
}

// WindSpeed returns the average wind speed data point as Speed.
//
// If the data point is not available in the CurrentWeather it will return Speed in which
// the "not available" field will be true.
func (cw CurrentWeather) WindSpeed() Speed {
	if cw.Data.WindSpeed == nil {
		return Speed{notAvailable: true}
	}
	speed := Speed{
		dateTime: cw.Data.WindSpeed.DateTime,
		name:     FieldWindSpeed,
		source:   SourceUnknown,
		floatVal: cw.Data.WindSpeed.Value,
	}
	if cw.Data.WindSpeed.Source != nil {
		speed.source = StringToSource(*cw.Data.WindSpeed.Source)
	}
	return speed
}
