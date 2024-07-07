// SPDX-FileCopyrightText: 2023 Winni Neessen <wn@neessen.dev>
//
// SPDX-License-Identifier: MIT

package meteologix

import (
	"fmt"
	"math"
	"time"
)

const (
	// MultiplierKnots is the multiplier for converting the base unit to knots
	MultiplierKnots = 1.9438444924
	// MultiplierKPH is the multiplier for converting the base unit to kilometers per hour
	MultiplierKPH = 3.6
	// MultiplierMPH is the multiplier for converting the base unit to miles per hour
	MultiplierMPH = 2.236936
)

// Speed is a type wrapper of an WeatherData for holding speed
// values in WeatherData
type Speed WeatherData

// IsAvailable returns true if an Speed value was
// available at time of query
func (s Speed) IsAvailable() bool {
	return !s.notAvailable
}

// DateTime returns the DateTime when the Speed was checked
func (s Speed) DateTime() time.Time {
	return s.dt
}

// Value returns the float64 value of an Speed in meters
// per second.
// If the Speed is not available in the WeatherData
// Vaule will return math.NaN instead.
func (s Speed) Value() float64 {
	if s.notAvailable {
		return math.NaN()
	}
	return s.floatVal
}

// String satisfies the fmt.Stringer interface for the Speed type
func (s Speed) String() string {
	return fmt.Sprintf("%.1fm/s", s.floatVal)
}

// Source returns the Source of Speed
// If the Source is not available it will return SourceUnknown
func (s Speed) Source() Source {
	return s.source
}

// KMH returns the Speed value in km/h
func (s Speed) KMH() float64 {
	return s.floatVal * MultiplierKPH
}

// KMHString returns the Speed value as formatted string in km/h
func (s Speed) KMHString() string {
	return fmt.Sprintf("%.1fkm/h", s.KMH())
}

// Knots returns the Speed value in kn
func (s Speed) Knots() float64 {
	return s.floatVal * MultiplierKnots
}

// KnotsString returns the Speed value as formatted string in kn
func (s Speed) KnotsString() string {
	return fmt.Sprintf("%.0fkn", s.Knots())
}

// MPH returns the Speed value in mi/h
func (s Speed) MPH() float64 {
	return s.floatVal * MultiplierMPH
}

// MPHString returns the Speed value as formatted string in mi/h
func (s Speed) MPHString() string {
	return fmt.Sprintf("%.1fmi/h", s.MPH())
}
