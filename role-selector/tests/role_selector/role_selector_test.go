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
	"colmena.bsc.es/role-selector/types"
)

type MockRoleRunner struct {
	RunCalled  bool
	StopCalled bool
	RunArgs    []struct {
		RoleId    string
		ServiceId string
		ImageId   string
	}
	StopArgs []struct {
		RoleId    string
		ServiceId string
		ImageId   string
	}
}

func (m *MockRoleRunner) Run(roleId string, serviceId string, imageId string) {
	m.RunCalled = true
	m.RunArgs = append(m.RunArgs, struct {
		RoleId    string
		ServiceId string
		ImageId   string
	}{roleId, serviceId, imageId})
}

func (m *MockRoleRunner) Stop(roleId string, serviceId string, imageId string) {
	m.StopCalled = true
	m.StopArgs = append(m.StopArgs, struct {
		RoleId    string
		ServiceId string
		ImageId   string
	}{roleId, serviceId, imageId})
}

type MockKpiRetrieverBroken struct {
	GetCalled bool
	GetArgs   []string
}

func (m *MockKpiRetrieverBroken) Get(serviceId string) ([]types.KPI, error) {
	m.GetCalled = true
	m.GetArgs = append(m.GetArgs, serviceId)
	return []types.KPI{
		{
			Query:          "avg_over_time(examplecontextdata_processing_time[5s]) < 1",
			Value:          2,
			Threshold:      1,
			Level:          "Broken",
			AssociatedRole: "Processing",
		},
	}, nil
}

type MockKpiRetrieverMet struct {
	GetCalled bool
	GetArgs   []string
}

func (m *MockKpiRetrieverMet) Get(serviceId string) ([]types.KPI, error) {
	m.GetCalled = true
	m.GetArgs = append(m.GetArgs, serviceId)
	return []types.KPI{
		{
			Query:          "avg_over_time(examplecontextdata_processing_time[5m]) < 1",
			Value:          0,
			Threshold:      1,
			Level:          "Met",
			AssociatedRole: "Processing",
		},
	}, nil
}

func TestRunEagerRole(t *testing.T) {
	hardware := "SENSOR"
	serviceDescriptionChan := make(chan types.ServiceDescription)
	lazyPolicy := &policy.LazyPolicy{}
	mockRunner := &MockRoleRunner{}
	selector := &roleselector.RoleSelector{
		ServiceDescriptionChan: serviceDescriptionChan,
		AlertsChan:             make(chan types.Alert),
		Hardware:               hardware,
		Policy:                 lazyPolicy,
		RoleRunner:             mockRunner,
		KpiRetriever:           &MockKpiRetrieverBroken{},
	}
	go selector.Run()

	fmt.Println("Reading service description...")
	sd, err := fileloader.LoadFromFile[types.ServiceDescription]("../resources/service_definition.json")
	if err != nil {
		fmt.Printf("Failed to load service definition: %v\n", err)
		return
	}
	serviceDescriptionChan <- *sd
	// Wait for the selector to process the service description. Now it's too long because of the sleep in the selector.
	time.Sleep(6 * time.Second)

	// Verify that Run was called with correct arguments
	if !mockRunner.RunCalled {
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

func TestRunRoleAfterReceivingAlert(t *testing.T) {
	hardware := "CPU"
	serviceDescriptionChan := make(chan types.ServiceDescription)
	alertsChan := make(chan types.Alert)
	lazyPolicy := &policy.LazyPolicy{}
	mockRunner := &MockRoleRunner{}
	mockKpiRetriever := &MockKpiRetrieverBroken{}
	selector := &roleselector.RoleSelector{
		ServiceDescriptionChan: serviceDescriptionChan,
		AlertsChan:             alertsChan,
		Hardware:               hardware,
		Policy:                 lazyPolicy,
		RoleRunner:             mockRunner,
		KpiRetriever:           mockKpiRetriever,
	}
	go selector.Run()

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

	serviceDescriptionChan <- *sd
	// Wait for the selector to process the service description
	time.Sleep(100 * time.Millisecond)
	alertsChan <- *alert
	// Wait for alert to process
	time.Sleep(100 * time.Millisecond)

	// Verify that Run was called with correct arguments
	if !mockRunner.RunCalled {
		t.Error("Run was not called")
	}

	if len(mockRunner.RunArgs) == 0 {
		t.Error("No Run arguments were recorded")
	} else {
		args := mockRunner.RunArgs[0]
		if args.RoleId != "Processing" {
			t.Errorf("Expected roleId to be 'Processing', got '%s'", args.RoleId)
		}
		if args.ServiceId != "ExampleSensorprocessor" {
			t.Errorf("Expected serviceId to be 'ExampleSensorprocessor', got '%s'", args.ServiceId)
		}
	}
}

func TestLazyPolicyKPIMet(t *testing.T) {
	hardware := "CPU"
	serviceDescriptionChan := make(chan types.ServiceDescription)
	alertsChan := make(chan types.Alert)
	lazyPolicy := &policy.LazyPolicy{}
	mockRunner := &MockRoleRunner{}
	mockKpiRetriever := &MockKpiRetrieverMet{}
	selector := &roleselector.RoleSelector{
		ServiceDescriptionChan: serviceDescriptionChan,
		AlertsChan:             alertsChan,
		Hardware:               hardware,
		Policy:                 lazyPolicy,
		RoleRunner:             mockRunner,
		KpiRetriever:           mockKpiRetriever,
	}
	go selector.Run()

	fmt.Println("Reading service description...")
	sd, err := fileloader.LoadFromFile[types.ServiceDescription]("../resources/service_definition.json")
	if err != nil {
		fmt.Printf("Failed to load service definition: %v\n", err)
		return
	}

	serviceDescriptionChan <- *sd

	time.Sleep(2 * time.Second)

	if mockRunner.RunCalled {
		t.Error("Run was called")
	}
}
