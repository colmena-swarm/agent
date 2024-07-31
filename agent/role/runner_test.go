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
	"os"
	"testing"
	"time"

	"colmena.bsc.es/agent/device"
	// "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockContainerEngine struct {
	mock.Mock
}
func (m *MockContainerEngine) RunContainer(name string, interfc string, device device.Info) (string, error) {
	args := m.Called(name)
	return args.String(0), nil
}
func (m *MockContainerEngine) StopContainer(containerId string) {
	m.Called(containerId)
}
func (m *MockContainerEngine) Subscribe(stopped chan string) {
	m.Called()
}

type MockKpiMatcher struct {
	mock.Mock
}

func (m *MockKpiMatcher) match(kpi KPI) bool {
	args := m.Called(kpi)
	kpisMet := args.Bool(0)
	return kpisMet
}

func TestEagerDeviceUsesRoleFeatures(t *testing.T) {
	tests := []struct {
		deviceFeatures []string
		expectContainerStarted bool
	}{
		{[]string{"CAMERA"}, true},
		{[]string{"CAMERA", "CPU"}, true},
		{[]string{}, false},
	}

	for _, tt := range tests {
		containerEngine := createMockContainerEngine(t, tt.expectContainerStarted)
		device := device.Info{Strategy: "EAGER", Features: tt.deviceFeatures}
		done := make(chan os.Signal, 1)
		found := make(chan []Role)
	
		go Run(device, done, found, containerEngine, nil)
		found <-foundRoles([]KPI{})
		time.Sleep(1 * time.Second)
	
		containerEngine.AssertExpectations(t)
	}
}

func TestLazyDeviceUsesDeviceFeaturesAndKpis(t *testing.T) {
	tests := []struct {
		deviceFeatures []string
		kpiMet bool
		expectContainerStarted bool
	}{
		{[]string{"CAMERA"}, true, false},
		{[]string{"CAMERA"}, false, true},
		{[]string{}, false, false},
	}
	for _, tt := range tests {
		containerEngine := createMockContainerEngine(t, tt.expectContainerStarted)
		device := device.Info{Strategy: "LAZY", Features: tt.deviceFeatures}
		done := make(chan os.Signal, 1)
		found := make(chan []Role)
		kpi := KPI{
			Key:           "second",
			Threshold:     51.55,
			FromUnit:      5,
			FromType:	  "m",
		}
	
		mockKpiMatcher := new (MockKpiMatcher)
		mockKpiMatcher.On("match", kpi).Return(tt.kpiMet)
	
		go Run(device, done, found, containerEngine, mockKpiMatcher)
		found <-foundRoles([]KPI{kpi})
		time.Sleep(1 * time.Second)
	
		containerEngine.AssertExpectations(t)
	}
}

func foundRoles(kpis []KPI) []Role {
	return []Role {
		{Id: "RoleId", ImageId: "ImageId", HardwareRequirements: []string{"CAMERA"}, Kpis: kpis},
	}
}

func createMockContainerEngine(t *testing.T, expectContainerStarted bool) *MockContainerEngine {
	containerEngine := new(MockContainerEngine)
	if (expectContainerStarted) {
		containerEngine.On("RunContainer", "ImageId").Return("ContainerId")
	} else {
		containerEngine.AssertNotCalled(t, "RunContainer")
	}
	containerEngine.On("Subscribe")
	return containerEngine
}
