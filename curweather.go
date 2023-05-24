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
	Dewpoint *APIValue `json:"dewpoint,omitempty"`
	// HumidityRelative represents the relative humidity in percent
	HumidityRelative *APIValue `json:"humidityRelative,omitempty"`
	// Precipitation represents the current amount of precipitation
	Precipitation *APIValue `json:"prec"`
	// Precipitation10m represents the amount of precipitation over the last 10 minutes
	Precipitation10m *APIValue `json:"prec10m"`
	// Precipitation1h represents the amount of precipitation over the last hour
	Precipitation1h *APIValue `json:"prec1h"`
	// Precipitation24h represents the amount of precipitation over the last 24 hours
	Precipitation24h *APIValue `json:"prec24h"`
	// PressureMSL represents the pressure at mean sea level (MSL) in hPa
	PressureMSL *APIValue `json:"pressureMsl"`
	// Temperature represents the temperature in °C
	Temperature *APIValue `json:"temp,omitempty"`
	// Windspeed represents the wind speed in knots
	Windspeed *APIValue `json:"windSpeed,omitempty"`
	// Winddirection represents the direction from which the wind
	// originates in degree (0=N, 90=E, 180=S, 270=W)
	Winddirection *APIValue `json:"windDirection,omitempty"`
	/*
		// DewPointMean represents the mean dewpoint in °C
		DewpointMean *APIValue `json:"dewpointMean,omitempty"`
		// GlobalRadiation10m represents the sum of global radiation over the last
		// 10 minutes in kJ/m²
		GlobalRadiation10m *APIValue `json:"globalRadiation10m,omitempty"`
		// GlobalRadiation1h represents the sum of global radiation over the last
		// 1 hour in kJ/m²
		GlobalRadiation1h *APIValue `json:"globalRadiation1h,omitempty"`
		// GlobalRadiation24h represents the sum of global radiation over the last
		// 24 hour in kJ/m²
		GlobalRadiation24h *APIValue `json:"globalRadiation24h,omitempty"`
		// PressureMSL represents the pressure at station level (QFE) in hPa
		PressureQFE *APIValue `json:"pressure"`
		// TemperatureMax represents the maximum temperature in °C
		TemperatureMax *APIValue `json:"tempMax,omitempty"`
		// TemperatureMean represents the mean temperature in °C
		TemperatureMean *APIValue `json:"tempMean,omitempty"`
		// TemperatureMin represents the minimum temperature in °C
		TemperatureMin *APIValue `json:"tempMin,omitempty"`
		// Temperature5cm represents the temperature 5cm above ground in °C
		Temperature5cm *APIValue `json:"temp5cm,omitempty"`
		// Temperature5cm represents the minimum temperature 5cm above
		// ground in °C
		Temperature5cmMin *APIValue `json:"temp5cmMin,omitempty"`

	*/
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

// Temperature returns the temperature data point as Temperature.
// If the data point is not available in the CurrentWeather it will return
// Temperature in which the "not available" field will be true.
func (cw CurrentWeather) Temperature() Temperature {
	if cw.Data.Temperature == nil {
		return Temperature{na: true}
	}
	v := Temperature{
		dt: cw.Data.Temperature.DateTime,
		n:  FieldTemperature,
		s:  SourceUnknown,
		v:  cw.Data.Temperature.Value,
	}
	if cw.Data.Temperature.Source != nil {
		v.s = StringToSource(*cw.Data.Temperature.Source)
	}
	return v
}

// Dewpoint returns the dewpoint data point as Temperature.
// If the data point is not available in the CurrentWeather it will return
// Temperature in which the "not available" field will be true.
func (cw CurrentWeather) Dewpoint() Temperature {
	if cw.Data.Dewpoint == nil {
		return Temperature{na: true}
	}
	v := Temperature{
		dt: cw.Data.Dewpoint.DateTime,
		n:  FieldDewpoint,
		s:  SourceUnknown,
		v:  cw.Data.Dewpoint.Value,
	}
	if cw.Data.Dewpoint.Source != nil {
		v.s = StringToSource(*cw.Data.Dewpoint.Source)
	}
	return v
}

// HumidityRelative returns the relative humidity data point as Humidity.
// If the data point is not available in the CurrentWeather it will return
// Humidity in which the "not available" field will be true.
func (cw CurrentWeather) HumidityRelative() Humidity {
	if cw.Data.Dewpoint == nil {
		return Humidity{na: true}
	}
	v := Humidity{
		dt: cw.Data.HumidityRelative.DateTime,
		n:  FieldHumidityRelative,
		s:  SourceUnknown,
		v:  cw.Data.HumidityRelative.Value,
	}
	if cw.Data.HumidityRelative.Source != nil {
		v.s = StringToSource(*cw.Data.HumidityRelative.Source)
	}
	return v
}

// Precipitation returns the current amount of precipitation (mm) as Precipitation
// If the data point is not available in the CurrentWeather it will return
// Precipitation in which the "not available" field will be true.
//
// At this point of development, it looks like currently only the 1 Hour value
// is returned by the endpoint, so expect non-availablity for any other Timespan
// at this point.
func (cw CurrentWeather) Precipitation(ts Timespan) Precipitation {
	var df *APIValue
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
		return Precipitation{na: true}
	}

	if df == nil {
		return Precipitation{na: true}
	}
	v := Precipitation{
		dt: df.DateTime,
		n:  fn,
		s:  SourceUnknown,
		v:  df.Value,
	}
	if df.Source != nil {
		v.s = StringToSource(*df.Source)
	}
	return v
}

// PressureMSL returns the pressure at mean sea level data point as Pressure.
// If the data point is not available in the CurrentWeather it will return
// Pressure in which the "not available" field will be true.
func (cw CurrentWeather) PressureMSL() Pressure {
	if cw.Data.PressureMSL == nil {
		return Pressure{na: true}
	}
	v := Pressure{
		dt: cw.Data.PressureMSL.DateTime,
		n:  FieldPressureMSL,
		s:  SourceUnknown,
		v:  cw.Data.PressureMSL.Value,
	}
	if cw.Data.PressureMSL.Source != nil {
		v.s = StringToSource(*cw.Data.PressureMSL.Source)
	}
	return v
}

// Winddirection returns the wind direction data point as Direction.
// If the data point is not available in the CurrentWeather it will return
// Direction in which the "not available" field will be true.
func (cw CurrentWeather) Winddirection() Direction {
	if cw.Data.Winddirection == nil {
		return Direction{na: true}
	}
	v := Direction{
		dt: cw.Data.Winddirection.DateTime,
		n:  FieldWinddirection,
		s:  SourceUnknown,
		v:  cw.Data.Winddirection.Value,
	}
	if cw.Data.Winddirection.Source != nil {
		v.s = StringToSource(*cw.Data.Winddirection.Source)
	}
	return v
}

// Windspeed returns the wind speed data point as Speed.
// If the data point is not available in the CurrentWeather it will return
// Speed in which the "not available" field will be true.
func (cw CurrentWeather) Windspeed() Speed {
	if cw.Data.Windspeed == nil {
		return Speed{na: true}
	}
	v := Speed{
		dt: cw.Data.Windspeed.DateTime,
		n:  FieldWindspeed,
		s:  SourceUnknown,
		v:  cw.Data.Windspeed.Value,
	}
	if cw.Data.Windspeed.Source != nil {
		v.s = StringToSource(*cw.Data.Windspeed.Source)
	}
	return v
}
