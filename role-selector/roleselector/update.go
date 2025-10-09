package roleselector

import (
	"log"

	"colmena.bsc.es/role-selector/types"
)

func Clean(serviceId string, old []*types.Role, new []*types.Role, roleRunner RoleRunner) {
	newRolesByRoleId := make(map[string]*types.Role)
	for _, newRole := range new {
		newRolesByRoleId[newRole.Id] = newRole
	}
	for _, oldRole := range old {
		if oldRole.State == types.Running {
			newRole, ok := newRolesByRoleId[oldRole.Id]
			if !ok {
				log.Printf("Role %v not found in new version of service %v, stopping...", oldRole.Id, serviceId)
				roleRunner.Stop(oldRole.Id, serviceId, oldRole.ImageId, true)
				oldRole.State = types.Stopped
				continue
			}

			if oldRole.ImageId == newRole.ImageId {
				log.Printf("Role %v is already up to date", oldRole.Id)
				continue
			}

			log.Printf("Stopping %v to update to %v", oldRole.ImageId, newRole.ImageId)
			roleRunner.Stop(oldRole.Id, serviceId, oldRole.ImageId, true)

			//prevent the role selector from starting either role until the previous has been stopped
			oldRole.State = types.Updating
			newRole.State = types.Updating
		}
	}
}

