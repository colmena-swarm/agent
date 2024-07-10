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
package es.bsc.colmena.dcp;

import es.bsc.colmena.agent.Device;
import es.bsc.colmena.app.routeguide.RouteGuide;
import es.bsc.colmena.library.DCPClient;
import es.bsc.colmena.library.Role;
import es.bsc.colmena.library.RoleDefinition;
import es.bsc.colmena.library.ServiceDescription;
import es.bsc.colmena.library.proto.servicedescription.ServiceDescriptionConverter;
import es.bsc.colmena.library.requirement.HardwareRequirement;
import es.bsc.colmena.library.roledefinition.JavaCodeRoleDefinition;
import es.bsc.colmena.util.BaseColmenaTest;
import es.bsc.colmena.util.DockerProperties;
import es.bsc.colmena.util.TestDeviceFactory;
import es.bsc.colmena.util.roles.PublisherRole;
import es.bsc.colmena.util.roles.SubscriptionItemRole;
import org.awaitility.Awaitility;
import org.junit.jupiter.api.Test;

import java.io.Serializable;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.concurrent.TimeUnit;

import static es.bsc.colmena.app.routeguide.RouteGuide.HardwareRequirement.CAMERA;
import static es.bsc.colmena.app.routeguide.RouteGuide.HardwareRequirement.CPU;
import static es.bsc.colmena.util.roles.PublisherRole.TEST_SUBSCRIPTION_KEY;
import static java.time.Instant.EPOCH;
import static org.awaitility.Awaitility.await;
import static es.bsc.colmena.dcp.metrics.MetricsStorage.SimpleOperation.SUM;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.is;

public class QueueSizeMetricsTest extends BaseColmenaTest {

    private final String queueSizeMetric = TEST_SUBSCRIPTION_KEY + "_queue_size";

    @Test
    public void when_publishing_getting_subscription_messages_then_metrics_are_incremented_and_decremented() throws InterruptedException {
        RoleDefinition publisherRoleDef = new JavaCodeRoleDefinition(PublisherRole.ROLE_ID, PublisherRole.class);
        RoleDefinition subscriptionItemRoleDef = new JavaCodeRoleDefinition(SubscriptionItemRole.ROLE_ID, SubscriptionItemRole.class);

        Device device = TestDeviceFactory.aDevice(Set.of(publisherRoleDef, subscriptionItemRoleDef));
        device.start();

        Awaitility.await().until(() -> device.getCurrentRoles().size() == 2);
        Map<String, Role> currentRoles = device.getCurrentRoles();
        PublisherRole publisherRole = (PublisherRole) currentRoles.get(PublisherRole.ROLE_ID);
        SubscriptionItemRole subscriptionItemRole = (SubscriptionItemRole) currentRoles.get(SubscriptionItemRole.ROLE_ID);

        publisherRole.publish("Hello");
        publisherRole.publish("World");
        assertThat(distributedColmenaPlatform.getMetricsStorage().get(queueSizeMetric, EPOCH, SUM), is(2d));

        Serializable first = subscriptionItemRole.getSubscriptionItem();
        assertThat(first, is("Hello"));
        Serializable second = subscriptionItemRole.getSubscriptionItem();
        assertThat(second, is("World"));
        assertThat(distributedColmenaPlatform.getMetricsStorage().get(queueSizeMetric, EPOCH, SUM), is(0d));

        device.disconnect();
    }

    @Test
    public void role_can_be_started_with_queue_size_qos() throws InterruptedException {
        ServiceDescription serviceDescription = ServiceDescriptionConverter.convert(basicServiceDescription());
        Device processor = TestDeviceFactory.aLazyDevice(Set.of(), Set.of(HardwareRequirement.CPU));
        processor.start();
        Device camera = TestDeviceFactory.aDevice(Set.of(), Set.of(HardwareRequirement.CAMERA));
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
                .setValue("buffer_queue_size[100000000s] < 5")
                .build();
        var dockerRoleDefinitions = List.of(
                RouteGuide.DockerRoleDefinition.newBuilder()
                        .setId("Sensing")
                        .setImageId(DockerProperties.formatImageName("colmena-sensing"))
                        .addHardwareRequirements(CAMERA)
                        .build(),
                RouteGuide.DockerRoleDefinition.newBuilder()
                        .setId("Processing")
                        .setImageId(DockerProperties.formatImageName("colmena-processing"))
                        .addHardwareRequirements(CPU)
                        .addAllKpis(Set.of(kpi))
                        .build()
        );
        return RouteGuide.ServiceDescription.newBuilder()
                .setId(RouteGuide.ServiceDescriptionId.newBuilder().setValue("Basic Service").build())
                .addAllDockerRoleDefinitions(dockerRoleDefinitions)
                .build();
    }

}
