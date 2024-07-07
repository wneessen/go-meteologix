// SPDX-FileCopyrightText: 2023 Winni Neessen <wn@neessen.dev>
//
// SPDX-License-Identifier: MIT

package meteologix

import (
	"time"
)

const (
	// CondCloudy represents cloudy weather conditions
	CondCloudy ConditionType = "cloudy"
	// CondFog represents foggy weather conditions
	CondFog ConditionType = "fog"
	// CondFreezingRain represents weather conditions with freezing rain
	CondFreezingRain ConditionType = "freezingrain"
	// CondOvercast represents overcast weather conditions
	CondOvercast ConditionType = "overcast"
	// CondPartlyCloudy represents partly cloudy weather conditions
	CondPartlyCloudy ConditionType = "partlycloudy"
	// CondRain represents rainy weather conditions.
	// Rain defines as following:
	// - Falls steadily
	// - Lasts for hours or days
	// - Typically widespread throughout your city or town
	CondRain ConditionType = "rain"
	// CondRainHeavy represents heavy rain weather conditions
	CondRainHeavy ConditionType = "rainheavy"
	// CondShowers represents weather conditions with showers.
	// Showers define as following:
	// - Lighter rainfall
	// - Shorter duration
	// - Can start and stop over a period of time
	// - Tends to be more scattered across an area
	CondShowers ConditionType = "showers"
	// CondShowersHeavy represents weather conditions with heavy showers
	CondShowersHeavy ConditionType = "showersheavy"
	// CondSnow represents snowy weather conditions
	CondSnow ConditionType = "snow"
	// CondSnowHeavy represents weather conditions with heavy snow
	CondSnowHeavy ConditionType = "snowheavy"
	// CondSnowRain represents weather conditions with snowy rain
	CondSnowRain ConditionType = "snowrain"
	// CondSunshine represents clear and sunny weather conditions
	CondSunshine ConditionType = "sunshine"
	// CondThunderStorm represents weather conditions with thunderstorms
	CondThunderStorm ConditionType = "thunderstorm"
	// CondUnknown represents a unknown weather condition
	CondUnknown ConditionType = "unknown"
)

// ConditionMap is a map to associate a specific ConditionType to a nicely
// formatted, human readable string
var ConditionMap = map[ConditionType]string{
	CondCloudy:       "Cloudy",
	CondFog:          "Fog",
	CondFreezingRain: "Freezing rain",
	CondOvercast:     "Overcast",
	CondPartlyCloudy: "Partly cloudy",
	CondRain:         "Rain",
	CondRainHeavy:    "Heavy rain",
	CondShowers:      "Showers",
	CondShowersHeavy: "Heavy showers",
	CondSnow:         "Snow",
	CondSnowHeavy:    "Heavy snow",
	CondSnowRain:     "Sleet",
	CondSunshine:     "Clear sky",
	CondThunderStorm: "Thunderstorm",
	CondUnknown:      "Unknown",
}

// Condition is a type wrapper of an WeatherData for holding
// a specific weather Condition value in the WeatherData
type Condition WeatherData

// ConditionType is a type wrapper for a string type
type ConditionType string

// IsAvailable returns true if a Condition value was available
// at time of query
func (c Condition) IsAvailable() bool {
	return !c.notAvailable
}

// DateTime returns the timestamp of a Condition value as time.Time
func (c Condition) DateTime() time.Time {
	return c.dt
}

// Value returns the raw value of a Condition as unformatted string
// as returned by the API
// If the Condition is not available in the WeatherData, Value will
// return DataUnavailable instead.
func (c Condition) Value() string {
	if c.notAvailable {
		return DataUnavailable
	}
	return c.sv
}

// Condition returns the actual value of that Condition as ConditionType.
// If the value is not available or not supported it will return a
// CondUnknown
func (c Condition) Condition() ConditionType {
	if c.notAvailable {
		return CondUnknown
	}
	if _, ok := ConditionMap[ConditionType(c.sv)]; ok {
		return ConditionType(c.sv)
	}
	return CondUnknown
}

// String returns the formatted, human readable string for a given
// Condition type and satisfies the fmt.Stringer interface
func (c Condition) String() string {
	return c.Condition().String()
}

// Source returns the Source of a Condition
// If the Source is not available it will return SourceUnknown
func (c Condition) Source() Source {
	return c.source
}

// String returns a human readable, formatted string for a ConditionType and
// satisfies the fmt.Stringer interface.
func (ct ConditionType) String() string {
	if cs, ok := ConditionMap[ct]; ok {
		return cs
	}
	return ConditionMap[CondUnknown]
}
