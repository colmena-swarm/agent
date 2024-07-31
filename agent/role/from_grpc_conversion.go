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
	"regexp"
	"strconv"

	"colmena.bsc.es/agent/grpc"
)

func parseGrpcRole(grpcRole *grpc.DockerRoleDefinition) Role {
	role := Role{
		Id:       				grpcRole.Id,
		ImageId:  				grpcRole.ImageId,
		HardwareRequirements: 	parseHardwareRequirements(grpcRole.HardwareRequirements),
		Kpis:     				parseKpis(grpcRole.Kpis),
	}
	return role
}

func parseHardwareRequirements(hardwareRequirements []grpc.HardwareRequirement) []string {
	features := []string{}
	for _, hardwareRequirement := range hardwareRequirements {
		features = append(features, hardwareRequirement.String())
	}
	return features
}

func parseKpis(kpis []*grpc.Kpi) []KPI {
	parsed := []KPI{}
	for _, kpi := range kpis {
		kpiString := kpi.Value
		parsed = append(parsed, parseKpi(kpiString))
	}
	return parsed
}

func parseKpi(kpi string) KPI {
	re := regexp.MustCompile(`(.*)\[(.*)\]\s*(>=|<=|>|<)\s*(.*)\s*`)
    match := re.FindStringSubmatch(kpi)
    key := match[1]
	comparison := match[3]
	thresholdValue, _ := strconv.ParseFloat(match[4], 64)

	matchedTime := match[2]
	timere := regexp.MustCompile(`(\d*)(\D*)`)
	timeMatch := timere.FindStringSubmatch(matchedTime)
	timeValue, _ := strconv.ParseFloat(timeMatch[1], 32)
	timeUnit := timeMatch[2]
	
	return KPI {
		Key: key,
		Threshold: thresholdValue,
		FromUnit: timeValue,
		FromType: timeUnit,
		Comparison: comparison,
	}
}
