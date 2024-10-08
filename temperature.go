// SPDX-FileCopyrightText: 2023 Winni Neessen <wn@neessen.dev>
//
// SPDX-License-Identifier: MIT

package meteologix

import (
	"fmt"
	"math"
	"time"
)

// Temperature is a type wrapper of an WeatherData for holding temperature values in WeatherData
type Temperature WeatherData

// IsAvailable returns true if an Temperature value was available at time of query
func (t Temperature) IsAvailable() bool {
	return !t.notAvailable
}

// DateTime returns the time at which the temperature data was generated or fetched
func (t Temperature) DateTime() time.Time {
	return t.dateTime
}

// Value returns the float64 value of an Temperature
//
// If the Temperature is not available in the WeatherData, Value will return math.NaN instead.
func (t Temperature) Value() float64 {
	if t.notAvailable {
		return math.NaN()
	}
	return t.floatVal
}

// Source returns the Source of an Temperature
//
// If the Source is not available it will return SourceUnknown
func (t Temperature) Source() Source {
	return t.source
}

// String satisfies the fmt.Stringer interface for the Temperature type
func (t Temperature) String() string {
	return fmt.Sprintf("%.1f°C", t.floatVal)
}

// Celsius returns the Temperature value in Celsius
func (t Temperature) Celsius() float64 {
	return t.floatVal
}

// CelsiusString returns the Temperature value as Celsius formated string.
//
// This is an alias for the fmt.Stringer interface
func (t Temperature) CelsiusString() string {
	return t.String()
}

// Fahrenheit returns the Temperature value in Fahrenheit
func (t Temperature) Fahrenheit() float64 {
	return t.floatVal*9/5 + 32
}

// FahrenheitString returns the Temperature value as Fahrenheit formated string.
func (t Temperature) FahrenheitString() string {
	return fmt.Sprintf("%.1f°F", t.Fahrenheit())
}
