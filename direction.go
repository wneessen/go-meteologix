// SPDX-FileCopyrightText: 2023 Winni Neessen <wn@neessen.dev>
//
// SPDX-License-Identifier: MIT

package meteologix

import (
	"fmt"
	"math"
	"sort"
	"time"
)

const (
	// DirectionMinAngel is the minimum angel for a direction
	DirectionMinAngel = 0
	// DirectionMaxAngel is the maximum angel for a direction
	DirectionMaxAngel = 360
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
func (d Direction) IsAvailable() bool {
	return !d.na
}

// DateTime returns true if an Direction value was
// available at time of query
func (d Direction) DateTime() time.Time {
	return d.dt
}

// Value returns the float64 value of an Direction in degrees
// If the Direction is not available in the WeatherData
// Vaule will return math.NaN instead.
func (d Direction) Value() float64 {
	if d.na {
		return math.NaN()
	}
	return d.fv
}

// String satisfies the fmt.Stringer interface for the Direction type
func (d Direction) String() string {
	return fmt.Sprintf("%.0f°", d.fv)
}

// Source returns the Source of a Direction
// If the Source is not available it will return SourceUnknown
func (d Direction) Source() Source {
	return d.s
}

// Direction returns the abbreviation string for a given Direction type
func (d Direction) Direction() string {
	if d.fv < DirectionMinAngel || d.fv > DirectionMaxAngel {
		return ErrUnsupportedDirection
	}
	if ds, ok := WindDirAbbrMap[d.fv]; ok {
		return ds
	}
	return findDirection(d.fv, WindDirAbbrMap)
}

// DirectionFull returns the full string for a given Direction type
func (d Direction) DirectionFull() string {
	if d.fv < DirectionMinAngel || d.fv > DirectionMaxAngel {
		return ErrUnsupportedDirection
	}
	if ds, ok := WindDirFullMap[d.fv]; ok {
		return ds
	}
	return findDirection(d.fv, WindDirFullMap)
}

// findDirection takes a Direction and tries to estimate the nearest
// direction string from a map
func findDirection(v float64, m map[float64]string) string {
	ks := make([]float64, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	sort.Float64s(ks)

	sv := 0.0
	ev := 0.0
	for i := range ks {
		if v > ks[i] {
			sv = ks[i]
			continue
		}
		if v < ks[i] {
			ev = ks[i]
			break
		}
	}
	sr := math.Mod(v, sv)
	er := math.Mod(ev, v)
	if er > sr {
		return m[sv]
	}
	return m[ev]
}
