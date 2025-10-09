package roleselector

import "colmena.bsc.es/role-selector/types"

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
		RemoveContainer bool
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

func (m *MockRoleRunner) Stop(roleId string, serviceId string, imageId string, removeContainer bool) {
	m.StopCalled = true
	m.StopArgs = append(m.StopArgs, struct {
		RoleId    string
		ServiceId string
		ImageId   string
		RemoveContainer bool
	}{roleId, serviceId, imageId, removeContainer})
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
			Query:          "avg_over_time(examplecontextdata_processing_time[5s]) < 15",
			Value:          2,
			Threshold:      1,
			Level:          "Broken",
			AssociatedRole: "Sensing",
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
			Query:          "avg_over_time(examplecontextdata_processing_time[5m]) < 15",
			Value:          0,
			Threshold:      1,
			Level:          "Met",
			AssociatedRole: "Sensing",
		},
	}, nil
}

type MockKpiRetrieverNoKpis struct {
	GetCalled bool
	GetArgs   []string
}

func (m *MockKpiRetrieverNoKpis) Get(serviceId string) ([]types.KPI, error) {
	m.GetCalled = true
	m.GetArgs = append(m.GetArgs, serviceId)
	return []types.KPI{}, nil
}