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
func (p Precipitation) IsAvailable() bool {
	return !p.notAvailable
}

// DateTime returns the DateTime when the Precipitation value was recorded
func (p Precipitation) DateTime() time.Time {
	return p.dt
}

// String satisfies the fmt.Stringer interface for the Precipitation type
func (p Precipitation) String() string {
	return fmt.Sprintf("%.1fmm", p.floatVal)
}

// Source returns the Source of Precipitation
// If the Source is not available it will return SourceUnknown
func (p Precipitation) Source() Source {
	return p.s
}

// Value returns the float64 value of an Precipitation
// If the Precipitation is not available in the WeatherData
// Vaule will return math.NaN instead.
func (p Precipitation) Value() float64 {
	if p.notAvailable {
		return math.NaN()
	}
	return p.floatVal
}
