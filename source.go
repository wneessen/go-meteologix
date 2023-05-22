// SPDX-FileCopyrightText: 2023 Winni Neessen <wn@neessen.dev>
//
// SPDX-License-Identifier: MIT

package meteologix

import "strings"

// Enum of different weather data sources
const (
	// SourceObservation represent observations from weather stations (high precision)
	SourceObservation = iota
	// SourceAnalysis represents weather data based on analysis (medium precision)
	SourceAnalysis
	// SourceForecast represents weather data based on weather forcecasts
	SourceForecast
	// SourceMixed represents weather data based on mixed sources
	SourceMixed
	// SourceUnknown represents weather data based on unknown sources
	SourceUnknown
)

// Source is a type wrapper for an int type to enum different weather sources
type Source int

// String satisfies the fmt.Stringer interface for the Source type
func (s Source) String() string {
	switch s {
	case SourceObservation:
		return "Observation"
	case SourceAnalysis:
		return "Analysis"
	case SourceForecast:
		return "Forecast"
	case SourceMixed:
		return "Mixed"
	case SourceUnknown:
		return "Unknown"
	default:
		return "Unknown"
	}
}

// StringToSource converts a given source string to a Source type
func StringToSource(s string) Source {
	switch strings.ToLower(s) {
	case "observation":
		return SourceObservation
	case "analysis":
		return SourceAnalysis
	case "forecast":
		return SourceForecast
	case "mixed":
		return SourceMixed
	default:
		return SourceUnknown
	}
}
