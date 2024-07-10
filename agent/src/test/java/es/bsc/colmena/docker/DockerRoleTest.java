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
import es.bsc.colmena.app.routeguide.RouteGuide;
import es.bsc.colmena.library.DCPClient;
import es.bsc.colmena.library.roledefinition.DockerImageRoleDefinition;
import es.bsc.colmena.library.proto.servicedescription.ServiceDescriptionConverter;
import es.bsc.colmena.util.DCPClientFactory;
import es.bsc.colmena.util.DockerColmenaTest;
import es.bsc.colmena.util.DockerProperties;
import es.bsc.colmena.util.TestDeviceFactory;
import org.junit.jupiter.api.Test;

import java.util.Set;
import java.util.concurrent.TimeUnit;

import static es.bsc.colmena.library.requirement.HardwareRequirement.CAMERA;
import static es.bsc.colmena.library.requirement.HardwareRequirement.CPU;
import static org.awaitility.Awaitility.await;

public class DockerRoleTest extends DockerColmenaTest {

    @Test
    public void docker_image_role_can_use_dcp_storage() throws InterruptedException {
        String imageName = DockerProperties.formatImageName("use-storage");
        String imageId = dockerClient.build(imageName, "../roles/storage/Dockerfile");

        DockerImageRoleDefinition roleDefinition = new DockerImageRoleDefinition("BasicRole", imageId);

        Device device = TestDeviceFactory.aDevice(Set.of(roleDefinition));
        device.start();

        await().untilTrue(device.getHasStarted());
        assert !device.getCurrentRoles().isEmpty();

        //wait until ByteString stored
        await().until(() -> distributedColmenaPlatform.getStorage().get(("World")) != null);

        DCPClient dcpClient = DCPClientFactory.client();
        await().until(()-> String.valueOf(dcpClient.getStored("World")).equals("Hello"));

        device.disconnect();
        dockerClient.deleteImage(imageName);
    }

    @Test
    public void storage_docker_image_role_can_be_added_to_a_running_system() throws InterruptedException {
        String imageName = DockerProperties.formatImageName("use-storage");
        String imageId = dockerClient.build(imageName, "../roles/storage/Dockerfile");

        Device device = TestDeviceFactory.aDeviceWithRequirements(Set.of(CAMERA));
        device.start();

        await().untilTrue(device.getHasStarted());
        assert device.getCurrentRoles().isEmpty();

        DCPClient dcpClient = new DCPClient(channel);
        dcpClient.addService(ServiceDescriptionConverter.convert(basicServiceDescription(imageId)));

        //my home internet is very slow...
        await().atMost(5, TimeUnit.MINUTES).until(() -> distributedColmenaPlatform.getStorage().get(("World")) != null);
        await().until(()-> String.valueOf(dcpClient.getStored("World")).equals("Hello"));

        device.disconnect();
        dockerClient.deleteImage(imageName);
    }

    @Test
    public void docker_hub_image_can_be_started_on_a_running_system() throws InterruptedException {
        Device device = TestDeviceFactory.aDeviceWithRequirements(Set.of(CAMERA));
        device.start();

        await().untilTrue(device.getHasStarted());
        assert device.getCurrentRoles().isEmpty();

        DCPClient dcpClient = new DCPClient(channel);
        String imageName = DockerProperties.formatImageName("use-storage");
        dcpClient.addService(ServiceDescriptionConverter.convert(basicServiceDescription(imageName)));

        //wait until ByteString stored
        await().atMost(5, TimeUnit.MINUTES).until(() -> distributedColmenaPlatform.getStorage().get(("World")) != null);
        await().until(()-> String.valueOf(dcpClient.getStored("World")).equals("Hello"));

        device.disconnect();
        dockerClient.deleteImage(imageName);
    }

    @Test
    public void pub_sub_using_two_docker_images() throws InterruptedException {
        String pubImageId = dockerClient.build(DockerProperties.formatImageName("pub"), "../roles/pub/Dockerfile");
        String subImageId = dockerClient.build(DockerProperties.formatImageName("sub"), "../roles/sub/Dockerfile");

        DockerImageRoleDefinition pubRoleDef = new DockerImageRoleDefinition("PUBROLE", pubImageId, Set.of(CAMERA));
        DockerImageRoleDefinition subRoleDef = new DockerImageRoleDefinition("SUBROLE", subImageId, Set.of(CPU));

        Device processor = TestDeviceFactory.aDevice(Set.of(pubRoleDef, subRoleDef), Set.of(CPU));
        processor.start();
        Device camera = TestDeviceFactory.aDevice(Set.of(pubRoleDef, subRoleDef), Set.of(CAMERA));
        camera.start();

        await().untilTrue(camera.getHasStarted());
        await().untilTrue(processor.getHasStarted());
        assert !camera.getCurrentRoles().isEmpty();
        assert !processor.getCurrentRoles().isEmpty();

        await().until(() -> distributedColmenaPlatform.getStorage().getQueues().get("ack") != null);
        await().until(() -> distributedColmenaPlatform.getStorage().getQueues().get("ack").size() > 5);

        camera.disconnect();
        processor.disconnect();

        //dockerClient.deleteImage(pubImageId);
        //dockerClient.deleteImage(subImageId);
    }

    public static RouteGuide.ServiceDescription basicServiceDescription(String imageId) {
        var dockerRoleDefinition = RouteGuide.DockerRoleDefinition.newBuilder()
                .setId("TestRole")
                .setImageId(imageId)
                .addHardwareRequirements(RouteGuide.HardwareRequirement.CAMERA)
                .addAllKpis(Set.of())
                .build();

        return RouteGuide.ServiceDescription.newBuilder()
                .setId(RouteGuide.ServiceDescriptionId.newBuilder().setValue("Basic Service").build())
                .addDockerRoleDefinitions(dockerRoleDefinition)
                .build();
    }
}
