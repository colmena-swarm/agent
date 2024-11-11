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

package service_description

import (
	"testing"

	"colmena.bsc.es/role-selector/fileloader"
	"colmena.bsc.es/role-selector/servicedescription"
	"colmena.bsc.es/role-selector/types"
)

func TestParseServiceDefinition(t *testing.T) {
	data, err := fileloader.LoadFromFile[types.ServiceDescription]("../resources/service_definition.json")
	if err != nil {
		t.Error(err)
		return
	}

	indicators, roleIds, err := servicedescription.Parse(data, "CPU")
	if err != nil || len(roleIds) != 1 {
		t.Fatalf("Error occurred: %v, Expected 1 role in roleIds, but got %d", err, len(roleIds))
	}

	expectedRoles := []string{"Processing"}
	for _, expectedRole := range expectedRoles {
		found := false
		for _, role := range roleIds {
			if role.Id == expectedRole {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected role '%s' not found in roleIds", expectedRole)
		}
	}

	if len(indicators) != 1 {
		t.Errorf("Expected 1 indicator, but got %d", len(indicators))
	}

	if indicators[0].Query != "avg_over_time(examplecontextdata_processing_time[5s]) < 15" {
		t.Errorf("Expected first indicator Query to be 'avg_over_time(examplecontextdata_processing_time[5s]) < 15', but got '%s'", indicators[0].Query)
	}
	if indicators[0].Threshold != 15.0 {
		t.Errorf("Expected first indicator Threshold to be 15, but got %f", indicators[0].Threshold)
	}
	if indicators[0].Operator != types.LessThan {
		t.Errorf("Expected first indicator Operator to be <, but got %s", indicators[0].Operator)
	}
	if indicators[0].AssociatedRole != "Processing" {
		t.Errorf("Expected first indicator AssociatedRole to be 'Processing', but got '%s'", indicators[0].AssociatedRole)
	}
}

func TestParseServiceDefinition_MultiVariableQuery(t *testing.T) {
	data, err := fileloader.LoadFromFile[types.ServiceDescription]("../resources/service_definition_2.json")
	if err != nil {
		t.Error(err)
		return
	}

	indicators, roleIds, err := servicedescription.Parse(data, "CPU")
	if err != nil || len(roleIds) != 1 {
		t.Fatalf("Error occurred: %v, Expected 1 role in roleIds, but got %d", err, len(roleIds))
	}

	expectedRoles := []string{"Processing"}
	for _, expectedRole := range expectedRoles {
		found := false
		for _, role := range roleIds {
			if role.Id == expectedRole {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected role '%s' not found in roleIds", expectedRole)
		}
	}

	if len(indicators) != 1 {
		t.Errorf("Expected 1 indicator, but got %d", len(indicators))
	}

	if indicators[0].Query != "avg_over_time(examplecontextdata_processing_time[5s]) + avg_over_time(examplecontextdata_processing_time[5s]) < 15" {
		t.Errorf("Expected first indicator Query to be 'avg_over_time(examplecontextdata_processing_time[5s]) < 15', but got '%s'", indicators[0].Query)
	}
	if indicators[0].Threshold != 15.0 {
		t.Errorf("Expected first indicator Threshold to be 15, but got %f", indicators[0].Threshold)
	}
	if indicators[0].Operator != types.LessThan {
		t.Errorf("Expected first indicator Operator to be <, but got %s", indicators[0].Operator)
	}
	if indicators[0].AssociatedRole != "Processing" {
		t.Errorf("Expected first indicator AssociatedRole to be 'Processing', but got '%s'", indicators[0].AssociatedRole)
	}
}
