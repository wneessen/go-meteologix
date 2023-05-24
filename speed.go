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
func (s Speed) IsAvailable() bool {
	return !s.na
}

// DateTime returns true if an Speed value was
// available at time of query
func (s Speed) DateTime() time.Time {
	return s.dt
}

// Value returns the float64 value of an Speed in knots
// If the Speed is not available in the Observation
// Vaule will return math.NaN instead.
func (s Speed) Value() float64 {
	if s.na {
		return math.NaN()
	}
	return s.fv
}

// String satisfies the fmt.Stringer interface for the Speed type
func (s Speed) String() string {
	return fmt.Sprintf("%.0fkn", s.fv)
}

// Source returns the Source of Speed
// If the Source is not available it will return SourceUnknown
func (s Speed) Source() Source {
	return s.s
}

// KMH returns the Speed value in km/h
func (s Speed) KMH() float64 {
	return s.fv * 1.852
}

// KMHString returns the Speed value as formatted string in km/h
func (s Speed) KMHString() string {
	return fmt.Sprintf("%.1fkm/h", s.KMH())
}

// MPH returns the Speed value in mi/h
func (s Speed) MPH() float64 {
	return s.fv * 1.151
}

// MPHString returns the Speed value as formatted string in mi/h
func (s Speed) MPHString() string {
	return fmt.Sprintf("%.1fmi/h", s.MPH())
}
