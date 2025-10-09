package roleselector

import (
	"testing"

	"colmena.bsc.es/role-selector/roleselector"
	"colmena.bsc.es/role-selector/types"
)

func TestClean(t *testing.T) {
	mockRunner := &MockRoleRunner{}	
	roleselector.Clean("serviceId", []*types.Role{
		{
			Id: "role1",
			ImageId: "image1",
			State: types.Running,
		},
	}, []*types.Role{}, mockRunner)

	if !mockRunner.StopCalled {
		t.Error("Stop was not called")
	}
}

func TestCleanWithNoRolesRunning(t *testing.T) {
	mockRunner := &MockRoleRunner{}
	roleselector.Clean("serviceId", []*types.Role{}, []*types.Role{
		{
			Id: "role1",
			ImageId: "image1",
			State: types.Stopped,
		},
	}, mockRunner)

	if mockRunner.StopCalled {
		t.Error("Stop was called")
	}
}

func TestCleanWithNewRole(t *testing.T) {
	mockRunner := &MockRoleRunner{}
	roleselector.Clean("serviceId", []*types.Role{
		{
			Id: "role1",
			ImageId: "image1",
			State: types.Running,
		},
	}, []*types.Role{
		{
			Id: "role1",
			ImageId: "image2",
			State: types.Stopped,
		},
	}, mockRunner)

	if !mockRunner.StopCalled {
		t.Error("Stop was not called")
	}

	if len(mockRunner.StopArgs) != 1 {
		t.Errorf("Expected 1 stop argument, got %d", len(mockRunner.StopArgs))
	}

	if mockRunner.StopArgs[0].ImageId != "image1" {
		t.Errorf("Expected imageId to be 'image1', got '%s'", mockRunner.StopArgs[0].ImageId)
	}
}