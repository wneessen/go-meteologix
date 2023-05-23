// SPDX-FileCopyrightText: 2023 Winni Neessen <wn@neessen.dev>
//
// SPDX-License-Identifier: MIT

package meteologix

import (
	"testing"
)

func TestSource_String(t *testing.T) {
	tt := []struct {
		// Original source
		os Source
		// Expected string
		es string
	}{
		{SourceObservation, "Observation"},
		{SourceAnalysis, "Analysis"},
		{SourceForecast, "Forecast"},
		{SourceMixed, "Mixed"},
		{SourceUnknown, "Unknown"},
		{999, "Unknown"},
	}
	for _, tc := range tt {
		t.Run(tc.os.String(), func(t *testing.T) {
			if tc.os.String() != tc.es {
				t.Errorf("String for Source failed, expected: %s, got: %s",
					tc.es, tc.os.String())
			}
		})
	}
}

func TestStringToSource(t *testing.T) {
	tt := []struct {
		// Original string
		os string
		// Expected source
		es Source
	}{
		{"Observation", SourceObservation},
		{"Analysis", SourceAnalysis},
		{"Forecast", SourceForecast},
		{"Mixed", SourceMixed},
		{"Unknown", SourceUnknown},
	}
	for _, tc := range tt {
		t.Run(tc.es.String(), func(t *testing.T) {
			if r := StringToSource(tc.os); r != tc.es {
				t.Errorf("StringToSource failed, expected: %s, got: %s",
					tc.es.String(), r.String())
			}
		})
	}
}
