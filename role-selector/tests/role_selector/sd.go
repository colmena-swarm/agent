package roleselector

import "colmena.bsc.es/role-selector/types"


type TestServiceDescriptionBuilder struct {
	ServiceId string
	Kpis []types.KpiDescription
	HardwareRequirements []string
	DockerRoleDefinitions []types.DockerRoleDefinition
}

func NewTestServiceDescriptionBuilder(serviceId string) *TestServiceDescriptionBuilder {
	return &TestServiceDescriptionBuilder{ServiceId: serviceId}
}

func (b *TestServiceDescriptionBuilder) AddDockerRoleDefinition(dockerRoleDefinition types.DockerRoleDefinition) *TestServiceDescriptionBuilder {
	b.DockerRoleDefinitions = append(b.DockerRoleDefinitions, dockerRoleDefinition)
	return b
}

func (b *TestServiceDescriptionBuilder) AddHardwareRequirement(hardwareRequirement string) *TestServiceDescriptionBuilder {
	b.HardwareRequirements = append(b.HardwareRequirements, hardwareRequirement)
	return b
}

func (b *TestServiceDescriptionBuilder) Build() *types.ServiceDescription {
	return &types.ServiceDescription{
		ServiceId: types.Id{
			Value: b.ServiceId,
		},
		Kpis: b.Kpis,
		DockerRoleDefinitions: b.DockerRoleDefinitions,
	}
}

func (b *TestServiceDescriptionBuilder) AddKpi(kpi types.KpiDescription) *TestServiceDescriptionBuilder {
	b.Kpis = append(b.Kpis, kpi)
	return b
}
