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

func TestClient_CurrentWeatherByLocation_CloudCoverage(t *testing.T) {
	tt := []struct {
		// Location name
		loc string
		// CurWeather coverage
		co *Coverage
	}{
		{"Ehrenfeld, Germany", &Coverage{
			dateTime: time.Date(2023, 5, 23, 7, 0, 0, 0, time.Local),
			source:   SourceAnalysis,
			floatVal: 95,
		}},
		{"Berlin, Germany", &Coverage{
			dateTime: time.Date(2023, 5, 23, 7, 0, 0, 0, time.Local),
			source:   SourceAnalysis,
			floatVal: 60,
		}},
		{"Neermoor, Germany", nil},
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
				t.Errorf("CurrentWeatherByLocation failed: %s", err)
				return
			}
			if tc.co != nil && tc.co.String() != cw.CloudCoverage().String() {
				t.Errorf("CurrentWeatherByLocation failed, expected cloud coverage "+
					"string: %s, got: %s", tc.co.String(), cw.CloudCoverage())
			}
			if tc.co != nil && tc.co.Value() != cw.CloudCoverage().Value() {
				t.Errorf("CurrentWeatherByLocation failed, expected cloud coverage "+
					"float: %f, got: %f", tc.co.Value(), cw.CloudCoverage().Value())
			}
			if tc.co != nil && cw.CloudCoverage().Source() != tc.co.source {
				t.Errorf("CurrentWeatherByLocation failed, expected source: %s, but got: %s",
					tc.co.source, cw.CloudCoverage().Source())
			}
			if tc.co == nil {
				if cw.CloudCoverage().IsAvailable() {
					t.Errorf("CurrentWeatherByLocation failed, expected cloud coverage "+
						"to have no data, but got: %s", cw.CloudCoverage())
				}
				if !math.IsNaN(cw.CloudCoverage().Value()) {
					t.Errorf("CurrentWeatherByLocation failed, expected cloud coverage "+
						"to return NaN, but got: %s", cw.CloudCoverage().String())
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
			dateTime: time.Date(2023, 5, 23, 7, 0, 0, 0, time.Local),
			source:   SourceObservation,
			floatVal: 11.5,
		}},
		{"Berlin, Germany", &Temperature{
			dateTime: time.Date(2023, 5, 23, 7, 0, 0, 0, time.Local),
			source:   SourceAnalysis,
			floatVal: 11.0,
		}},
		{"Neermoor, Germany", nil},
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
				t.Errorf("CurrentWeatherByLocation failed: %s", err)
				return
			}
			if tc.t != nil && tc.t.String() != cw.Dewpoint().String() {
				t.Errorf("CurrentWeatherByLocation failed, expected dewpoint "+
					"string: %s, got: %s", tc.t.String(), cw.Dewpoint())
			}
			if tc.t != nil && tc.t.Value() != cw.Dewpoint().Value() {
				t.Errorf("CurrentWeatherByLocation failed, expected dewpoint "+
					"float: %f, got: %f", tc.t.Value(), cw.Dewpoint().Value())
			}
			if tc.t != nil && cw.Dewpoint().Source() != tc.t.source {
				t.Errorf("CurrentWeatherByLocation failed, expected source: %s, but got: %s",
					tc.t.source, cw.Dewpoint().Source())
			}
			if tc.t == nil {
				if cw.Dewpoint().IsAvailable() {
					t.Errorf("CurrentWeatherByLocation failed, expected dewpoint "+
						"to have no data, but got: %s", cw.Dewpoint())
				}
				if !math.IsNaN(cw.Dewpoint().Value()) {
					t.Errorf("CurrentWeatherByLocation failed, expected dewpoint "+
						"to return NaN, but got: %s", cw.Dewpoint().String())
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
			dateTime: time.Date(2023, 5, 23, 7, 0, 0, 0, time.Local),
			source:   SourceObservation,
			floatVal: 82,
		}},
		{"Berlin, Germany", &Humidity{
			dateTime: time.Date(2023, 5, 23, 7, 0, 0, 0, time.Local),
			source:   SourceAnalysis,
			floatVal: 64,
		}},
		{"Neermoor, Germany", nil},
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
				t.Errorf("CurrentWeatherByLocation failed: %s", err)
				return
			}
			if tc.h != nil && tc.h.String() != cw.HumidityRelative().String() {
				t.Errorf("CurrentWeatherByLocation failed, expected humidity "+
					"string: %s, got: %s", tc.h.String(), cw.HumidityRelative())
			}
			if tc.h != nil && tc.h.Value() != cw.HumidityRelative().Value() {
				t.Errorf("CurrentWeatherByLocation failed, expected humidity "+
					"float: %f, got: %f", tc.h.Value(), cw.HumidityRelative().Value())
			}
			if tc.h != nil && cw.HumidityRelative().Source() != tc.h.source {
				t.Errorf("CurrentWeatherByLocation failed, expected source: %s, but got: %s",
					tc.h.source, cw.HumidityRelative().Source())
			}
			if tc.h == nil {
				if cw.HumidityRelative().IsAvailable() {
					t.Errorf("CurrentWeatherByLocation failed, expected humidity "+
						"to have no data, but got: %s", cw.HumidityRelative())
				}
				if !math.IsNaN(cw.HumidityRelative().Value()) {
					t.Errorf("CurrentWeatherByLocation failed, expected humidity "+
						"to return NaN, but got: %s", cw.HumidityRelative().String())
				}
			}
		})
	}
}

func TestClient_CurrentWeatherByLocation_IsDay(t *testing.T) {
	tt := []struct {
		// Location name
		loc string
		// CurWeather IsDay
		d bool
	}{
		{"Ehrenfeld, Germany", false},
		{"Berlin, Germany", true},
		{"Neermoor, Germany", false},
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
				t.Errorf("CurrentWeatherByLocation failed: %s", err)
				return
			}
			if cw.IsDay() != tc.d {
				t.Errorf("CurrentWeather IsDay failed, expected: %t, got: %t", cw.IsDay(), tc.d)
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
		{"Neermoor, Germany", nil},
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
				t.Errorf("CurrentWeatherByLocation failed: %s", err)
				return
			}
			if tc.p != nil && tc.p.String() != cw.Precipitation(TimespanCurrent).String() {
				t.Errorf("CurrentWeatherByLocation failed, expected precipitation "+
					"string: %s, got: %s", tc.p.String(), cw.Precipitation(TimespanCurrent))
			}
			if tc.p != nil && tc.p.Value() != cw.Precipitation(TimespanCurrent).Value() {
				t.Errorf("CurrentWeatherByLocation failed, expected precipitation "+
					"float: %f, got: %f", tc.p.Value(), cw.Precipitation(TimespanCurrent).Value())
			}
			if tc.p != nil && cw.Precipitation(TimespanCurrent).Source() != tc.p.source {
				t.Errorf("CurrentWeatherByLocation failed, expected source: %s, but got: %s",
					tc.p.source, cw.Precipitation(TimespanCurrent).Source())
			}
			if tc.p == nil {
				if cw.Precipitation(TimespanCurrent).IsAvailable() {
					t.Errorf("CurrentWeatherByLocation failed, expected precipitation "+
						"to have no data, but got: %s", cw.Precipitation(TimespanCurrent))
				}
				if !math.IsNaN(cw.Precipitation(TimespanCurrent).Value()) {
					t.Errorf("CurrentWeatherByLocation failed, expected precipitation "+
						"to return NaN, but got: %s", cw.Precipitation(TimespanCurrent).String())
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
		{"Neermoor, Germany", nil},
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
				t.Errorf("CurrentWeatherByLocation failed: %s", err)
				return
			}
			if tc.p != nil && tc.p.String() != cw.Precipitation(Timespan10Min).String() {
				t.Errorf("CurrentWeatherByLocation failed, expected precipitation "+
					"string: %s, got: %s", tc.p.String(), cw.Precipitation(Timespan10Min))
			}
			if tc.p != nil && tc.p.Value() != cw.Precipitation(Timespan10Min).Value() {
				t.Errorf("CurrentWeatherByLocation failed, expected precipitation "+
					"float: %f, got: %f", tc.p.Value(), cw.Precipitation(Timespan10Min).Value())
			}
			if tc.p != nil && cw.Precipitation(Timespan10Min).Source() != tc.p.source {
				t.Errorf("CurrentWeatherByLocation failed, expected source: %s, but got: %s",
					tc.p.source, cw.Precipitation(Timespan10Min).Source())
			}
			if tc.p == nil {
				if cw.Precipitation(Timespan10Min).IsAvailable() {
					t.Errorf("CurrentWeatherByLocation failed, expected precipitation "+
						"to have no data, but got: %s", cw.Precipitation(Timespan10Min))
				}
				if !math.IsNaN(cw.Precipitation(Timespan10Min).Value()) {
					t.Errorf("CurrentWeatherByLocation failed, expected precipitation "+
						"to return NaN, but got: %s", cw.Precipitation(Timespan10Min).String())
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
			dateTime: time.Date(2023, 5, 23, 7, 0, 0, 0, time.Local),
			source:   SourceObservation,
			floatVal: 0,
		}},
		{"Berlin, Germany", &Precipitation{
			dateTime: time.Date(2023, 5, 23, 7, 0, 0, 0, time.Local),
			source:   SourceAnalysis,
			floatVal: 0.0092,
		}},
		{"Neermoor, Germany", nil},
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
				t.Errorf("CurrentWeatherByLocation failed: %s", err)
				return
			}
			if tc.p != nil && tc.p.String() != cw.Precipitation(Timespan1Hour).String() {
				t.Errorf("CurrentWeatherByLocation failed, expected precipitation "+
					"string: %s, got: %s", tc.p.String(), cw.Precipitation(Timespan1Hour))
			}
			if tc.p != nil && tc.p.Value() != cw.Precipitation(Timespan1Hour).Value() {
				t.Errorf("CurrentWeatherByLocation failed, expected precipitation "+
					"float: %f, got: %f", tc.p.Value(), cw.Precipitation(Timespan1Hour).Value())
			}
			if tc.p != nil && cw.Precipitation(Timespan1Hour).Source() != tc.p.source {
				t.Errorf("CurrentWeatherByLocation failed, expected source: %s, but got: %s",
					tc.p.source, cw.Precipitation(Timespan1Hour).Source())
			}
			if tc.p == nil {
				if cw.Precipitation(Timespan1Hour).IsAvailable() {
					t.Errorf("CurrentWeatherByLocation failed, expected precipitation "+
						"to have no data, but got: %s", cw.Precipitation(Timespan1Hour))
				}
				if !math.IsNaN(cw.Precipitation(Timespan1Hour).Value()) {
					t.Errorf("CurrentWeatherByLocation failed, expected precipitation "+
						"to return NaN, but got: %s", cw.Precipitation(Timespan1Hour).String())
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
		{"Neermoor, Germany", nil},
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
				t.Errorf("CurrentWeatherByLocation failed: %s", err)
				return
			}
			if tc.p != nil && tc.p.String() != cw.Precipitation(Timespan24Hours).String() {
				t.Errorf("CurrentWeatherByLocation failed, expected precipitation "+
					"string: %s, got: %s", tc.p.String(), cw.Precipitation(Timespan24Hours))
			}
			if tc.p != nil && tc.p.Value() != cw.Precipitation(Timespan24Hours).Value() {
				t.Errorf("CurrentWeatherByLocation failed, expected precipitation "+
					"float: %f, got: %f", tc.p.Value(), cw.Precipitation(Timespan24Hours).Value())
			}
			if tc.p != nil && cw.Precipitation(Timespan24Hours).Source() != tc.p.source {
				t.Errorf("CurrentWeatherByLocation failed, expected source: %s, but got: %s",
					tc.p.source, cw.Precipitation(Timespan24Hours).Source())
			}
			if tc.p == nil {
				if cw.Precipitation(Timespan24Hours).IsAvailable() {
					t.Errorf("CurrentWeatherByLocation failed, expected precipitation "+
						"to have no data, but got: %s", cw.Precipitation(Timespan24Hours))
				}
				if !math.IsNaN(cw.Precipitation(Timespan24Hours).Value()) {
					t.Errorf("CurrentWeatherByLocation failed, expected precipitation "+
						"to return NaN, but got: %s", cw.Precipitation(Timespan24Hours).String())
				}
			}
		})
	}
}

func TestClient_CurrentWeatherByLocation_PressureMSL(t *testing.T) {
	tt := []struct {
		// Location name
		loc string
		// CurWeather pressure
		p *Pressure
	}{
		{"Ehrenfeld, Germany", &Pressure{
			dateTime: time.Date(2023, 5, 23, 7, 0, 0, 0, time.Local),
			source:   SourceAnalysis,
			floatVal: 1018.9,
		}},
		{"Berlin, Germany", &Pressure{
			dateTime: time.Date(2023, 5, 23, 7, 0, 0, 0, time.Local),
			source:   SourceAnalysis,
			floatVal: 1011.5,
		}},
		{"Neermoor, Germany", nil},
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
				t.Errorf("CurrentWeatherByLocation failed: %s", err)
				return
			}
			if tc.p != nil && tc.p.String() != cw.PressureMSL().String() {
				t.Errorf("CurrentWeatherByLocation failed, expected pressure "+
					"string: %s, got: %s", tc.p.String(), cw.PressureMSL())
			}
			if tc.p != nil && tc.p.Value() != cw.PressureMSL().Value() {
				t.Errorf("CurrentWeatherByLocation failed, expected pressure "+
					"float: %f, got: %f", tc.p.Value(), cw.PressureMSL().Value())
			}
			if tc.p != nil && cw.PressureMSL().Source() != tc.p.source {
				t.Errorf("CurrentWeatherByLocation failed, expected source: %s, but got: %s",
					tc.p.source, cw.PressureMSL().Source())
			}
			if tc.p == nil {
				if cw.PressureMSL().IsAvailable() {
					t.Errorf("CurrentWeatherByLocation failed, expected pressure "+
						"to have no data, but got: %s", cw.PressureMSL())
				}
				if !math.IsNaN(cw.PressureMSL().Value()) {
					t.Errorf("CurrentWeatherByLocation failed, expected pressure "+
						"to return NaN, but got: %s", cw.PressureMSL().String())
				}
			}
		})
	}
}

func TestClient_CurrentWeatherByLocation_PressureQFE(t *testing.T) {
	tt := []struct {
		// Location name
		loc string
		// CurWeather pressure
		p *Pressure
	}{
		{"Ehrenfeld, Germany", &Pressure{
			dateTime: time.Date(2023, 5, 23, 7, 0, 0, 0, time.Local),
			source:   SourceAnalysis,
			floatVal: 1011.7,
		}},
		{"Berlin, Germany", nil},
		{"Neermoor, Germany", nil},
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
				t.Errorf("CurrentWeatherByLocation failed: %s", err)
				return
			}
			if tc.p != nil && tc.p.String() != cw.PressureQFE().String() {
				t.Errorf("CurrentWeatherByLocation failed, expected pressure "+
					"string: %s, got: %s", tc.p.String(), cw.PressureQFE())
			}
			if tc.p != nil && tc.p.Value() != cw.PressureQFE().Value() {
				t.Errorf("CurrentWeatherByLocation failed, expected pressure "+
					"float: %f, got: %f", tc.p.Value(), cw.PressureQFE().Value())
			}
			if tc.p != nil && cw.PressureQFE().Source() != tc.p.source {
				t.Errorf("CurrentWeatherByLocation failed, expected source: %s, but got: %s",
					tc.p.source, cw.PressureQFE().Source())
			}
			if tc.p == nil {
				if cw.PressureQFE().IsAvailable() {
					t.Errorf("CurrentWeatherByLocation failed, expected pressure "+
						"to have no data, but got: %s", cw.PressureQFE())
				}
				if !math.IsNaN(cw.PressureQFE().Value()) {
					t.Errorf("CurrentWeatherByLocation failed, expected pressure "+
						"to return NaN, but got: %s", cw.PressureQFE().String())
				}
			}
		})
	}
}

func TestClient_CurrentWeatherByLocation_SnowAmount(t *testing.T) {
	tt := []struct {
		// Location name
		loc string
		// CurWeather pressure
		d *Density
	}{
		{"Ehrenfeld, Germany", &Density{
			dateTime: time.Date(2023, 5, 23, 6, 0, 0, 0, time.UTC),
			source:   SourceAnalysis,
			floatVal: 0,
		}},
		{"Berlin, Germany", &Density{
			dateTime: time.Date(2023, 5, 23, 8, 0, 0, 0, time.UTC),
			source:   SourceAnalysis,
			floatVal: 21.1,
		}},
		{"Neermoor, Germany", nil},
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
				t.Errorf("CurrentWeatherByLocation failed: %s", err)
				return
			}
			if tc.d != nil && tc.d.String() != cw.SnowAmount().String() {
				t.Errorf("CurrentWeatherByLocation failed, expected snow amount "+
					"string: %s, got: %s", tc.d.String(), cw.SnowAmount())
			}
			if tc.d != nil && tc.d.Value() != cw.SnowAmount().Value() {
				t.Errorf("CurrentWeatherByLocation failed, expected snow amount "+
					"float: %f, got: %f", tc.d.Value(), cw.SnowAmount().Value())
			}
			if tc.d != nil && cw.SnowAmount().Source() != tc.d.source {
				t.Errorf("CurrentWeatherByLocation failed, expected source: %s, but got: %s",
					tc.d.source, cw.SnowAmount().Source())
			}
			if tc.d != nil && tc.d.dateTime.Unix() != cw.SnowAmount().DateTime().Unix() {
				t.Errorf("CurrentWeatherByLocation failed, expected datetime: %s, got: %s",
					tc.d.dateTime.Format(time.RFC3339), cw.SnowAmount().DateTime().Format(time.RFC3339))
			}
			if tc.d == nil {
				if cw.SnowAmount().IsAvailable() {
					t.Errorf("CurrentWeatherByLocation failed, expected snow amount "+
						"to have no data, but got: %s", cw.SnowAmount())
				}
				if !math.IsNaN(cw.SnowAmount().Value()) {
					t.Errorf("CurrentWeatherByLocation failed, expected snow amount "+
						"to return NaN, but got: %s", cw.SnowAmount().String())
				}
			}
		})
	}
}

func TestClient_CurrentWeatherByLocation_SnowHeight(t *testing.T) {
	tt := []struct {
		// Location name
		loc string
		// CurWeather height
		h *Height
	}{
		{"Ehrenfeld, Germany", &Height{
			dateTime: time.Date(2023, 5, 23, 6, 0, 0, 0, time.UTC),
			source:   SourceAnalysis,
			floatVal: 1.23,
		}},
		{"Berlin, Germany", &Height{
			dateTime: time.Date(2023, 5, 23, 6, 0, 0, 0, time.UTC),
			source:   SourceAnalysis,
			floatVal: 0.003,
		}},
		{"Neermoor, Germany", nil},
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
				t.Errorf("CurrentWeatherByLocation failed: %s", err)
				return
			}
			if tc.h != nil && tc.h.String() != cw.SnowHeight().String() {
				t.Errorf("CurrentWeatherByLocation failed, expected snow height "+
					"string: %s, got: %s", tc.h.String(), cw.SnowHeight())
			}
			if tc.h != nil && tc.h.Value() != cw.SnowHeight().Value() {
				t.Errorf("CurrentWeatherByLocation failed, expected snow height "+
					"float: %f, got: %f", tc.h.Value(), cw.SnowHeight().Value())
			}
			if tc.h != nil && tc.h.MeterString() != cw.SnowHeight().MeterString() {
				t.Errorf("CurrentWeatherByLocation failed, expected snow height "+
					"string: %s, got: %s", tc.h.MeterString(), cw.SnowHeight().MeterString())
			}
			if tc.h != nil && tc.h.Meter() != cw.SnowHeight().Meter() {
				t.Errorf("CurrentWeatherByLocation failed, expected snow height "+
					"float: %f, got: %f", tc.h.Meter(), cw.SnowHeight().Meter())
			}
			if tc.h != nil && tc.h.CentiMeterString() != cw.SnowHeight().CentiMeterString() {
				t.Errorf("CurrentWeatherByLocation failed, expected snow height "+
					"string: %s, got: %s", tc.h.CentiMeterString(), cw.SnowHeight().CentiMeterString())
			}
			if tc.h != nil && tc.h.CentiMeter() != cw.SnowHeight().CentiMeter() {
				t.Errorf("CurrentWeatherByLocation failed, expected snow height "+
					"float: %f, got: %f", tc.h.CentiMeter(), cw.SnowHeight().CentiMeter())
			}
			if tc.h != nil && tc.h.MilliMeterString() != cw.SnowHeight().MilliMeterString() {
				t.Errorf("CurrentWeatherByLocation failed, expected snow height "+
					"string: %s, got: %s", tc.h.MilliMeterString(), cw.SnowHeight().MilliMeterString())
			}
			if tc.h != nil && tc.h.MilliMeter() != cw.SnowHeight().MilliMeter() {
				t.Errorf("CurrentWeatherByLocation failed, expected snow height "+
					"float: %f, got: %f", tc.h.MilliMeter(), cw.SnowHeight().MilliMeter())
			}
			if tc.h != nil && cw.SnowHeight().Source() != tc.h.source {
				t.Errorf("CurrentWeatherByLocation failed, expected source: %s, but got: %s",
					tc.h.source, cw.SnowHeight().Source())
			}
			if tc.h != nil && tc.h.dateTime.Unix() != cw.SnowHeight().DateTime().Unix() {
				t.Errorf("CurrentWeatherByLocation failed, expected datetime: %s, got: %s",
					tc.h.dateTime.Format(time.RFC3339), cw.SnowHeight().DateTime().Format(time.RFC3339))
			}
			if tc.h == nil {
				if cw.SnowHeight().IsAvailable() {
					t.Errorf("CurrentWeatherByLocation failed, expected snow height "+
						"to have no data, but got: %s", cw.SnowHeight())
				}
				if !math.IsNaN(cw.SnowHeight().Value()) {
					t.Errorf("CurrentWeatherByLocation failed, expected snow height "+
						"to return NaN, but got: %s", cw.SnowHeight().String())
				}
				if !math.IsNaN(cw.SnowHeight().Meter()) {
					t.Errorf("CurrentWeatherByLocation failed, expected snow height "+
						"to return NaN, but got: %f", cw.SnowHeight().Meter())
				}
				if !math.IsNaN(cw.SnowHeight().CentiMeter()) {
					t.Errorf("CurrentWeatherByLocation failed, expected snow height "+
						"to return NaN, but got: %f", cw.SnowHeight().CentiMeter())
				}
				if !math.IsNaN(cw.SnowHeight().MilliMeter()) {
					t.Errorf("CurrentWeatherByLocation failed, expected snow height "+
						"to return NaN, but got: %f", cw.SnowHeight().MilliMeter())
				}
			}
		})
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
			dateTime: time.Date(2023, 5, 23, 7, 0, 0, 0, time.Local),
			source:   SourceObservation,
			floatVal: 14.6,
		}},
		{"Berlin, Germany", &Temperature{
			dateTime: time.Date(2023, 5, 23, 7, 0, 0, 0, time.Local),
			source:   SourceAnalysis,
			floatVal: 17.8,
		}},
		{"Neermoor, Germany", nil},
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
				t.Errorf("CurrentWeatherByLocation failed: %s", err)
				return
			}
			if tc.t != nil && tc.t.String() != cw.Temperature().String() {
				t.Errorf("CurrentWeatherByLocation failed, expected temperature "+
					"string: %s, got: %s", tc.t.String(), cw.Temperature())
			}
			if tc.t != nil && tc.t.Value() != cw.Temperature().Value() {
				t.Errorf("CurrentWeatherByLocation failed, expected temperature "+
					"float: %f, got: %f", tc.t.Value(), cw.Temperature().Value())
			}
			if tc.t != nil && cw.Temperature().Source() != tc.t.source {
				t.Errorf("CurrentWeatherByLocation failed, expected source: %s, but got: %s",
					tc.t.source, cw.Temperature().Source())
			}
			if tc.t == nil {
				if cw.Temperature().IsAvailable() {
					t.Errorf("CurrentWeatherByLocation failed, expected temperature "+
						"to have no data, but got: %s", cw.Temperature())
				}
				if !math.IsNaN(cw.Temperature().Value()) {
					t.Errorf("CurrentWeatherByLocation failed, expected temperature "+
						"to return NaN, but got: %s", cw.Temperature().String())
				}
			}
		})
	}
}

func TestClient_CurrentWeatherByLocation_WeatherSymbol(t *testing.T) {
	tt := []struct {
		// Location name
		loc string
		// CurWeather generic string
		gs *Condition
	}{
		{"Ehrenfeld, Germany", &Condition{
			dateTime:  time.Date(2023, 5, 23, 7, 30, 0, 0, time.UTC),
			source:    SourceAnalysis,
			stringVal: "overcast",
		}},
		{"Berlin, Germany", &Condition{
			dateTime:  time.Date(2023, 5, 23, 8, 50, 0, 0, time.UTC),
			source:    SourceAnalysis,
			stringVal: "cloudy",
		}},
		{"Neermoor, Germany", nil},
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
				t.Errorf("CurrentWeatherByLocation failed: %s", err)
				return
			}
			if tc.gs != nil && tc.gs.String() != cw.WeatherSymbol().String() {
				t.Errorf("CurrentWeatherByLocation failed, expected weathersymbol "+
					"string: %s, got: %s", tc.gs.String(), cw.WeatherSymbol())
			}
			if tc.gs != nil && tc.gs.Condition() != cw.WeatherSymbol().Condition() {
				t.Errorf("CurrentWeatherByLocation failed, expected condition "+
					"string: %s, got: %s", tc.gs.Condition(), cw.WeatherSymbol().Condition())
			}
			if tc.gs != nil && tc.gs.Value() != cw.WeatherSymbol().Value() {
				t.Errorf("CurrentWeatherByLocation failed, expected weathersymbol "+
					"string: %s, got: %s", tc.gs.Value(), cw.WeatherSymbol().Value())
			}
			if tc.gs != nil && cw.WeatherSymbol().Source() != tc.gs.source {
				t.Errorf("CurrentWeatherByLocation failed, expected source: %s, but got: %s",
					tc.gs.source, cw.WeatherSymbol().Source())
			}
			if tc.gs != nil && tc.gs.dateTime.Unix() != cw.WeatherSymbol().DateTime().Unix() {
				t.Errorf("CurrentWeatherByLocation failed, expected datetime: %s, got: %s",
					tc.gs.dateTime.Format(time.RFC3339), cw.WeatherSymbol().DateTime().Format(time.RFC3339))
			}
			if tc.gs == nil {
				if cw.WeatherSymbol().IsAvailable() {
					t.Errorf("CurrentWeatherByLocation failed, expected weathersymbol "+
						"to have no data, but got: %s", cw.WeatherSymbol())
				}
				if cw.WeatherSymbol().Value() != DataUnavailable {
					t.Errorf("CurrentWeatherByLocation failed, expected weathersymbol "+
						"to return DataUnavailable, but got: %s", cw.WeatherSymbol().String())
				}
			}
		})
	}
}

func TestClient_CurrentWeatherByLocation_WindDirection(t *testing.T) {
	tt := []struct {
		// Location name
		loc string
		// CurWeather direction
		d *Direction
		// Direction abbr. string
		da string
		// Direction full string
		df string
	}{
		{"Ehrenfeld, Germany", &Direction{
			dateTime: time.Date(2023, 5, 23, 7, 0, 0, 0, time.Local),
			source:   SourceAnalysis,
			floatVal: 302,
		}, "NWbW", "Northwest by West"},
		{"Berlin, Germany", &Direction{
			dateTime: time.Date(2023, 5, 23, 7, 0, 0, 0, time.Local),
			source:   SourceAnalysis,
			floatVal: 286,
		}, "WbN", "West by North"},
		{"Neermoor, Germany", nil, "", ""},
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
				t.Errorf("CurrentWeatherByLocation failed: %s", err)
				return
			}
			if tc.d != nil && tc.d.String() != cw.WindDirection().String() {
				t.Errorf("CurrentWeatherByLocation failed, expected wind direction "+
					"string: %s, got: %s", tc.d.String(), cw.WindDirection())
			}
			if tc.d != nil && tc.d.Value() != cw.WindDirection().Value() {
				t.Errorf("CurrentWeatherByLocation failed, expected wind direction "+
					"float: %f, got: %f", tc.d.Value(), cw.WindDirection().Value())
			}
			if tc.d != nil && cw.WindDirection().Source() != tc.d.source {
				t.Errorf("CurrentWeatherByLocation failed, expected source: %s, but got: %s",
					tc.d.source, cw.WindDirection().Source())
			}
			if tc.d != nil && cw.WindDirection().Direction() != tc.da {
				t.Errorf("CurrentWeatherByLocation failed, expected direction abbr.: %s, but got: %s",
					tc.da, cw.WindDirection().Direction())
			}
			if tc.d != nil && cw.WindDirection().DirectionFull() != tc.df {
				t.Errorf("CurrentWeatherByLocation failed, expected direction full: %s, but got: %s",
					tc.df, cw.WindDirection().DirectionFull())
			}
			if tc.d == nil {
				if cw.WindDirection().IsAvailable() {
					t.Errorf("CurrentWeatherByLocation failed, expected wind direction "+
						"to have no data, but got: %s", cw.WindDirection())
				}
				if !math.IsNaN(cw.WindDirection().Value()) {
					t.Errorf("CurrentWeatherByLocation failed, expected wind direction "+
						"to return NaN, but got: %s", cw.WindSpeed().String())
				}
			}
		})
	}
}

func TestClient_CurrentWeatherByLocation_WindGust(t *testing.T) {
	tt := []struct {
		// Location name
		loc string
		// CurWeather speed
		s *Speed
	}{
		{"Ehrenfeld, Germany", &Speed{
			dateTime: time.Date(2023, 5, 23, 7, 0, 0, 0, time.Local),
			source:   SourceAnalysis,
			floatVal: 7.770000,
		}},
		{"Berlin, Germany", &Speed{
			dateTime: time.Date(2023, 5, 23, 7, 0, 0, 0, time.Local),
			source:   SourceAnalysis,
			floatVal: 5.570000,
		}},
		{"Neermoor, Germany", nil},
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
				t.Errorf("CurrentWeatherByLocation failed: %s", err)
				return
			}
			if tc.s != nil && tc.s.String() != cw.WindGust().String() {
				t.Errorf("CurrentWeatherByLocation failed, expected wind gust "+
					"string: %s, got: %s", tc.s.String(), cw.WindGust())
			}
			if tc.s != nil && tc.s.Value() != cw.WindGust().Value() {
				t.Errorf("CurrentWeatherByLocation failed, expected wind gust "+
					"float: %f, got: %f", tc.s.Value(), cw.WindGust().Value())
			}
			if tc.s != nil && cw.WindGust().Source() != tc.s.source {
				t.Errorf("CurrentWeatherByLocation failed, expected source: %s, but got: %s",
					tc.s.source, cw.WindGust().Source())
			}
			if tc.s == nil {
				if cw.WindGust().IsAvailable() {
					t.Errorf("CurrentWeatherByLocation failed, expected wind gust "+
						"to have no data, but got: %s", cw.WindGust())
				}
				if !math.IsNaN(cw.WindGust().Value()) {
					t.Errorf("CurrentWeatherByLocation failed, expected wind gust "+
						"to return NaN, but got: %s", cw.WindGust().String())
				}
			}
		})
	}
}

func TestClient_CurrentWeatherByLocation_WindSpeed(t *testing.T) {
	tt := []struct {
		// Location name
		loc string
		// CurWeather speed
		s *Speed
	}{
		{"Ehrenfeld, Germany", &Speed{
			dateTime: time.Date(2023, 5, 23, 7, 0, 0, 0, time.Local),
			source:   SourceAnalysis,
			floatVal: 3.94,
		}},
		{"Berlin, Germany", &Speed{
			dateTime: time.Date(2023, 5, 23, 7, 0, 0, 0, time.Local),
			source:   SourceAnalysis,
			floatVal: 3.19,
		}},
		{"Neermoor, Germany", nil},
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
				t.Errorf("CurrentWeatherByLocation failed: %s", err)
				return
			}
			if tc.s != nil && tc.s.String() != cw.WindSpeed().String() {
				t.Errorf("CurrentWeatherByLocation failed, expected wind speed "+
					"string: %s, got: %s", tc.s.String(), cw.WindSpeed())
			}
			if tc.s != nil && tc.s.Value() != cw.WindSpeed().Value() {
				t.Errorf("CurrentWeatherByLocation failed, expected wind speed "+
					"float: %f, got: %f", tc.s.Value(), cw.WindSpeed().Value())
			}
			if tc.s != nil && cw.WindSpeed().Source() != tc.s.source {
				t.Errorf("CurrentWeatherByLocation failed, expected source: %s, but got: %s",
					tc.s.source, cw.WindSpeed().Source())
			}
			if tc.s == nil {
				if cw.WindSpeed().IsAvailable() {
					t.Errorf("CurrentWeatherByLocation failed, expected wind speed "+
						"to have no data, but got: %s", cw.WindSpeed())
				}
				if !math.IsNaN(cw.WindSpeed().Value()) {
					t.Errorf("CurrentWeatherByLocation failed, expected wind speed "+
						"to return NaN, but got: %s", cw.WindSpeed().String())
				}
			}
		})
	}
}
