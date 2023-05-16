// SPDX-FileCopyrightText: 2023 Winni Neessen <winni@neessen.dev>
//
// SPDX-License-Identifier: MIT

package meteologix

import (
	"encoding/json"
	"fmt"
	"time"
)

// DataNotAvailable is returned if a requested data point returned no data
const DataNotAvailable = "data not available"

// Observation represents the observation API response for a Station
type Observation struct {
	// Altitude is the altitude of the station providing the Observation
	Altitude *int `json:"ele,omitempty"`
	// Data holds the different ObservationData points
	Data ObservationData `json:"data"`
	// Name is the name of the Station providing the Observation
	Name string `json:"name"`
	// Latitude represents the GeoLocation latitude coordinates for the Station
	Latitude float64 `json:"lat"`
	// Longitude represents the GeoLocation longitude coordinates for the Station
	Longitude float64 `json:"lon"`
	// StationID is the ID of the Station providing the Observation
	StationID string `json:"stationId"`
}

// ObservationData holds the different data points of the Observation.
//
// Please keep in mind that different Station types return different values, therefore
// all values are represented as pointer type returning nil if the data point in question
// is not returned for the requested Station.
type ObservationData struct {
	// DewPoint represents the dewpoint in °C
	DewPoint *ObservationTemperature `json:"dewpoint,omitempty"`
	// HumidityRelative represents the relative humidity in percent
	HumidityRelative *ObservationHumidity `json:"humidityRelative,omitempty"`
	// Precipitation represents the current amount of precipitation
	Precipitation *ObservationPrecipitation `json:"prec"`
	// Precipitation10m represents the amount of precipitation over the last 10 minutes
	Precipitation10m *ObservationPrecipitation `json:"prec10m"`
	// Precipitation1h represents the amount of precipitation over the last hour
	Precipitation1h *ObservationPrecipitation `json:"prec1h"`
	// Precipitation24h represents the amount of precipitation over the last 24 hours
	Precipitation24h *ObservationPrecipitation `json:"prec24h"`
	// PressureMSL represents the pressure at mean sea level (MSL) in hPa
	PressureMSL *ObservationPressure `json:"pressureMsl"`
	// PressureMSL represents the pressure at station level (QFE) in hPa
	PressureQFE *ObservationPressure `json:"pressure"`
	// Temperature represents the temperature in °C
	Temperature *ObservationTemperature `json:"temp,omitempty"`
	// TemperatureMax represents the maximum temperature in °C
	TemperatureMax *ObservationTemperature `json:"tempMax,omitempty"`
	// TemperatureMin represents the minimum temperature in °C
	TemperatureMin *ObservationTemperature `json:"tempMin,omitempty"`
	// Temperature5cm represents the temperature 5cm above ground in °C
	Temperature5cm *ObservationTemperature `json:"temp5cm,omitempty"`
	// Temperature5cm represents the minimum temperature 5cm above
	// ground in °C
	Temperature5cmMin *ObservationTemperature `json:"temp5cmMin,omitempty"`
}

// ObservationTemperature is a type wrapper for a temperature value
// in an Observation
type ObservationTemperature ObservationValueFloat

// ObservationHumidity is a type wrapper for a humidity value
// in an Observation
type ObservationHumidity ObservationValueFloat

// ObservationPrecipitation is a type wrapper for a precipitation value
// in an Observation
type ObservationPrecipitation ObservationValueFloat

// ObservationPressure is a type wrapper for a pressure value
// in an Observation
type ObservationPressure ObservationValueFloat

// ObservationValueFloat represents a observation value returning a
// Float type
type ObservationValueFloat struct {
	DateTime time.Time `json:"dateTime"`
	Value    float64   `json:"value"`
}

// ObservationLatestByStationID returns the latest Observation values from the
// given Station
func (c *Client) ObservationLatestByStationID(si string) (Observation, error) {
	var o Observation
	u := fmt.Sprintf("%s/station/%s/observations/latest", c.config.apiURL, si)
	r, err := c.httpClient.Get(u)
	if err != nil {
		return o, fmt.Errorf("API request failed: %w", err)
	}

	if err := json.Unmarshal(r, &o); err != nil {
		return o, fmt.Errorf("failed to unmarshal API response JSON: %w", err)
	}

	return o, nil
}

// Dewpoint returns the dewpoint data point as formatted string.
// If the data point is not available in the Observation it will return a
// corresponding DataNotAvailable string
func (o Observation) Dewpoint() string {
	if o.Data.DewPoint == nil {
		return DataNotAvailable
	}
	return o.Data.DewPoint.String()
}

// Temperature returns the temperature data point as formatted string.
// If the data point is not available in the Observation it will return a
// corresponding DataNotAvailable string
func (o Observation) Temperature() string {
	if o.Data.Temperature == nil {
		return DataNotAvailable
	}
	return o.Data.Temperature.String()
}

// TemperatureAtGround returns the temperature at ground level (5cm)
// data point as formatted string.
// If the data point is not available in the Observation it will return a
// corresponding DataNotAvailable string
func (o Observation) TemperatureAtGround() string {
	if o.Data.Temperature5cm == nil {
		return DataNotAvailable
	}
	return o.Data.Temperature5cm.String()
}

// TemperatureMin returns the minimum temperature so far data point as
// formatted string.
// If the data point is not available in the Observation it will return a
// corresponding DataNotAvailable string
func (o Observation) TemperatureMin() string {
	if o.Data.TemperatureMin == nil {
		return DataNotAvailable
	}
	return o.Data.TemperatureMin.String()
}

// TemperatureMax returns the maximum temperature so far data point as
// formatted string.
// If the data point is not available in the Observation it will return a
// corresponding DataNotAvailable string
func (o Observation) TemperatureMax() string {
	if o.Data.TemperatureMax == nil {
		return DataNotAvailable
	}
	return o.Data.TemperatureMax.String()
}

// TemperatureAtGroundMin returns the minimum temperature so far
// at ground level (5cm) data point as formatted string.
// If the data point is not available in the Observation it will return a
// corresponding DataNotAvailable string
func (o Observation) TemperatureAtGroundMin() string {
	if o.Data.Temperature5cmMin == nil {
		return DataNotAvailable
	}
	return o.Data.Temperature5cmMin.String()
}

// String satisfies the fmt.Stringer interface for the ObservationTemperature type
func (t ObservationTemperature) String() string {
	return fmt.Sprintf("%.1f°C", t.Value)
}

func (t ObservationTemperature) Timestamp() time.Time {
	return t.DateTime
}

// Celsius returns the ObservationTemperature value in Celsius
func (t ObservationTemperature) Celsius() float64 {
	return t.Value
}

// CelsiusString returns the ObservationTemperature value as Celsius
// formated string.
//
// This is an alias for the fmt.Stringer interface
func (t ObservationTemperature) CelsiusString() string {
	return t.String()
}

// Fahrenheit returns the ObservationTemperature value in Fahrenheit
func (t ObservationTemperature) Fahrenheit() float64 {
	return t.Value*9/5 + 32
}

// FahrenheitString returns the ObservationTemperature value as Fahrenheit
// formated string.
func (t ObservationTemperature) FahrenheitString() string {
	return fmt.Sprintf("%.1f°F", t.Fahrenheit())
}

// String satisfies the fmt.Stringer interface for the ObservationHumidity type
func (t ObservationHumidity) String() string {
	return fmt.Sprintf("%.1f%%", t.Value)
}

// String satisfies the fmt.Stringer interface for the ObservationPrecipitation type
func (t ObservationPrecipitation) String() string {
	return fmt.Sprintf("%.1fmm", t.Value)
}

// String satisfies the fmt.Stringer interface for the ObservationPressure type
func (t ObservationPressure) String() string {
	return fmt.Sprintf("%.1fhPa", t.Value)
}
