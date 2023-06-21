// SPDX-FileCopyrightText: 2023 Winni Neessen <wn@neessen.dev>
//
// SPDX-License-Identifier: MIT

package meteologix

import (
	"fmt"
	"strings"
	"time"
)

// DataUnavailable is a constant string that is returned if a
// data point is not available
const DataUnavailable = "Data unavailable"

// DateFormat is the parsing format that is used for datetime strings
// that only hold the date but no time
const DateFormat = "2006-01-02"

// Enum for different Fieldname values
const (
	// FieldDewpoint represents the Dewpoint data point
	FieldDewpoint Fieldname = iota
	// FieldDewpointMean represents the TemperatureMean data point
	FieldDewpointMean
	// FieldGlobalRadiation10m represents the GlobalRadiation10m data point
	FieldGlobalRadiation10m
	// FieldGlobalRadiation1h represents the GlobalRadiation1h data point
	FieldGlobalRadiation1h
	// FieldGlobalRadiation24h represents the GlobalRadiation24h data point
	FieldGlobalRadiation24h
	// FieldHumidityRelative represents the HumidityRelative data point
	FieldHumidityRelative
	// FieldPrecipitation represents the Precipitation data point
	FieldPrecipitation
	// FieldPrecipitation10m represents the Precipitation10m data point
	FieldPrecipitation10m
	// FieldPrecipitation1h represents the Precipitation1h data point
	FieldPrecipitation1h
	// FieldPrecipitation24h represents the Precipitation24h data point
	FieldPrecipitation24h
	// FieldPressureMSL represents the PressureMSL data point
	FieldPressureMSL
	// FieldPressureQFE represents the PressureQFE data point
	FieldPressureQFE
	// FieldSunrise represents the Sunrise data point
	FieldSunrise
	// FieldSunset represents the Sunset data point
	FieldSunset
	// FieldTemperature represents the Temperature data point
	FieldTemperature
	// FieldTemperatureAtGround represents the TemperatureAtGround data point
	FieldTemperatureAtGround
	// FieldTemperatureAtGroundMin represents the TemperatureAtGroundMin data point
	FieldTemperatureAtGroundMin
	// FieldTemperatureMax represents the TemperatureMax data point
	FieldTemperatureMax
	// FieldTemperatureMean represents the TemperatureMean data point
	FieldTemperatureMean
	// FieldTemperatureMin represents the TemperatureMin data point
	FieldTemperatureMin
	// FieldWeatherSymbol represents the weather symbol data point
	FieldWeatherSymbol
	// FieldWinddirection represents the WindDirection data point
	FieldWinddirection
	// FieldWindspeed represents the WindSpeed data point
	FieldWindspeed
)

// Enum for different Timespan values
const (
	// TimespanCurrent represents the moment of the last observation
	TimespanCurrent Timespan = iota
	// Timespan10Min represents the last 10 minutes
	Timespan10Min
	// Timespan1Hour represents the last hour
	Timespan1Hour
	// Timespan24Hours represents the last 24 hours
	Timespan24Hours
)

// APIDate is type wrapper for datestamp (without time) returned by
// the API endpoints
type APIDate struct {
	time.Time
}

// APIBool is the JSON structure of the weather data that is
// returned by the API endpoints in which the value is a boolean
type APIBool struct {
	DateTime time.Time `json:"dateTime"`
	Source   *string   `json:"source,omitempty"`
	Value    bool      `json:"value"`
}

// APIFloat is the JSON structure of the weather data that is
// returned by the API endpoints in which the value is a float
type APIFloat struct {
	DateTime time.Time `json:"dateTime"`
	Source   *string   `json:"source,omitempty"`
	Value    float64   `json:"value"`
}

// APIString is the JSON structure of the weather data that is
// returned by the API endpoints in which the value is a string
type APIString struct {
	DateTime time.Time `json:"dateTime"`
	Source   *string   `json:"source,omitempty"`
	Value    string    `json:"value"`
}

// Timespan is a type wrapper for an int type
type Timespan int

// WeatherData is a type that holds weather (Observation, Current
// Weather) data and can be wrapped into other types to provide type
// specific receiver methods
type WeatherData struct {
	bv bool
	dt time.Time
	dv time.Time
	fv float64
	n  Fieldname
	na bool
	s  Source
	sv string
}

// Fieldname is a type wrapper for an int for field names
// of an Observation
type Fieldname int

// UnmarshalJSON interprets the API datestamp and converts it into a
// time.Time type
func (a *APIDate) UnmarshalJSON(s []byte) error {
	d := string(s)
	d = strings.ReplaceAll(d, `"`, ``)
	if d == "null" {
		return nil
	}

	pd, err := time.Parse(DateFormat, d)
	if err != nil {
		return fmt.Errorf("failed to parse JSON string as APIDate string: %w", err)
	}
	a.Time = pd
	return nil
}
