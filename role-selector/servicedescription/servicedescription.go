/*
 *  Copyright 2002-2025 Barcelona Supercomputing Center (www.bsc.es)
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

package servicedescription

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"colmena.bsc.es/role-selector/types"
)

func Parse(serviceDescription *types.ServiceDescription, hardware string) ([]types.KPI, []types.DockerRoleDefinition, error) {
	var kpis []types.KPI
	var roles []types.DockerRoleDefinition

	re := regexp.MustCompile(`[<>]\s*(\d+(\.\d+)?)`)

	for _, kpiStr := range serviceDescription.Kpis {
		kpi, err := createKpi(kpiStr.Query, "", re)
		if err != nil {
			return nil, nil, err
		}
		kpis = append(kpis, kpi)
	}

	for _, role := range serviceDescription.DockerRoleDefinitions {
		if role.HardwareRequirements[0] != hardware {
			log.Printf("RoleId: %v does not match hardware: %v, required: %v", role.Id, hardware, role.HardwareRequirements)
			continue
		}

		log.Printf("RoleId: %v matches hardware %v", role.Id, hardware)
		roles = append(roles, role)
		for _, kpiStr := range role.Kpis {
			kpi, err := createKpi(kpiStr.Query, role.Id, re)
			if err != nil {
				return nil, nil, err
			}
			kpis = append(kpis, kpi)
		}
	}

	log.Printf("Parsed kpis: %v", kpis)
	log.Printf("Parsed roles: %v", roles)
	return kpis, roles, nil
}

func createKpi(kpiStr string, associatedRole string, re *regexp.Regexp) (types.KPI, error) {
	kpi := types.KPI{
		Query:          kpiStr,
		Value:          0,
		Threshold:      0,
		Operator:       "",
		AssociatedRole: associatedRole,
		Level:          "",
	}
	operator, threshold, err := FetchOperatorAndThreshold(kpiStr)
	kpi.Operator = operator
	kpi.Threshold = threshold

	if err != nil {
		fmt.Printf("Could not parse KPI with Operator and Threshold")
		return kpi, err
	}

	return kpi, nil
}

func FetchOperatorAndThreshold(kpi string) (types.Operator, float64, error) {
	var op types.Operator

	// Walk the string once, looking for '<' or '>'
	for i, r := range kpi {
		if r == '<' || r == '>' {
			op = types.Operator(string(r))

			// Everything after the operator (trim leading spaces) should be the threshold
			after := strings.TrimSpace(kpi[i+1:])

			// Grab the substring up to the first space (if any)
			end := strings.IndexByte(after, ' ')
			if end == -1 {
				end = len(after)
			}
			numStr := after[:end]

			val, err := strconv.ParseFloat(numStr, 64)
			if err != nil {
				return "", 0, err
			}
			return op, val, nil
		}
	}
	return "", 0, fmt.Errorf("no comparison operator found in KPI")
}
