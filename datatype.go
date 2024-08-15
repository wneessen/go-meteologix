// SPDX-FileCopyrightText: 2023 Winni Neessen <wn@neessen.dev>
//
// SPDX-License-Identifier: MIT

package meteologix

import (
	"fmt"
	"time"
)

// DataUnavailable is a constant string that is returned if a data point is not available
const DataUnavailable = "Data unavailable"

// DateFormat is the parsing format that is used for datetime strings that only hold
// the date but no time
const DateFormat = "2006-01-02"

// Enum for different Fieldname values
const (
	// FieldCloudCoverage represents the CloudCoverage data point
	FieldCloudCoverage Fieldname = iota
	// FieldDewpoint represents the Dewpoint data point
	FieldDewpoint
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
	// FieldSnowAmount represents the SnowAmount data point
	FieldSnowAmount
	// FieldSnowHeight represents the SnowHeight data point
	FieldSnowHeight
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
	// FieldWindDirection represents the WindDirection data point
	FieldWindDirection
	// FieldWindGust represents the WindGust data point
	FieldWindGust
	// FieldWindGust3h represents the WindGust (over the last 3h) data point
	FieldWindGust3h
	// FieldWindSpeed represents the WindSpeed data point
	FieldWindSpeed
)

// Enum for different Timespan values
const (
	// TimespanCurrent represents the moment of the last observation
	TimespanCurrent Timespan = iota
	// Timespan10Min represents the last 10 minutes
	Timespan10Min
	// Timespan1Hour represents the 1 hour
	Timespan1Hour
	// Timespan3Hours represents the 3 hours
	Timespan3Hours
	// Timespan6Hours represents the 6 hours
	Timespan6Hours
	// Timespan24Hours represents the last 24 hours
	Timespan24Hours
)

// APIDate is type wrapper for datestamp (without time) returned by the API endpoints
type APIDate struct {
	time.Time
}

// APIBool is the JSON structure of the weather data that is returned by the API endpoints
// in which the value is a boolean
type APIBool struct {
	DateTime time.Time `json:"dateTime"`
	Source   *string   `json:"source,omitempty"`
	Value    bool      `json:"value"`
}

// APIFloat is the JSON structure of the weather data that is returned by the API endpoints
// in which the value is a float
type APIFloat struct {
	DateTime time.Time `json:"dateTime"`
	Source   *string   `json:"source,omitempty"`
	Value    float64   `json:"value"`
}

// APIString is the JSON structure of the weather data that is returned by the API endpoints
// in which the value is a string
type APIString struct {
	DateTime time.Time `json:"dateTime"`
	Source   *string   `json:"source,omitempty"`
	Value    string    `json:"value"`
}

// Timespan is a type wrapper for an int type
type Timespan int

// WeatherData is a type that holds weather (Observation, Current Weather) data and can be wrapped
// into other types to provide type specific receiver methods
type WeatherData struct {
	// bv bool
	dateTime     time.Time
	dateTimeVal  time.Time
	floatVal     float64
	name         Fieldname
	notAvailable bool
	source       Source
	stringVal    string
}

// Fieldname is a type wrapper for an int for field names of an Observation
type Fieldname int

// String returns a string representation of the Timespan value and satisfies the fmt.Stringer interface.
func (t Timespan) String() string {
	switch t {
	case TimespanCurrent:
		return "current"
	case Timespan10Min:
		return "10m"
	case Timespan1Hour:
		return "1h"
	case Timespan3Hours:
		return "3h"
	case Timespan6Hours:
		return "6h"
	case Timespan24Hours:
		return "24h"
	default:
		return "unknown"
	}
}

// UnmarshalJSON interprets the API datestamp and converts it into a time.Time type
func (a *APIDate) UnmarshalJSON(data []byte) error {
	date := string(data)
	if date == "null" {
		return nil
	}
	date = date[1 : len(date)-1]
	parsedDate, err := time.Parse(DateFormat, date)
	if err != nil {
		return fmt.Errorf("failed to parse JSON string as APIDate string: %w", err)
	}
	a.Time = parsedDate
	return nil
}
