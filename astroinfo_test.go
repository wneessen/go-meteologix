// SPDX-FileCopyrightText: 2023 Winni Neessen <wn@neessen.dev>
//
// SPDX-License-Identifier: MIT
package meteologix

import (
	"testing"
	"time"
)

func TestClient_AstronomicalInfoByCoordinates(t *testing.T) {
	la := 52.5067296
	lo := 13.2599306
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		t.Errorf("failed to load time location data for Europe/Berlin: %s", err)
		return
	}
	rt := time.Date(2023, 5, 28, 15, 8, 33, 0, loc)
	nfmt := time.Date(2023, 6, 4, 5, 43, 56, 0, loc)
	nnmt := time.Date(2023, 6, 18, 6, 39, 10, 0, loc)
	c := New(withMockAPI())
	ai, err := c.AstronomicalInfoByCoordinates(la, lo)
	if err != nil {
		t.Errorf("failed to fetch astronomical information: %s", err)
		return
	}
	if ai.Latitude != la {
		t.Errorf("AstronomicalInfoByCoordinates failed, expected lat: %f, got: %f", la,
			ai.Latitude)
	}
	if ai.Longitude != lo {
		t.Errorf("AstronomicalInfoByCoordinates failed, expected lon: %f, got: %f", lo,
			ai.Longitude)
	}
	if ai.Run.UnixMilli() != rt.UnixMilli() {
		t.Errorf("AstronomicalInfoByCoordinates failed, expected run time: %s, got: %s",
			rt.String(), ai.Run.String())
	}
	if ai.TimeZone != loc.String() {
		t.Errorf("AstronomicalInfoByCoordinates failed, expected time zone: %s, got: %s",
			loc.String(), ai.TimeZone)
	}
	if ai.NextFullMoon.UnixMilli() != nfmt.UnixMilli() {
		t.Errorf("AstronomicalInfoByCoordinates failed, expected next full moon: %s, got: %s",
			nfmt.String(), ai.NextFullMoon.String())
	}
	if ai.NextNewMoon.UnixMilli() != nnmt.UnixMilli() {
		t.Errorf("AstronomicalInfoByCoordinates failed, expected next new moon: %s, got: %s",
			nnmt.String(), ai.NextNewMoon.String())
	}
}

func TestClient_AstronomicalInfoByLocation(t *testing.T) {
	la := 52.5067296
	lo := 13.2599306
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		t.Errorf("failed to load time location data for Europe/Berlin: %s", err)
		return
	}
	rt := time.Date(2023, 5, 28, 15, 8, 33, 0, loc)
	nfmt := time.Date(2023, 6, 4, 5, 43, 56, 0, loc)
	nnmt := time.Date(2023, 6, 18, 6, 39, 10, 0, loc)
	c := New(withMockAPI())
	ai, err := c.AstronomicalInfoByLocation("Berlin, Germany")
	if err != nil {
		t.Errorf("failed to fetch astronomical information: %s", err)
		return
	}
	if ai.Latitude != la {
		t.Errorf("AstronomicalInfoByCoordinates failed, expected lat: %f, got: %f", la,
			ai.Latitude)
	}
	if ai.Longitude != lo {
		t.Errorf("AstronomicalInfoByCoordinates failed, expected lon: %f, got: %f", lo,
			ai.Longitude)
	}
	if ai.Run.UnixMilli() != rt.UnixMilli() {
		t.Errorf("AstronomicalInfoByCoordinates failed, expected run time: %s, got: %s",
			rt.String(), ai.Run.String())
	}
	if ai.TimeZone != loc.String() {
		t.Errorf("AstronomicalInfoByCoordinates failed, expected time zone: %s, got: %s",
			loc.String(), ai.TimeZone)
	}
	if ai.NextFullMoon.UnixMilli() != nfmt.UnixMilli() {
		t.Errorf("AstronomicalInfoByCoordinates failed, expected next full moon: %s, got: %s",
			nfmt.String(), ai.NextFullMoon.String())
	}
	if ai.NextNewMoon.UnixMilli() != nnmt.UnixMilli() {
		t.Errorf("AstronomicalInfoByCoordinates failed, expected next new moon: %s, got: %s",
			nnmt.String(), ai.NextNewMoon.String())
	}
}

func TestAstronomicalInfo_SunsetByDateString(t *testing.T) {
	la := 52.5067296
	lo := 13.2599306
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		t.Errorf("failed to load time location data for Europe/Berlin: %s", err)
		return
	}
	ti := time.Date(2023, 5, 28, 21, 16, 37, 0, loc)
	c := New(withMockAPI())
	ai, err := c.AstronomicalInfoByCoordinates(la, lo)
	if err != nil {
		t.Errorf("failed to fetch astronomical information: %s", err)
		return
	}
	if !ai.SunsetByTime(ti).IsAvailable() {
		t.Errorf("SunsetByTime failed, expected entry, but got 'not available'")
		return
	}
	if ai.SunsetByTime(ti).Value().UnixMilli() != ti.UnixMilli() {
		t.Errorf("SunsetByTime failed, expected sunset: %s, got: %s",
			ti.String(), ai.SunsetByTime(ti).Value().String())
	}
	if !ai.SunsetByDateString(ti.Format(DateFormat)).IsAvailable() {
		t.Errorf("SunsetByDateString failed, expected entry, but got 'not available'")
		return
	}
	if ai.SunsetByTime(ti).String() != ti.Format(time.RFC3339) {
		t.Errorf("SunsetByTime failed, expected sunset: %s, got: %s",
			ti.Format(time.RFC3339), ai.SunsetByTime(ti).String())
	}
	if ai.SunsetByDateString(ti.Format(DateFormat)).Value().UnixMilli() != ti.UnixMilli() {
		t.Errorf("SunsetByDateString failed, expected sunset: %s, got: %s",
			ti.String(), ai.SunsetByDateString(ti.Format(DateFormat)).Value().String())
	}
	ti = time.Time{}
	if ai.SunsetByTime(ti).IsAvailable() {
		t.Errorf("SunsetByTime failed, expected no entry, but got: %s",
			ai.SunsetByTime(ti).Value().String())
	}
	if !ai.SunsetByTime(ti).Value().IsZero() {
		t.Errorf("SunsetByTime failed, expected no entry, but got: %s",
			ai.SunsetByTime(ti).Value().String())
	}
}
