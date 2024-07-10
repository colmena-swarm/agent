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
import es.bsc.colmena.library.DCPClient;
import es.bsc.colmena.library.roledefinition.DockerImageRoleDefinition;
import es.bsc.colmena.util.*;
import org.junit.jupiter.api.Test;

import java.util.Set;

import static org.awaitility.Awaitility.await;

public class ContainerLifecycleTest extends DockerColmenaTest {

    @Test
    public void when_a_docker_container_stops_then_it_is_removed_from_device_current_roles() throws InterruptedException {
        String imageName = DockerProperties.formatImageName("exits");
        String imageId = dockerClient.build(imageName, "../roles/exits/Dockerfile");

        DockerImageRoleDefinition roleDefinition = new DockerImageRoleDefinition("BasicRole", imageId);

        Device device = TestDeviceFactory.aDevice(Set.of(roleDefinition));
        device.start();

        await().untilTrue(device.getHasStarted());
        assert !device.getCurrentRoles().isEmpty();

        //wait until ByteString stored
        await().until(() -> distributedColmenaPlatform.getStorage().get(("exits")) != null);

        DCPClient dcpClient = DCPClientFactory.client();
        await().until(()-> String.valueOf(dcpClient.getStored("exits")).equals("hasStarted"));

        await().until(() -> device.getCurrentRoles().isEmpty());

        device.disconnect();
        device.blockUntilShutdown();
        dockerClient.deleteImage(imageName);
    }
}
