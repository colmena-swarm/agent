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
package es.bsc.colmena.rolerunner;

import com.google.common.net.HostAndPort;
import es.bsc.colmena.docker.DockerClient;
import es.bsc.colmena.docker.DockerRole;
import es.bsc.colmena.library.Role;
import lombok.extern.slf4j.Slf4j;

import java.util.HashSet;
import java.util.Set;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ConcurrentMap;

@Slf4j
public class DockerRoleRunner {
    private final DockerClient dockerClient;
    private final ConcurrentMap<String, DockerRole> runningRoles = new ConcurrentHashMap<>();

    public DockerRoleRunner(HostAndPort dcpHostPort) {
        dockerClient = new DockerClient(dcpHostPort, this::roleStopped);
    }

    public synchronized void startRole(DockerRole dockerRole) {
        log.info("Starting roleId: {}, containerId: {}", dockerRole.getRoleId(), dockerRole.getContainerId());
        String containerId = dockerClient.run(dockerRole.getImageId());
        dockerRole.setContainerId(containerId);
        runningRoles.put(dockerRole.getRoleId(), dockerRole);
        log.info("Started roleId: {}, containerId: {}", dockerRole.getRoleId(), dockerRole.getContainerId());
    }

    public synchronized void stopRole(DockerRole dockerRole) {
        log.info("Stopping roleId: {}, containerId: {}", dockerRole.getRoleId(), dockerRole.getContainerId());
        dockerClient.stop(dockerRole.getContainerId());
        dockerClient.removeContainer(dockerRole.getContainerId());
        runningRoles.remove(dockerRole.getRoleId());
        log.info("Stopped roleId: {}, containerId: {}", dockerRole.getRoleId(), dockerRole.getContainerId());
    }

    public synchronized void roleStopped(String containerId) {
        runningRoles.values().stream()
                .filter(role -> role.getContainerId().equals(containerId))
                .forEach(role -> {
                        log.info("Role stopped. Removing container. roleId: {}, containerId: {}", role.getRoleId(), role.getContainerId());
                        dockerClient.removeContainer(role.getContainerId());
                        runningRoles.remove(role.getRoleId());
                });
    }

    public synchronized void disconnect() throws InterruptedException {
        log.info("Disconnecting. Stopping {} roles", runningRoles.size());
        runningRoles.values().forEach(this::stopRole);
    }

    public synchronized Set<Role> getRunningRoles() {
        return new HashSet<>(runningRoles.values());
    }
}
