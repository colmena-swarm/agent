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
package es.bsc.colmena.util;

import com.github.dockerjava.api.exception.DockerClientException;
import es.bsc.colmena.docker.DockerClient;
import lombok.extern.slf4j.Slf4j;

import java.io.File;
import java.util.Set;

@Slf4j
public class TestDockerClient {
    com.github.dockerjava.api.DockerClient dockerClient;

    public TestDockerClient() {
        this.dockerClient = DockerClient.createDockerClient();
    }

    public String build(String tag, String dockerfileLocation) {
        try {
            return dockerClient.buildImageCmd(new File(dockerfileLocation))
                    .withTags(Set.of(tag))
                    .withNoCache(true)
                    .start()
                    .awaitImageId();
        }
        catch (DockerClientException dce) {
            log.error("Caught DockerClientException. " +
                    "This might be because containerd is being used for storing/pulling images for multi-platform build." +
                    "More info: https://gitlab.bsc.es/wdc/projects/colmena/-/issues/2", dce);
            throw new RuntimeException();
        }
    }

    public void deleteImage(String imageId) {
        dockerClient.removeImageCmd(imageId).exec();
    }
}
