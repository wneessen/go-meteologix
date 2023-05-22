// SPDX-FileCopyrightText: 2023 Winni Neessen <wn@neessen.dev>
//
// SPDX-License-Identifier: MIT

package meteologix

import (
	"fmt"
	"math"
	"time"
)

// Pressure is a type wrapper of an WeatherData for holding pressure
// values in WeatherData
type Pressure WeatherData

// IsAvailable returns true if an Pressure value was
// available at time of query
func (t Pressure) IsAvailable() bool {
	return !t.na
}

// DateTime returns true if an Pressure value was
// available at time of query
func (t Pressure) DateTime() time.Time {
	return t.dt
}

// String satisfies the fmt.Stringer interface for the Pressure type
func (t Pressure) String() string {
	return fmt.Sprintf("%.1fhPa", t.v)
}

// Value returns the float64 value of an Pressure
// If the Pressure is not available in the Observation
// Vaule will return math.NaN instead.
func (t Pressure) Value() float64 {
	if t.na {
		return math.NaN()
	}
	return t.v
}
