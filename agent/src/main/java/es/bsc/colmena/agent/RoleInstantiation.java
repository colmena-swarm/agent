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
package es.bsc.colmena.agent;

import es.bsc.colmena.docker.DockerRole;
import es.bsc.colmena.library.*;
import es.bsc.colmena.library.roledefinition.DockerImageRoleDefinition;
import es.bsc.colmena.library.roledefinition.JavaCodeRoleDefinition;
import lombok.SneakyThrows;
import org.apache.commons.lang3.NotImplementedException;

public class RoleInstantiation {

    @SneakyThrows
    public static Role instantiate(RoleDefinition roleDefinition, ColmenaAgent colmenaAgent) {
        if (roleDefinition instanceof DockerImageRoleDefinition dockerImageRoleDefinition) {
            return new DockerRole(roleDefinition.getRoleId(), dockerImageRoleDefinition.getImageId());
        }

        if (roleDefinition instanceof JavaCodeRoleDefinition javaCodeRoleDefinition) {
            Class<? extends JavaRole> clazz = javaCodeRoleDefinition.getRole();
            return JavaRoleInstantiation.instantiate(clazz, colmenaAgent.getDcpClient());
        }

        throw new NotImplementedException();
    }
}
