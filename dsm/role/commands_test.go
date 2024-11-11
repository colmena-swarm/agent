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
package role

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

// MockContainerEngine implements the ContainerEngine interface for testing
type MockContainerEngine struct {
	mu                 sync.RWMutex
	containers         map[string]bool // roleId -> isRunning
	runContainerCalls  []RunContainerCall
	stopContainerCalls []StopContainerCall
	runContainerError  error
	stopContainerError error
}

type RunContainerCall struct {
	RoleId  string
	ImageId string
	AgentId string
	Interfc string
}

type StopContainerCall struct {
	ContainerId string
}

func NewMockContainerEngine() *MockContainerEngine {
	return &MockContainerEngine{
		containers:         make(map[string]bool),
		runContainerCalls:  make([]RunContainerCall, 0),
		stopContainerCalls: make([]StopContainerCall, 0),
	}
}

func (m *MockContainerEngine) RunContainer(roleId string, imageId string, agentId string, interfc string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.runContainerCalls = append(m.runContainerCalls, RunContainerCall{
		RoleId:  roleId,
		ImageId: imageId,
		AgentId: agentId,
		Interfc: interfc,
	})

	if m.runContainerError != nil {
		return "", m.runContainerError
	}

	m.containers[roleId] = true
	return fmt.Sprintf("container-%s", roleId), nil
}

func (m *MockContainerEngine) StopContainer(containerId string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.stopContainerCalls = append(m.stopContainerCalls, StopContainerCall{
		ContainerId: containerId,
	})

	if m.stopContainerError != nil {
		return m.stopContainerError
	}

	// Find and stop the container by roleId (assuming containerId is roleId for simplicity)
	for roleId := range m.containers {
		if strings.Contains(containerId, roleId) {
			m.containers[roleId] = false
			break
		}
	}

	return nil
}

func (m *MockContainerEngine) Subscribe(stopped chan string) {
	// Mock implementation - not used in these tests
}

// Test helper methods
func (m *MockContainerEngine) GetRunContainerCalls() []RunContainerCall {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]RunContainerCall{}, m.runContainerCalls...)
}

func (m *MockContainerEngine) GetStopContainerCalls() []StopContainerCall {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]StopContainerCall{}, m.stopContainerCalls...)
}

func (m *MockContainerEngine) SetRunContainerError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.runContainerError = err
}

func (m *MockContainerEngine) SetStopContainerError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stopContainerError = err
}

func (m *MockContainerEngine) IsContainerRunning(roleId string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.containers[roleId]
}

func TestCommandListener_RegisterEndpoints(t *testing.T) {
	mockEngine := NewMockContainerEngine()
	commandListener := CommandListener{
		AgentId:         "test-agent-123",
		Interfc:         "eth0",
		ContainerEngine: mockEngine,
	}

	handler := commandListener.Endpoints()

	// Test that we get a valid handler
	if handler == nil {
		t.Fatal("RegisterEndpoints() returned nil handler")
	}
}

func TestStartEndpoint_Success(t *testing.T) {
	mockEngine := NewMockContainerEngine()
	commandListener := CommandListener{
		AgentId:         "test-agent-123",
		Interfc:         "eth0",
		ContainerEngine: mockEngine,
	}

	server := httptest.NewServer(commandListener.Endpoints())
	defer server.Close()

	// Test data
	roleCmd := RoleCommand{
		ServiceId: "service-123",
		RoleId:    "role-456",
		ImageId:   "test-image:latest",
	}

	jsonData, err := json.Marshal(roleCmd)
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}

	// Make request
	resp, err := http.Post(server.URL+"/start", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Wait a bit for the goroutine to execute
	time.Sleep(100 * time.Millisecond)

	// Check that RunContainer was called with correct parameters
	calls := mockEngine.GetRunContainerCalls()
	if len(calls) != 1 {
		t.Fatalf("Expected 1 RunContainer call, got %d", len(calls))
	}

	call := calls[0]
	if call.RoleId != "role-456" {
		t.Errorf("Expected RoleId 'role-456', got '%s'", call.RoleId)
	}
	if call.ImageId != "test-image:latest" {
		t.Errorf("Expected ImageId 'test-image:latest', got '%s'", call.ImageId)
	}
	if call.AgentId != "test-agent-123" {
		t.Errorf("Expected AgentId 'test-agent-123', got '%s'", call.AgentId)
	}
	if call.Interfc != "eth0" {
		t.Errorf("Expected Interfc 'eth0', got '%s'", call.Interfc)
	}

	// Check that container is marked as running
	if !mockEngine.IsContainerRunning("role-456") {
		t.Error("Expected container to be running")
	}
}

func TestStartEndpoint_InvalidJSON(t *testing.T) {
	mockEngine := NewMockContainerEngine()
	commandListener := CommandListener{
		AgentId:         "test-agent-123",
		Interfc:         "eth0",
		ContainerEngine: mockEngine,
	}

	server := httptest.NewServer(commandListener.Endpoints())
	defer server.Close()

	// Send invalid JSON
	invalidJSON := `{"serviceId": "service-123", "roleId": "role-456", "imageId": "test-image:latest"` // Missing closing brace
	resp, err := http.Post(server.URL+"/start", "application/json", bytes.NewBufferString(invalidJSON))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}

	// Check that RunContainer was not called
	calls := mockEngine.GetRunContainerCalls()
	if len(calls) != 0 {
		t.Errorf("Expected 0 RunContainer calls, got %d", len(calls))
	}
}

func TestStopEndpoint_Success(t *testing.T) {
	mockEngine := NewMockContainerEngine()
	commandListener := CommandListener{
		AgentId:         "test-agent-123",
		Interfc:         "eth0",
		ContainerEngine: mockEngine,
	}

	server := httptest.NewServer(commandListener.Endpoints())
	defer server.Close()

	// Test data
	roleCmd := RoleCommand{
		ServiceId: "service-123",
		RoleId:    "role-456",
		ImageId:   "test-image:latest",
	}

	jsonData, err := json.Marshal(roleCmd)
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}

	// Make request
	resp, err := http.Post(server.URL+"/stop", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Wait a bit for the goroutine to execute
	time.Sleep(100 * time.Millisecond)

	// Check that StopContainer was called with correct parameters
	calls := mockEngine.GetStopContainerCalls()
	if len(calls) != 1 {
		t.Fatalf("Expected 1 StopContainer call, got %d", len(calls))
	}

	call := calls[0]
	if call.ContainerId != "role-456" {
		t.Errorf("Expected ContainerId 'role-456', got '%s'", call.ContainerId)
	}
}

func TestStopEndpoint_InvalidJSON(t *testing.T) {
	mockEngine := NewMockContainerEngine()
	commandListener := CommandListener{
		AgentId:         "test-agent-123",
		Interfc:         "eth0",
		ContainerEngine: mockEngine,
	}

	server := httptest.NewServer(commandListener.Endpoints())
	defer server.Close()

	// Send invalid JSON
	invalidJSON := `{"serviceId": "service-123", "roleId": "role-456"` // Missing closing brace and imageId
	resp, err := http.Post(server.URL+"/stop", "application/json", bytes.NewBufferString(invalidJSON))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}

	// Check that StopContainer was not called
	calls := mockEngine.GetStopContainerCalls()
	if len(calls) != 0 {
		t.Errorf("Expected 0 StopContainer calls, got %d", len(calls))
	}
}

func TestMultipleRequests(t *testing.T) {
	mockEngine := NewMockContainerEngine()
	commandListener := CommandListener{
		AgentId:         "test-agent-123",
		Interfc:         "eth0",
		ContainerEngine: mockEngine,
	}

	server := httptest.NewServer(commandListener.Endpoints())
	defer server.Close()

	// Test multiple start requests
	roleCommands := []RoleCommand{
		{ServiceId: "service-1", RoleId: "role-1", ImageId: "image-1"},
		{ServiceId: "service-2", RoleId: "role-2", ImageId: "image-2"},
		{ServiceId: "service-3", RoleId: "role-3", ImageId: "image-3"},
	}

	for _, cmd := range roleCommands {
		jsonData, err := json.Marshal(cmd)
		if err != nil {
			t.Fatalf("Failed to marshal test data: %v", err)
		}

		resp, err := http.Post(server.URL+"/start", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200 for role %s, got %d", cmd.RoleId, resp.StatusCode)
		}
	}

	// Wait for all goroutines to complete
	time.Sleep(200 * time.Millisecond)

	// Check that all containers were started
	calls := mockEngine.GetRunContainerCalls()
	if len(calls) != 3 {
		t.Fatalf("Expected 3 RunContainer calls, got %d", len(calls))
	}

	// Check that all containers are running
	for _, cmd := range roleCommands {
		if !mockEngine.IsContainerRunning(cmd.RoleId) {
			t.Errorf("Expected container %s to be running", cmd.RoleId)
		}
	}
}

func TestConcurrentRequests(t *testing.T) {
	mockEngine := NewMockContainerEngine()
	commandListener := CommandListener{
		AgentId:         "test-agent-123",
		Interfc:         "eth0",
		ContainerEngine: mockEngine,
	}

	server := httptest.NewServer(commandListener.Endpoints())
	defer server.Close()

	// Test concurrent requests
	numRequests := 10
	done := make(chan bool, numRequests)

	for i := 0; i < numRequests; i++ {
		go func(i int) {
			defer func() { done <- true }()

			roleCmd := RoleCommand{
				ServiceId: fmt.Sprintf("service-%d", i),
				RoleId:    fmt.Sprintf("role-%d", i),
				ImageId:   fmt.Sprintf("image-%d", i),
			}

			jsonData, err := json.Marshal(roleCmd)
			if err != nil {
				t.Errorf("Failed to marshal test data: %v", err)
				return
			}

			resp, err := http.Post(server.URL+"/start", "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				t.Errorf("Failed to make request: %v", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status 200 for role %d, got %d", i, resp.StatusCode)
			}
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numRequests; i++ {
		<-done
	}

	// Wait a bit more for all container operations to complete
	time.Sleep(200 * time.Millisecond)

	// Check that all containers were started
	calls := mockEngine.GetRunContainerCalls()
	if len(calls) != numRequests {
		t.Fatalf("Expected %d RunContainer calls, got %d", numRequests, len(calls))
	}
}

func TestParseRoleCommand(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		expected RoleCommand
		hasError bool
	}{
		{
			name:     "Valid JSON",
			jsonData: `{"serviceId": "service-123", "roleId": "role-456", "imageId": "image-789"}`,
			expected: RoleCommand{ServiceId: "service-123", RoleId: "role-456", ImageId: "image-789"},
			hasError: false,
		},
		{
			name:     "Empty JSON",
			jsonData: `{}`,
			expected: RoleCommand{},
			hasError: false,
		},
		{
			name:     "Invalid JSON",
			jsonData: `{"serviceId": "service-123", "roleId": "role-456"`,
			expected: RoleCommand{},
			hasError: true,
		},
		{
			name:     "Missing fields",
			jsonData: `{"serviceId": "service-123"}`,
			expected: RoleCommand{ServiceId: "service-123"},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(tt.jsonData))
			req.Header.Set("Content-Type", "application/json")

			result, err := parseRoleCommand(req)

			if tt.hasError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Expected %+v, got %+v", tt.expected, result)
				}
			}
		})
	}
}
