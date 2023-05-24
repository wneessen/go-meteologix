package meteologix

import (
	"time"
)

// GenericString is a type wrapper of an WeatherData for holding
// a generic string value in the WeatherData
type GenericString WeatherData

// IsAvailable returns true if a GenericString value was available
// at time of query
func (gs GenericString) IsAvailable() bool {
	return !gs.na
}

// DateTime returns the timestamp of a GenericString value as time.Time
func (gs GenericString) DateTime() time.Time {
	return gs.dt
}

// Value returns the string value of a GenericString as simple
// unformatted string
// If the GenericSString is not available in the WeatherData
// Value will return DataUnavailable instead.
func (gs GenericString) Value() string {
	if gs.na {
		return DataUnavailable
	}
	return gs.sv
}

// String satisfies the fmt.Stringer interface for the GenericString type
func (gs GenericString) String() string {
	return gs.Value()
}

// Source returns the Source of a GenericString
// If the Source is not available it will return SourceUnknown
func (gs GenericString) Source() Source {
	return gs.s
}
