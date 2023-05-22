// SPDX-FileCopyrightText: 2023 Winni Neessen <wn@neessen.dev>
//
// SPDX-License-Identifier: MIT

package meteologix

import (
	"fmt"
	"math"
	"time"
)

// Humidity is a type wrapper of an WeatherData for holding humidity
// values in WeatherData
type Humidity WeatherData

// IsAvailable returns true if an Humidity value was
// available at time of query
func (t Humidity) IsAvailable() bool {
	return !t.na
}

// DateTime returns true if an Humidity value was
// available at time of query
func (t Humidity) DateTime() time.Time {
	return t.dt
}

// String satisfies the fmt.Stringer interface for the Humidity type
func (t Humidity) String() string {
	return fmt.Sprintf("%.1f%%", t.v)
}

// Value returns the float64 value of an Humidity
// If the Humidity is not available in the Observation
// Vaule will return math.NaN instead.
func (t Humidity) Value() float64 {
	if t.na {
		return math.NaN()
	}
	return t.v
}
