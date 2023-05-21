// SPDX-FileCopyrightText: 2023 Winni Neessen <wn@neessen.dev>
//
// SPDX-License-Identifier: MIT

package meteologix

import (
	"encoding/json"
	"fmt"
	"math"
	"time"
)

const (
	// FieldDewpoint represents the Dewpoint data point
	FieldDewpoint ObservationFieldName = iota
	// FieldDewpointMean represents the TemperatureMean data point
	FieldDewpointMean
	// FieldGlobalRadiation10m represents the GlobalRadiation10m data point
	FieldGlobalRadiation10m
	// FieldGlobalRadiation1h represents the GlobalRadiation1h data point
	FieldGlobalRadiation1h
	// FieldGlobalRadiation24h represents the GlobalRadiation24h data point
	FieldGlobalRadiation24h
	// FieldHumidityRelative represents the HumidityRelative data point
	FieldHumidityRelative
	// FieldPrecipitation represents the Precipitation data point
	FieldPrecipitation
	// FieldPrecipitation10m represents the Precipitation10m data point
	FieldPrecipitation10m
	// FieldPrecipitation1h represents the Precipitation1h data point
	FieldPrecipitation1h
	// FieldPrecipitation24h represents the Precipitation24h data point
	FieldPrecipitation24h
	// FieldPressureMSL represents the PressureMSL data point
	FieldPressureMSL
	// FieldPressureQFE represents the PressureQFE data point
	FieldPressureQFE
	// FieldTemperature represents the Temperature data point
	FieldTemperature
	// FieldTemperatureAtGround represents the TemperatureAtGround data point
	FieldTemperatureAtGround
	// FieldTemperatureAtGroundMin represents the TemperatureAtGroundMin data point
	FieldTemperatureAtGroundMin
	// FieldTemperatureMax represents the TemperatureMax data point
	FieldTemperatureMax
	// FieldTemperatureMean represents the TemperatureMean data point
	FieldTemperatureMean
	// FieldTemperatureMin represents the TemperatureMin data point
	FieldTemperatureMin
	// FieldWinddirection represents the Winddirection data point
	FieldWinddirection
	// FieldWindspeed represents the Windspeed data point
	FieldWindspeed
)

const (
	// TimespanCurrent represents the moment of the last observation
	TimespanCurrent Timespan = iota
	// Timespan10Min represents the last 10 minutes
	Timespan10Min
	// Timespan1Hour represents the last hour
	Timespan1Hour
	// Timespan24Hours represents the last 24 hours
	Timespan24Hours
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

// ErrUnsupportedDirection is returned when a direction degree is given,
// that is not resolvable
var ErrUnsupportedDirection = "Unsupported direction"

// Observation represents the observation API response for a Station
type Observation struct {
	// Altitude is the altitude of the station providing the Observation
	Altitude *int `json:"ele,omitempty"`
	// Data holds the different ObservationData points
	Data ObservationData `json:"data"`
	// Name is the name of the Station providing the Observation
	Name string `json:"name"`
	// Latitude represents the GeoLocation latitude coordinates for the Station
	Latitude float64 `json:"lat"`
	// Longitude represents the GeoLocation longitude coordinates for the Station
	Longitude float64 `json:"lon"`
	// StationID is the ID of the Station providing the Observation
	StationID string `json:"stationId"`
}

// ObservationData holds the different data points of the Observation.
//
// Please keep in mind that different Station types return different values, therefore
// all values are represented as pointer type returning nil if the data point in question
// is not returned for the requested Station.
type ObservationData struct {
	// Dewpoint represents the dewpoint in °C
	Dewpoint *ObservationValue `json:"dewpoint,omitempty"`
	// DewPointMean represents the mean dewpoint in °C
	DewpointMean *ObservationValue `json:"dewpointMean,omitempty"`
	// GlobalRadiation10m represents the sum of global radiation over the last
	// 10 minutes in kJ/m²
	GlobalRadiation10m *ObservationValue `json:"globalRadiation10m,omitempty"`
	// GlobalRadiation1h represents the sum of global radiation over the last
	// 1 hour in kJ/m²
	GlobalRadiation1h *ObservationValue `json:"globalRadiation1h,omitempty"`
	// GlobalRadiation24h represents the sum of global radiation over the last
	// 24 hour in kJ/m²
	GlobalRadiation24h *ObservationValue `json:"globalRadiation24h,omitempty"`
	// HumidityRelative represents the relative humidity in percent
	HumidityRelative *ObservationValue `json:"humidityRelative,omitempty"`
	// Precipitation represents the current amount of precipitation
	Precipitation *ObservationValue `json:"prec"`
	// Precipitation10m represents the amount of precipitation over the last 10 minutes
	Precipitation10m *ObservationValue `json:"prec10m"`
	// Precipitation1h represents the amount of precipitation over the last hour
	Precipitation1h *ObservationValue `json:"prec1h"`
	// Precipitation24h represents the amount of precipitation over the last 24 hours
	Precipitation24h *ObservationValue `json:"prec24h"`
	// PressureMSL represents the pressure at mean sea level (MSL) in hPa
	PressureMSL *ObservationValue `json:"pressureMsl"`
	// PressureMSL represents the pressure at station level (QFE) in hPa
	PressureQFE *ObservationValue `json:"pressure"`
	// Temperature represents the temperature in °C
	Temperature *ObservationValue `json:"temp,omitempty"`
	// TemperatureMax represents the maximum temperature in °C
	TemperatureMax *ObservationValue `json:"tempMax,omitempty"`
	// TemperatureMean represents the mean temperature in °C
	TemperatureMean *ObservationValue `json:"tempMean,omitempty"`
	// TemperatureMin represents the minimum temperature in °C
	TemperatureMin *ObservationValue `json:"tempMin,omitempty"`
	// Temperature5cm represents the temperature 5cm above ground in °C
	Temperature5cm *ObservationValue `json:"temp5cm,omitempty"`
	// Temperature5cm represents the minimum temperature 5cm above
	// ground in °C
	Temperature5cmMin *ObservationValue `json:"temp5cmMin,omitempty"`
	// Winddirection represents the direction from which the wind
	// originates in degree (0=N, 90=E, 180=S, 270=W)
	Winddirection *ObservationValue `json:"windDirection,omitempty"`
	// Windspeed represents the wind speed in knots
	Windspeed *ObservationValue `json:"windSpeed,omitempty"`
}

// ObservationValue is the JSON structure of the Observation data that is
// returned by the API endpoints
type ObservationValue struct {
	DateTime time.Time `json:"dateTime"`
	Value    float64   `json:"value"`
}

// ObservationField is a type that holds Observation data and can be wrapped
// into other types to provide type specific receiver methods
type ObservationField struct {
	dt time.Time
	n  ObservationFieldName
	na bool
	v  float64
}

// ObservationFieldName is a type wrapper for an int for field names
// of an Observation
type ObservationFieldName int

// ObservationDirection is a type wrapper of an ObservationField for
// holding directional values
type ObservationDirection ObservationField

// ObservationHumidity is a type wrapper of an ObservationField for
// holding humidity values
type ObservationHumidity ObservationField

// ObservationPrecipitation is a type wrapper for a precipitation value
// in an Observation
type ObservationPrecipitation ObservationField

// ObservationPressure is a type wrapper for a pressure value
// in an Observation
type ObservationPressure ObservationField

// ObservationRadiation is a type wrapper of an ObservationField for
// holding radiation values
type ObservationRadiation ObservationField

// ObservationSpeed is a type wrapper of an ObservationField for
// holding speed values
type ObservationSpeed ObservationField

// ObservationTemperature is a type wrapper of an ObservationField for
// holding temperature values
type ObservationTemperature ObservationField

// Timespan is a type wrapper for an int type
type Timespan int

// ObservationLatestByStationID returns the latest Observation values from the
// given Station
func (c *Client) ObservationLatestByStationID(si string) (Observation, error) {
	var o Observation
	u := fmt.Sprintf("%s/station/%s/observations/latest", c.config.apiURL, si)
	r, err := c.httpClient.Get(u)
	if err != nil {
		return o, fmt.Errorf("API request failed: %w", err)
	}

	if err := json.Unmarshal(r, &o); err != nil {
		return o, fmt.Errorf("failed to unmarshal API response JSON: %w", err)
	}

	return o, nil
}

// ObservationLatestByLocation performs a GeoLocation lookup of the location string, checks for any
// nearby weather stations (25 km radius) and returns the latest Observation values from the
// Stations with the shortest distance. It will also return the Station that was used for the query.
// It will throw an error if no station could be found in that queried location.
func (c *Client) ObservationLatestByLocation(l string) (Observation, Station, error) {
	sl, err := c.StationSearchByLocationWithinRadius(l, 25)
	if err != nil {
		return Observation{}, Station{}, fmt.Errorf("failed search locations at given location: %w", err)
	}
	s := sl[0]
	o, err := c.ObservationLatestByStationID(s.ID)
	return o, s, err
}

// Dewpoint returns the dewpoint data point as ObservationTemperature
// If the data point is not available in the Observation it will return
// ObservationTemperature in which the "not available" field will be
// true.
func (o Observation) Dewpoint() ObservationTemperature {
	if o.Data.Dewpoint == nil {
		return ObservationTemperature{na: true}
	}
	return ObservationTemperature{
		dt: o.Data.Dewpoint.DateTime,
		n:  FieldDewpoint,
		v:  o.Data.Dewpoint.Value,
	}
}

// DewpointMean returns the mean dewpoint data point as ObservationTemperature.
// If the data point is not available in the Observation it will return
// ObservationTemperature in which the "not available" field will be
// true.
func (o Observation) DewpointMean() ObservationTemperature {
	if o.Data.DewpointMean == nil {
		return ObservationTemperature{na: true}
	}
	return ObservationTemperature{
		dt: o.Data.DewpointMean.DateTime,
		n:  FieldDewpointMean,
		v:  o.Data.DewpointMean.Value,
	}
}

// Temperature returns the temperature data point as ObservationTemperature.
// If the data point is not available in the Observation it will return
// ObservationTemperature in which the "not available" field will be
// true.
func (o Observation) Temperature() ObservationTemperature {
	if o.Data.Temperature == nil {
		return ObservationTemperature{na: true}
	}
	return ObservationTemperature{
		dt: o.Data.Temperature.DateTime,
		n:  FieldTemperature,
		v:  o.Data.Temperature.Value,
	}
}

// TemperatureAtGround returns the temperature at ground level (5cm)
// data point as ObservationTemperature.
// If the data point is not available in the Observation it will return
// ObservationTemperature in which the "not available" field will be
// true.
func (o Observation) TemperatureAtGround() ObservationTemperature {
	if o.Data.Temperature5cm == nil {
		return ObservationTemperature{na: true}
	}
	return ObservationTemperature{
		dt: o.Data.Temperature5cm.DateTime,
		n:  FieldTemperatureAtGround,
		v:  o.Data.Temperature5cm.Value,
	}
}

// TemperatureMax returns the maximum temperature so far data point as
// ObservationTemperature.
// If the data point is not available in the Observation it will return
// ObservationTemperature in which the "not available" field will be
// true.
func (o Observation) TemperatureMax() ObservationTemperature {
	if o.Data.TemperatureMax == nil {
		return ObservationTemperature{na: true}
	}
	return ObservationTemperature{
		dt: o.Data.TemperatureMax.DateTime,
		n:  FieldTemperatureMax,
		v:  o.Data.TemperatureMax.Value,
	}
}

// TemperatureMin returns the minimum temperature so far data point as
// ObservationTemperature.
// If the data point is not available in the Observation it will return
// ObservationTemperature in which the "not available" field will be
// true.
func (o Observation) TemperatureMin() ObservationTemperature {
	if o.Data.TemperatureMin == nil {
		return ObservationTemperature{na: true}
	}
	return ObservationTemperature{
		dt: o.Data.TemperatureMin.DateTime,
		n:  FieldTemperatureMin,
		v:  o.Data.TemperatureMin.Value,
	}
}

// TemperatureAtGroundMin returns the minimum temperature so far
// at ground level (5cm) data point as ObservationTemperature
// If the data point is not available in the Observation it will return
// ObservationTemperature in which the "not available" field will be
// true.
func (o Observation) TemperatureAtGroundMin() ObservationTemperature {
	if o.Data.Temperature5cmMin == nil {
		return ObservationTemperature{na: true}
	}
	return ObservationTemperature{
		dt: o.Data.Temperature5cmMin.DateTime,
		n:  FieldTemperatureAtGroundMin,
		v:  o.Data.Temperature5cmMin.Value,
	}
}

// TemperatureMean returns the mean temperature data point as ObservationTemperature.
// If the data point is not available in the Observation it will return
// ObservationTemperature in which the "not available" field will be
// true.
func (o Observation) TemperatureMean() ObservationTemperature {
	if o.Data.TemperatureMean == nil {
		return ObservationTemperature{na: true}
	}
	return ObservationTemperature{
		dt: o.Data.TemperatureMean.DateTime,
		n:  FieldTemperatureMean,
		v:  o.Data.TemperatureMean.Value,
	}
}

// HumidityRelative returns the relative humidity data point as float64.
// If the data point is not available in the Observation it will return
// ObservationHumidity in which the "not available" field will be
// true.
func (o Observation) HumidityRelative() ObservationHumidity {
	if o.Data.HumidityRelative == nil {
		return ObservationHumidity{na: true}
	}
	return ObservationHumidity{
		dt: o.Data.HumidityRelative.DateTime,
		n:  FieldHumidityRelative,
		v:  o.Data.HumidityRelative.Value,
	}
}

// PressureMSL returns the relative pressure at mean seal level data point
// as ObservationPressure.
// If the data point is not available in the Observation it will return
// ObservationPressure in which the "not available" field will be
// true.
func (o Observation) PressureMSL() ObservationPressure {
	if o.Data.PressureMSL == nil {
		return ObservationPressure{na: true}
	}
	return ObservationPressure{
		dt: o.Data.PressureMSL.DateTime,
		n:  FieldPressureMSL,
		v:  o.Data.PressureMSL.Value,
	}
}

// PressureQFE returns the relative pressure at mean seal level data point
// as ObservationPressure.
// If the data point is not available in the Observation it will return
// ObservationPressure in which the "not available" field will be
// true.
func (o Observation) PressureQFE() ObservationPressure {
	if o.Data.PressureQFE == nil {
		return ObservationPressure{na: true}
	}
	return ObservationPressure{
		dt: o.Data.PressureQFE.DateTime,
		n:  FieldPressureQFE,
		v:  o.Data.PressureQFE.Value,
	}
}

// Precipitation returns the current amount of precipitation (mm) as
// ObservationPrecipitation
// If the data point is not available in the Observation it will return
// ObservationPrecipitation in which the "not available" field will be
// true.
func (o Observation) Precipitation(ts Timespan) ObservationPrecipitation {
	var df *ObservationValue
	var fn ObservationFieldName
	switch ts {
	case TimespanCurrent:
		df = o.Data.Precipitation
		fn = FieldPrecipitation
	case Timespan10Min:
		df = o.Data.Precipitation10m
		fn = FieldPrecipitation10m
	case Timespan1Hour:
		df = o.Data.Precipitation1h
		fn = FieldPrecipitation1h
	case Timespan24Hours:
		df = o.Data.Precipitation24h
		fn = FieldPrecipitation24h
	default:
		return ObservationPrecipitation{na: true}
	}

	if df == nil {
		return ObservationPrecipitation{na: true}
	}
	return ObservationPrecipitation{
		dt: df.DateTime,
		n:  fn,
		v:  df.Value,
	}
}

// GlobalRadiation returns the current amount of global radiation as
// ObservationRadiation
// If the data point is not available in the Observation it will return
// ObservationRadiation in which the "not available" field will be
// true.
func (o Observation) GlobalRadiation(ts Timespan) ObservationRadiation {
	var df *ObservationValue
	var fn ObservationFieldName
	switch ts {
	case Timespan10Min:
		df = o.Data.GlobalRadiation10m
		fn = FieldGlobalRadiation10m
	case Timespan1Hour:
		df = o.Data.GlobalRadiation1h
		fn = FieldGlobalRadiation1h
	case Timespan24Hours:
		df = o.Data.GlobalRadiation24h
		fn = FieldGlobalRadiation24h
	default:
		return ObservationRadiation{na: true}
	}

	if df == nil {
		return ObservationRadiation{na: true}
	}
	return ObservationRadiation{
		dt: df.DateTime,
		n:  fn,
		v:  df.Value,
	}
}

// Winddirection returns the current direction from which the wind
// originates in degree (0=N, 90=E, 180=S, 270=W) as ObservationDirection.
// If the data point is not available in the Observation it will return
// ObservationDirection in which the "not available" field will be true.
func (o Observation) Winddirection() ObservationDirection {
	if o.Data.Winddirection == nil {
		return ObservationDirection{na: true}
	}
	return ObservationDirection{
		dt: o.Data.Winddirection.DateTime,
		n:  FieldWinddirection,
		v:  o.Data.Winddirection.Value,
	}
}

// Windspeed returns the current windspeed data point as ObservationSpeed.
// If the data point is not available in the Observation it will return
// ObservationSpeed in which the "not available" field will be true.
func (o Observation) Windspeed() ObservationSpeed {
	if o.Data.Windspeed == nil {
		return ObservationSpeed{na: true}
	}
	return ObservationSpeed{
		dt: o.Data.Windspeed.DateTime,
		n:  FieldWindspeed,
		v:  o.Data.Windspeed.Value,
	}
}

// IsAvailable returns true if an ObservationTemperature value was
// available at time of query
func (t ObservationTemperature) IsAvailable() bool {
	return !t.na
}

// DateTime returns true if an ObservationTemperature value was
// available at time of query
func (t ObservationTemperature) DateTime() time.Time {
	return t.dt
}

// Value returns the float64 value of an ObservationTemperature
// If the ObservationTemperature is not available in the Observation
// Vaule will return math.NaN instead.
func (t ObservationTemperature) Value() float64 {
	if t.na {
		return math.NaN()
	}
	return t.v
}

// String satisfies the fmt.Stringer interface for the ObservationTemperature type
func (t ObservationTemperature) String() string {
	return fmt.Sprintf("%.1f°C", t.v)
}

// Celsius returns the ObservationTemperature value in Celsius
func (t ObservationTemperature) Celsius() float64 {
	return t.v
}

// CelsiusString returns the ObservationTemperature value as Celsius
// formated string.
//
// This is an alias for the fmt.Stringer interface
func (t ObservationTemperature) CelsiusString() string {
	return t.String()
}

// Fahrenheit returns the ObservationTemperature value in Fahrenheit
func (t ObservationTemperature) Fahrenheit() float64 {
	return t.v*9/5 + 32
}

// FahrenheitString returns the ObservationTemperature value as Fahrenheit
// formated string.
func (t ObservationTemperature) FahrenheitString() string {
	return fmt.Sprintf("%.1f°F", t.Fahrenheit())
}

// IsAvailable returns true if an ObservationHumidity value was
// available at time of query
func (t ObservationHumidity) IsAvailable() bool {
	return !t.na
}

// DateTime returns true if an ObservationHumidity value was
// available at time of query
func (t ObservationHumidity) DateTime() time.Time {
	return t.dt
}

// String satisfies the fmt.Stringer interface for the ObservationHumidity type
func (t ObservationHumidity) String() string {
	return fmt.Sprintf("%.1f%%", t.v)
}

// Value returns the float64 value of an ObservationHumidity
// If the ObservationHumidity is not available in the Observation
// Vaule will return math.NaN instead.
func (t ObservationHumidity) Value() float64 {
	if t.na {
		return math.NaN()
	}
	return t.v
}

// IsAvailable returns true if an ObservationPrecipitation value was
// available at time of query
func (t ObservationPrecipitation) IsAvailable() bool {
	return !t.na
}

// DateTime returns true if an ObservationPrecipitation value was
// available at time of query
func (t ObservationPrecipitation) DateTime() time.Time {
	return t.dt
}

// String satisfies the fmt.Stringer interface for the ObservationPrecipitation type
func (t ObservationPrecipitation) String() string {
	return fmt.Sprintf("%.1fmm", t.v)
}

// Value returns the float64 value of an ObservationPrecipitation
// If the ObservationPrecipitation is not available in the Observation
// Vaule will return math.NaN instead.
func (t ObservationPrecipitation) Value() float64 {
	if t.na {
		return math.NaN()
	}
	return t.v
}

// IsAvailable returns true if an ObservationPressure value was
// available at time of query
func (t ObservationPressure) IsAvailable() bool {
	return !t.na
}

// DateTime returns true if an ObservationPressure value was
// available at time of query
func (t ObservationPressure) DateTime() time.Time {
	return t.dt
}

// String satisfies the fmt.Stringer interface for the ObservationPressure type
func (t ObservationPressure) String() string {
	return fmt.Sprintf("%.1fhPa", t.v)
}

// Value returns the float64 value of an ObservationPressure
// If the ObservationPressure is not available in the Observation
// Vaule will return math.NaN instead.
func (t ObservationPressure) Value() float64 {
	if t.na {
		return math.NaN()
	}
	return t.v
}

// IsAvailable returns true if an ObservationRadiation value was
// available at time of query
func (t ObservationRadiation) IsAvailable() bool {
	return !t.na
}

// DateTime returns true if an ObservationRadiation value was
// available at time of query
func (t ObservationRadiation) DateTime() time.Time {
	return t.dt
}

// Value returns the float64 value of an ObservationRadiation
// If the ObservationRadiation is not available in the Observation
// Vaule will return math.NaN instead.
func (t ObservationRadiation) Value() float64 {
	if t.na {
		return math.NaN()
	}
	return t.v
}

// String satisfies the fmt.Stringer interface for the ObservationRadiation type
func (t ObservationRadiation) String() string {
	return fmt.Sprintf("%.0fkJ/m²", t.v)
}

// IsAvailable returns true if an ObservationSpeed value was
// available at time of query
func (t ObservationSpeed) IsAvailable() bool {
	return !t.na
}

// DateTime returns true if an ObservationSpeed value was
// available at time of query
func (t ObservationSpeed) DateTime() time.Time {
	return t.dt
}

// Value returns the float64 value of an ObservationSpeed in knots
// If the ObservationSpeed is not available in the Observation
// Vaule will return math.NaN instead.
func (t ObservationSpeed) Value() float64 {
	if t.na {
		return math.NaN()
	}
	return t.v
}

// String satisfies the fmt.Stringer interface for the ObservationSpeed type
func (t ObservationSpeed) String() string {
	return fmt.Sprintf("%.0fkn", t.v)
}

// KMH returns the ObservationSpeed value in km/h
func (t ObservationSpeed) KMH() float64 {
	return t.v * 1.852
}

// KMHString returns the ObservationSpeed value as formatted string in km/h
func (t ObservationSpeed) KMHString() string {
	return fmt.Sprintf("%.1fkm/h", t.KMH())
}

// MPH returns the ObservationSpeed value in mi/h
func (t ObservationSpeed) MPH() float64 {
	return t.v * 1.151
}

// MPHString returns the ObservationSpeed value as formatted string in mi/h
func (t ObservationSpeed) MPHString() string {
	return fmt.Sprintf("%.1fmi/h", t.MPH())
}

// IsAvailable returns true if an ObservationDirection value was
// available at time of query
func (t ObservationDirection) IsAvailable() bool {
	return !t.na
}

// DateTime returns true if an ObservationDirection value was
// available at time of query
func (t ObservationDirection) DateTime() time.Time {
	return t.dt
}

// Value returns the float64 value of an ObservationDirection in degrees
// If the ObservationDirection is not available in the Observation
// Vaule will return math.NaN instead.
func (t ObservationDirection) Value() float64 {
	if t.na {
		return math.NaN()
	}
	return t.v
}

// String satisfies the fmt.Stringer interface for the ObservationDirection type
func (t ObservationDirection) String() string {
	return fmt.Sprintf("%.0f°", t.v)
}

// Direction returns the abbreviation string for a given ObservationDirection type
func (t ObservationDirection) Direction() string {
	if ds, ok := WindDirAbbrMap[t.v]; ok {
		return ds
	}
	return ErrUnsupportedDirection
}

// DirectionFull returns the full string for a given ObservationDirection type
func (t ObservationDirection) DirectionFull() string {
	if ds, ok := WindDirFullMap[t.v]; ok {
		return ds
	}
	return ErrUnsupportedDirection
}
