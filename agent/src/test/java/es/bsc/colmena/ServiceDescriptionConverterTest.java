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
import es.bsc.colmena.library.metrics.ThresholdType;
import es.bsc.colmena.util.TestServiceDescriptionFactory;
import es.bsc.colmena.library.proto.servicedescription.ServiceDescriptionConverter;
import es.bsc.colmena.library.requirement.HardwareRequirement;
import es.bsc.colmena.app.routeguide.RouteGuide;
import es.bsc.colmena.library.roledefinition.JavaCodeRoleDefinition;
import es.bsc.colmena.library.ServiceDescription;
import lombok.extern.slf4j.Slf4j;
import org.hamcrest.Matchers;
import org.junit.jupiter.api.Test;

import java.time.temporal.ChronoUnit;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;

@Slf4j
public class ServiceDescriptionConverterTest {

    @Test
    public void when_agent_server_running_then_service_description_can_be_added() {
        RouteGuide.ServiceDescription serviceDescription = TestServiceDescriptionFactory.create();
        ServiceDescription capturedServiceDescription = ServiceDescriptionConverter.convert(serviceDescription);

        assertThat(capturedServiceDescription.getRoles().size(), equalTo(1));
        JavaCodeRoleDefinition capturedRoleDefinition = (JavaCodeRoleDefinition) capturedServiceDescription.getRoles().toArray()[0];
        assertThat(capturedRoleDefinition.getRoleId(), equalTo("TestRole"));
        assertThat(capturedRoleDefinition.getRequirements().size(), equalTo(1));
        assertThat(capturedRoleDefinition.getRequirements().contains(HardwareRequirement.CAMERA), Matchers.is(true));
        assertThat(capturedRoleDefinition.getMetrics().size(), equalTo(1));
        RoleDefinition.Metric capturedMetric = (RoleDefinition.Metric) capturedRoleDefinition.getMetrics().toArray()[0];
        assertThat(capturedMetric.getKey(), equalTo("test"));
        assertThat(capturedMetric.getUnit(), equalTo(ChronoUnit.HOURS));
        assertThat(capturedMetric.getThreshold(), equalTo(12.0));
        assertThat(capturedMetric.getThresholdType(), equalTo(ThresholdType.GREATER_THAN_OR_EQUAL_TO));
        assertThat(capturedMetric.getAmountOfTime(), equalTo(5));
    }
}
