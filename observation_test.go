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
	}{
		{"Koeln-Botanischer Garten", "199942", 44, 50.9667, 6.9667},
		{"Koeln-Stammheim", "H744", 43, 50.9833, 6.9833},
		{"All data fields", "all", 123, 1.234, -1.234},
		{"No data fields", "none", 123, 1.234, -1.234},
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
		})
	}
}

func TestClient_ObservationLatestByStationID_Dewpoint(t *testing.T) {
	tt := []struct {
		// Test name
		n string
		// Station ID
		sid string
		// Observation dewpoint
		dp *ObservationTemperature
	}{
		{"K-Botanischer Garten", "199942", &ObservationTemperature{Value: 10.1}},
		{"K-Stammheim", "H744", &ObservationTemperature{Value: 9.7}},
		{"All data fields", "all", &ObservationTemperature{Value: 6.5}},
		{"No data fields", "none", nil},
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
			if tc.dp != nil && tc.dp.String() != o.DewpointString() {
				t.Errorf("ObservationLatestByStationID failed, expected dewpoint "+
					"string: %s, got: %s", tc.dp.String(), o.DewpointString())
			}
			if tc.dp != nil && tc.dp.Value != o.Dewpoint() {
				t.Errorf("ObservationLatestByStationID failed, expected dewpoint "+
					"float: %f, got: %f", tc.dp.Value, o.Dewpoint())
			}
			if tc.dp == nil {
				if o.DewpointString() != DataNotAvailable {
					t.Errorf("ObservationLatestByStationID failed, expected dewpoint "+
						"to have no data, but got: %s", o.DewpointString())
				}
				if !math.IsNaN(o.Dewpoint()) {
					t.Errorf("ObservationLatestByStationID failed, expected dewpoint "+
						"to return NaN, but got: %f", o.Dewpoint())
				}
			}
		})
	}
}

func TestClient_ObservationLatestByStationID_HumidityRealtive(t *testing.T) {
	tt := []struct {
		// Test name
		n string
		// Station ID
		sid string
		// Observation dewpoint
		h *ObservationHumidity
	}{
		{"K-Botanischer Garten", "199942", &ObservationHumidity{Value: 80}},
		{"K-Stammheim", "H744", &ObservationHumidity{Value: 73}},
		{"All data fields", "all", &ObservationHumidity{Value: 72}},
		{"No data fields", "none", nil},
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
			if tc.h != nil && tc.h.String() != o.HumidityRelativeString() {
				t.Errorf("ObservationLatestByStationID failed, expected humidity "+
					"string: %s, got: %s", tc.h.String(), o.HumidityRelativeString())
			}
			if tc.h != nil && tc.h.Value != o.HumidityRelative() {
				t.Errorf("ObservationLatestByStationID failed, expected humidity "+
					"float: %f, got: %f", tc.h.Value, o.HumidityRelative())
			}
			if tc.h == nil {
				if o.HumidityRelativeString() != DataNotAvailable {
					t.Errorf("ObservationLatestByStationID failed, expected humidity "+
						"to have no data, but got: %s", o.HumidityRelativeString())
				}
				if !math.IsNaN(o.HumidityRelative()) {
					t.Errorf("ObservationLatestByStationID failed, expected humidity "+
						"to return NaN, but got: %f", o.HumidityRelative())
				}
			}
		})
	}
}

func TestClient_ObservationLatestByStationID_PrecipitationCurrent(t *testing.T) {
	tt := []struct {
		// Test name
		n string
		// Station ID
		sid string
		// Observation precipitation
		p *ObservationPrecipitation
	}{
		{"K-Botanischer Garten", "199942", &ObservationPrecipitation{Value: 0}},
		{"K-Stammheim", "H744", &ObservationPrecipitation{Value: 0}},
		{"All data fields", "all", &ObservationPrecipitation{Value: 0.1}},
		{"No data fields", "none", nil},
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
				t.Errorf("ObservationLatestByStationID with station %s "+
					"failed: %s", tc.sid, err)
				return
			}
			if tc.p != nil && tc.p.String() != o.PrecipitationString(PrecipitationCurrent) {
				t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
					"string: %s, got: %s", tc.p.String(), o.PrecipitationString(PrecipitationCurrent))
			}
			if tc.p != nil && tc.p.Value != o.Precipitation(PrecipitationCurrent) {
				t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
					"float: %f, got: %f", tc.p.Value, o.Precipitation(PrecipitationCurrent))
			}
			if tc.p == nil {
				if o.PrecipitationString(PrecipitationCurrent) != DataNotAvailable {
					t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
						"to have no data, but got: %s", o.PrecipitationString(PrecipitationCurrent))
				}
				if !math.IsNaN(o.Precipitation(PrecipitationCurrent)) {
					t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
						"to return NaN, but got: %f", o.Precipitation(PrecipitationCurrent))
				}
			}
		})
	}
}

func TestClient_ObservationLatestByStationID_Precipitation10m(t *testing.T) {
	tt := []struct {
		// Test name
		n string
		// Station ID
		sid string
		// Observation precipitation
		p *ObservationPrecipitation
	}{
		{"K-Botanischer Garten", "199942", &ObservationPrecipitation{Value: 0}},
		{"K-Stammheim", "H744", &ObservationPrecipitation{Value: 0}},
		{"All data fields", "all", &ObservationPrecipitation{Value: 0.5}},
		{"No data fields", "none", nil},
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
				t.Errorf("ObservationLatestByStationID with station %s "+
					"failed: %s", tc.sid, err)
				return
			}
			if tc.p != nil && tc.p.String() != o.PrecipitationString(Precipitation10Min) {
				t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
					"string: %s, got: %s", tc.p.String(), o.PrecipitationString(Precipitation10Min))
			}
			if tc.p != nil && tc.p.Value != o.Precipitation(Precipitation10Min) {
				t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
					"float: %f, got: %f", tc.p.Value, o.Precipitation(Precipitation10Min))
			}
			if tc.p == nil {
				if o.PrecipitationString(Precipitation10Min) != DataNotAvailable {
					t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
						"to have no data, but got: %s", o.PrecipitationString(Precipitation10Min))
				}
				if !math.IsNaN(o.Precipitation(Precipitation10Min)) {
					t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
						"to return NaN, but got: %f", o.Precipitation(Precipitation10Min))
				}
			}
		})
	}
}

func TestClient_ObservationLatestByStationID_Precipitation1h(t *testing.T) {
	tt := []struct {
		// Test name
		n string
		// Station ID
		sid string
		// Observation precipitation
		p *ObservationPrecipitation
	}{
		{"K-Botanischer Garten", "199942", &ObservationPrecipitation{Value: 0}},
		{"K-Stammheim", "H744", &ObservationPrecipitation{Value: 0}},
		{"All data fields", "all", &ObservationPrecipitation{Value: 10.3}},
		{"No data fields", "none", nil},
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
				t.Errorf("ObservationLatestByStationID with station %s "+
					"failed: %s", tc.sid, err)
				return
			}
			if tc.p != nil && tc.p.String() != o.PrecipitationString(Precipitation1Hour) {
				t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
					"string: %s, got: %s", tc.p.String(), o.PrecipitationString(Precipitation1Hour))
			}
			if tc.p != nil && tc.p.Value != o.Precipitation(Precipitation1Hour) {
				t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
					"float: %f, got: %f", tc.p.Value, o.Precipitation(Precipitation1Hour))
			}
			if tc.p == nil {
				if o.PrecipitationString(Precipitation1Hour) != DataNotAvailable {
					t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
						"to have no data, but got: %s", o.PrecipitationString(Precipitation1Hour))
				}
				if !math.IsNaN(o.Precipitation(Precipitation1Hour)) {
					t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
						"to return NaN, but got: %f", o.Precipitation(Precipitation1Hour))
				}
			}
		})
	}
}

func TestClient_ObservationLatestByStationID_Precipitation24h(t *testing.T) {
	tt := []struct {
		// Test name
		n string
		// Station ID
		sid string
		// Observation precipitation
		p *ObservationPrecipitation
	}{
		{"K-Botanischer Garten", "199942", &ObservationPrecipitation{Value: 0}},
		{"K-Stammheim", "H744", &ObservationPrecipitation{Value: 0}},
		{"All data fields", "all", &ObservationPrecipitation{Value: 32.12}},
		{"No data fields", "none", nil},
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
				t.Errorf("ObservationLatestByStationID with station %s "+
					"failed: %s", tc.sid, err)
				return
			}
			if tc.p != nil && tc.p.String() != o.PrecipitationString(Precipitation24Hours) {
				t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
					"string: %s, got: %s", tc.p.String(), o.PrecipitationString(Precipitation24Hours))
			}
			if tc.p != nil && tc.p.Value != o.Precipitation(Precipitation24Hours) {
				t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
					"float: %f, got: %f", tc.p.Value, o.Precipitation(Precipitation24Hours))
			}
			if tc.p == nil {
				if o.PrecipitationString(Precipitation24Hours) != DataNotAvailable {
					t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
						"to have no data, but got: %s", o.PrecipitationString(Precipitation24Hours))
				}
				if !math.IsNaN(o.Precipitation(Precipitation24Hours)) {
					t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
						"to return NaN, but got: %f", o.Precipitation(Precipitation24Hours))
				}
			}
		})
	}
}

func TestClient_ObservationLatestByStationID_PrecipitationUnknown(t *testing.T) {
	c := New(withMockAPI())
	if c == nil {
		t.Errorf("failed to create new Client, got nil")
		return
	}
	o, err := c.ObservationLatestByStationID("all")
	if err != nil {
		t.Errorf("ObservationLatestByStationID with station %s "+
			"failed: %s", "all", err)
		return
	}
	if o.PrecipitationString(999) != ErrTimespanUnsupported {
		t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
			"to have no data, but got: %s", o.PrecipitationString(999))
	}
	if !math.IsNaN(o.Precipitation(999)) {
		t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
			"to return NaN, but got: %f", o.Precipitation(999))
	}
}

func TestClient_ObservationLatestByStationID_Temperature(t *testing.T) {
	tt := []struct {
		// Test name
		n string
		// Station ID
		sid string
		// Observation dewpoint
		t *ObservationTemperature
	}{
		{"K-Botanischer Garten", "199942", &ObservationTemperature{Value: 13.4}},
		{"K-Stammheim", "H744", &ObservationTemperature{Value: 14.4}},
		{"All data fields", "all", &ObservationTemperature{Value: 10.8}},
		{"No data fields", "none", nil},
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
			if tc.t != nil && tc.t.String() != o.TemperatureString() {
				t.Errorf("ObservationLatestByStationID failed, expected temperature "+
					"string: %s, got: %s", tc.t.String(), o.TemperatureString())
			}
			if tc.t != nil && tc.t.Value != o.Temperature() {
				t.Errorf("ObservationLatestByStationID failed, expected temperature "+
					"float: %f, got: %f", tc.t.Value, o.Temperature())
			}
			if tc.t == nil {
				if o.TemperatureString() != DataNotAvailable {
					t.Errorf("ObservationLatestByStationID failed, expected temperature "+
						"to have no data, but got: %s", o.TemperatureString())
				}
				if !math.IsNaN(o.Temperature()) {
					t.Errorf("ObservationLatestByStationID failed, expected temperature "+
						"to return NaN, but got: %f", o.Temperature())
				}
			}
		})
	}
}

func TestClient_ObservationLatestByStationID_PressureMSL(t *testing.T) {
	tt := []struct {
		// Test name
		n string
		// Station ID
		sid string
		// Observation dewpoint
		p *ObservationPressure
	}{
		{"K-Botanischer Garten", "199942", &ObservationPressure{Value: 1015.5}},
		{"K-Stammheim", "H744", nil},
		{"All data fields", "all", &ObservationPressure{Value: 1026.3}},
		{"No data fields", "none", nil},
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
			if tc.p != nil && tc.p.String() != o.PresureMSLString() {
				t.Errorf("ObservationLatestByStationID failed, expected pressure MSL "+
					"string: %s, got: %s", tc.p.String(), o.PresureMSLString())
			}
			if tc.p != nil && tc.p.Value != o.PressureMSL() {
				t.Errorf("ObservationLatestByStationID failed, expected pressure MSL "+
					"float: %f, got: %f", tc.p.Value, o.PressureMSL())
			}
			if tc.p == nil {
				if o.PresureMSLString() != DataNotAvailable {
					t.Errorf("ObservationLatestByStationID failed, expected pressure MSL "+
						"to have no data, but got: %s", o.PresureMSLString())
				}
				if !math.IsNaN(o.PressureMSL()) {
					t.Errorf("ObservationLatestByStationID failed, expected pressure MSL "+
						"to return NaN, but got: %f", o.PressureMSL())
				}
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
