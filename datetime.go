// SPDX-FileCopyrightText: 2023 Winni Neessen <wn@neessen.dev>
//
// SPDX-License-Identifier: MIT

package meteologix

import (
	"time"
)

// DateTime is a type wrapper of an WeatherData for holding datetime
// values in WeatherData
type DateTime WeatherData

// IsAvailable returns true if an Direction value was
// available at time of query
func (dt DateTime) IsAvailable() bool {
	return !dt.na
}

// DateTime returns the timestamp for that specific DateTime type
func (dt DateTime) DateTime() time.Time {
	return dt.dt
}

// Value returns the time.Time value of an DateTime value
// If the DateTime is not available in the WeatherData
// Value will return time.Time with a zero value instead.
func (dt DateTime) Value() time.Time {
	if dt.na {
		return time.Time{}
	}
	return dt.dv
}

// String satisfies the fmt.Stringer interface for the DateTime type
// The date will returned as time.RFC3339 format
func (dt DateTime) String() string {
	return dt.Value().Format(time.RFC3339)
}
