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

import es.bsc.colmena.util.BaseColmenaTest;
import es.bsc.colmena.util.TestDeviceFactory;
import es.bsc.colmena.util.TestServiceDescriptionFactory;
import es.bsc.colmena.agent.Device;
import es.bsc.colmena.library.DCPClient;
import es.bsc.colmena.library.proto.servicedescription.ServiceDescriptionConverter;
import es.bsc.colmena.library.requirement.HardwareRequirement;
import org.junit.jupiter.api.Test;

import java.util.Set;

import static org.awaitility.Awaitility.await;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.is;

public class ServiceDescriptionTest extends BaseColmenaTest {

    /*
        add a service description and see the client take
        add a service description see the client take, disconnect and then add a second and see the subscription managed
     */

    private Device device;

    @Test
    public void when_device_starts_then_register_service_description_subscription() throws InterruptedException {
        device = TestDeviceFactory.aDeviceWithRequirements(Set.of(HardwareRequirement.CAMERA));
        device.start();
        await().untilTrue(device.getHasStarted());
        assertThat(distributedColmenaPlatform.getServiceStorage().getSubscribers().size(), is(1));
        device.disconnect();
    }

    @Test
    public void when_device_subscribed_then_service_can_be_added_and_added_to_device() throws InterruptedException {
        device = TestDeviceFactory.aDeviceWithRequirements(Set.of(HardwareRequirement.CAMERA));
        device.start();
        await().untilTrue(device.getHasStarted());
        assertThat(distributedColmenaPlatform.getServiceStorage().getSubscribers().size(), is(1));

        DCPClient dcpClient = new DCPClient(channel);
        dcpClient.addService(ServiceDescriptionConverter.convert(TestServiceDescriptionFactory.create()));

        assertThat(distributedColmenaPlatform.getServiceStorage().getServices().size(), is(1));
        await().until(() -> !device.getCurrentRoles().isEmpty());

        device.disconnect();
    }

    @Test
    public void given_an_added_service_description_when_an_agent_started_then_starts_with_role() throws InterruptedException {
        DCPClient dcpClient = new DCPClient(channel);
        dcpClient.addService(ServiceDescriptionConverter.convert(TestServiceDescriptionFactory.create()));

        device = TestDeviceFactory.aDeviceWithRequirements(Set.of(HardwareRequirement.CAMERA));
        device.start();
        await().untilTrue(device.getHasStarted());
        assertThat(distributedColmenaPlatform.getServiceStorage().getSubscribers().size(), is(1));
        assertThat(distributedColmenaPlatform.getServiceStorage().getServices().size(), is(1));
        await().until(() -> !device.getCurrentRoles().isEmpty());

        device.disconnect();
    }

    @Test
    public void when_subscribed_device_disconnects_then_subscription_disconnects_gracefully() throws InterruptedException {
        device = TestDeviceFactory.aDeviceWithRequirements(Set.of(HardwareRequirement.CAMERA));
        device.start();
        await().untilTrue(device.getHasStarted());
        assertThat(distributedColmenaPlatform.getServiceStorage().getSubscribers().size(), is(1));

        device.disconnect();

        //gRPC doesn't give access to the underlying cxn so at this point does not know that the connection has dropped
        assertThat(distributedColmenaPlatform.getServiceStorage().getSubscribers().size(), is(1));

        DCPClient dcpClient = new DCPClient(channel);
        dcpClient.addService(ServiceDescriptionConverter.convert(TestServiceDescriptionFactory.create()));

        assertThat(distributedColmenaPlatform.getServiceStorage().getSubscribers().size(), is(0));
        assertThat(distributedColmenaPlatform.getServiceStorage().getServices().size(), is(1));
    }
}
