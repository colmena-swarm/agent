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

import es.bsc.colmena.agent.Device;
import es.bsc.colmena.library.roledefinition.DockerImageRoleDefinition;
import es.bsc.colmena.util.DockerColmenaTest;
import es.bsc.colmena.util.DockerProperties;
import es.bsc.colmena.util.TestDeviceFactory;
import lombok.extern.slf4j.Slf4j;
import org.junit.jupiter.api.Test;

import java.util.Set;

import static es.bsc.colmena.library.requirement.HardwareRequirement.CAMERA;
import static es.bsc.colmena.library.requirement.HardwareRequirement.CPU;
import static org.awaitility.Awaitility.await;

@Slf4j
public class SingleSubTest extends DockerColmenaTest {

    @Test
    public void pub_sub_using_two_docker_images() throws InterruptedException {
        String pubImageId = dockerClient.build(DockerProperties.formatImageName("pub"), "../roles/pub/Dockerfile");
        String subImageId = dockerClient.build(DockerProperties.formatImageName("sub"), "../roles/singlesub/Dockerfile");

        DockerImageRoleDefinition pubRoleDef = new DockerImageRoleDefinition("PUBROLE", pubImageId, Set.of(CAMERA));
        DockerImageRoleDefinition subRoleDef = new DockerImageRoleDefinition("SUBROLE", subImageId, Set.of(CPU));

        Device processor = TestDeviceFactory.aDevice(Set.of(pubRoleDef, subRoleDef), Set.of(CPU));
        processor.start();
        Device camera = TestDeviceFactory.aDevice(Set.of(pubRoleDef, subRoleDef), Set.of(CAMERA));
        camera.start();

        await().until(() -> distributedColmenaPlatform.getStorage().getQueues().get("ack") != null);
        await().until(() -> distributedColmenaPlatform.getStorage().getQueues().get("ack").size() > 5);

        log.info("disconnecting devices");
        camera.disconnect();
        processor.disconnect();

        //dockerClient.deleteImage(pubImageId);
        //dockerClient.deleteImage(subImageId);
    }

}
