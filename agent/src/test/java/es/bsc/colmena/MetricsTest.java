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

import es.bsc.colmena.agent.Device;
import es.bsc.colmena.library.RoleDefinition;
import es.bsc.colmena.library.roledefinition.JavaCodeRoleDefinition;
import es.bsc.colmena.util.roles.ThroughputRole;
import es.bsc.colmena.util.BaseColmenaTest;
import es.bsc.colmena.util.TestDeviceFactory;
import org.junit.jupiter.api.Test;

import java.time.Instant;
import java.time.temporal.ChronoUnit;
import java.util.Set;
import java.util.concurrent.TimeUnit;

import static es.bsc.colmena.util.roles.ThroughputRole.METRICS_KEY;
import static org.awaitility.Awaitility.await;
import static es.bsc.colmena.library.metrics.ThresholdType.GREATER_THAN_OR_EQUAL_TO;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;

public class MetricsTest extends BaseColmenaTest {

    @Test
    public void when_role_increments_metrics_then_afterwards_retrieved_through_dcp() throws InterruptedException {
        Set<RoleDefinition> roleDefinitions = Set.of(
                new JavaCodeRoleDefinition(ThroughputRole.ROLE_ID, ThroughputRole.class)
        );
        Device device = TestDeviceFactory.aDevice(roleDefinitions);
        device.start();
        await().atMost(1, TimeUnit.SECONDS).untilTrue(device.getHasStarted());
        Instant beforeTestStarted = Instant.now().minus(100, ChronoUnit.SECONDS);
        await().atMost(2, TimeUnit.SECONDS).until(
                () -> distributedColmenaPlatform.getMonitor().metricIsMet(METRICS_KEY, 1, GREATER_THAN_OR_EQUAL_TO, beforeTestStarted));
        device.disconnect();
    }

    @Test
    public void when_eager_device_starts_and_metric_not_violated_then_start_role() throws InterruptedException {
        Set<RoleDefinition> roleDefinitions = Set.of(
                new JavaCodeRoleDefinition(ThroughputRole.ROLE_ID, ThroughputRole.class, Set.of(),
                        Set.of(new RoleDefinition.Metric(METRICS_KEY, 0, GREATER_THAN_OR_EQUAL_TO, 100, ChronoUnit.SECONDS)))
        );
        Device device = TestDeviceFactory.aDevice(roleDefinitions);
        device.start();
        await().atMost(1, TimeUnit.SECONDS).untilTrue(device.getHasStarted());
        assertThat(device.getCurrentRoles().size(), equalTo(1));
        device.disconnect();
    }

    @Test
    public void when_lazy_device_starts_and_metric_not_violated_then_do_not_start_role() throws InterruptedException {
        Set<RoleDefinition> roleDefinitions = Set.of(
                new JavaCodeRoleDefinition(ThroughputRole.ROLE_ID, ThroughputRole.class, Set.of(),
                        Set.of(new RoleDefinition.Metric(METRICS_KEY, 0, GREATER_THAN_OR_EQUAL_TO, 100, ChronoUnit.SECONDS)))
        );
        Device device = TestDeviceFactory.aLazyDevice(roleDefinitions);
        device.start();
        await().atMost(1, TimeUnit.SECONDS).untilTrue(device.getHasStarted());
        assertThat(device.getCurrentRoles().size(), equalTo(0));
        device.disconnect();
    }

    @Test
    public void when_lazy_device_starts_and_metric_violated_then_start_role() throws InterruptedException {
        Set<RoleDefinition> roleDefinitions = Set.of(
                new JavaCodeRoleDefinition(ThroughputRole.ROLE_ID, ThroughputRole.class, Set.of(),
                        Set.of(new RoleDefinition.Metric(METRICS_KEY, 100, GREATER_THAN_OR_EQUAL_TO, 100, ChronoUnit.SECONDS)))
        );
        Device device = TestDeviceFactory.aLazyDevice(roleDefinitions);
        device.start();
        await().atMost(1, TimeUnit.SECONDS).untilTrue(device.getHasStarted());
        assertThat(device.getCurrentRoles().size(), equalTo(1));
        device.disconnect();
    }
}
