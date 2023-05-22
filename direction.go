// SPDX-FileCopyrightText: 2023 Winni Neessen <wn@neessen.dev>
//
// SPDX-License-Identifier: MIT

package meteologix

import (
	"fmt"
	"math"
	"time"
)

// WindDirAbbrMap is a map to associate a wind direction degree value with
// the abbreviated direction string
var WindDirAbbrMap = map[float64]string{
	0: "N", 11.25: "NbE", 22.5: "NNE", 33.75: "NEbN", 45: "NE", 56.25: "NEbE",
	67.5: "ENE", 78.75: "EbN", 90: "E", 101.25: "EbS", 112.5: "ESE", 123.75: "SEbE",
	135: "SE", 146.25: "SEbS", 157.5: "SSE", 168.75: "SbE", 180: "S",
	191.25: "SbW", 202.5: "SSW", 213.75: "SWbS", 225: "SW", 236.25: "SWbW",
	247.5: "WSW", 258.75: "WbS", 270: "W", 281.25: "WbN", 292.5: "WNW",
	303.75: "NWbW", 315: "NW", 326.25: "NWbN", 337.5: "NNW", 348.75: "NbW",
}

// WindDirFullMap is a map to associate a wind direction degree value with
// the full direction string
var WindDirFullMap = map[float64]string{
	0: "North", 11.25: "North by East", 22.5: "North-Northeast",
	33.75: "Northeast by North", 45: "Northeast", 56.25: "Northeast by East",
	67.5: "East-Northeast", 78.75: "East by North", 90: "East",
	101.25: "East by South", 112.5: "East-Southeast", 123.75: "Southeast by East",
	135: "Southeast", 146.25: "Southeast by South", 157.5: "South-Southeast",
	168.75: "South by East", 180: "South", 191.25: "South by West",
	202.5: "South-Southwest", 213.75: "Southwest by South", 225: "Southwest",
	236.25: "Southwest by West", 247.5: "West-Southwest", 258.75: "West by South",
	270: "West", 281.25: "West by North", 292.5: "West-Northwest",
	303.75: "Northwest by West", 315: "Northwest", 326.25: "Northwest by North",
	337.5: "North-Northwest", 348.75: "North by West",
}

// Direction is a type wrapper of an WeatherData for holding directional
// values in WeatherData
type Direction WeatherData

// IsAvailable returns true if an Direction value was
// available at time of query
func (t Direction) IsAvailable() bool {
	return !t.na
}

// DateTime returns true if an Direction value was
// available at time of query
func (t Direction) DateTime() time.Time {
	return t.dt
}

// Value returns the float64 value of an Direction in degrees
// If the Direction is not available in the Observation
// Vaule will return math.NaN instead.
func (t Direction) Value() float64 {
	if t.na {
		return math.NaN()
	}
	return t.v
}

// String satisfies the fmt.Stringer interface for the Direction type
func (t Direction) String() string {
	return fmt.Sprintf("%.0f°", t.v)
}

// Direction returns the abbreviation string for a given Direction type
func (t Direction) Direction() string {
	if ds, ok := WindDirAbbrMap[t.v]; ok {
		return ds
	}
	return ErrUnsupportedDirection
}

// DirectionFull returns the full string for a given Direction type
func (t Direction) DirectionFull() string {
	if ds, ok := WindDirFullMap[t.v]; ok {
		return ds
	}
	return ErrUnsupportedDirection
}
