// SPDX-FileCopyrightText: 2023 Winni Neessen <wn@neessen.dev>
//
// SPDX-License-Identifier: MIT

package meteologix

import (
	"testing"
)

func TestFindDirection(t *testing.T) {
	// Prepare test cases
	tt := []struct {
		v  float64
		dm map[float64]string
		er string
	}{
		{15, WindDirAbbrMap, "NbE"},
		{47, WindDirAbbrMap, "NE"},
		{200, WindDirAbbrMap, "SSW"},
		{330, WindDirAbbrMap, "NWbN"},
		{15, WindDirFullMap, "North by East"},
		{47, WindDirFullMap, "Northeast"},
		{200, WindDirFullMap, "South-Southwest"},
		{330, WindDirFullMap, "Northwest by North"},
	}

	// Run tests
	for _, tc := range tt {
		t.Run("", func(t *testing.T) {
			r := findDirection(tc.v, tc.dm)
			if tc.er != r {
				t.Errorf("findDirection failed, expected: %s, got: %s", tc.er, r)
			}
		})
	}
}
