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
package es.bsc.colmena.docker.python;

import com.google.protobuf.util.JsonFormat;
import es.bsc.colmena.agent.Device;
import es.bsc.colmena.app.routeguide.RouteGuide;
import es.bsc.colmena.library.DCPClient;
import es.bsc.colmena.library.ServiceDescription;
import es.bsc.colmena.library.roledefinition.DockerImageRoleDefinition;
import es.bsc.colmena.library.proto.servicedescription.ServiceDescriptionConverter;
import es.bsc.colmena.util.DockerColmenaTest;
import es.bsc.colmena.util.DockerProperties;
import es.bsc.colmena.util.TestDeviceFactory;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.io.FileUtils;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;

import java.io.File;
import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.util.Set;
import java.util.concurrent.TimeUnit;

import static es.bsc.colmena.library.requirement.HardwareRequirement.CAMERA;
import static es.bsc.colmena.library.requirement.HardwareRequirement.CPU;
import static org.awaitility.Awaitility.await;

@Slf4j
public class PythonDockerRoleTest extends DockerColmenaTest {

    private static final String EXAMPLE_APPLICATION = "../../service-tools/test/examples/example_application/build/service_description.json";
    private static final String EXAMPLE_INTERFACES = "../../service-tools/test/examples/example_interfaces/build/service_description.json";

    @Test
    // this is used to test the service description built by the service tools
    public void service_tools_service_description_can_be_started_on_a_running_system() throws InterruptedException, IOException {
        ServiceDescription serviceDescription = ServiceDescriptionConverter.convert(serviceDescriptionFromJson(EXAMPLE_APPLICATION));
        Device processor = TestDeviceFactory.aDevice(Set.of(), Set.of(CPU));
        processor.start();
        Device camera = TestDeviceFactory.aDevice(Set.of(), Set.of(CAMERA));
        camera.start();

        await().untilTrue(processor.getHasStarted());
        await().untilTrue(camera.getHasStarted());

        DCPClient dcpClient = new DCPClient(channel);
        dcpClient.addService(serviceDescription);

        await().atMost(5, TimeUnit.MINUTES).until(() -> distributedColmenaPlatform.getStorage().getQueues().get("result") != null);
        await().atMost(5, TimeUnit.MINUTES).until(() -> distributedColmenaPlatform.getStorage().getQueues().get("result").size() > 5);

        camera.disconnect();
        processor.disconnect();

    }

    @Test
    public void pub_sub_using_two_python_docker_images() throws InterruptedException {
        String pubImageId = dockerClient.build(DockerProperties.formatImageName("colmena-sensing"), "../../service-tools/test/examples/example_application/build/Sensing/Dockerfile");
        String subImageId = dockerClient.build(DockerProperties.formatImageName("colmena-processing"), "../../service-tools/test/examples/example_application/build/Processing/Dockerfile");

        DockerImageRoleDefinition pubRoleDef = new DockerImageRoleDefinition("PUBROLE", pubImageId, Set.of(CAMERA));
        DockerImageRoleDefinition subRoleDef = new DockerImageRoleDefinition("SUBROLE", subImageId, Set.of(CPU));

        Device processor = TestDeviceFactory.aDevice(Set.of(pubRoleDef, subRoleDef), Set.of(CPU));
        processor.start();
        Device camera = TestDeviceFactory.aDevice(Set.of(pubRoleDef, subRoleDef), Set.of(CAMERA));
        camera.start();

        await().atMost(5, TimeUnit.MINUTES).until(() -> distributedColmenaPlatform.getStorage().getQueues().get("result") != null);
        await().atMost(5, TimeUnit.MINUTES).until(() -> distributedColmenaPlatform.getStorage().getQueues().get("result").size() > 5);

        camera.disconnect();
        processor.disconnect();

        dockerClient.deleteImage(pubImageId);
        dockerClient.deleteImage(subImageId);
    }

    @Test
    public void pub_sub_using_kpi_parsing() throws InterruptedException, IOException {
        ServiceDescription serviceDescription = ServiceDescriptionConverter.convert(serviceDescriptionFromJson(EXAMPLE_APPLICATION));
        DCPClient dcpClient = new DCPClient(channel);
        dcpClient.addService(serviceDescription);

        Device processor = TestDeviceFactory.aLazyDevice(Set.of(), Set.of(CPU));
        processor.start();
        await().untilTrue(processor.getHasStarted());
        Assertions.assertTrue(processor.getCurrentRoles().isEmpty());

        Device camera = TestDeviceFactory.aDevice(Set.of(), Set.of(CAMERA));
        camera.start();

        await().atMost(5, TimeUnit.MINUTES).until(() -> distributedColmenaPlatform.getStorage().getQueues().get("result") != null);
        await().atMost(5, TimeUnit.MINUTES).until(() -> distributedColmenaPlatform.getStorage().getQueues().get("result").size() > 5);
        log.info("test completed. Disconnecting devices");

        camera.disconnect();

        await().atMost(30, TimeUnit.SECONDS).until(() -> processor.getCurrentRoles().isEmpty());
        processor.disconnect();
    }

    @Test
    public void test_interfaces() throws InterruptedException, IOException {
        ServiceDescription serviceDescription = ServiceDescriptionConverter.convert(serviceDescriptionFromJson(EXAMPLE_INTERFACES));
        Device processor = TestDeviceFactory.aDevice(Set.of(), Set.of(CPU));
        processor.start();
        Device camera = TestDeviceFactory.aDevice(Set.of(), Set.of(CAMERA));
        camera.start();

        await().untilTrue(processor.getHasStarted());
        await().untilTrue(camera.getHasStarted());

        DCPClient dcpClient = new DCPClient(channel);
        dcpClient.addService(serviceDescription);

        await().atMost(5, TimeUnit.MINUTES).until(() -> distributedColmenaPlatform.getStorage().get("evaluate") != null);

        camera.disconnect();
        processor.disconnect();
    }


    public static RouteGuide.ServiceDescription serviceDescriptionFromJson(String serviceDescriptionFileLocation) throws IOException {
        String content = FileUtils.readFileToString(new File(serviceDescriptionFileLocation), StandardCharsets.UTF_8);
        RouteGuide.ServiceDescription.Builder sdBuilder = RouteGuide.ServiceDescription.newBuilder();
        JsonFormat.parser().ignoringUnknownFields().merge(content, sdBuilder);
        return sdBuilder.build();
    }
}
