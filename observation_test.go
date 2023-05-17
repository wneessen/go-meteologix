// SPDX-FileCopyrightText: 2023 Winni Neessen <winni@neessen.dev>
//
// SPDX-License-Identifier: MIT

package meteologix

import (
	"fmt"
	"math"
	"testing"
)

func TestClient_ObservationLatestByStationID_Mock(t *testing.T) {
	tt := []struct {
		// Test name
		n string
		// Station ID
		sid string
		// Altitude
		alt int
		// Latitude
		lat float64
		// Longitude
		lon float64
		// Observation data points
		dp *ObservationData
	}{
		{"Koeln-Botanischer Garten", "199942", 44, 50.9667, 6.9667, &ObservationData{
			DewPoint:         &ObservationTemperature{Value: 10.1},
			HumidityRelative: &ObservationHumidity{Value: 80},
			Precipitation:    &ObservationPrecipitation{Value: 0},
			Precipitation10m: &ObservationPrecipitation{Value: 0},
			Precipitation1h:  &ObservationPrecipitation{Value: 0},
			Precipitation24h: &ObservationPrecipitation{Value: 0},
			Temperature:      &ObservationTemperature{Value: 13.4},
		}},
		{"Koeln-Stammheim", "H744", 43, 50.9833, 6.9833, &ObservationData{
			DewPoint:         &ObservationTemperature{Value: 9.7},
			HumidityRelative: &ObservationHumidity{Value: 73},
			Precipitation:    &ObservationPrecipitation{Value: 0},
			Precipitation10m: &ObservationPrecipitation{Value: 0},
			Precipitation1h:  &ObservationPrecipitation{Value: 0},
			Precipitation24h: &ObservationPrecipitation{Value: 0},
			Temperature:      &ObservationTemperature{Value: 14.4},
		}},
		{"All data fields", "all", 123, 1.234, -1.234, &ObservationData{
			DewPoint:         &ObservationTemperature{Value: 6.5},
			HumidityRelative: &ObservationHumidity{Value: 72},
			Precipitation:    &ObservationPrecipitation{Value: 0.1},
			Precipitation10m: &ObservationPrecipitation{Value: 0.5},
			Precipitation1h:  &ObservationPrecipitation{Value: 10.3},
			Precipitation24h: &ObservationPrecipitation{Value: 32.12},
			Temperature:      &ObservationTemperature{Value: 10.8},
		}},
	}
	c := New(withMockAPI())
	if c == nil {
		t.Errorf("failed to create new Client, got nil")
		return
	}
	for _, tc := range tt {
		t.Run(tc.n, func(t *testing.T) {
			o, err := c.ObservationLatestByStationID(tc.sid)
			if err != nil {
				t.Errorf("ObservationLatestByStationID with station %s failed: %s", tc.sid, err)
				return
			}
			if o.StationID != tc.sid {
				t.Errorf("ObservationLatestByStationID failed, expected station id: %s, got: %s",
					tc.sid, o.StationID)
			}
			if o.Name != tc.n {
				t.Errorf("ObservationLatestByStationID failed, expected name: %s, got: %s",
					tc.n, o.Name)
			}
			if o.Altitude == nil {
				t.Errorf("ObservationLatestByStationID failed, expected altitude but got nil")
			}
			if o.Altitude != nil && *o.Altitude != tc.alt {
				t.Errorf("ObservationLatestByStationID failed, expected altitude: %d, got: %d",
					tc.alt, *o.Altitude)
			}
			if tc.dp == nil {
				t.Skip("No data points received, this might be intentionally. Skipping data point validation.")
			}
			if tc.dp.DewPoint != nil && tc.dp.DewPoint.String() != o.Dewpoint() {
				t.Errorf("ObservationLatestByStationID failed, expected dewpoint string: %s, got: %s",
					tc.dp.DewPoint.String(), o.Dewpoint())
			}
			if tc.dp.DewPoint == nil {
				if o.Dewpoint() != DataNotAvailable {
					t.Errorf("ObservationLatestByStationID failed, expected dewpoint to have "+
						"no data, but got: %s", o.Dewpoint())
				}
			}
			if tc.dp.HumidityRelative != nil && tc.dp.HumidityRelative.String() != o.HumidityRelative() {
				t.Errorf("ObservationLatestByStationID failed, expected humidity string: %s, got: %s",
					tc.dp.HumidityRelative.String(), o.HumidityRelative())
			}
			if tc.dp.HumidityRelative == nil {
				if o.HumidityRelative() != DataNotAvailable {
					t.Errorf("ObservationLatestByStationID failed, expected relative humidity to have "+
						"no data, but got: %s", o.HumidityRelative())
				}
			}
			if tc.dp.Precipitation != nil && tc.dp.Precipitation.String() != o.Precipitation(PrecipitationCurrent) {
				t.Errorf("ObservationLatestByStationID failed, expected precipitation string: %s, got: %s",
					tc.dp.Precipitation.String(), o.Precipitation(PrecipitationCurrent))
			}
			if tc.dp.Precipitation == nil {
				if o.Precipitation(PrecipitationCurrent) != DataNotAvailable {
					t.Errorf("ObservationLatestByStationID failed, expected precipitation (current) to have "+
						"no data, but got: %s", o.Precipitation(PrecipitationCurrent))
				}
			}
			if tc.dp.Precipitation10m != nil && tc.dp.Precipitation10m.String() != o.Precipitation(Precipitation10Min) {
				t.Errorf("ObservationLatestByStationID failed, expected precipitation (10m) string: %s, got: %s",
					tc.dp.Precipitation10m.String(), o.Precipitation(Precipitation10Min))
			}
			if tc.dp.Precipitation10m == nil {
				if o.Precipitation(Precipitation10Min) != DataNotAvailable {
					t.Errorf("ObservationLatestByStationID failed, expected precipitation (10m) to have "+
						"no data, but got: %s", o.Precipitation(Precipitation10Min))
				}
			}
			if tc.dp.Precipitation1h != nil && tc.dp.Precipitation1h.String() != o.Precipitation(Precipitation1Hour) {
				t.Errorf("ObservationLatestByStationID failed, expected precipitation (1h) string: %s, got: %s",
					tc.dp.Precipitation1h.String(), o.Precipitation(Precipitation1Hour))
			}
			if tc.dp.Precipitation1h == nil {
				if o.Precipitation(Precipitation1Hour) != DataNotAvailable {
					t.Errorf("ObservationLatestByStationID failed, expected precipitation (1h) to have "+
						"no data, but got: %s", o.Precipitation(Precipitation1Hour))
				}
			}
			if tc.dp.Precipitation24h != nil && tc.dp.Precipitation24h.String() != o.Precipitation(Precipitation24Hours) {
				t.Errorf("ObservationLatestByStationID failed, expected precipitation (24h) string: %s, got: %s",
					tc.dp.Precipitation24h.String(), o.Precipitation(Precipitation24Hours))
			}
			if tc.dp.Precipitation24h == nil {
				if o.Precipitation(Precipitation24Hours) != DataNotAvailable {
					t.Errorf("ObservationLatestByStationID failed, expected precipitation (24h) to have "+
						"no data, but got: %s", o.Precipitation(Precipitation24Hours))
				}
			}
			if tc.dp.Temperature != nil && tc.dp.Temperature.String() != o.TemperatureString() {
				t.Errorf("ObservationLatestByStationID failed, expected temperature string: %s, got: %s",
					tc.dp.Temperature.String(), o.TemperatureString())
			}
			if tc.dp.Temperature != nil && tc.dp.Temperature.Value != o.Temperature() {
				t.Errorf("ObservationLatestByStationID failed, expected temperature value: %f, got: %f",
					tc.dp.Temperature.Value, o.Temperature())
			}
			if tc.dp.Temperature == nil {
				if o.TemperatureString() != DataNotAvailable {
					t.Errorf("ObservationLatestByStationID failed, expected temperature to have "+
						"no data, but got: %s", o.TemperatureString())
				}
				if !math.IsNaN(o.Temperature()) {
					t.Errorf("ObservationLatestByStationID failed, expected temperature to be NaN, "+
						"but got: %f", o.Temperature())
				}
			}
			if tc.dp.Temperature != nil && tc.dp.Temperature.Celsius() != o.Data.Temperature.Celsius() {
				t.Errorf("ObservationLatestByStationID failed, expected temperature float: %f, got: %f",
					tc.dp.Temperature.Celsius(), o.Data.Temperature.Celsius())
			}
			if tc.dp.Temperature != nil && tc.dp.Temperature.Fahrenheit() != o.Data.Temperature.Fahrenheit() {
				t.Errorf("ObservationLatestByStationID failed, expected temperature (F) float: %f, got: %f",
					tc.dp.Temperature.Fahrenheit(), o.Data.Temperature.Fahrenheit())
			}
		})
	}
}

func TestObservationTemperature_String(t *testing.T) {
	tt := []struct {
		// Original celsius value
		c float64
		// Fahrenheit value
		f float64
	}{
		{-273.15, -459.66999999999996},
		{-50, -58},
		{-40, -40},
		{-30, -22},
		{-20, -4},
		{-10, 14},
		{-9, 15.8},
		{-8, 17.6},
		{-7, 19.4},
		{-6, 21.2},
		{-5, 23},
		{-4, 24.8},
		{-3, 26.6},
		{-2, 28.4},
		{-1, 30.2},
		{0, 32},
		{1, 33.8},
		{2, 35.6},
		{3, 37.4},
		{4, 39.2},
		{5, 41},
		{6, 42.8},
		{7, 44.6},
		{8, 46.4},
		{9, 48.2},
		{10, 50},
		{20, 68},
		{30, 86},
		{40, 104},
		{50, 122},
		{100, 212},
	}
	cf := "%.1f°C"
	ff := "%.1f°F"
	for _, tc := range tt {
		t.Run(fmt.Sprintf("%.2f°C", tc.c), func(t *testing.T) {
			ot := ObservationTemperature{Value: tc.c}
			if ot.Celsius() != tc.c {
				t.Errorf("ObservationTemperature.Celsius failed, expected: %f, got: %f", tc.c,
					ot.Celsius())
			}
			if ot.CelsiusString() != fmt.Sprintf(cf, tc.c) {
				t.Errorf("ObservationTemperature.CelsiusString failed, expected: %s, got: %s",
					fmt.Sprintf(cf, tc.c), ot.CelsiusString())
			}
			if ot.Fahrenheit() != tc.f {
				t.Errorf("ObservationTemperature.Fahrenheit failed, expected: %f, got: %f", tc.f,
					ot.Fahrenheit())
			}
			if ot.FahrenheitString() != fmt.Sprintf(ff, tc.f) {
				t.Errorf("ObservationTemperature.FahrenheitString failed, expected: %s, got: %s",
					fmt.Sprintf(ff, tc.f), ot.FahrenheitString())
			}
		})
	}
}
