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

import com.google.common.collect.Sets;
import com.google.common.net.HostAndPort;
import es.bsc.colmena.docker.DockerRole;
import es.bsc.colmena.library.JavaRole;
import es.bsc.colmena.library.Role;
import es.bsc.colmena.rolerunner.DockerRoleRunner;
import es.bsc.colmena.rolerunner.JavaRoleRunner;
import lombok.extern.slf4j.Slf4j;

import java.util.Map;
import java.util.Set;
import java.util.function.Function;
import java.util.stream.Collectors;

@Slf4j
public class RoleRunner {
    private final JavaRoleRunner javaRoleRunner = new JavaRoleRunner();
    private final DockerRoleRunner dockerRoleRunner;

    public RoleRunner(HostAndPort dcpHostPort) {
        dockerRoleRunner = new DockerRoleRunner(dcpHostPort);
    }

    public Map<String, Role> getRunningRoles() {
        Set<Role> dockerRoles = dockerRoleRunner.getRunningRoles();
        Set<Role> javaRoles = javaRoleRunner.getRunningRoles();
        return Sets.union(dockerRoles, javaRoles).stream()
                .collect(Collectors.toMap(Role::getRoleId, Function.identity()));
    }

    public void startRole(Role role) {
        if (role instanceof DockerRole dockerRole) {
            dockerRoleRunner.startRole(dockerRole);
        }
        if (role instanceof JavaRole javaRole) {
            javaRoleRunner.startRole(javaRole);
        }
    }

    public void stopRole(Role role) {
        if (role instanceof DockerRole dockerRole) {
            dockerRoleRunner.stopRole(dockerRole);
        }
        if (role instanceof JavaRole javaRole) {
            javaRoleRunner.stopRole(javaRole);
        }
    }

    public void disconnect() throws InterruptedException {
        dockerRoleRunner.disconnect();
        javaRoleRunner.disconnect();
    }

    public void blockUntilShutdown() throws InterruptedException {
        javaRoleRunner.blockUntilShutdown();
    }
}
