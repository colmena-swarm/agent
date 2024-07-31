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

func TestKpiConversion(t *testing.T) {
	kpi := KPI{
		Key:			"key",
		Threshold:		51.55,
		FromUnit:		5,
		FromType:		"m",
		Comparison:		"<=",
	}
	converted := *convertKpi(kpi)
	assert.Equal(t, kpi.Key, converted.Key)
	assert.Equal(t, float32(kpi.Threshold), converted.Threshold)
	assert.Equal(t, float32(kpi.FromUnit), converted.FromUnit)
	assert.Equal(t, kpi.FromType, converted.FromType)
	assert.Equal(t, kpi.Comparison, converted.Comparison)
}