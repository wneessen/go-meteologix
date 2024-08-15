// SPDX-FileCopyrightText: 2023 Winni Neessen <wn@neessen.dev>
//
// SPDX-License-Identifier: MIT

package meteologix

import (
	"fmt"
	"math"
	"time"
)

// Coverage is a type wrapper of WeatherData for holding coverage values in %
type Coverage WeatherData

// IsAvailable returns true if a Coverage value was available at time of query
func (c Coverage) IsAvailable() bool {
	return !c.notAvailable
}

// DateTime returns the DateTime of the queried Coverage value
func (c Coverage) DateTime() time.Time {
	return c.dateTime
}

// String satisfies the fmt.Stringer interface for the Coverage type
func (c Coverage) String() string {
	return fmt.Sprintf("%.0f%%", c.floatVal)
}

// Source returns the Source of Coverage
//
// If the Source is not available it will return SourceUnknown
func (c Coverage) Source() Source {
	return c.source
}

// Value returns the float64 value of a Coverage
//
// If the Coverage is not available in the WeatherData, Value will return math.NaN instead.
func (c Coverage) Value() float64 {
	if c.notAvailable {
		return math.NaN()
	}
	return c.floatVal
}

// Description returns a descriptive string value of a Coverage
//
// If the Coverage is not available in the WeatherData, Description will return Unknown instead.
func (c Coverage) Description() string {
	switch {
	case c.floatVal <= 10:
		return "Clear sky"
	case c.floatVal <= 30:
		return "Mostly clear"
	case c.floatVal <= 50:
		return "Partly cloudy"
	case c.floatVal <= 70:
		return "Mostly cloudy"
	case c.floatVal <= 90:
		return "Overcast"
	case c.floatVal <= 100:
		return "Very cloudy"
	default:
		return "Unknown"
	}
}
