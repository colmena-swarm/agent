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

import es.bsc.colmena.agent.Device;
import es.bsc.colmena.app.routeguide.RouteGuide;
import es.bsc.colmena.library.DCPClient;
import es.bsc.colmena.library.ServiceDescription;
import es.bsc.colmena.library.proto.servicedescription.ServiceDescriptionConverter;
import es.bsc.colmena.util.DockerColmenaTest;
import es.bsc.colmena.util.DockerProperties;
import es.bsc.colmena.util.TestDeviceFactory;
import org.junit.jupiter.api.Test;

import java.util.List;
import java.util.Set;
import java.util.concurrent.TimeUnit;

import static es.bsc.colmena.library.requirement.HardwareRequirement.CAMERA;
import static es.bsc.colmena.library.requirement.HardwareRequirement.CPU;
import static org.awaitility.Awaitility.await;

public class NumExecutionsTest extends DockerColmenaTest {

    @Test
    public void given_a_role_with_metrics_that_are_met_with_the_role_doing_work_then_role_is_not_stopped() throws InterruptedException {
        ServiceDescription serviceDescription = ServiceDescriptionConverter.convert(basicServiceDescription());
        Device processor = TestDeviceFactory.aLazyDevice(Set.of(), Set.of(CPU));
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

    public static RouteGuide.ServiceDescription basicServiceDescription() {
        var kpi = RouteGuide.Kpi.newBuilder()
                .setValue("processed[3s] >= 1")
                .build();
        var dockerRoleDefinitions = List.of(
                    RouteGuide.DockerRoleDefinition.newBuilder()
                    .setId("Sensing")
                    .setImageId(DockerProperties.formatImageName("colmena-sensing"))
                    .addHardwareRequirements(RouteGuide.HardwareRequirement.CAMERA)
                    .build(),
                    RouteGuide.DockerRoleDefinition.newBuilder()
                    .setId("Processing")
                    .setImageId(DockerProperties.formatImageName("colmena-processing"))
                    .addHardwareRequirements(RouteGuide.HardwareRequirement.CPU)
                    .addAllKpis(Set.of(kpi))
                    .build()
                );
        return RouteGuide.ServiceDescription.newBuilder()
                .setId(RouteGuide.ServiceDescriptionId.newBuilder().setValue("Basic Service").build())
                .addAllDockerRoleDefinitions(dockerRoleDefinitions)
                .build();
    }
}
