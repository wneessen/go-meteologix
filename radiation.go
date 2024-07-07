// SPDX-FileCopyrightText: 2023 Winni Neessen <wn@neessen.dev>
//
// SPDX-License-Identifier: MIT

package meteologix

import (
	"fmt"
	"math"
	"time"
)

// Radiation is a type wrapper of an WeatherData for holding radiation
// values in WeatherData
type Radiation WeatherData

// IsAvailable returns true if an Radiation value was
// available at time of query
func (r Radiation) IsAvailable() bool {
	return !r.notAvailable
}

// DateTime returns the time.Time object representing the date and time
// at which the Radiation value was queried
func (r Radiation) DateTime() time.Time {
	return r.dt
}

// Value returns the float64 value of an Radiation
// If the Radiation is not available in the WeatherData
// Vaule will return math.NaN instead.
func (r Radiation) Value() float64 {
	if r.notAvailable {
		return math.NaN()
	}
	return r.floatVal
}

// String satisfies the fmt.Stringer interface for the Radiation type
func (r Radiation) String() string {
	return fmt.Sprintf("%.0fkJ/m²", r.floatVal)
}

// Source returns the Source of Pressure
// If the Source is not available it will return SourceUnknown
func (r Radiation) Source() Source {
	return r.source
}
