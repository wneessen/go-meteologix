// SPDX-FileCopyrightText: 2023 Winni Neessen <wn@neessen.dev>
//
// SPDX-License-Identifier: MIT

package meteologix

import (
	"testing"
	"time"
)

func TestWeatherForecast_All(t *testing.T) {
	tests := []struct {
		name         string
		lat          float64
		lon          float64
		timespan     Timespan
		fcastdetails ForecastDetails
		datapoints   int
	}{
		{"1h Standard", 50.9586327, 6.9685969, Timespan1Hour, ForecastDetailStandard, 24},
		{"3h Standard", 50.9586327, 6.9685969, Timespan3Hours, ForecastDetailStandard, 37},
		{"6h Standard", 50.9586327, 6.9685969, Timespan6Hours, ForecastDetailStandard, 38},
		{"1h Advanced", 50.9586327, 6.9685969, Timespan1Hour, ForecastDetailAdvanced, 24},
	}
	client := New(withMockAPI())
	if client == nil {
		t.Errorf("failed to create new Client, got nil")
		return
	}
	for _, testcase := range tests {
		t.Run(testcase.name, func(t *testing.T) {
			forecast, err := client.ForecastByCoordinates(testcase.lat, testcase.lon, testcase.timespan,
				testcase.fcastdetails)
			if err != nil {
				t.Errorf("ForecastByLocation failed: %s", err)
				return
			}
			data := forecast.All()
			if len(data) <= 0 {
				t.Errorf("ForecastByLocation failed, expected at least one forecast, got: %d", len(data))
			}
			if len(data) != testcase.datapoints {
				t.Errorf("ForecastByLocation failed, expected %d datapoints, got: %d", testcase.datapoints,
					len(data))
			}
		})
	}
}

func TestWeatherForecast_At(t *testing.T) {
	tests := []struct {
		name         string
		lat          float64
		lon          float64
		timespan     Timespan
		fcastdetails ForecastDetails
		timestamp    time.Time
		expectedTime time.Time
	}{
		{"1h Standard", 50.9586327, 6.9685969, Timespan1Hour, ForecastDetailStandard,
			time.Date(2024, 8, 13, 12, 22, 29, 0, time.UTC),
			time.Date(2024, 8, 13, 12, 0, 0, 0, time.UTC)},
		{"3h Standard", 50.9586327, 6.9685969, Timespan3Hours, ForecastDetailStandard,
			time.Date(2024, 8, 16, 8, 49, 12, 0, time.UTC),
			time.Date(2024, 8, 16, 9, 33, 03, 0, time.UTC)},
		{"6h Standard", 50.9586327, 6.9685969, Timespan6Hours, ForecastDetailStandard,
			time.Date(2024, 9, 2, 23, 3, 47, 0, time.UTC),
			time.Date(2024, 9, 2, 21, 54, 03, 0, time.UTC)},
		{"1h Advanced", 50.9586327, 6.9685969, Timespan1Hour, ForecastDetailAdvanced,
			time.Date(2024, 8, 30, 2, 9, 0, 0, time.UTC),
			time.Date(2024, 8, 30, 2, 0, 0, 0, time.UTC)},
	}
	client := New(withMockAPI())
	if client == nil {
		t.Errorf("failed to create new Client, got nil")
		return
	}
	for _, testcase := range tests {
		t.Run(testcase.name, func(t *testing.T) {
			forecast, err := client.ForecastByCoordinates(testcase.lat, testcase.lon, testcase.timespan,
				testcase.fcastdetails)
			if err != nil {
				t.Errorf("ForecastByLocation failed: %s", err)
				return
			}
			data := forecast.At(testcase.timestamp)
			if !data.DateTime().Equal(testcase.expectedTime) {
				t.Errorf("ForecastByCoordinates failed, expected %s, got: %s", testcase.expectedTime,
					data.DateTime())
			}
		})
	}
}
