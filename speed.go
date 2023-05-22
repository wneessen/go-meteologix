// SPDX-FileCopyrightText: 2023 Winni Neessen <wn@neessen.dev>
//
// SPDX-License-Identifier: MIT

package meteologix

import (
	"fmt"
	"math"
	"time"
)

// Speed is a type wrapper of an WeatherData for holding speed
// values in WeatherData
type Speed WeatherData

// IsAvailable returns true if an Speed value was
// available at time of query
func (t Speed) IsAvailable() bool {
	return !t.na
}

// DateTime returns true if an Speed value was
// available at time of query
func (t Speed) DateTime() time.Time {
	return t.dt
}

// Value returns the float64 value of an Speed in knots
// If the Speed is not available in the Observation
// Vaule will return math.NaN instead.
func (t Speed) Value() float64 {
	if t.na {
		return math.NaN()
	}
	return t.v
}

// String satisfies the fmt.Stringer interface for the Speed type
func (t Speed) String() string {
	return fmt.Sprintf("%.0fkn", t.v)
}

// KMH returns the Speed value in km/h
func (t Speed) KMH() float64 {
	return t.v * 1.852
}

// KMHString returns the Speed value as formatted string in km/h
func (t Speed) KMHString() string {
	return fmt.Sprintf("%.1fkm/h", t.KMH())
}

// MPH returns the Speed value in mi/h
func (t Speed) MPH() float64 {
	return t.v * 1.151
}

// MPHString returns the Speed value as formatted string in mi/h
func (t Speed) MPHString() string {
	return fmt.Sprintf("%.1fmi/h", t.MPH())
}
