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

package roleselector

import (
	"fmt"
	"testing"
	"time"

	"colmena.bsc.es/role-selector/fileloader"
	"colmena.bsc.es/role-selector/policy"
	"colmena.bsc.es/role-selector/roleselector"
	"colmena.bsc.es/role-selector/sla"
	"colmena.bsc.es/role-selector/types"
	"github.com/go-awaitility/awaitility"
)


func roleRunner(policy policy.Policy, kpiRetriever sla.KpiRetriever, roleRunner roleselector.RoleRunner) *roleselector.RoleSelector {
	hardware := "SENSOR"
	serviceDescriptionChan := make(chan types.ServiceDescription)
	lazyPolicy := policy
	selector := &roleselector.RoleSelector{
		ServiceDescriptionChan: serviceDescriptionChan,
		AlertsChan:             make(chan types.Alert),
		Hardware:               hardware,
		Policy:                 lazyPolicy,
		RoleRunner:             roleRunner,
		KpiRetriever:           kpiRetriever,
	}
	return selector
}

/*
Given an eager role selector, when a service description is received, then the role selector should run the role.
*/
func TestRunEagerRole(t *testing.T) {	
	mockRunner := &MockRoleRunner{}
	selector := roleRunner(&policy.EagerPolicy{}, &MockKpiRetrieverBroken{}, mockRunner)
	go selector.Run(make(chan roleselector.StoppedEvent))

	fmt.Println("Reading service description...")
	sd, err := fileloader.LoadFromFile[types.ServiceDescription]("../resources/service_definition.json")
	if err != nil {
		fmt.Printf("Failed to load service definition: %v\n", err)
		return
	}
	selector.ServiceDescriptionChan <- *sd
	err = awaitility.Await(10 * time.Millisecond, 100 * time.Millisecond, func() bool {
		return mockRunner.RunCalled
	})
	if err != nil {
		t.Error("Run was not called")
	}

	if len(mockRunner.RunArgs) == 0 {
		t.Error("No Run arguments were recorded")
	} else {
		args := mockRunner.RunArgs[0]
		if args.RoleId != "Sensing" {
			t.Errorf("Expected roleId to be 'Sensing', got '%s'", args.RoleId)
		}
		if args.ServiceId != "ExampleSensorprocessor" {
			t.Errorf("Expected serviceId to be 'ExampleSensorprocessor', got '%s'", args.ServiceId)
		}
	}
}

/*
Given a lazy role selector, when a service description is received, then the role selector should run the role if it has no associated KPIs.
*/
func TestLazyWithoutKpis(t *testing.T) {
	mockRunner := &MockRoleRunner{}
	selector := roleRunner(&policy.LazyPolicy{}, &MockKpiRetrieverNoKpis{}, mockRunner)
	go selector.Run(make(chan roleselector.StoppedEvent))
	
	sdBuilder := NewTestServiceDescriptionBuilder("ExampleSensorprocessor")
	sdBuilder.AddDockerRoleDefinition(types.DockerRoleDefinition{
		Id: "Sensing",
		ImageId: "sensing:latest",
		HardwareRequirements: []string{"SENSOR"},
		Kpis: []types.KpiDescription{},
		})
	sd := sdBuilder.Build()
	
	selector.ServiceDescriptionChan <- *sd

	err := awaitility.Await(100 * time.Millisecond, 100 * time.Millisecond, func() bool {
		return mockRunner.RunCalled
	})
	if err != nil {
		t.Error("Run was not called")
	}
}

/*
Given a lazy role selector, when a service description is received, then the role selector should not run the role until an alert is received.
*/
func TestLazyRunRoleAfterReceivingAlert(t *testing.T) {
	mockRunner := &MockRoleRunner{}
	selector := roleRunner(&policy.LazyPolicy{}, &MockKpiRetrieverMet{}, mockRunner)
	go selector.Run(make(chan roleselector.StoppedEvent))

	fmt.Println("Reading service description...")
	sd, err := fileloader.LoadFromFile[types.ServiceDescription]("../resources/service_definition.json")
	if err != nil {
		fmt.Printf("Failed to load service definition: %v\n", err)
		return
	}

	fmt.Println("Reading alert...")
	alert, err := fileloader.LoadFromFile[types.Alert]("../resources/alert.json")
	if err != nil {
		fmt.Printf("Failed to read alert: %v\n", err)
	}

	selector.ServiceDescriptionChan <- *sd
	// Wait for the selector to process the service description, it should not run the role yet because the KPI is met
	time.Sleep(100 * time.Millisecond)

	selector.AlertsChan <- *alert

	err = awaitility.Await(10 * time.Millisecond, 100 * time.Millisecond, func() bool {
		return mockRunner.RunCalled
	})
	if err != nil {
		t.Error("Run was not called")
	}

	// Verify that Run was called with correct arguments
	if len(mockRunner.RunArgs) == 0 {
		t.Error("No Run arguments were recorded")
	} else {
		args := mockRunner.RunArgs[0]
		if args.RoleId != "Sensing" {
			t.Errorf("Expected roleId to be 'Sensing', got '%s'", args.RoleId)
		}
		if args.ServiceId != "ExampleSensorprocessor" {
			t.Errorf("Expected serviceId to be 'ExampleSensorprocessor', got '%s'", args.ServiceId)
		}
	}
}

/*
Given a lazy role selector, when a service description is received, then the role selector should not run if the KPI is met.
*/
func TestLazyPolicyKPIMet(t *testing.T) {
	mockRunner := &MockRoleRunner{}
	selector := roleRunner(&policy.LazyPolicy{}, &MockKpiRetrieverMet{}, mockRunner)
	go selector.Run(make(chan roleselector.StoppedEvent))

	selector.ServiceDescriptionChan <- *ServiceDescription()

	err := awaitility.Await(100 * time.Millisecond, 100 * time.Millisecond, func() bool {
		return !mockRunner.RunCalled
	})
	if err != nil {
		t.Error("Run was called")
	}
}

/*
Given a lazy role selector, when a service description is received, then the role selector should run if the KPI is broken.
*/
func TestLazyPolicyKPIBroken(	t *testing.T) {
	mockRunner := &MockRoleRunner{}
	selector := roleRunner(&policy.LazyPolicy{}, &MockKpiRetrieverBroken{}, mockRunner)
	go selector.Run(make(chan roleselector.StoppedEvent))

	selector.ServiceDescriptionChan <- *ServiceDescription()

	err := awaitility.Await(100 * time.Millisecond, 100 * time.Millisecond, func() bool {
		return mockRunner.RunCalled
	})
	if err != nil {
		t.Error("Run was called")
	}
}

/*
Given an eager role selector, when an updated service description is received, 
then the role selector should stop the previous version of the role, wait until it is stopped, and then start the new version of the role.
*/
func TestEagerPolicyUpdateService(t *testing.T) {
	mockRunner := &MockRoleRunner{}
	selector := roleRunner(&policy.EagerPolicy{}, &MockKpiRetrieverNoKpis{}, mockRunner)
	stoppedChannel := make(chan roleselector.StoppedEvent)
	go selector.Run(stoppedChannel)

	initialServiceDescription := ServiceDescription()
	selector.ServiceDescriptionChan <- *initialServiceDescription

	err := awaitility.Await(100 * time.Millisecond, 100 * time.Millisecond, func() bool {
		return mockRunner.RunCalled
	})
	if err != nil {
		t.Error("Run was called")
	}
	
	updatedServiceDescription := UpdatedServiceDescription()
	selector.ServiceDescriptionChan <- *updatedServiceDescription

	// the previous version of the role is now stopped
	err = awaitility.Await(100 * time.Millisecond, 100 * time.Millisecond, func() bool {
		return mockRunner.StopCalled
	})
	if err != nil {
		t.Error("Stop was called")
	}

	// the new version of the role is now updating to stop the roleselector from starting it
	status, ok := selector.GetRoleStatus("Sensing", "ExampleSensorprocessor")
	if !ok {
		t.Error("Role not found")
	}
	if status != types.Updating {
		t.Errorf("Expected role to be updating, got %v", status)
	}

	// we now stop the previous version of the role
	stoppedChannel <- roleselector.StoppedEvent{
		ServiceId: initialServiceDescription.ServiceId.Value,
		RoleId: initialServiceDescription.DockerRoleDefinitions[0].Id,
		ImageId: initialServiceDescription.DockerRoleDefinitions[0].ImageId,
	}

	// run should be called with the new image id of the new version of the role
	err = awaitility.Await(100 * time.Millisecond, 100 * time.Millisecond, func() bool {
		return mockRoleRunnerExpectedCall(
			mockRunner, 
			updatedServiceDescription.ServiceId.Value, 
			updatedServiceDescription.DockerRoleDefinitions[0].Id, 
			updatedServiceDescription.DockerRoleDefinitions[0].ImageId)
	})
	if err != nil {
		t.Error("Run was not called")
	}
}

func ServiceDescription() *types.ServiceDescription {
	sd := NewTestServiceDescriptionBuilder("ExampleSensorprocessor")
	sd.AddDockerRoleDefinition(types.DockerRoleDefinition{
		Id: "Sensing",
		ImageId: "sensing:latest",
		HardwareRequirements: []string{"SENSOR"},
		Kpis: []types.KpiDescription{
			{
				Query: "avg_over_time(examplecontextdata_processing_time[5s]) < 15",
				Scope: "company_premises/building = .",
			},
		},
	})
	return sd.Build()
}

func UpdatedServiceDescription() *types.ServiceDescription {
	sd := NewTestServiceDescriptionBuilder("ExampleSensorprocessor")
	sd.AddDockerRoleDefinition(types.DockerRoleDefinition{
		Id: "Sensing",
		ImageId: "sensing:latest2",
		HardwareRequirements: []string{"SENSOR"},
		Kpis: []types.KpiDescription{
			{
				Query: "avg_over_time(examplecontextdata_processing_time[5s]) < 15",
				Scope: "company_premises/building = .",
			},
		},
	})
	return sd.Build()
}

func mockRoleRunnerExpectedCall(
	mockRunner *MockRoleRunner, 
	serviceId string, 
	roleId string, 
	imageId string) bool {
	runCalled := mockRunner.RunCalled
	found := false
	for _, arg := range mockRunner.RunArgs {
		if arg.RoleId == roleId && arg.ServiceId == serviceId && arg.ImageId == imageId {
			found = true
			break
		}
	}
	return runCalled && found
}