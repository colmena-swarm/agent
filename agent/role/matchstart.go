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
	"slices"

	"colmena.bsc.es/agent/device"
)

func matchstart(device device.Info, roles []Role, kpiMatcher KpiMatcher) []Role {
	matching := []Role{}
	for _, role := range roles {
		if (matchDeviceFeatures(device, role) && qosNotMet(device, role, kpiMatcher)) {
			matching = append(matching, role)
		}
	}
	return matching
}

func qosNotMet(device device.Info, role Role, kpiMatcher KpiMatcher) bool {
	if (device.Strategy == "EAGER") {
		return true;
	}
	met := kpisMet(role.Kpis, kpiMatcher)
	return !met
}

func matchDeviceFeatures(device device.Info, role Role) bool {
	for _, roleFeature := range role.HardwareRequirements {
		if !slices.Contains(device.Features, roleFeature) {
			return false
		}
	}
	return true;
}