// SPDX-FileCopyrightText: 2023 Winni Neessen <wn@neessen.dev>
//
// SPDX-License-Identifier: MIT

package meteologix

import (
	"fmt"
	"math"
	"time"
)

// Height is a type wrapper of an WeatherData for holding height
// values in WeatherData (based on meters a default unit)
type Height WeatherData

// IsAvailable returns true if an Height value was
// available at time of query
func (h Height) IsAvailable() bool {
	return !h.na
}

// DateTime returns true if an Height value was
// available at time of query
func (h Height) DateTime() time.Time {
	return h.dt
}

// String satisfies the fmt.Stringer interface for the Height type
func (h Height) String() string {
	return fmt.Sprintf("%.3fm", h.fv)
}

// Source returns the Source of Height
// If the Source is not available it will return SourceUnknown
func (h Height) Source() Source {
	return h.s
}

// Value returns the float64 value of an Height
// If the Height is not available in the WeatherData
// Vaule will return math.NaN instead.
func (h Height) Value() float64 {
	if h.na {
		return math.NaN()
	}
	return h.fv
}

// Meter returns the Height type value as float64 in meters.
// This is an alias for the Value() method
func (h Height) Meter() float64 {
	return h.Value()
}

// MeterString returns the Height type as formatted string in meters
// This is an alias for the String() method
func (h Height) MeterString() string {
	return h.String()
}

// CentiMeter returns the Height type value as float64 in centimeters.
func (h Height) CentiMeter() float64 {
	if h.na {
		return math.NaN()
	}
	return h.fv / 100
}

// CentiMeterString returns the Height type as formatted string in centimeters
func (h Height) CentiMeterString() string {
	return fmt.Sprintf("%.3fcm", h.CentiMeter())
}

// MilliMeter returns the Height type value as float64 in milliimeters.
func (h Height) MilliMeter() float64 {
	if h.na {
		return math.NaN()
	}
	return h.fv / 1000
}

// MilliMeterString returns the Height type as formatted string in millimeters
func (h Height) MilliMeterString() string {
	return fmt.Sprintf("%.3fmm", h.MilliMeter())
}
