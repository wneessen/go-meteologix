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

func TestClient_CurrentWeatherByCoordinates_Mock(t *testing.T) {
	tt := []struct {
		// Latitude
		lat float64
		// Longitude
		lon float64
		// us
		us string
	}{
		{50.9833, 6.9833, "metric"},
	}
	c := New(withMockAPI())
	if c == nil {
		t.Errorf("failed to create new Client, got nil")
		return
	}
	for _, tc := range tt {
		t.Run(fmt.Sprintf("%.3f/%.3f", tc.lat, tc.lon), func(t *testing.T) {
			cw, err := c.CurrentWeatherByCoordinates(tc.lat, tc.lon)
			if err != nil {
				t.Errorf("CurrentWeatherByCoordinates failed: %s", err)
				return
			}
			if cw.Latitude != tc.lat {
				t.Errorf("CurrentWeatherByCoordinates failed, expected latitude: %f, got: %f", tc.lat,
					cw.Latitude)
			}
			if cw.Longitude != tc.lon {
				t.Errorf("CurrentWeatherByCoordinates failed, expected longitude: %f, got: %f", tc.lon,
					cw.Longitude)
			}
			if cw.UnitSystem != tc.us {
				t.Errorf("CurrentWeatherByCoordinates failed, expected unit system: %s, got: %s", tc.us,
					cw.UnitSystem)
			}
		})
	}
}

func TestClient_CurrentWeatherByLocation(t *testing.T) {
	tt := []struct {
		// Location string
		loc string
		// Latitude
		lat float64
		// Longitude
		lon float64
		// us
		us string
	}{
		{"Ehrenfeld, Germany", 50.9833, 6.9833, "metric"},
	}
	c := New(withMockAPI())
	if c == nil {
		t.Errorf("failed to create new Client, got nil")
		return
	}
	for _, tc := range tt {
		t.Run(tc.loc, func(t *testing.T) {
			cw, err := c.CurrentWeatherByLocation(tc.loc)
			if err != nil {
				t.Errorf("CurrentWeatherByCoordinates failed: %s", err)
				return
			}
			if cw.Latitude != tc.lat {
				t.Errorf("CurrentWeatherByCoordinates failed, expected latitude: %f, got: %f", tc.lat,
					cw.Latitude)
			}
			if cw.Longitude != tc.lon {
				t.Errorf("CurrentWeatherByCoordinates failed, expected longitude: %f, got: %f", tc.lon,
					cw.Longitude)
			}
			if cw.UnitSystem != tc.us {
				t.Errorf("CurrentWeatherByCoordinates failed, expected unit system: %s, got: %s", tc.us,
					cw.UnitSystem)
			}
		})
	}
}

func TestClient_CurrentWeatherByLocation_Fail(t *testing.T) {
	c := New(withMockAPI())
	if c == nil {
		t.Errorf("failed to create new Client, got nil")
		return
	}
	_, err := c.CurrentWeatherByLocation("Timbucktu, Atlantis")
	if err == nil {
		t.Errorf("CurrentWeatherByCoordinates was supposed to fail, but didn't")
	}
	_, err = c.CurrentWeatherByLocation("")
	if err == nil {
		t.Errorf("CurrentWeatherByCoordinates was supposed to fail, but didn't")
	}
}

func TestClient_CurrentWeatherByLocation_Temperature(t *testing.T) {
	tt := []struct {
		// Location name
		loc string
		// CurWeather temperature
		t *Temperature
	}{
		{"Ehrenfeld, Germany", &Temperature{
			dt: time.Date(2023, 5, 23, 7, 0, 0, 0, time.Local),
			s:  SourceObservation,
			v:  14.6,
		}},
		{"Berlin, Germany", &Temperature{
			dt: time.Date(2023, 5, 23, 7, 0, 0, 0, time.Local),
			s:  SourceAnalysis,
			v:  17.8,
		}},
	}
	c := New(withMockAPI())
	if c == nil {
		t.Errorf("failed to create new Client, got nil")
		return
	}
	for _, tc := range tt {
		t.Run(tc.loc, func(t *testing.T) {
			o, err := c.CurrentWeatherByLocation(tc.loc)
			if err != nil {
				t.Errorf("CurrentWeatherByLocation failed: %s", err)
				return
			}
			if tc.t != nil && tc.t.String() != o.Temperature().String() {
				t.Errorf("CurrentWeatherByLocation failed, expected temperature "+
					"string: %s, got: %s", tc.t.String(), o.Temperature())
			}
			if tc.t != nil && tc.t.Value() != o.Temperature().Value() {
				t.Errorf("CurrentWeatherByLocation failed, expected temperature "+
					"float: %f, got: %f", tc.t.Value(), o.Temperature().Value())
			}
			if tc.t != nil && o.Temperature().Source() != tc.t.s {
				t.Errorf("CurrentWeatherByLocation failed, expected source: %s, but got: %s",
					tc.t.s, o.Temperature().Source())
			}
			if tc.t == nil {
				if o.Temperature().IsAvailable() {
					t.Errorf("CurrentWeatherByLocation failed, expected temperature "+
						"to have no data, but got: %s", o.Temperature())
				}
				if !math.IsNaN(o.Temperature().Value()) {
					t.Errorf("CurrentWeatherByLocation failed, expected temperature "+
						"to return NaN, but got: %s", o.Temperature().String())
				}
			}
		})
	}
}

func TestClient_CurrentWeatherByLocation_Dewpoint(t *testing.T) {
	tt := []struct {
		// Location name
		loc string
		// CurWeather dewpoint
		t *Temperature
	}{
		{"Ehrenfeld, Germany", &Temperature{
			dt: time.Date(2023, 5, 23, 7, 0, 0, 0, time.Local),
			s:  SourceObservation,
			v:  11.5,
		}},
		{"Berlin, Germany", &Temperature{
			dt: time.Date(2023, 5, 23, 7, 0, 0, 0, time.Local),
			s:  SourceAnalysis,
			v:  11.0,
		}},
	}
	c := New(withMockAPI())
	if c == nil {
		t.Errorf("failed to create new Client, got nil")
		return
	}
	for _, tc := range tt {
		t.Run(tc.loc, func(t *testing.T) {
			o, err := c.CurrentWeatherByLocation(tc.loc)
			if err != nil {
				t.Errorf("CurrentWeatherByLocation failed: %s", err)
				return
			}
			if tc.t != nil && tc.t.String() != o.Dewpoint().String() {
				t.Errorf("CurrentWeatherByLocation failed, expected dewpoint "+
					"string: %s, got: %s", tc.t.String(), o.Dewpoint())
			}
			if tc.t != nil && tc.t.Value() != o.Dewpoint().Value() {
				t.Errorf("CurrentWeatherByLocation failed, expected dewpoint "+
					"float: %f, got: %f", tc.t.Value(), o.Dewpoint().Value())
			}
			if tc.t != nil && o.Dewpoint().Source() != tc.t.s {
				t.Errorf("CurrentWeatherByLocation failed, expected source: %s, but got: %s",
					tc.t.s, o.Dewpoint().Source())
			}
			if tc.t == nil {
				if o.Dewpoint().IsAvailable() {
					t.Errorf("CurrentWeatherByLocation failed, expected dewpoint "+
						"to have no data, but got: %s", o.Dewpoint())
				}
				if !math.IsNaN(o.Dewpoint().Value()) {
					t.Errorf("CurrentWeatherByLocation failed, expected dewpoint "+
						"to return NaN, but got: %s", o.Dewpoint().String())
				}
			}
		})
	}
}

func TestClient_CurrentWeatherByLocation_HumidityRelative(t *testing.T) {
	tt := []struct {
		// Location name
		loc string
		// CurWeather humidity
		h *Humidity
	}{
		{"Ehrenfeld, Germany", &Humidity{
			dt: time.Date(2023, 5, 23, 7, 0, 0, 0, time.Local),
			s:  SourceObservation,
			v:  82,
		}},
		{"Berlin, Germany", &Humidity{
			dt: time.Date(2023, 5, 23, 7, 0, 0, 0, time.Local),
			s:  SourceAnalysis,
			v:  64,
		}},
	}
	c := New(withMockAPI())
	if c == nil {
		t.Errorf("failed to create new Client, got nil")
		return
	}
	for _, tc := range tt {
		t.Run(tc.loc, func(t *testing.T) {
			o, err := c.CurrentWeatherByLocation(tc.loc)
			if err != nil {
				t.Errorf("CurrentWeatherByLocation failed: %s", err)
				return
			}
			if tc.h != nil && tc.h.String() != o.HumidityRelative().String() {
				t.Errorf("CurrentWeatherByLocation failed, expected humidity "+
					"string: %s, got: %s", tc.h.String(), o.HumidityRelative())
			}
			if tc.h != nil && tc.h.Value() != o.HumidityRelative().Value() {
				t.Errorf("CurrentWeatherByLocation failed, expected humidity "+
					"float: %f, got: %f", tc.h.Value(), o.HumidityRelative().Value())
			}
			if tc.h != nil && o.HumidityRelative().Source() != tc.h.s {
				t.Errorf("CurrentWeatherByLocation failed, expected source: %s, but got: %s",
					tc.h.s, o.HumidityRelative().Source())
			}
			if tc.h == nil {
				if o.HumidityRelative().IsAvailable() {
					t.Errorf("CurrentWeatherByLocation failed, expected humidity "+
						"to have no data, but got: %s", o.HumidityRelative())
				}
				if !math.IsNaN(o.HumidityRelative().Value()) {
					t.Errorf("CurrentWeatherByLocation failed, expected humidity "+
						"to return NaN, but got: %s", o.HumidityRelative().String())
				}
			}
		})
	}
}

func TestClient_CurrentWeatherByLocation_PrecipitationCurrent(t *testing.T) {
	tt := []struct {
		// Location name
		loc string
		// CurWeather precipitation
		p *Precipitation
	}{
		{"Ehrenfeld, Germany", nil},
		{"Berlin, Germany", nil},
	}
	c := New(withMockAPI())
	if c == nil {
		t.Errorf("failed to create new Client, got nil")
		return
	}
	for _, tc := range tt {
		t.Run(tc.loc, func(t *testing.T) {
			o, err := c.CurrentWeatherByLocation(tc.loc)
			if err != nil {
				t.Errorf("CurrentWeatherByLocation failed: %s", err)
				return
			}
			if tc.p != nil && tc.p.String() != o.Precipitation(TimespanCurrent).String() {
				t.Errorf("CurrentWeatherByLocation failed, expected precipitation "+
					"string: %s, got: %s", tc.p.String(), o.Precipitation(TimespanCurrent))
			}
			if tc.p != nil && tc.p.Value() != o.Precipitation(TimespanCurrent).Value() {
				t.Errorf("CurrentWeatherByLocation failed, expected precipitation "+
					"float: %f, got: %f", tc.p.Value(), o.Precipitation(TimespanCurrent).Value())
			}
			if tc.p != nil && o.Precipitation(TimespanCurrent).Source() != tc.p.s {
				t.Errorf("CurrentWeatherByLocation failed, expected source: %s, but got: %s",
					tc.p.s, o.Precipitation(TimespanCurrent).Source())
			}
			if tc.p == nil {
				if o.Precipitation(TimespanCurrent).IsAvailable() {
					t.Errorf("CurrentWeatherByLocation failed, expected precipitation "+
						"to have no data, but got: %s", o.Precipitation(TimespanCurrent))
				}
				if !math.IsNaN(o.Precipitation(TimespanCurrent).Value()) {
					t.Errorf("CurrentWeatherByLocation failed, expected precipitation "+
						"to return NaN, but got: %s", o.Precipitation(TimespanCurrent).String())
				}
			}
		})
	}
}

func TestClient_CurrentWeatherByLocation_Precipitation10m(t *testing.T) {
	tt := []struct {
		// Location name
		loc string
		// CurWeather precipitation
		p *Precipitation
	}{
		{"Ehrenfeld, Germany", nil},
		{"Berlin, Germany", nil},
	}
	c := New(withMockAPI())
	if c == nil {
		t.Errorf("failed to create new Client, got nil")
		return
	}
	for _, tc := range tt {
		t.Run(tc.loc, func(t *testing.T) {
			o, err := c.CurrentWeatherByLocation(tc.loc)
			if err != nil {
				t.Errorf("CurrentWeatherByLocation failed: %s", err)
				return
			}
			if tc.p != nil && tc.p.String() != o.Precipitation(Timespan10Min).String() {
				t.Errorf("CurrentWeatherByLocation failed, expected precipitation "+
					"string: %s, got: %s", tc.p.String(), o.Precipitation(Timespan10Min))
			}
			if tc.p != nil && tc.p.Value() != o.Precipitation(Timespan10Min).Value() {
				t.Errorf("CurrentWeatherByLocation failed, expected precipitation "+
					"float: %f, got: %f", tc.p.Value(), o.Precipitation(Timespan10Min).Value())
			}
			if tc.p != nil && o.Precipitation(Timespan10Min).Source() != tc.p.s {
				t.Errorf("CurrentWeatherByLocation failed, expected source: %s, but got: %s",
					tc.p.s, o.Precipitation(Timespan10Min).Source())
			}
			if tc.p == nil {
				if o.Precipitation(Timespan10Min).IsAvailable() {
					t.Errorf("CurrentWeatherByLocation failed, expected precipitation "+
						"to have no data, but got: %s", o.Precipitation(Timespan10Min))
				}
				if !math.IsNaN(o.Precipitation(Timespan10Min).Value()) {
					t.Errorf("CurrentWeatherByLocation failed, expected precipitation "+
						"to return NaN, but got: %s", o.Precipitation(Timespan10Min).String())
				}
			}
		})
	}
}

func TestClient_CurrentWeatherByLocation_Precipitation1h(t *testing.T) {
	tt := []struct {
		// Location name
		loc string
		// CurWeather precipitation
		p *Precipitation
	}{
		{"Ehrenfeld, Germany", &Precipitation{
			dt: time.Date(2023, 5, 23, 7, 0, 0, 0, time.Local),
			s:  SourceObservation,
			v:  0,
		}},
		{"Berlin, Germany", &Precipitation{
			dt: time.Date(2023, 5, 23, 7, 0, 0, 0, time.Local),
			s:  SourceAnalysis,
			v:  0.0092,
		}},
	}
	c := New(withMockAPI())
	if c == nil {
		t.Errorf("failed to create new Client, got nil")
		return
	}
	for _, tc := range tt {
		t.Run(tc.loc, func(t *testing.T) {
			o, err := c.CurrentWeatherByLocation(tc.loc)
			if err != nil {
				t.Errorf("CurrentWeatherByLocation failed: %s", err)
				return
			}
			if tc.p != nil && tc.p.String() != o.Precipitation(Timespan1Hour).String() {
				t.Errorf("CurrentWeatherByLocation failed, expected precipitation "+
					"string: %s, got: %s", tc.p.String(), o.Precipitation(Timespan1Hour))
			}
			if tc.p != nil && tc.p.Value() != o.Precipitation(Timespan1Hour).Value() {
				t.Errorf("CurrentWeatherByLocation failed, expected precipitation "+
					"float: %f, got: %f", tc.p.Value(), o.Precipitation(Timespan1Hour).Value())
			}
			if tc.p != nil && o.Precipitation(Timespan1Hour).Source() != tc.p.s {
				t.Errorf("CurrentWeatherByLocation failed, expected source: %s, but got: %s",
					tc.p.s, o.Precipitation(Timespan1Hour).Source())
			}
			if tc.p == nil {
				if o.Precipitation(Timespan1Hour).IsAvailable() {
					t.Errorf("CurrentWeatherByLocation failed, expected precipitation "+
						"to have no data, but got: %s", o.Precipitation(Timespan1Hour))
				}
				if !math.IsNaN(o.Precipitation(Timespan1Hour).Value()) {
					t.Errorf("CurrentWeatherByLocation failed, expected precipitation "+
						"to return NaN, but got: %s", o.Precipitation(Timespan1Hour).String())
				}
			}
		})
	}
}

func TestClient_CurrentWeatherByLocation_Precipitation24h(t *testing.T) {
	tt := []struct {
		// Location name
		loc string
		// CurWeather precipitation
		p *Precipitation
	}{
		{"Ehrenfeld, Germany", nil},
		{"Berlin, Germany", nil},
	}
	c := New(withMockAPI())
	if c == nil {
		t.Errorf("failed to create new Client, got nil")
		return
	}
	for _, tc := range tt {
		t.Run(tc.loc, func(t *testing.T) {
			o, err := c.CurrentWeatherByLocation(tc.loc)
			if err != nil {
				t.Errorf("CurrentWeatherByLocation failed: %s", err)
				return
			}
			if tc.p != nil && tc.p.String() != o.Precipitation(Timespan24Hours).String() {
				t.Errorf("CurrentWeatherByLocation failed, expected precipitation "+
					"string: %s, got: %s", tc.p.String(), o.Precipitation(Timespan24Hours))
			}
			if tc.p != nil && tc.p.Value() != o.Precipitation(Timespan24Hours).Value() {
				t.Errorf("CurrentWeatherByLocation failed, expected precipitation "+
					"float: %f, got: %f", tc.p.Value(), o.Precipitation(Timespan24Hours).Value())
			}
			if tc.p != nil && o.Precipitation(Timespan24Hours).Source() != tc.p.s {
				t.Errorf("CurrentWeatherByLocation failed, expected source: %s, but got: %s",
					tc.p.s, o.Precipitation(Timespan24Hours).Source())
			}
			if tc.p == nil {
				if o.Precipitation(Timespan24Hours).IsAvailable() {
					t.Errorf("CurrentWeatherByLocation failed, expected precipitation "+
						"to have no data, but got: %s", o.Precipitation(Timespan24Hours))
				}
				if !math.IsNaN(o.Precipitation(Timespan24Hours).Value()) {
					t.Errorf("CurrentWeatherByLocation failed, expected precipitation "+
						"to return NaN, but got: %s", o.Precipitation(Timespan24Hours).String())
				}
			}
		})
	}
}