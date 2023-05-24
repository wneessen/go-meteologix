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
		// PressureMSL represents the pressure at station level (QFE) in hPa
		PressureQFE *APIValue `json:"pressure"`

	*/
	// Temperature represents the temperature in °C
	Temperature *APIValue `json:"temp,omitempty"`
	/*
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
		// Winddirection represents the direction from which the wind
		// originates in degree (0=N, 90=E, 180=S, 270=W)
		Winddirection *APIValue `json:"windDirection,omitempty"`
		// Windspeed represents the wind speed in knots
		Windspeed *APIValue `json:"windSpeed,omitempty"`

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
