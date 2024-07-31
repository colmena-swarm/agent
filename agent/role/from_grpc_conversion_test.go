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

	"colmena.bsc.es/agent/grpc"
	"github.com/stretchr/testify/assert"
)

func TestKpiParse(t *testing.T) {
	parsed := parseKpi("processing_time[2h] >= 60.0")
	assert.Equal(t, "processing_time", parsed.Key)
	assert.Equal(t, float64(2), parsed.FromUnit)
	assert.Equal(t, "h", parsed.FromType)
	assert.Equal(t, 60.0, parsed.Threshold)
	assert.Equal(t, ">=", parsed.Comparison)
}

func TestParseGrpcRole(t *testing.T) {
	testKpi := grpc.Kpi{Value: "x[1s] < 22.5"}
	testRole := grpc.DockerRoleDefinition{
		Id: "RoleId",
		ImageId: "ImageId",
		HardwareRequirements: []grpc.HardwareRequirement{grpc.HardwareRequirement_CAMERA},
		Kpis: []*grpc.Kpi{&testKpi},
	}
	converted := parseGrpcRole(&testRole)
	assert.Equal(t, "RoleId", converted.Id)
	assert.Equal(t, "ImageId", converted.ImageId)
	assert.Equal(t, []string{"CAMERA"}, converted.HardwareRequirements)
	convertedKpi := converted.Kpis[0]
	assert.Equal(t, "x", convertedKpi.Key)
	assert.Equal(t, 22.5, convertedKpi.Threshold)
	assert.Equal(t, "<", convertedKpi.Comparison)
}