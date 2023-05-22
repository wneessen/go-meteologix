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
func (t Radiation) IsAvailable() bool {
	return !t.na
}

// DateTime returns true if an Radiation value was
// available at time of query
func (t Radiation) DateTime() time.Time {
	return t.dt
}

// Value returns the float64 value of an Radiation
// If the Radiation is not available in the Observation
// Vaule will return math.NaN instead.
func (t Radiation) Value() float64 {
	if t.na {
		return math.NaN()
	}
	return t.v
}

// String satisfies the fmt.Stringer interface for the Radiation type
func (t Radiation) String() string {
	return fmt.Sprintf("%.0fkJ/m²", t.v)
}
