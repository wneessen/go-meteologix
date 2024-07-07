// SPDX-FileCopyrightText: 2023 Winni Neessen <wn@neessen.dev>
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

func TestClient_ObservationLatestByStationID_MockFail(t *testing.T) {
	c := New(withMockAPI())
	if c == nil {
		t.Errorf("failed to create new Client, got nil")
		return
	}
	_, err := c.ObservationLatestByStationID(" ")
	if err == nil {
		t.Errorf("ObservationLatestByStationID with non-sense station ID was supposed to fail, but didn't")
	}
}

func TestClient_ObservationLatestByLocation(t *testing.T) {
	ak := getAPIKeyFromEnv(t)
	if ak == "" {
		t.Skip("no API_KEY found in environment, skipping test")
	}
	c := New(WithAPIKey(ak))
	if c == nil {
		t.Errorf("failed to create new Client, got nil")
		return
	}
	o, s, err := c.ObservationLatestByLocation("Ehrenfeld, Germany")
	if err != nil {
		t.Errorf("ObservationLatestByLocation failed: %s", err)
		return
	}
	if o.Name != "Koeln-Botanischer Garten" {
		t.Errorf("ObservationLatestByLocation failed, expected name: %s, got: %s",
			"Koeln-Botanischer Garten", o.Name)
	}
	if o.StationID != s.ID {
		t.Errorf("ObservationLatestByLocation failed, expected ID: %s, got: %s",
			"Köln-Botanischer Garten", o.StationID)
	}
	if o.Altitude != nil && *o.Altitude != s.Altitude {
		t.Errorf("ObservationLatestByLocation failed, expected altitude: %d, got: %d",
			s.Altitude, *o.Altitude)
	}
	if o.Altitude == nil {
		t.Errorf("ObservationLatestByLocation failed, expected altitude, got nil")
	}
	if o.Latitude != 50.966700 {
		t.Errorf("ObservationLatestByLocation failed, expected latitude: %f, got: %f",
			50.966700, o.Latitude)
	}
	if o.Longitude != 6.966700 {
		t.Errorf("ObservationLatestByLocation failed, expected longitude: %f, got: %f",
			6.966700, o.Longitude)
	}
}

func TestClient_ObservationLatestByLocation_Fail(t *testing.T) {
	ak := getAPIKeyFromEnv(t)
	if ak == "" {
		t.Skip("no API_KEY found in environment, skipping test")
	}
	c := New(WithAPIKey(ak))
	if c == nil {
		t.Errorf("failed to create new Client, got nil")
		return
	}
	_, _, err := c.ObservationLatestByLocation("Timbugtu")
	if err == nil {
		t.Errorf("ObservationLatestByLocation with non-sense location was supposed to fail, but didn't")
	}
}

func TestClient_ObservationLatestByStationID_Dewpoint(t *testing.T) {
	tt := []struct {
		// Test name
		n string
		// Station ID
		sid string
		// Observation dewpoint
		dp *Temperature
	}{
		{"K-Botanischer Garten", "199942", &Temperature{
			dt:       time.Date(2023, 0o5, 15, 20, 10, 0, 0, time.UTC),
			floatVal: 10.1,
		}},
		{"K-Stammheim", "H744", &Temperature{
			dt:       time.Date(2023, 0o5, 15, 19, 30, 0, 0, time.UTC),
			floatVal: 9.7,
		}},
		{"All data fields", "all", &Temperature{
			dt:       time.Date(2023, 0o5, 17, 7, 40, 0, 0, time.UTC),
			floatVal: 6.5,
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
			if tc.dp != nil && tc.dp.dt.Unix() != o.Dewpoint().DateTime().Unix() {
				t.Errorf("ObservationLatestByStationID failed, expected datetime: %s, got: %s",
					tc.dp.dt.Format(time.RFC3339), o.Dewpoint().DateTime().Format(time.RFC3339))
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

func TestClient_ObservationLatestByStationID_DewpointMean(t *testing.T) {
	tt := []struct {
		// Test name
		n string
		// Station ID
		sid string
		// Observation dewpoint
		t *Temperature
	}{
		{"K-Botanischer Garten", "199942", nil},
		{"K-Stammheim", "H744", nil},
		{"All data fields", "all", &Temperature{floatVal: 8.3}},
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
			if tc.t != nil && tc.t.String() != o.DewpointMean().String() {
				t.Errorf("ObservationLatestByStationID failed, expected mean dewpoint "+
					"string: %s, got: %s", tc.t.String(), o.DewpointMean())
			}
			if tc.t != nil && tc.t.Value() != o.DewpointMean().Value() {
				t.Errorf("ObservationLatestByStationID failed, expected mean dewpoint "+
					"float: %f, got: %f", tc.t.Value(), o.DewpointMean().Value())
			}
			if tc.t == nil {
				if o.DewpointMean().IsAvailable() {
					t.Errorf("ObservationLatestByStationID failed, expected mean dewpoint "+
						"to have no data, but got: %s", o.DewpointMean())
				}
				if !math.IsNaN(o.DewpointMean().Value()) {
					t.Errorf("ObservationLatestByStationID failed, expected mean dewpoint "+
						"to return NaN, but got: %s", o.DewpointMean().String())
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
		h *Humidity
	}{
		{"K-Botanischer Garten", "199942", &Humidity{
			dt:       time.Date(2023, 0o5, 15, 20, 10, 0, 0, time.UTC),
			floatVal: 80,
		}},
		{"K-Stammheim", "H744", &Humidity{
			dt:       time.Date(2023, 0o5, 15, 19, 30, 0, 0, time.UTC),
			floatVal: 73,
		}},
		{"All data fields", "all", &Humidity{
			dt:       time.Date(2023, 0o5, 17, 7, 40, 0, 0, time.UTC),
			floatVal: 72,
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
			if tc.h != nil && tc.h.dt.Unix() != o.HumidityRelative().DateTime().Unix() {
				t.Errorf("ObservationLatestByStationID failed, expected datetime: %s, got: %s",
					tc.h.dt.Format(time.RFC3339), o.HumidityRelative().DateTime().Format(time.RFC3339))
			}
			if o.HumidityRelative().Source() != SourceObservation {
				t.Errorf("ObservationLatestByStationID failed, expected observation source, but got: %s",
					o.HumidityRelative().Source())
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
		p *Precipitation
	}{
		{"K-Botanischer Garten", "199942", &Precipitation{
			dt:       time.Date(2023, 0o5, 15, 18, 0, 0, 0, time.UTC),
			floatVal: 0,
		}},
		{"K-Stammheim", "H744", &Precipitation{
			dt:       time.Date(2023, 0o5, 15, 19, 30, 0, 0, time.UTC),
			floatVal: 0,
		}},
		{"All data fields", "all", &Precipitation{
			dt:       time.Date(2023, 0o5, 17, 7, 30, 0, 0, time.UTC),
			floatVal: 0.1,
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
			if tc.p != nil && tc.p.String() != o.Precipitation(TimespanCurrent).String() {
				t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
					"string: %s, got: %s", tc.p.String(), o.Precipitation(TimespanCurrent))
			}
			if tc.p != nil && tc.p.Value() != o.Precipitation(TimespanCurrent).Value() {
				t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
					"float: %f, got: %f", tc.p.Value(), o.Precipitation(TimespanCurrent).Value())
			}
			if tc.p != nil && tc.p.dt.Unix() != o.Precipitation(TimespanCurrent).DateTime().Unix() {
				t.Errorf("ObservationLatestByStationID failed, expected datetime: %s, got: %s",
					tc.p.dt.Format(time.RFC3339),
					o.Precipitation(TimespanCurrent).DateTime().Format(time.RFC3339))
			}
			if o.Precipitation(TimespanCurrent).Source() != SourceObservation {
				t.Errorf("ObservationLatestByStationID failed, expected observation source, but got: %s",
					o.Precipitation(TimespanCurrent).Source())
			}
			if tc.p == nil {
				if o.Precipitation(TimespanCurrent).IsAvailable() {
					t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
						"to have no data, but got: %s", o.Precipitation(TimespanCurrent))
				}
				if !math.IsNaN(o.Precipitation(TimespanCurrent).Value()) {
					t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
						"to return NaN, but got: %s", o.Precipitation(TimespanCurrent).String())
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
		p *Precipitation
	}{
		{"K-Botanischer Garten", "199942", &Precipitation{floatVal: 0}},
		{"K-Stammheim", "H744", &Precipitation{floatVal: 0}},
		{"All data fields", "all", &Precipitation{floatVal: 0.5}},
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
			if tc.p != nil && tc.p.String() != o.Precipitation(Timespan10Min).String() {
				t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
					"string: %s, got: %s", tc.p.String(), o.Precipitation(Timespan10Min))
			}
			if tc.p != nil && tc.p.Value() != o.Precipitation(Timespan10Min).Value() {
				t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
					"float: %f, got: %f", tc.p.Value(), o.Precipitation(Timespan10Min).Value())
			}
			if tc.p == nil {
				if o.Precipitation(Timespan10Min).IsAvailable() {
					t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
						"to have no data, but got: %s", o.Precipitation(Timespan10Min))
				}
				if !math.IsNaN(o.Precipitation(Timespan10Min).Value()) {
					t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
						"to return NaN, but got: %s", o.Precipitation(Timespan10Min).String())
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
		p *Precipitation
	}{
		{"K-Botanischer Garten", "199942", &Precipitation{floatVal: 0}},
		{"K-Stammheim", "H744", &Precipitation{floatVal: 0}},
		{"All data fields", "all", &Precipitation{floatVal: 10.3}},
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
			if tc.p != nil && tc.p.String() != o.Precipitation(Timespan1Hour).String() {
				t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
					"string: %s, got: %s", tc.p.String(), o.Precipitation(Timespan1Hour))
			}
			if tc.p != nil && tc.p.Value() != o.Precipitation(Timespan1Hour).Value() {
				t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
					"float: %f, got: %f", tc.p.Value(), o.Precipitation(Timespan1Hour).Value())
			}
			if tc.p == nil {
				if o.Precipitation(Timespan1Hour).IsAvailable() {
					t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
						"to have no data, but got: %s", o.Precipitation(Timespan1Hour))
				}
				if !math.IsNaN(o.Precipitation(Timespan1Hour).Value()) {
					t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
						"to return NaN, but got: %s", o.Precipitation(Timespan1Hour).String())
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
		p *Precipitation
	}{
		{"K-Botanischer Garten", "199942", &Precipitation{floatVal: 0}},
		{"K-Stammheim", "H744", &Precipitation{floatVal: 0}},
		{"All data fields", "all", &Precipitation{floatVal: 32.12}},
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
			if tc.p != nil && tc.p.String() != o.Precipitation(Timespan24Hours).String() {
				t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
					"string: %s, got: %s", tc.p.String(), o.Precipitation(Timespan24Hours))
			}
			if tc.p != nil && tc.p.Value() != o.Precipitation(Timespan24Hours).Value() {
				t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
					"float: %f, got: %f", tc.p.Value(), o.Precipitation(Timespan24Hours).Value())
			}
			if tc.p == nil {
				if o.Precipitation(Timespan24Hours).IsAvailable() {
					t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
						"to have no data, but got: %s", o.Precipitation(Timespan24Hours))
				}
				if !math.IsNaN(o.Precipitation(Timespan24Hours).Value()) {
					t.Errorf("ObservationLatestByStationID failed, expected precipitation "+
						"to return NaN, but got: %s", o.Precipitation(Timespan24Hours).String())
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
		t *Temperature
	}{
		{"K-Botanischer Garten", "199942", &Temperature{floatVal: 13.4}},
		{"K-Stammheim", "H744", &Temperature{floatVal: 14.4}},
		{"All data fields", "all", &Temperature{floatVal: 10.8}},
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
			if o.Temperature().Source() != SourceObservation {
				t.Errorf("ObservationLatestByStationID failed, expected observation source, but got: %s",
					o.Temperature().Source())
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
		t *Temperature
	}{
		{"K-Botanischer Garten", "199942", nil},
		{"K-Stammheim", "H744", &Temperature{floatVal: 14.3}},
		{"All data fields", "all", &Temperature{floatVal: 15.4}},
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
		t *Temperature
	}{
		{"K-Botanischer Garten", "199942", &Temperature{floatVal: 12.3}},
		{"K-Stammheim", "H744", &Temperature{floatVal: 11.9}},
		{"All data fields", "all", &Temperature{floatVal: 6.2}},
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
		t *Temperature
	}{
		{"K-Botanischer Garten", "199942", &Temperature{floatVal: 20.5}},
		{"K-Stammheim", "H744", &Temperature{floatVal: 20.7}},
		{"All data fields", "all", &Temperature{floatVal: 12.4}},
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
		t *Temperature
	}{
		{"K-Botanischer Garten", "199942", nil},
		{"K-Stammheim", "H744", &Temperature{floatVal: 12.8}},
		{"All data fields", "all", &Temperature{floatVal: 3.7}},
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

func TestClient_ObservationLatestByStationID_TemperatureMean(t *testing.T) {
	tt := []struct {
		// Test name
		n string
		// Station ID
		sid string
		// Observation dewpoint
		t *Temperature
	}{
		{"K-Botanischer Garten", "199942", nil},
		{"K-Stammheim", "H744", nil},
		{"All data fields", "all", &Temperature{floatVal: 16.3}},
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
			if tc.t != nil && tc.t.String() != o.TemperatureMean().String() {
				t.Errorf("ObservationLatestByStationID failed, expected mean temperature "+
					"string: %s, got: %s", tc.t.String(), o.TemperatureMean())
			}
			if tc.t != nil && tc.t.Value() != o.TemperatureMean().Value() {
				t.Errorf("ObservationLatestByStationID failed, expected mean temperature "+
					"float: %f, got: %f", tc.t.Value(), o.TemperatureMean().Value())
			}
			if tc.t == nil {
				if o.TemperatureMean().IsAvailable() {
					t.Errorf("ObservationLatestByStationID failed, expected mean temperature "+
						"to have no data, but got: %s", o.TemperatureMean())
				}
				if !math.IsNaN(o.TemperatureMean().Value()) {
					t.Errorf("ObservationLatestByStationID failed, expected mean temperature "+
						"to return NaN, but got: %s", o.TemperatureMean().String())
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
		p *Pressure
	}{
		{"K-Botanischer Garten", "199942", &Pressure{
			dt:       time.Date(2023, 0o5, 15, 20, 10, 0, 0, time.UTC),
			floatVal: 1015.5,
		}},
		{"K-Stammheim", "H744", nil},
		{"All data fields", "all", &Pressure{
			dt:       time.Date(2023, 0o5, 17, 7, 40, 0, 0, time.UTC),
			floatVal: 1026.3,
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
			if tc.p != nil && tc.p.dt.Unix() != o.PressureMSL().DateTime().Unix() {
				t.Errorf("ObservationLatestByStationID failed, expected datetime: %s, got: %s",
					tc.p.dt.Format(time.RFC3339), o.PressureMSL().DateTime().Format(time.RFC3339))
			}
			if o.PressureMSL().Source() != SourceObservation {
				t.Errorf("ObservationLatestByStationID failed, expected observation source, but got: %s",
					o.PressureMSL().Source())
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
		p *Pressure
	}{
		{"K-Botanischer Garten", "199942", &Pressure{floatVal: 1010.2}},
		{"K-Stammheim", "H744", nil},
		{"All data fields", "all", &Pressure{floatVal: 1020.9}},
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

func TestClient_ObservationLatestByStationID_GlobalRadiationCurrent(t *testing.T) {
	c := New(withMockAPI())
	if c == nil {
		t.Errorf("failed to create new Client, got nil")
		return
	}
	o, err := c.ObservationLatestByStationID("199942")
	if err != nil {
		t.Errorf("ObservationLatestByStationID with station %s "+
			"failed: %s", "199942", err)
		return
	}
	if o.GlobalRadiation(TimespanCurrent).IsAvailable() {
		t.Errorf("ObservationLatestByStationID failed, expected glob. radiation "+
			"to have no data, but got: %s", o.GlobalRadiation(TimespanCurrent))
	}
	if !math.IsNaN(o.GlobalRadiation(TimespanCurrent).Value()) {
		t.Errorf("ObservationLatestByStationID failed, expected glob. radiation "+
			"to return NaN, but got: %s", o.GlobalRadiation(TimespanCurrent).String())
	}
}

func TestClient_ObservationLatestByStationID_GlobalRadiation10m(t *testing.T) {
	tt := []struct {
		// Test name
		n string
		// Station ID
		sid string
		// Observation radiation
		p *Radiation
	}{
		{"K-Botanischer Garten", "199942", &Radiation{
			dt:       time.Date(2023, 0o5, 15, 20, 10, 0, 0, time.UTC),
			floatVal: 0,
		}},
		{"K-Stammheim", "H744", nil},
		{"All data fields", "all", &Radiation{
			dt:       time.Date(2023, 0o5, 17, 7, 40, 0, 0, time.UTC),
			floatVal: 62,
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
			if tc.p != nil && tc.p.String() != o.GlobalRadiation(Timespan10Min).String() {
				t.Errorf("ObservationLatestByStationID failed, expected glob. radiation "+
					"string: %s, got: %s", tc.p.String(), o.GlobalRadiation(Timespan10Min))
			}
			if tc.p != nil && tc.p.Value() != o.GlobalRadiation(Timespan10Min).Value() {
				t.Errorf("ObservationLatestByStationID failed, expected glob. radiation "+
					"float: %f, got: %f", tc.p.Value(), o.GlobalRadiation(Timespan10Min).Value())
			}
			if tc.p != nil && tc.p.dt.Unix() != o.GlobalRadiation(Timespan10Min).DateTime().Unix() {
				t.Errorf("ObservationLatestByStationID failed, expected datetime: %s, got: %s",
					tc.p.dt.Format(time.RFC3339), o.GlobalRadiation(Timespan10Min).DateTime().Format(time.RFC3339))
			}
			if o.GlobalRadiation(Timespan10Min).Source() != SourceObservation {
				t.Errorf("ObservationLatestByStationID failed, expected observation source, but got: %s",
					o.GlobalRadiation(Timespan10Min).Source())
			}
			if tc.p == nil {
				if o.GlobalRadiation(Timespan10Min).IsAvailable() {
					t.Errorf("ObservationLatestByStationID failed, expected glob. radiation "+
						"to have no data, but got: %s", o.GlobalRadiation(Timespan10Min))
				}
				if !math.IsNaN(o.GlobalRadiation(Timespan10Min).Value()) {
					t.Errorf("ObservationLatestByStationID failed, expected glob. radiation "+
						"to return NaN, but got: %s", o.GlobalRadiation(Timespan10Min).String())
				}
			}
		})
	}
}

func TestClient_ObservationLatestByStationID_GlobalRadiation1h(t *testing.T) {
	tt := []struct {
		// Test name
		n string
		// Station ID
		sid string
		// Observation radiation
		p *Radiation
	}{
		{"K-Botanischer Garten", "199942", &Radiation{floatVal: 0}},
		{"K-Stammheim", "H744", nil},
		{"All data fields", "all", &Radiation{floatVal: 200}},
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
			if tc.p != nil && tc.p.String() != o.GlobalRadiation(Timespan1Hour).String() {
				t.Errorf("ObservationLatestByStationID failed, expected glob. radiation "+
					"string: %s, got: %s", tc.p.String(), o.GlobalRadiation(Timespan1Hour))
			}
			if tc.p != nil && tc.p.Value() != o.GlobalRadiation(Timespan1Hour).Value() {
				t.Errorf("ObservationLatestByStationID failed, expected glob. radiation "+
					"float: %f, got: %f", tc.p.Value(), o.GlobalRadiation(Timespan1Hour).Value())
			}
			if tc.p == nil {
				if o.GlobalRadiation(Timespan1Hour).IsAvailable() {
					t.Errorf("ObservationLatestByStationID failed, expected glob. radiation "+
						"to have no data, but got: %s", o.GlobalRadiation(Timespan1Hour))
				}
				if !math.IsNaN(o.GlobalRadiation(Timespan1Hour).Value()) {
					t.Errorf("ObservationLatestByStationID failed, expected glob. radiation "+
						"to return NaN, but got: %s", o.GlobalRadiation(Timespan1Hour).String())
				}
			}
		})
	}
}

func TestClient_ObservationLatestByStationID_GlobalRadiation24h(t *testing.T) {
	tt := []struct {
		// Test name
		n string
		// Station ID
		sid string
		// Observation radiation
		p *Radiation
	}{
		{"K-Botanischer Garten", "199942", &Radiation{floatVal: 774}},
		{"K-Stammheim", "H744", nil},
		{"All data fields", "all", &Radiation{floatVal: 756}},
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
			if tc.p != nil && tc.p.String() != o.GlobalRadiation(Timespan24Hours).String() {
				t.Errorf("ObservationLatestByStationID failed, expected glob. radiation "+
					"string: %s, got: %s", tc.p.String(), o.GlobalRadiation(Timespan24Hours))
			}
			if tc.p != nil && tc.p.Value() != o.GlobalRadiation(Timespan24Hours).Value() {
				t.Errorf("ObservationLatestByStationID failed, expected glob. radiation "+
					"float: %f, got: %f", tc.p.Value(), o.GlobalRadiation(Timespan24Hours).Value())
			}
			if tc.p == nil {
				if o.GlobalRadiation(Timespan24Hours).IsAvailable() {
					t.Errorf("ObservationLatestByStationID failed, expected glob. radiation "+
						"to have no data, but got: %s", o.GlobalRadiation(Timespan24Hours))
				}
				if !math.IsNaN(o.GlobalRadiation(Timespan24Hours).Value()) {
					t.Errorf("ObservationLatestByStationID failed, expected glob. radiation "+
						"to return NaN, but got: %s", o.GlobalRadiation(Timespan24Hours).String())
				}
			}
		})
	}
}

func TestClient_ObservationLatestByStationID_WindDirection(t *testing.T) {
	tt := []struct {
		// Test name
		n string
		// Station ID
		sid string
		// Observation dewpoint
		p *Direction
	}{
		{"K-Botanischer Garten", "199942", nil},
		{"K-Stammheim", "H744", nil},
		{"All data fields", "all", &Direction{
			dt:       time.Date(2023, 0o5, 21, 11, 30, 0, 0, time.UTC),
			floatVal: 90,
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
			if tc.p != nil && tc.p.String() != o.WindDirection().String() {
				t.Errorf("ObservationLatestByStationID failed, expected wind direction "+
					"string: %s, got: %s", tc.p.String(), o.WindDirection())
			}
			if tc.p != nil && tc.p.Value() != o.WindDirection().Value() {
				t.Errorf("ObservationLatestByStationID failed, expected wind direction "+
					"float: %f, got: %f", tc.p.Value(), o.WindDirection().Value())
			}
			if tc.p != nil && tc.p.dt.Unix() != o.WindDirection().DateTime().Unix() {
				t.Errorf("ObservationLatestByStationID failed, expected datetime: %s, got: %s",
					tc.p.dt.Format(time.RFC3339), o.WindDirection().DateTime().Format(time.RFC3339))
			}
			if o.WindDirection().Source() != SourceObservation {
				t.Errorf("ObservationLatestByStationID failed, expected observation source, but got: %s",
					o.WindDirection().Source())
			}
			if tc.p == nil {
				if o.WindDirection().IsAvailable() {
					t.Errorf("ObservationLatestByStationID failed, expected wind direction "+
						"to have no data, but got: %s", o.WindDirection())
				}
				if !math.IsNaN(o.WindDirection().Value()) {
					t.Errorf("ObservationLatestByStationID failed, expected wind direction "+
						"to return NaN, but got: %s", o.WindDirection().String())
				}
			}
		})
	}
}

func TestClient_ObservationLatestByStationID_WindSpeed(t *testing.T) {
	tt := []struct {
		// Test name
		n string
		// Station ID
		sid string
		// Observation dewpoint
		p *Speed
	}{
		{"K-Botanischer Garten", "199942", nil},
		{"K-Stammheim", "H744", nil},
		{"All data fields", "all", &Speed{
			dt:       time.Date(2023, 0o5, 21, 11, 30, 0, 0, time.UTC),
			floatVal: 7.716666666,
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
			if tc.p != nil && tc.p.String() != o.WindSpeed().String() {
				t.Errorf("ObservationLatestByStationID failed, expected windspeed "+
					"string: %s, got: %s", tc.p.String(), o.WindSpeed())
			}
			if tc.p != nil && tc.p.Value() != o.WindSpeed().Value() {
				t.Errorf("ObservationLatestByStationID failed, expected windspeed "+
					"float: %f, got: %f, %+v", tc.p.Value(), o.WindSpeed().Value(), o.Data.WindSpeed)
			}
			if tc.p != nil && tc.p.dt.Unix() != o.WindSpeed().DateTime().Unix() {
				t.Errorf("ObservationLatestByStationID failed, expected datetime: %s, got: %s",
					tc.p.dt.Format(time.RFC3339), o.WindSpeed().DateTime().Format(time.RFC3339))
			}
			if o.WindSpeed().Source() != SourceObservation {
				t.Errorf("ObservationLatestByStationID failed, expected observation source, but got: %s",
					o.WindSpeed().Source())
			}
			if tc.p == nil {
				if o.WindSpeed().IsAvailable() {
					t.Errorf("ObservationLatestByStationID failed, expected windspeed "+
						"to have no data, but got: %s", o.WindSpeed())
				}
				if !math.IsNaN(o.WindSpeed().Value()) {
					t.Errorf("ObservationLatestByStationID failed, expected windspeed "+
						"to return NaN, but got: %s", o.WindSpeed().String())
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
			ot := Temperature{floatVal: tc.c}
			if ot.Celsius() != tc.c {
				t.Errorf("Temperature.Celsius failed, expected: %f, got: %f", tc.c,
					ot.Celsius())
			}
			if ot.CelsiusString() != fmt.Sprintf(cf, tc.c) {
				t.Errorf("Temperature.CelsiusString failed, expected: %s, got: %s",
					fmt.Sprintf(cf, tc.c), ot.CelsiusString())
			}
			if ot.Fahrenheit() != tc.f {
				t.Errorf("Temperature.Fahrenheit failed, expected: %f, got: %f", tc.f,
					ot.Fahrenheit())
			}
			if ot.FahrenheitString() != fmt.Sprintf(ff, tc.f) {
				t.Errorf("Temperature.FahrenheitString failed, expected: %s, got: %s",
					fmt.Sprintf(ff, tc.f), ot.FahrenheitString())
			}
		})
	}
}

func TestObservationSpeed_Conversion(t *testing.T) {
	tt := []struct {
		// Original m/s value
		ms float64
		// km/h value
		kmh float64
		// mi/h value
		mph float64
		// knots value
		kn float64
	}{
		{0, 0, 0, 0},
		{1, 3.6, 2.236936, 1.9438444924},
		{10, 36, 22.369360, 19.438444924},
		{15, 54, 33.554040, 29.157667386},
		{30, 108, 67.108080, 58.315334772},
	}
	msf := "%.1fm/s"
	knf := "%.0fkn"
	kmhf := "%.1fkm/h"
	mphf := "%.1fmi/h"
	for _, tc := range tt {
		t.Run(fmt.Sprintf("%.0fm/s", tc.ms), func(t *testing.T) {
			os := Speed{floatVal: tc.ms}
			if os.Value() != tc.ms {
				t.Errorf("Speed.Value failed, expected: %f, got: %f", tc.ms,
					os.Value())
			}
			if os.String() != fmt.Sprintf(msf, tc.ms) {
				t.Errorf("Speed.String failed, expected: %s, got: %s",
					fmt.Sprintf(msf, tc.ms), os.String())
			}
			if os.KMH() != tc.kmh {
				t.Errorf("Speed.KMH failed, expected: %f, got: %f", tc.kmh,
					os.KMH())
			}
			if os.KMHString() != fmt.Sprintf(kmhf, tc.kmh) {
				t.Errorf("Speed.KMHString failed, expected: %s, got: %s",
					fmt.Sprintf(kmhf, tc.kmh), os.KMHString())
			}
			if os.MPH() != tc.mph {
				t.Errorf("Speed.MPH failed, expected: %f, got: %f", tc.mph,
					os.MPH())
			}
			if os.MPHString() != fmt.Sprintf(mphf, tc.mph) {
				t.Errorf("Speed.MPHString failed, expected: %s, got: %s",
					fmt.Sprintf(mphf, tc.mph), os.MPHString())
			}
			if os.Knots() != tc.kn {
				t.Errorf("Speed.Knots failed, expected: %f, got: %f", tc.kn,
					os.Knots())
			}
			if os.KnotsString() != fmt.Sprintf(knf, tc.kn) {
				t.Errorf("Speed.KnotsString failed, expected: %s, got: %s",
					fmt.Sprintf(knf, tc.kn), os.KnotsString())
			}
		})
	}
}

func TestObservationDirection_Direction(t *testing.T) {
	tt := []struct {
		// Original direction in degree
		d float64
		// Direction string
		ds string
	}{
		{0, "N"},
		{11.25, "NbE"},
		{22.5, "NNE"},
		{33.75, "NEbN"},
		{45, "NE"},
		{56.25, "NEbE"},
		{67.5, "ENE"},
		{78.75, "EbN"},
		{90, "E"},
		{101.25, "EbS"},
		{112.5, "ESE"},
		{123.75, "SEbE"},
		{135, "SE"},
		{146.25, "SEbS"},
		{157.5, "SSE"},
		{168.75, "SbE"},
		{180, "S"},
		{191.25, "SbW"},
		{202.5, "SSW"},
		{213.75, "SWbS"},
		{225, "SW"},
		{236.25, "SWbW"},
		{247.5, "WSW"},
		{258.75, "WbS"},
		{270, "W"},
		{281.25, "WbN"},
		{292.5, "WNW"},
		{303.75, "NWbW"},
		{315, "NW"},
		{326.25, "NWbN"},
		{337.5, "NNW"},
		{348.75, "NbW"},
		{999, ErrUnsupportedDirection},
	}
	for _, tc := range tt {
		t.Run(fmt.Sprintf("%.2f° => %s", tc.d, tc.ds), func(t *testing.T) {
			d := Direction{floatVal: tc.d}
			if d.Direction() != tc.ds {
				t.Errorf("Direction.Direction failed, expected: %s, got: %s",
					tc.ds, d.Direction())
			}
		})
	}
}

func TestObservationDirection_DirectionFull(t *testing.T) {
	tt := []struct {
		// Original direction in degree
		d float64
		// Direction string
		ds string
	}{
		{0, "North"},
		{11.25, "North by East"},
		{22.5, "North-Northeast"},
		{33.75, "Northeast by North"},
		{45, "Northeast"},
		{56.25, "Northeast by East"},
		{67.5, "East-Northeast"},
		{78.75, "East by North"},
		{90, "East"},
		{101.25, "East by South"},
		{112.5, "East-Southeast"},
		{123.75, "Southeast by East"},
		{135, "Southeast"},
		{146.25, "Southeast by South"},
		{157.5, "South-Southeast"},
		{168.75, "South by East"},
		{180, "South"},
		{191.25, "South by West"},
		{202.5, "South-Southwest"},
		{213.75, "Southwest by South"},
		{225, "Southwest"},
		{236.25, "Southwest by West"},
		{247.5, "West-Southwest"},
		{258.75, "West by South"},
		{270, "West"},
		{281.25, "West by North"},
		{292.5, "West-Northwest"},
		{303.75, "Northwest by West"},
		{315, "Northwest"},
		{326.25, "Northwest by North"},
		{337.5, "North-Northwest"},
		{348.75, "North by West"},
		{999, ErrUnsupportedDirection},
	}
	for _, tc := range tt {
		t.Run(fmt.Sprintf("%.2f° => %s", tc.d, tc.ds), func(t *testing.T) {
			d := Direction{floatVal: tc.d}
			if d.DirectionFull() != tc.ds {
				t.Errorf("Direction.Direction failed, expected: %s, got: %s",
					tc.ds, d.DirectionFull())
			}
		})
	}
}
