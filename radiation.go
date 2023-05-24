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
	return !r.na
}

// DateTime returns true if an Radiation value was
// available at time of query
func (r Radiation) DateTime() time.Time {
	return r.dt
}

// Value returns the float64 value of an Radiation
// If the Radiation is not available in the Observation
// Vaule will return math.NaN instead.
func (r Radiation) Value() float64 {
	if r.na {
		return math.NaN()
	}
	return r.fv
}

// String satisfies the fmt.Stringer interface for the Radiation type
func (r Radiation) String() string {
	return fmt.Sprintf("%.0fkJ/m²", r.fv)
}

// Source returns the Source of Pressure
// If the Source is not available it will return SourceUnknown
func (r Radiation) Source() Source {
	return r.s
}
