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

import es.bsc.colmena.agent.Device;
import es.bsc.colmena.library.RoleDefinition;
import es.bsc.colmena.library.metrics.ThresholdType;
import es.bsc.colmena.library.roledefinition.DockerImageRoleDefinition;
import es.bsc.colmena.util.DockerColmenaTest;
import es.bsc.colmena.util.DockerProperties;
import es.bsc.colmena.util.TestDeviceFactory;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;

import java.time.temporal.ChronoUnit;
import java.util.Set;

import static es.bsc.colmena.library.requirement.HardwareRequirement.CAMERA;
import static es.bsc.colmena.library.requirement.HardwareRequirement.CPU;
import static org.awaitility.Awaitility.await;

public class EagerLazyTest extends DockerColmenaTest {

    @Test
    public void lazy_processor_eager_producer() throws InterruptedException {
        String pubImageId = dockerClient.build(DockerProperties.formatImageName("pub"), "../roles/pub/Dockerfile");
        String subImageId = dockerClient.build(DockerProperties.formatImageName("sub"), "../roles/sub/Dockerfile");

        DockerImageRoleDefinition pubRoleDef = new DockerImageRoleDefinition("PUBROLE", pubImageId, Set.of(CAMERA));
        DockerImageRoleDefinition subRoleDef = new DockerImageRoleDefinition("SUBROLE", subImageId, Set.of(CPU),
                Set.of(new RoleDefinition.Metric("message_published", 1, ThresholdType.LESS_THAN, 5, ChronoUnit.SECONDS))
        );

        Device lazyProcessor = TestDeviceFactory.aLazyDevice(Set.of(pubRoleDef, subRoleDef), Set.of(CPU));
        lazyProcessor.start();
        await().untilTrue(lazyProcessor.getHasStarted());
        Assertions.assertTrue(lazyProcessor.getCurrentRoles().isEmpty());

        Device eagerCamera = TestDeviceFactory.aDevice(Set.of(pubRoleDef, subRoleDef), Set.of(CAMERA));
        eagerCamera.start();
        await().untilTrue(eagerCamera.getHasStarted());
        Assertions.assertFalse(eagerCamera.getCurrentRoles().isEmpty());

        lazyProcessor.tryRunRoles();

        await().until(() -> distributedColmenaPlatform.getStorage().getQueues().get("ack") != null);
        await().until(() -> distributedColmenaPlatform.getStorage().getQueues().get("ack").size() > 5);

        eagerCamera.disconnect();

        //wait until queue is empty and for processor role to be stopped
        await().until(() -> lazyProcessor.getCurrentRoles().isEmpty());
        lazyProcessor.disconnect();

        dockerClient.deleteImage(pubImageId);
        dockerClient.deleteImage(subImageId);
    }
}
