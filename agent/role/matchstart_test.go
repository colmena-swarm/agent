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

	"colmena.bsc.es/agent/device"
	"github.com/stretchr/testify/assert"
)

func TestDeviceFeaturesMatch(t *testing.T) {
	testrole := Role{"id", "imageId", []string{"CAMERA"}, []KPI{}}
	var tests = []struct {
		deviceFeatures []string
		expected       bool
	}{
		{[]string{"CAMERA"}, true},
		{[]string{"CAMERA", "CPU"}, true},
		{[]string{"CPU"}, false},
	}
	for _, tt := range tests {
		device := device.Info{Strategy: "eager", Features: tt.deviceFeatures}
		match := matchDeviceFeatures(device, testrole)
		assert.Equal(t, tt.expected, match, "with deviceFeatures %v, was %v", device.Features, tt.expected)
	}
}