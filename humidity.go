// SPDX-FileCopyrightText: 2023 Winni Neessen <wn@neessen.dev>
//
// SPDX-License-Identifier: MIT

package meteologix

import (
	"fmt"
	"math"
	"time"
)

// Humidity is a type wrapper of an WeatherData for holding humidity
// values in WeatherData
type Humidity WeatherData

// IsAvailable returns true if an Humidity value was
// available at time of query
func (h Humidity) IsAvailable() bool {
	return !h.notAvailable
}

// DateTime returns the timestamp of when the humidity
// measurement was taken.
func (h Humidity) DateTime() time.Time {
	return h.dt
}

// String satisfies the fmt.Stringer interface for the Humidity type
func (h Humidity) String() string {
	return fmt.Sprintf("%.1f%%", h.floatVal)
}

// Source returns the Source of Humidity
// If the Source is not available it will return SourceUnknown
func (h Humidity) Source() Source {
	return h.s
}

// Value returns the float64 value of an Humidity
// If the Humidity is not available in the WeatherData
// Value will return math.NaN instead.
func (h Humidity) Value() float64 {
	if h.notAvailable {
		return math.NaN()
	}
	return h.floatVal
}
