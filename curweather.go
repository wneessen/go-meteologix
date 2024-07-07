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

// APICurrentWeatherData holds the different data points of the CurrentWeather as
// returned by the current weather API endpoints.
//
// Please keep in mind that different Station types return different values, therefore
// all values are represented as pointer type returning nil if the data point in question
// is not returned for the requested Station.
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
func (c *Client) CurrentWeatherByCoordinates(la, lo float64) (CurrentWeather, error) {
	var cw CurrentWeather
	lat := strconv.FormatFloat(la, 'f', -1, 64)
	lon := strconv.FormatFloat(lo, 'f', -1, 64)
	u, err := url.Parse(fmt.Sprintf("%s/current/%s/%s", c.config.apiURL, lat, lon))
	if err != nil {
		return cw, fmt.Errorf("failed to parse current weather URL: %w", err)
	}
	uq := u.Query()
	uq.Add("units", "metric")
	u.RawQuery = uq.Encode()

	r, err := c.httpClient.Get(u.String())
	if err != nil {
		return cw, fmt.Errorf("API request failed: %w", err)
	}

	if err := json.Unmarshal(r, &cw); err != nil {
		return cw, fmt.Errorf("failed to unmarshal API response JSON: %w", err)
	}

	return cw, nil
}

// CurrentWeatherByLocation returns the CurrentWeather values for the given location
func (c *Client) CurrentWeatherByLocation(lo string) (CurrentWeather, error) {
	gl, err := c.GetGeoLocationByName(lo)
	if err != nil {
		return CurrentWeather{}, fmt.Errorf("failed too look up geolocation: %w", err)
	}
	return c.CurrentWeatherByCoordinates(gl.Latitude, gl.Longitude)
}

// Dewpoint returns the dewpoint data point as Temperature.
// If the data point is not available in the CurrentWeather it will return
// Temperature in which the "not available" field will be true.
func (cw CurrentWeather) Dewpoint() Temperature {
	if cw.Data.Dewpoint == nil {
		return Temperature{notAvailable: true}
	}
	v := Temperature{
		dateTime: cw.Data.Dewpoint.DateTime,
		name:     FieldDewpoint,
		source:   SourceUnknown,
		floatVal: cw.Data.Dewpoint.Value,
	}
	if cw.Data.Dewpoint.Source != nil {
		v.source = StringToSource(*cw.Data.Dewpoint.Source)
	}
	return v
}

// HumidityRelative returns the relative humidity data point as Humidity.
// If the data point is not available in the CurrentWeather it will return
// Humidity in which the "not available" field will be true.
func (cw CurrentWeather) HumidityRelative() Humidity {
	if cw.Data.HumidityRelative == nil {
		return Humidity{notAvailable: true}
	}
	v := Humidity{
		dateTime: cw.Data.HumidityRelative.DateTime,
		name:     FieldHumidityRelative,
		source:   SourceUnknown,
		floatVal: cw.Data.HumidityRelative.Value,
	}
	if cw.Data.HumidityRelative.Source != nil {
		v.source = StringToSource(*cw.Data.HumidityRelative.Source)
	}
	return v
}

// IsDay returns true if it is currently day at queried location
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
func (cw CurrentWeather) Precipitation(ts Timespan) Precipitation {
	var df *APIFloat
	var fn Fieldname
	switch ts {
	case TimespanCurrent:
		df = cw.Data.Precipitation
		fn = FieldPrecipitation
	case Timespan10Min:
		df = cw.Data.Precipitation10m
		fn = FieldPrecipitation10m
	case Timespan1Hour:
		df = cw.Data.Precipitation1h
		fn = FieldPrecipitation1h
	case Timespan24Hours:
		df = cw.Data.Precipitation24h
		fn = FieldPrecipitation24h
	default:
		return Precipitation{notAvailable: true}
	}

	if df == nil {
		return Precipitation{notAvailable: true}
	}
	v := Precipitation{
		dateTime: df.DateTime,
		name:     fn,
		source:   SourceUnknown,
		floatVal: df.Value,
	}
	if df.Source != nil {
		v.source = StringToSource(*df.Source)
	}
	return v
}

// PressureMSL returns the pressure at mean sea level data point as Pressure.
// If the data point is not available in the CurrentWeather it will return
// Pressure in which the "not available" field will be true.
func (cw CurrentWeather) PressureMSL() Pressure {
	if cw.Data.PressureMSL == nil {
		return Pressure{notAvailable: true}
	}
	v := Pressure{
		dateTime: cw.Data.PressureMSL.DateTime,
		name:     FieldPressureMSL,
		source:   SourceUnknown,
		floatVal: cw.Data.PressureMSL.Value,
	}
	if cw.Data.PressureMSL.Source != nil {
		v.source = StringToSource(*cw.Data.PressureMSL.Source)
	}
	return v
}

// PressureQFE returns the pressure at mean sea level data point as Pressure.
// If the data point is not available in the CurrentWeather it will return
// Pressure in which the "not available" field will be true.
func (cw CurrentWeather) PressureQFE() Pressure {
	if cw.Data.PressureQFE == nil {
		return Pressure{notAvailable: true}
	}
	v := Pressure{
		dateTime: cw.Data.PressureQFE.DateTime,
		name:     FieldPressureQFE,
		source:   SourceUnknown,
		floatVal: cw.Data.PressureQFE.Value,
	}
	if cw.Data.PressureQFE.Source != nil {
		v.source = StringToSource(*cw.Data.PressureQFE.Source)
	}
	return v
}

// SnowAmount returns the amount of snow data point as Density.
// If the data point is not available in the CurrentWeather it will return
// Density in which the "not available" field will be true.
func (cw CurrentWeather) SnowAmount() Density {
	if cw.Data.SnowAmount == nil {
		return Density{notAvailable: true}
	}
	v := Density{
		dateTime: cw.Data.SnowAmount.DateTime,
		name:     FieldSnowAmount,
		source:   SourceUnknown,
		floatVal: cw.Data.SnowAmount.Value,
	}
	if cw.Data.SnowAmount.Source != nil {
		v.source = StringToSource(*cw.Data.SnowAmount.Source)
	}
	return v
}

// SnowHeight returns the snow height data point as Height.
// If the data point is not available in the CurrentWeather it will return
// Height in which the "not available" field will be true.
func (cw CurrentWeather) SnowHeight() Height {
	if cw.Data.SnowHeight == nil {
		return Height{notAvailable: true}
	}
	v := Height{
		dateTime: cw.Data.SnowHeight.DateTime,
		name:     FieldSnowHeight,
		source:   SourceUnknown,
		floatVal: cw.Data.SnowHeight.Value,
	}
	if cw.Data.SnowHeight.Source != nil {
		v.source = StringToSource(*cw.Data.SnowHeight.Source)
	}
	return v
}

// Temperature returns the temperature data point as Temperature.
// If the data point is not available in the CurrentWeather it will return
// Temperature in which the "not available" field will be true.
func (cw CurrentWeather) Temperature() Temperature {
	if cw.Data.Temperature == nil {
		return Temperature{notAvailable: true}
	}
	v := Temperature{
		dateTime: cw.Data.Temperature.DateTime,
		name:     FieldTemperature,
		source:   SourceUnknown,
		floatVal: cw.Data.Temperature.Value,
	}
	if cw.Data.Temperature.Source != nil {
		v.source = StringToSource(*cw.Data.Temperature.Source)
	}
	return v
}

// WeatherSymbol returns a text representation of the current weather
// as Condition.
// If the data point is not available in the CurrentWeather it will return
// Condition in which the "not available" field will be true.
func (cw CurrentWeather) WeatherSymbol() Condition {
	if cw.Data.WeatherSymbol == nil {
		return Condition{notAvailable: true}
	}
	v := Condition{
		dateTime:  cw.Data.WeatherSymbol.DateTime,
		name:      FieldWeatherSymbol,
		source:    SourceUnknown,
		stringVal: cw.Data.WeatherSymbol.Value,
	}
	if cw.Data.WeatherSymbol.Source != nil {
		v.source = StringToSource(*cw.Data.WeatherSymbol.Source)
	}
	return v
}

// WindDirection returns the wind direction data point as Direction.
// If the data point is not available in the CurrentWeather it will return
// Direction in which the "not available" field will be true.
func (cw CurrentWeather) WindDirection() Direction {
	if cw.Data.WindDirection == nil {
		return Direction{notAvailable: true}
	}
	v := Direction{
		dateTime: cw.Data.WindDirection.DateTime,
		name:     FieldWindDirection,
		source:   SourceUnknown,
		floatVal: cw.Data.WindDirection.Value,
	}
	if cw.Data.WindDirection.Source != nil {
		v.source = StringToSource(*cw.Data.WindDirection.Source)
	}
	return v
}

// WindGust returns the wind gust data point as Speed.
// If the data point is not available in the CurrentWeather it will return
// Speed in which the "not available" field will be true.
func (cw CurrentWeather) WindGust() Speed {
	if cw.Data.WindGust == nil {
		return Speed{notAvailable: true}
	}
	v := Speed{
		dateTime: cw.Data.WindGust.DateTime,
		name:     FieldWindGust,
		source:   SourceUnknown,
		floatVal: cw.Data.WindGust.Value,
	}
	if cw.Data.WindGust.Source != nil {
		v.source = StringToSource(*cw.Data.WindGust.Source)
	}
	return v
}

// WindSpeed returns the average wind speed data point as Speed.
// If the data point is not available in the CurrentWeather it will return
// Speed in which the "not available" field will be true.
func (cw CurrentWeather) WindSpeed() Speed {
	if cw.Data.WindSpeed == nil {
		return Speed{notAvailable: true}
	}
	v := Speed{
		dateTime: cw.Data.WindSpeed.DateTime,
		name:     FieldWindSpeed,
		source:   SourceUnknown,
		floatVal: cw.Data.WindSpeed.Value,
	}
	if cw.Data.WindSpeed.Source != nil {
		v.source = StringToSource(*cw.Data.WindSpeed.Source)
	}
	return v
}
