// SPDX-FileCopyrightText: 2023 Winni Neessen <wn@neessen.dev>
//
// SPDX-License-Identifier: MIT

package meteologix

import (
	"fmt"
	"math"
	"time"
)

// Density is a type wrapper of WeatherData for holding density
// values in kg/m³ in WeatherData
type Density WeatherData

// IsAvailable returns true if an Density value was
// available at time of query
func (d Density) IsAvailable() bool {
	return !d.notAvailable
}

// DateTime returns the DateTime of the queried Density value
func (d Density) DateTime() time.Time {
	return d.dateTime
}

// String satisfies the fmt.Stringer interface for the Density type
func (d Density) String() string {
	return fmt.Sprintf("%.1fkg/m³", d.floatVal)
}

// Source returns the Source of Density
// If the Source is not available it will return SourceUnknown
func (d Density) Source() Source {
	return d.source
}

// Value returns the float64 value of an Density
// If the Density is not available in the WeatherData
// Vaule will return math.NaN instead.
func (d Density) Value() float64 {
	if d.notAvailable {
		return math.NaN()
	}
	return d.floatVal
}
