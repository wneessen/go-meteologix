// SPDX-FileCopyrightText: 2023 Winni Neessen <winni@neessen.dev>
//
// SPDX-License-Identifier: MIT

package meteologix

import (
	"fmt"
	"math"
	"testing"
	"time"
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
		{"K-Botanischer Garten", "199942", &ObservationTemperature{
			dt: time.Date(2023, 0o5, 15, 20, 10, 0, 0, time.UTC),
			v:  10.1,
		}},
		{"K-Stammheim", "H744", &ObservationTemperature{
			dt: time.Date(2023, 0o5, 15, 19, 30, 0, 0, time.UTC),
			v:  9.7,
		}},
		{"All data fields", "all", &ObservationTemperature{
			dt: time.Date(2023, 0o5, 17, 7, 40, 0, 0, time.UTC),
			v:  6.5,
		}},
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
			if tc.dp != nil && tc.dp.String() != o.Dewpoint().String() {
				t.Errorf("ObservationLatestByStationID failed, expected dewpoint "+
					"string: %s, got: %s", tc.dp.String(), o.Dewpoint())
			}
			if tc.dp != nil && tc.dp.Value() != o.Dewpoint().Value() {
				t.Errorf("ObservationLatestByStationID failed, expected dewpoint "+
					"float: %f, got: %f", tc.dp.Value(), o.Dewpoint().Value())
			}
			if tc.dp != nil && tc.dp.dt.Unix() != o.Dewpoint().Datetime().Unix() {
				t.Errorf("ObservationLatestByStationID failed, expected datetime: %s, got: %s",
					tc.dp.dt.Format(time.RFC3339), o.Dewpoint().Datetime().Format(time.RFC3339))
			}
			if tc.dp == nil {
				if o.Dewpoint().IsAvailable() {
					t.Errorf("ObservationLatestByStationID failed, expected dewpoint "+
						"to have no data, but got: %s", o.Dewpoint().String())
				}
				if !math.IsNaN(o.Dewpoint().Value()) {
					t.Errorf("ObservationLatestByStationID failed, expected dewpoint "+
						"to return NaN, but got: %s", o.Dewpoint().String())
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
		{"K-Botanischer Garten", "199942", &ObservationHumidity{
			dt: time.Date(2023, 0o5, 15, 20, 10, 0, 0, time.UTC),
			v:  80,
		}},
		{"K-Stammheim", "H744", &ObservationHumidity{
			dt: time.Date(2023, 0o5, 15, 19, 30, 0, 0, time.UTC),
			v:  73,
		}},
		{"All data fields", "all", &ObservationHumidity{
			dt: time.Date(2023, 0o5, 17, 7, 40, 0, 0, time.UTC),
			v:  72,
		}},
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
			if tc.h != nil && tc.h.String() != o.HumidityRelative().String() {
				t.Errorf("ObservationLatestByStationID failed, expected humidity "+
					"string: %s, got: %s", tc.h.String(), o.HumidityRelative())
			}
			if tc.h != nil && tc.h.Value() != o.HumidityRelative().Value() {
				t.Errorf("ObservationLatestByStationID failed, expected humidity "+
					"float: %f, got: %f", tc.h.Value(), o.HumidityRelative().Value())
			}
			if tc.h != nil && tc.h.dt.Unix() != o.HumidityRelative().Datetime().Unix() {
				t.Errorf("ObservationLatestByStationID failed, expected datetime: %s, got: %s",
					tc.h.dt.Format(time.RFC3339), o.HumidityRelative().Datetime().Format(time.RFC3339))
			}
			if tc.h == nil {
				if o.HumidityRelative().IsAvailable() {
					t.Errorf("ObservationLatestByStationID failed, expected humidity "+
						"to have no data, but got: %s", o.HumidityRelative())
				}
				if !math.IsNaN(o.HumidityRelative().Value()) {
					t.Errorf("ObservationLatestByStationID failed, expected humidity "+
						"to return NaN, but got: %s", o.HumidityRelative().String())
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
		{"K-Botanischer Garten", "199942", &ObservationPrecipitation{
			dt: time.Date(2023, 0o5, 15, 18, 0, 0, 0, time.UTC),
			v:  0,
		}},
		{"K-Stammheim", "H744", &ObservationPrecipitation{
			dt: time.Date(2023, 0o5, 15, 19, 30, 0, 0, time.UTC),
			v:  0,
		}},
		{"All data fields", "all", &ObservationPrecipitation{
			dt: time.Date(2023, 0o5, 17, 7, 30, 0, 0, time.UTC),
			v:  0.1,
		}},
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
			if tc.p != nil && tc.p.String() != o.Precipitation(PrecipitationCurrent).String() {
				t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
					"string: %s, got: %s", tc.p.String(), o.Precipitation(PrecipitationCurrent))
			}
			if tc.p != nil && tc.p.Value() != o.Precipitation(PrecipitationCurrent).Value() {
				t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
					"float: %f, got: %f", tc.p.Value(), o.Precipitation(PrecipitationCurrent).Value())
			}
			if tc.p != nil && tc.p.dt.Unix() != o.Precipitation(PrecipitationCurrent).Datetime().Unix() {
				t.Errorf("ObservationLatestByStationID failed, expected datetime: %s, got: %s",
					tc.p.dt.Format(time.RFC3339),
					o.Precipitation(PrecipitationCurrent).Datetime().Format(time.RFC3339))
			}
			if tc.p == nil {
				if o.Precipitation(PrecipitationCurrent).IsAvailable() {
					t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
						"to have no data, but got: %s", o.Precipitation(PrecipitationCurrent))
				}
				if !math.IsNaN(o.Precipitation(PrecipitationCurrent).Value()) {
					t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
						"to return NaN, but got: %s", o.Precipitation(PrecipitationCurrent).String())
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
		{"K-Botanischer Garten", "199942", &ObservationPrecipitation{v: 0}},
		{"K-Stammheim", "H744", &ObservationPrecipitation{v: 0}},
		{"All data fields", "all", &ObservationPrecipitation{v: 0.5}},
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
			if tc.p != nil && tc.p.String() != o.Precipitation(Precipitation10Min).String() {
				t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
					"string: %s, got: %s", tc.p.String(), o.Precipitation(Precipitation10Min))
			}
			if tc.p != nil && tc.p.Value() != o.Precipitation(Precipitation10Min).Value() {
				t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
					"float: %f, got: %f", tc.p.Value(), o.Precipitation(Precipitation10Min).Value())
			}
			if tc.p == nil {
				if o.Precipitation(Precipitation10Min).IsAvailable() {
					t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
						"to have no data, but got: %s", o.Precipitation(Precipitation10Min))
				}
				if !math.IsNaN(o.Precipitation(Precipitation10Min).Value()) {
					t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
						"to return NaN, but got: %s", o.Precipitation(Precipitation10Min).String())
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
		{"K-Botanischer Garten", "199942", &ObservationPrecipitation{v: 0}},
		{"K-Stammheim", "H744", &ObservationPrecipitation{v: 0}},
		{"All data fields", "all", &ObservationPrecipitation{v: 10.3}},
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
			if tc.p != nil && tc.p.String() != o.Precipitation(Precipitation1Hour).String() {
				t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
					"string: %s, got: %s", tc.p.String(), o.Precipitation(Precipitation1Hour))
			}
			if tc.p != nil && tc.p.Value() != o.Precipitation(Precipitation1Hour).Value() {
				t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
					"float: %f, got: %f", tc.p.Value(), o.Precipitation(Precipitation1Hour).Value())
			}
			if tc.p == nil {
				if o.Precipitation(Precipitation1Hour).IsAvailable() {
					t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
						"to have no data, but got: %s", o.Precipitation(Precipitation1Hour))
				}
				if !math.IsNaN(o.Precipitation(Precipitation1Hour).Value()) {
					t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
						"to return NaN, but got: %s", o.Precipitation(Precipitation1Hour).String())
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
		{"K-Botanischer Garten", "199942", &ObservationPrecipitation{v: 0}},
		{"K-Stammheim", "H744", &ObservationPrecipitation{v: 0}},
		{"All data fields", "all", &ObservationPrecipitation{v: 32.12}},
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
			if tc.p != nil && tc.p.String() != o.Precipitation(Precipitation24Hours).String() {
				t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
					"string: %s, got: %s", tc.p.String(), o.Precipitation(Precipitation24Hours))
			}
			if tc.p != nil && tc.p.Value() != o.Precipitation(Precipitation24Hours).Value() {
				t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
					"float: %f, got: %f", tc.p.Value(), o.Precipitation(Precipitation24Hours).Value())
			}
			if tc.p == nil {
				if o.Precipitation(Precipitation24Hours).IsAvailable() {
					t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
						"to have no data, but got: %s", o.Precipitation(Precipitation24Hours))
				}
				if !math.IsNaN(o.Precipitation(Precipitation24Hours).Value()) {
					t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
						"to return NaN, but got: %s", o.Precipitation(Precipitation24Hours).String())
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
	if o.Precipitation(999).IsAvailable() {
		t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
			"to have no data, but got: %s", o.Precipitation(999))
	}
	if !math.IsNaN(o.Precipitation(999).Value()) {
		t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
			"to return NaN, but got: %s", o.Precipitation(999).String())
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
		{"K-Botanischer Garten", "199942", &ObservationTemperature{v: 13.4}},
		{"K-Stammheim", "H744", &ObservationTemperature{v: 14.4}},
		{"All data fields", "all", &ObservationTemperature{v: 10.8}},
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
			if tc.t != nil && tc.t.String() != o.Temperature().String() {
				t.Errorf("ObservationLatestByStationID failed, expected temperature "+
					"string: %s, got: %s", tc.t.String(), o.Temperature())
			}
			if tc.t != nil && tc.t.Value() != o.Temperature().Value() {
				t.Errorf("ObservationLatestByStationID failed, expected temperature "+
					"float: %f, got: %f", tc.t.Value(), o.Temperature().Value())
			}
			if tc.t == nil {
				if o.Temperature().IsAvailable() {
					t.Errorf("ObservationLatestByStationID failed, expected temperature "+
						"to have no data, but got: %s", o.Temperature())
				}
				if !math.IsNaN(o.Temperature().Value()) {
					t.Errorf("ObservationLatestByStationID failed, expected temperature "+
						"to return NaN, but got: %s", o.Temperature().String())
				}
			}
		})
	}
}

func TestClient_ObservationLatestByStationID_TemperatureAtGround(t *testing.T) {
	tt := []struct {
		// Test name
		n string
		// Station ID
		sid string
		// Observation dewpoint
		t *ObservationTemperature
	}{
		{"K-Botanischer Garten", "199942", nil},
		{"K-Stammheim", "H744", &ObservationTemperature{v: 14.3}},
		{"All data fields", "all", &ObservationTemperature{v: 15.4}},
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
			if tc.t != nil && tc.t.String() != o.TemperatureAtGround().String() {
				t.Errorf("ObservationLatestByStationID failed, expected temperature (5cm) "+
					"string: %s, got: %s", tc.t.String(), o.TemperatureAtGround().String())
			}
			if tc.t != nil && tc.t.Value() != o.TemperatureAtGround().Value() {
				t.Errorf("ObservationLatestByStationID failed, expected temperature (5cm) "+
					"float: %f, got: %f", tc.t.Value(), o.TemperatureAtGround().Value())
			}
			if tc.t == nil {
				if o.TemperatureAtGround().IsAvailable() {
					t.Errorf("ObservationLatestByStationID failed, expected temperature (5cm) "+
						"to have no data, but got: %s", o.TemperatureAtGround())
				}
				if !math.IsNaN(o.TemperatureAtGround().Value()) {
					t.Errorf("ObservationLatestByStationID failed, expected temperature (5cm) "+
						"to return NaN, but got: %s", o.TemperatureAtGround().String())
				}
			}
		})
	}
}

func TestClient_ObservationLatestByStationID_TemperatureMin(t *testing.T) {
	tt := []struct {
		// Test name
		n string
		// Station ID
		sid string
		// Observation dewpoint
		t *ObservationTemperature
	}{
		{"K-Botanischer Garten", "199942", &ObservationTemperature{v: 12.3}},
		{"K-Stammheim", "H744", &ObservationTemperature{v: 11.9}},
		{"All data fields", "all", &ObservationTemperature{v: 6.2}},
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
			if tc.t != nil && tc.t.String() != o.TemperatureMin().String() {
				t.Errorf("ObservationLatestByStationID failed, expected temperature (min) "+
					"string: %s, got: %s", tc.t.String(), o.TemperatureMin())
			}
			if tc.t != nil && tc.t.Value() != o.TemperatureMin().Value() {
				t.Errorf("ObservationLatestByStationID failed, expected temperature (min) "+
					"float: %f, got: %f", tc.t.Value(), o.TemperatureMin().Value())
			}
			if tc.t == nil {
				if o.TemperatureMin().IsAvailable() {
					t.Errorf("ObservationLatestByStationID failed, expected temperature (min) "+
						"to have no data, but got: %s", o.TemperatureMin())
				}
				if !math.IsNaN(o.TemperatureMin().Value()) {
					t.Errorf("ObservationLatestByStationID failed, expected temperature (min) "+
						"to return NaN, but got: %s", o.TemperatureMin().String())
				}
			}
		})
	}
}

func TestClient_ObservationLatestByStationID_TemperatureMax(t *testing.T) {
	tt := []struct {
		// Test name
		n string
		// Station ID
		sid string
		// Observation dewpoint
		t *ObservationTemperature
	}{
		{"K-Botanischer Garten", "199942", &ObservationTemperature{v: 20.5}},
		{"K-Stammheim", "H744", &ObservationTemperature{v: 20.7}},
		{"All data fields", "all", &ObservationTemperature{v: 12.4}},
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
			if tc.t != nil && tc.t.String() != o.TemperatureMax().String() {
				t.Errorf("ObservationLatestByStationID failed, expected temperature (max) "+
					"string: %s, got: %s", tc.t.String(), o.TemperatureMax())
			}
			if tc.t != nil && tc.t.Value() != o.TemperatureMax().Value() {
				t.Errorf("ObservationLatestByStationID failed, expected temperature (max) "+
					"float: %f, got: %f", tc.t.Value(), o.TemperatureMax().Value())
			}
			if tc.t == nil {
				if o.TemperatureMax().IsAvailable() {
					t.Errorf("ObservationLatestByStationID failed, expected temperature (max) "+
						"to have no data, but got: %s", o.TemperatureMax())
				}
				if !math.IsNaN(o.TemperatureMax().Value()) {
					t.Errorf("ObservationLatestByStationID failed, expected temperature (max) "+
						"to return NaN, but got: %s", o.TemperatureMax().String())
				}
			}
		})
	}
}

func TestClient_ObservationLatestByStationID_TemperatureAtGroundMin(t *testing.T) {
	tt := []struct {
		// Test name
		n string
		// Station ID
		sid string
		// Observation dewpoint
		t *ObservationTemperature
	}{
		{"K-Botanischer Garten", "199942", nil},
		{"K-Stammheim", "H744", &ObservationTemperature{v: 12.8}},
		{"All data fields", "all", &ObservationTemperature{v: 3.7}},
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
			if tc.t != nil && tc.t.String() != o.TemperatureAtGroundMin().String() {
				t.Errorf("ObservationLatestByStationID failed, expected temperature (5cm-min) "+
					"string: %s, got: %s", tc.t.String(), o.TemperatureAtGroundMin())
			}
			if tc.t != nil && tc.t.Value() != o.TemperatureAtGroundMin().Value() {
				t.Errorf("ObservationLatestByStationID failed, expected temperature (5cm-min) "+
					"float: %f, got: %f", tc.t.Value(), o.TemperatureAtGroundMin().Value())
			}
			if tc.t == nil {
				if o.TemperatureAtGroundMin().IsAvailable() {
					t.Errorf("ObservationLatestByStationID failed, expected temperature (5cm-min) "+
						"to have no data, but got: %s", o.TemperatureAtGroundMin())
				}
				if !math.IsNaN(o.TemperatureAtGroundMin().Value()) {
					t.Errorf("ObservationLatestByStationID failed, expected temperature (5cm-min) "+
						"to return NaN, but got: %s", o.TemperatureAtGroundMin().String())
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
		{"K-Botanischer Garten", "199942", &ObservationPressure{
			dt: time.Date(2023, 0o5, 15, 20, 10, 0, 0, time.UTC),
			v:  1015.5,
		}},
		{"K-Stammheim", "H744", nil},
		{"All data fields", "all", &ObservationPressure{
			dt: time.Date(2023, 0o5, 17, 7, 40, 0, 0, time.UTC),
			v:  1026.3,
		}},
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
			if tc.p != nil && tc.p.String() != o.PressureMSL().String() {
				t.Errorf("ObservationLatestByStationID failed, expected pressure MSL "+
					"string: %s, got: %s", tc.p.String(), o.PressureMSL())
			}
			if tc.p != nil && tc.p.Value() != o.PressureMSL().Value() {
				t.Errorf("ObservationLatestByStationID failed, expected pressure MSL "+
					"float: %f, got: %f", tc.p.Value(), o.PressureMSL().Value())
			}
			if tc.p != nil && tc.p.dt.Unix() != o.PressureMSL().Datetime().Unix() {
				t.Errorf("ObservationLatestByStationID failed, expected datetime: %s, got: %s",
					tc.p.dt.Format(time.RFC3339), o.PressureMSL().Datetime().Format(time.RFC3339))
			}
			if tc.p == nil {
				if o.PressureMSL().IsAvailable() {
					t.Errorf("ObservationLatestByStationID failed, expected pressure MSL "+
						"to have no data, but got: %s", o.PressureMSL())
				}
				if !math.IsNaN(o.PressureMSL().Value()) {
					t.Errorf("ObservationLatestByStationID failed, expected pressure MSL "+
						"to return NaN, but got: %s", o.PressureMSL().String())
				}
			}
		})
	}
}

func TestClient_ObservationLatestByStationID_PressureQFE(t *testing.T) {
	tt := []struct {
		// Test name
		n string
		// Station ID
		sid string
		// Observation dewpoint
		p *ObservationPressure
	}{
		{"K-Botanischer Garten", "199942", &ObservationPressure{v: 1010.2}},
		{"K-Stammheim", "H744", nil},
		{"All data fields", "all", &ObservationPressure{v: 1020.9}},
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
			if tc.p != nil && tc.p.String() != o.PressureQFE().String() {
				t.Errorf("ObservationLatestByStationID failed, expected pressure QFE "+
					"string: %s, got: %s", tc.p.String(), o.PressureQFE())
			}
			if tc.p != nil && tc.p.Value() != o.PressureQFE().Value() {
				t.Errorf("ObservationLatestByStationID failed, expected pressure QFE "+
					"float: %f, got: %f", tc.p.Value(), o.PressureQFE().Value())
			}
			if tc.p == nil {
				if o.PressureQFE().IsAvailable() {
					t.Errorf("ObservationLatestByStationID failed, expected pressure QFE "+
						"to have no data, but got: %s", o.PressureMSL())
				}
				if !math.IsNaN(o.PressureQFE().Value()) {
					t.Errorf("ObservationLatestByStationID failed, expected pressure QFE "+
						"to return NaN, but got: %s", o.PressureQFE().String())
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
			ot := ObservationTemperature{v: tc.c}
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
