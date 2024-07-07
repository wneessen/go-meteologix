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
func (p Pressure) IsAvailable() bool {
	return !p.notAvailable
}

// DateTime returns the date and time of the Pressure reading
func (p Pressure) DateTime() time.Time {
	return p.dt
}

// String satisfies the fmt.Stringer interface for the Pressure type
func (p Pressure) String() string {
	return fmt.Sprintf("%.1fhPa", p.floatVal)
}

// Source returns the Source of Pressure
// If the Source is not available it will return SourceUnknown
func (p Pressure) Source() Source {
	return p.s
}

// Value returns the float64 value of an Pressure
// If the Pressure is not available in the WeatherData
// Vaule will return math.NaN instead.
func (p Pressure) Value() float64 {
	if p.notAvailable {
		return math.NaN()
	}
	return p.floatVal
}
