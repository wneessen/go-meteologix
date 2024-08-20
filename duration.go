// SPDX-FileCopyrightText: 2023 Winni Neessen <wn@neessen.dev>
//
// SPDX-License-Identifier: MIT

package meteologix

import (
	"fmt"
	"math"
	"time"
)

// Duration is a type wrapper of an WeatherData for holding height values in WeatherData
// (based on meters a default unit)
type Duration WeatherData

// IsAvailable returns true if an Duration value was available at time of query
func (d Duration) IsAvailable() bool {
	return !d.notAvailable
}

// DateTime returns the timestamp associated with the Duration value
func (d Duration) DateTime() time.Time {
	return d.dateTime
}

// String satisfies the fmt.Stringer interface for the Duration type
func (d Duration) String() string {
	return fmt.Sprintf("%.2fh", d.floatVal)
}

// Source returns the Source of Duration
//
// If the Source is not available it will return SourceUnknown
func (d Duration) Source() Source {
	return d.source
}

// Value returns the float64 value of an Duration
//
// If the Duration is not available in the WeatherData, Value will return math.NaN instead.
func (d Duration) Value() float64 {
	if d.notAvailable {
		return math.NaN()
	}
	return d.floatVal
}
