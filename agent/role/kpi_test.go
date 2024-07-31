/*
 *  Copyright 2002-2024 Barcelona Supercomputing Center (www.bsc.es)
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */
package role

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var first = KPI{
	Key:        "first",
	Threshold:  51.55,
	FromUnit:   5,
	FromType:   "m",
	Comparison: "<=",
}
var second = KPI{
	Key:        "second",
	Threshold:  51.55,
	FromUnit:   5,
	FromType:   "m",
	Comparison: "<=",
}

func TestQoSResultsAreCorrectlyAggregated(t *testing.T) {
	var tests = []struct {
		a, b bool
		want bool
	}{
		{true, true, true},
		{true, false, false},
		{false, true, false},
		{false, false, false},
	}
	for _, tt := range tests {
		resultMapper := func(kpi KPI) bool {
			if kpi.Key == "first" {
				return tt.a
			}
			return tt.b
		}
		kpisMet := kpisMet([]KPI{first, second}, KpiMatcherFunc(resultMapper))
		assert.Equal(t, tt.want, kpisMet, "First: %v, Second: %v, Result: %v, Expected: %v", tt.a, tt.b, kpisMet, tt.want)
	}
}

func TestKpisMetWhenNoKpis(t *testing.T) {
	kpisMet := kpisMet([]KPI{}, returns(false))
	assert.True(t, kpisMet)
}

func returns(returning bool) KpiMatcher {
	return KpiMatcherFunc(func(kpi KPI) bool { return returning })
}
