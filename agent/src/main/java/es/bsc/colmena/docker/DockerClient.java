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
package es.bsc.colmena.docker;

import com.github.dockerjava.api.command.CreateContainerResponse;
import com.github.dockerjava.api.command.PullImageResultCallback;
import com.github.dockerjava.api.exception.DockerClientException;
import com.github.dockerjava.api.exception.NotFoundException;
import com.github.dockerjava.api.exception.NotModifiedException;
import com.github.dockerjava.api.model.HostConfig;
import com.github.dockerjava.api.model.PullResponseItem;
import com.github.dockerjava.core.DefaultDockerClientConfig;
import com.github.dockerjava.core.DockerClientConfig;
import com.github.dockerjava.core.DockerClientImpl;
import com.github.dockerjava.httpclient5.ApacheDockerHttpClient;
import com.github.dockerjava.transport.DockerHttpClient;
import com.google.common.net.HostAndPort;
import es.bsc.colmena.infrastructure.Host;
import lombok.extern.slf4j.Slf4j;

import java.time.Duration;
import java.util.List;
import java.util.function.Consumer;

@Slf4j
public class DockerClient {
    private static final boolean REMOVE_CONTAINER = true;
    private final HostAndPort dcpHostPort;
    private final com.github.dockerjava.api.DockerClient dockerClient;
    private final DockerStorageLimit dockerStorageLimit;

    public DockerClient(HostAndPort dcpHostPort, Consumer<String> roleStopped) {
        this.dcpHostPort = dcpHostPort;
        dockerClient = createDockerClient();
        dockerStorageLimit = new DockerStorageLimit(dockerClient);
        DockerMonitor dockerMonitor = new DockerMonitor(roleStopped);
        dockerClient.eventsCmd().exec(dockerMonitor);
    }

    public static com.github.dockerjava.api.DockerClient createDockerClient() {
        DockerClientConfig config = DefaultDockerClientConfig.createDefaultConfigBuilder().build();
        DockerHttpClient httpClient = new ApacheDockerHttpClient.Builder()
                .dockerHost(config.getDockerHost())
                .sslConfig(config.getSSLConfig())
                .maxConnections(100)
                .connectionTimeout(Duration.ofSeconds(30))
                .responseTimeout(Duration.ofSeconds(45))
                .build();
        return DockerClientImpl.getInstance(config, httpClient);
    }

    public String run(String imageId) {
        try {
            String containerId = createContainer(imageId);
            dockerClient.startContainerCmd(containerId).exec();
            dockerStorageLimit.manage();
            log.info("Started containerId: {}, imageId: {}", containerId, imageId);
            return containerId;
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

    private synchronized String createContainer(String imageId) {
        List<String> environmentVariables = List.of(
                "DCP_IP_ADDRESS=" + dcpHostPort.getHost(),
                "HOSTNAME=" + Host.hostname());
        try {
            CreateContainerResponse response = dockerClient.createContainerCmd(imageId)
                    .withEnv(environmentVariables)
                    .withHostConfig(HostConfig.newHostConfig().withNetworkMode("host"))
                    .exec();
            return response.getId();
        }
        catch (NotFoundException notFoundException) {
            try {
                log.info("Image not found locally. Trying Docker Hub. imageId: {}", imageId);
                dockerClient.pullImageCmd(imageId).withTag("latest")
                        .exec(new PullImageResultCallback() {
                            @Override
                            public void onNext(PullResponseItem item) {super.onNext(item);
                            }
                        })
                        .awaitCompletion();
                log.info("Image pulled from docker hub, creating container. imageId: {}", imageId);
                return createContainer(imageId);
            } catch (DockerClientException e) {
                log.info("Caught DockerClientException when trying to pull imageId: {}. exceptionMessage: {}" +
                        "This can be because the image already exists." +
                        "Retrying container creation...", imageId, e.getMessage());
                return createContainer(imageId);
            } catch (InterruptedException e) {
                log.error("Interrupted while pulling image from Docker Hub. imageId: {}", imageId, e);
                throw new RuntimeException();
            }
        }
    }

    public void stop(String containerId) {
        try {
            dockerClient.stopContainerCmd(containerId).exec();
            log.info("Stopped containerId: {}", containerId);
        } catch (NotModifiedException nme) {
            log.error("Could not stop containerId: {}. Already stopped?", containerId, nme);
        }
    }

    public void removeContainer(String containerId) {
        log.info("Removing container {}, containerId: {}", REMOVE_CONTAINER, containerId);
        if (REMOVE_CONTAINER) {
            dockerClient.removeContainerCmd(containerId).exec();
        }
    }
}
