// SPDX-FileCopyrightText: 2023 Winni Neessen <wn@neessen.dev>
//
// SPDX-License-Identifier: MIT

package meteologix

import (
	"encoding/json"
	"fmt"
)

// ErrUnsupportedDirection is returned when a direction degree is given,
// that is not resolvable
var ErrUnsupportedDirection = "Unsupported direction"

// Observation represents the observation API response for a Station
type Observation struct {
	// Altitude is the altitude of the station providing the Observation
	Altitude *int `json:"ele,omitempty"`
	// Data holds the different APIObservationData points
	Data APIObservationData `json:"data"`
	// Name is the name of the Station providing the Observation
	Name string `json:"name"`
	// Latitude represents the GeoLocation latitude coordinates for the Station
	Latitude float64 `json:"lat"`
	// Longitude represents the GeoLocation longitude coordinates for the Station
	Longitude float64 `json:"lon"`
	// StationID is the ID of the Station providing the Observation
	StationID string `json:"stationId"`
}

// APIObservationData holds the different data points of the Observation as
// returned by the station observation API endpoints.
//
// Please keep in mind that different Station types return different values, therefore
// all values are represented as pointer type returning nil if the data point in question
// is not returned for the requested Station.
type APIObservationData struct {
	// Dewpoint represents the dewpoint in °C
	Dewpoint *APIFloat `json:"dewpoint,omitempty"`
	// DewPointMean represents the mean dewpoint in °C
	DewpointMean *APIFloat `json:"dewpointMean,omitempty"`
	// GlobalRadiation10m represents the sum of global radiation over the last
	// 10 minutes in kJ/m²
	GlobalRadiation10m *APIFloat `json:"globalRadiation10m,omitempty"`
	// GlobalRadiation1h represents the sum of global radiation over the last
	// 1 hour in kJ/m²
	GlobalRadiation1h *APIFloat `json:"globalRadiation1h,omitempty"`
	// GlobalRadiation24h represents the sum of global radiation over the last
	// 24 hour in kJ/m²
	GlobalRadiation24h *APIFloat `json:"globalRadiation24h,omitempty"`
	// HumidityRelative represents the relative humidity in percent
	HumidityRelative *APIFloat `json:"humidityRelative,omitempty"`
	// Precipitation represents the current amount of precipitation
	Precipitation *APIFloat `json:"prec"`
	// Precipitation10m represents the amount of precipitation over the last 10 minutes
	Precipitation10m *APIFloat `json:"prec10m"`
	// Precipitation1h represents the amount of precipitation over the last hour
	Precipitation1h *APIFloat `json:"prec1h"`
	// Precipitation24h represents the amount of precipitation over the last 24 hours
	Precipitation24h *APIFloat `json:"prec24h"`
	// PressureMSL represents the pressure at mean sea level (MSL) in hPa
	PressureMSL *APIFloat `json:"pressureMsl"`
	// PressureMSL represents the pressure at station level (QFE) in hPa
	PressureQFE *APIFloat `json:"pressure"`
	// Temperature represents the temperature in °C
	Temperature *APIFloat `json:"temp,omitempty"`
	// TemperatureMax represents the maximum temperature in °C
	TemperatureMax *APIFloat `json:"tempMax,omitempty"`
	// TemperatureMean represents the mean temperature in °C
	TemperatureMean *APIFloat `json:"tempMean,omitempty"`
	// TemperatureMin represents the minimum temperature in °C
	TemperatureMin *APIFloat `json:"tempMin,omitempty"`
	// Temperature5cm represents the temperature 5cm above ground in °C
	Temperature5cm *APIFloat `json:"temp5cm,omitempty"`
	// Temperature5cm represents the minimum temperature 5cm above
	// ground in °C
	Temperature5cmMin *APIFloat `json:"temp5cmMin,omitempty"`
	// Winddirection represents the direction from which the wind
	// originates in degree (0=N, 90=E, 180=S, 270=W)
	Winddirection *APIFloat `json:"windDirection,omitempty"`
	// Windspeed represents the wind speed in knots
	Windspeed *APIFloat `json:"windSpeed,omitempty"`
}

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

// Dewpoint returns the dewpoint data point as Temperature
// If the data point is not available in the Observation it will return
// Temperature in which the "not available" field will be
// true.
func (o Observation) Dewpoint() Temperature {
	if o.Data.Dewpoint == nil {
		return Temperature{na: true}
	}
	return Temperature{
		dt: o.Data.Dewpoint.DateTime,
		n:  FieldDewpoint,
		s:  SourceObservation,
		fv: o.Data.Dewpoint.Value,
	}
}

// DewpointMean returns the mean dewpoint data point as Temperature.
// If the data point is not available in the Observation it will return
// Temperature in which the "not available" field will be
// true.
func (o Observation) DewpointMean() Temperature {
	if o.Data.DewpointMean == nil {
		return Temperature{na: true}
	}
	return Temperature{
		dt: o.Data.DewpointMean.DateTime,
		n:  FieldDewpointMean,
		s:  SourceObservation,
		fv: o.Data.DewpointMean.Value,
	}
}

// Temperature returns the temperature data point as Temperature.
// If the data point is not available in the Observation it will return
// Temperature in which the "not available" field will be
// true.
func (o Observation) Temperature() Temperature {
	if o.Data.Temperature == nil {
		return Temperature{na: true}
	}
	return Temperature{
		dt: o.Data.Temperature.DateTime,
		n:  FieldTemperature,
		s:  SourceObservation,
		fv: o.Data.Temperature.Value,
	}
}

// TemperatureAtGround returns the temperature at ground level (5cm)
// data point as Temperature.
// If the data point is not available in the Observation it will return
// Temperature in which the "not available" field will be
// true.
func (o Observation) TemperatureAtGround() Temperature {
	if o.Data.Temperature5cm == nil {
		return Temperature{na: true}
	}
	return Temperature{
		dt: o.Data.Temperature5cm.DateTime,
		n:  FieldTemperatureAtGround,
		s:  SourceObservation,
		fv: o.Data.Temperature5cm.Value,
	}
}

// TemperatureMax returns the maximum temperature so far data point as
// Temperature.
// If the data point is not available in the Observation it will return
// Temperature in which the "not available" field will be
// true.
func (o Observation) TemperatureMax() Temperature {
	if o.Data.TemperatureMax == nil {
		return Temperature{na: true}
	}
	return Temperature{
		dt: o.Data.TemperatureMax.DateTime,
		n:  FieldTemperatureMax,
		s:  SourceObservation,
		fv: o.Data.TemperatureMax.Value,
	}
}

// TemperatureMin returns the minimum temperature so far data point as
// Temperature.
// If the data point is not available in the Observation it will return
// Temperature in which the "not available" field will be
// true.
func (o Observation) TemperatureMin() Temperature {
	if o.Data.TemperatureMin == nil {
		return Temperature{na: true}
	}
	return Temperature{
		dt: o.Data.TemperatureMin.DateTime,
		n:  FieldTemperatureMin,
		s:  SourceObservation,
		fv: o.Data.TemperatureMin.Value,
	}
}

// TemperatureAtGroundMin returns the minimum temperature so far
// at ground level (5cm) data point as Temperature
// If the data point is not available in the Observation it will return
// Temperature in which the "not available" field will be
// true.
func (o Observation) TemperatureAtGroundMin() Temperature {
	if o.Data.Temperature5cmMin == nil {
		return Temperature{na: true}
	}
	return Temperature{
		dt: o.Data.Temperature5cmMin.DateTime,
		n:  FieldTemperatureAtGroundMin,
		s:  SourceObservation,
		fv: o.Data.Temperature5cmMin.Value,
	}
}

// TemperatureMean returns the mean temperature data point as Temperature.
// If the data point is not available in the Observation it will return
// Temperature in which the "not available" field will be
// true.
func (o Observation) TemperatureMean() Temperature {
	if o.Data.TemperatureMean == nil {
		return Temperature{na: true}
	}
	return Temperature{
		dt: o.Data.TemperatureMean.DateTime,
		n:  FieldTemperatureMean,
		s:  SourceObservation,
		fv: o.Data.TemperatureMean.Value,
	}
}

// HumidityRelative returns the relative humidity data point as float64.
// If the data point is not available in the Observation it will return
// Humidity in which the "not available" field will be
// true.
func (o Observation) HumidityRelative() Humidity {
	if o.Data.HumidityRelative == nil {
		return Humidity{na: true}
	}
	return Humidity{
		dt: o.Data.HumidityRelative.DateTime,
		n:  FieldHumidityRelative,
		s:  SourceObservation,
		fv: o.Data.HumidityRelative.Value,
	}
}

// PressureMSL returns the relative pressure at mean seal level data point
// as Pressure.
// If the data point is not available in the Observation it will return
// Pressure in which the "not available" field will be
// true.
func (o Observation) PressureMSL() Pressure {
	if o.Data.PressureMSL == nil {
		return Pressure{na: true}
	}
	return Pressure{
		dt: o.Data.PressureMSL.DateTime,
		n:  FieldPressureMSL,
		s:  SourceObservation,
		fv: o.Data.PressureMSL.Value,
	}
}

// PressureQFE returns the relative pressure at mean seal level data point
// as Pressure.
// If the data point is not available in the Observation it will return
// Pressure in which the "not available" field will be
// true.
func (o Observation) PressureQFE() Pressure {
	if o.Data.PressureQFE == nil {
		return Pressure{na: true}
	}
	return Pressure{
		dt: o.Data.PressureQFE.DateTime,
		n:  FieldPressureQFE,
		s:  SourceObservation,
		fv: o.Data.PressureQFE.Value,
	}
}

// Precipitation returns the current amount of precipitation (mm) as
// Precipitation
// If the data point is not available in the Observation it will return
// Precipitation in which the "not available" field will be
// true.
func (o Observation) Precipitation(ts Timespan) Precipitation {
	var df *APIFloat
	var fn Fieldname
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
		return Precipitation{na: true}
	}

	if df == nil {
		return Precipitation{na: true}
	}
	return Precipitation{
		dt: df.DateTime,
		n:  fn,
		s:  SourceObservation,
		fv: df.Value,
	}
}

// GlobalRadiation returns the current amount of global radiation as
// Radiation
// If the data point is not available in the Observation it will return
// Radiation in which the "not available" field will be
// true.
func (o Observation) GlobalRadiation(ts Timespan) Radiation {
	var df *APIFloat
	var fn Fieldname
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
		return Radiation{na: true}
	}

	if df == nil {
		return Radiation{na: true}
	}
	return Radiation{
		dt: df.DateTime,
		n:  fn,
		s:  SourceObservation,
		fv: df.Value,
	}
}

// Winddirection returns the current direction from which the wind
// originates in degree (0=N, 90=E, 180=S, 270=W) as Direction.
// If the data point is not available in the Observation it will return
// Direction in which the "not available" field will be true.
func (o Observation) Winddirection() Direction {
	if o.Data.Winddirection == nil {
		return Direction{na: true}
	}
	return Direction{
		dt: o.Data.Winddirection.DateTime,
		n:  FieldWinddirection,
		s:  SourceObservation,
		fv: o.Data.Winddirection.Value,
	}
}

// Windspeed returns the current windspeed data point as Speed.
// If the data point is not available in the Observation it will return
// Speed in which the "not available" field will be true.
func (o Observation) Windspeed() Speed {
	if o.Data.Windspeed == nil {
		return Speed{na: true}
	}
	return Speed{
		dt: o.Data.Windspeed.DateTime,
		n:  FieldWindspeed,
		s:  SourceObservation,
		fv: o.Data.Windspeed.Value,
	}
}
