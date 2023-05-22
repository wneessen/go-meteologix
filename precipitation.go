// SPDX-FileCopyrightText: 2023 Winni Neessen <wn@neessen.dev>
//
// SPDX-License-Identifier: MIT

package meteologix

import (
	"fmt"
	"math"
	"time"
)

// Precipitation is a type wrapper of an WeatherData for holding precipitation
// values in WeatherData
type Precipitation WeatherData

// IsAvailable returns true if an Precipitation value was
// available at time of query
func (t Precipitation) IsAvailable() bool {
	return !t.na
}

// DateTime returns true if an Precipitation value was
// available at time of query
func (t Precipitation) DateTime() time.Time {
	return t.dt
}

// String satisfies the fmt.Stringer interface for the Precipitation type
func (t Precipitation) String() string {
	return fmt.Sprintf("%.1fmm", t.v)
}

// Value returns the float64 value of an Precipitation
// If the Precipitation is not available in the Observation
// Vaule will return math.NaN instead.
func (t Precipitation) Value() float64 {
	if t.na {
		return math.NaN()
	}
	return t.v
}
