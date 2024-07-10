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
package es.bsc.colmena;

import es.bsc.colmena.library.RoleDefinition;
import es.bsc.colmena.util.BaseColmenaTest;
import es.bsc.colmena.util.TestDeviceFactory;
import es.bsc.colmena.agent.Device;
import es.bsc.colmena.util.roles.PublisherRole;
import es.bsc.colmena.util.roles.SubscriberRole;
import es.bsc.colmena.library.roledefinition.JavaCodeRoleDefinition;
import lombok.extern.slf4j.Slf4j;
import es.bsc.colmena.dcp.storage.Subscribers;
import org.junit.jupiter.api.Test;

import java.util.Set;

import static es.bsc.colmena.library.requirement.HardwareRequirement.CAMERA;
import static es.bsc.colmena.library.requirement.HardwareRequirement.CPU;
import static java.util.concurrent.TimeUnit.SECONDS;
import static org.awaitility.Awaitility.await;

@Slf4j
public class PublishSubscribeTest extends BaseColmenaTest {
    private static final String[] TEST_MESSAGES = new String[]{"One", "Two", "Three"};
    private static final Set<RoleDefinition> CAMERA_CPU_ROLE_DEFINITIONS = Set.of(
            new JavaCodeRoleDefinition(PublisherRole.ROLE_ID, PublisherRole.class, Set.of(CAMERA)),
            new JavaCodeRoleDefinition(SubscriberRole.ROLE_ID, SubscriberRole.class, Set.of(CPU)));

    @Test
    public void given_one_device_when_messages_published_to_storage_brokered_pub_sub_then_messages_are_processed_by_subscription() throws InterruptedException {
        Device device = TestDeviceFactory.aDevice(Set.of(
                new JavaCodeRoleDefinition(PublisherRole.ROLE_ID, PublisherRole.class),
                new JavaCodeRoleDefinition(SubscriberRole.ROLE_ID, SubscriberRole.class)));

        device.start();
        await().atMost(1, SECONDS).untilTrue(device.getHasStarted());

        SubscriberRole subscriberRole = getSubscriptionRole(device);
        await().until(this::subscribedAsExpected);

        subscriberRole.setExpectedMessages(TEST_MESSAGES);
        for (String testMessage : TEST_MESSAGES) {
            getPublisherRole(device).publish(testMessage);
        }

        await().atMost(3, SECONDS).untilTrue(subscriberRole.getHasProcessedData());
        device.disconnect();
    }

    @Test
    public void given_producer_and_subscriber_devices_when_publisher_starts_first_then_subscription_can_process_messages() throws InterruptedException {
        Device camera = TestDeviceFactory.aDevice(CAMERA_CPU_ROLE_DEFINITIONS, Set.of(CAMERA));
        Device cpu = TestDeviceFactory.aDevice(CAMERA_CPU_ROLE_DEFINITIONS, Set.of(CPU));

        camera.start();
        await().atMost(1, SECONDS).untilTrue(camera.getHasStarted());

        cpu.start();
        await().atMost(1, SECONDS).untilTrue(cpu.getHasStarted());

        SubscriberRole subscriberRole = getSubscriptionRole(cpu);
        await().until(this::subscribedAsExpected);

        subscriberRole.setExpectedMessages(TEST_MESSAGES);
        for (String testMessage : TEST_MESSAGES) {
            getPublisherRole(camera).publish(testMessage);
        }

        await().atMost(3, SECONDS).untilTrue(subscriberRole.getHasProcessedData());
        cpu.disconnect();
        camera.disconnect();
    }

    @Test
    public void given_producer_and_subscriber_devices_when_subscriber_starts_first_then_subscription_can_process_messages() throws InterruptedException {
        Device camera = TestDeviceFactory.aDevice(CAMERA_CPU_ROLE_DEFINITIONS, Set.of(CAMERA));
        Device cpu = TestDeviceFactory.aDevice(CAMERA_CPU_ROLE_DEFINITIONS, Set.of(CPU));

        cpu.start();
        await().atMost(1, SECONDS).untilTrue(cpu.getHasStarted());

        camera.start();
        await().atMost(1, SECONDS).untilTrue(camera.getHasStarted());

        SubscriberRole subscriberRole = getSubscriptionRole(cpu);
        await().until(this::subscribedAsExpected);

        subscriberRole.setExpectedMessages(TEST_MESSAGES);
        for (String testMessage : TEST_MESSAGES) {
            getPublisherRole(camera).publish(testMessage);
        }

        await().atMost(3, SECONDS).untilTrue(subscriberRole.getHasProcessedData());
        cpu.disconnect();
        camera.disconnect();
    }

    private SubscriberRole getSubscriptionRole(Device device) {
        return (SubscriberRole) device.getCurrentRoles().get(SubscriberRole.ROLE_ID);
    }

    private PublisherRole getPublisherRole(Device device) {
        return (PublisherRole) device.getCurrentRoles().get(PublisherRole.ROLE_ID);
    }

    private boolean subscribedAsExpected() {
        //once subscriber device has started then check DCP storage to see subscriber role has registered correctly
        Subscribers subscribers = distributedColmenaPlatform.getStorage().getQueueSubscribers().get(PublisherRole.TEST_SUBSCRIPTION_KEY);
        return subscribers.numberOfSubscribers() > 0;
    }

}
