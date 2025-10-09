package roleselector

import (
	"colmena.bsc.es/role-selector/types"
)

func parse(roleDescriptions []types.DockerRoleDefinition, serviceId string) []*types.Role {
	var roles []*types.Role
	for _, roleDescription := range roleDescriptions {
		//role initialised as stopped to be considered for starting
		toRun := types.Stopped
		role := &types.Role{
			Id:        roleDescription.Id,
			ImageId:   roleDescription.ImageId,
			ServiceId: serviceId,
			Kpis:      roleDescription.Kpis,
			State:     toRun,
			Resources: DefaultResources,
		}
		roles = append(roles, role)
	}
	return roles
}
