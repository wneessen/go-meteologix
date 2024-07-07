// SPDX-FileCopyrightText: 2023 Winni Neessen <wn@neessen.dev>
//
// SPDX-License-Identifier: MIT

package meteologix

import (
	"testing"
	"time"
)

func TestCondition_Condition(t *testing.T) {
	tc := Condition{
		dateTime:  time.Date(2023, 5, 23, 8, 50, 0, 0, time.UTC),
		source:    SourceAnalysis,
		stringVal: "cloudy",
	}
	if tc.Condition() != CondCloudy {
		t.Errorf("Condition failed, expected: %s, got: %s", CondCloudy.String(),
			tc.Condition().String())
	}
	tc = Condition{
		dateTime:  time.Date(2023, 5, 23, 8, 50, 0, 0, time.UTC),
		source:    SourceAnalysis,
		stringVal: "non-existing",
	}
	if tc.Condition() != CondUnknown {
		t.Errorf("Condition failed, expected: %s, got: %s", CondUnknown.String(),
			tc.Condition().String())
	}
	tc = Condition{notAvailable: true}
	if tc.Condition() != CondUnknown {
		t.Errorf("Condition failed, expected: %s, got: %s", CondUnknown.String(),
			tc.Condition().String())
	}
	ct := ConditionType("foo")
	if ct.String() != CondUnknown.String() {
		t.Errorf("Condition.String for unknown type failed, expected: %s, got: %s",
			CondUnknown.String(), ct.String())
	}
}
