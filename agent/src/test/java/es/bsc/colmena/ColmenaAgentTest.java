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
import es.bsc.colmena.agent.Device;
import es.bsc.colmena.library.Requirement;
import es.bsc.colmena.library.roledefinition.JavaCodeRoleDefinition;
import es.bsc.colmena.util.roles.BasicRole;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.Arguments;
import org.junit.jupiter.params.provider.MethodSource;

import java.util.Set;
import java.util.stream.Stream;

import static es.bsc.colmena.library.requirement.HardwareRequirement.CAMERA;
import static org.awaitility.Awaitility.await;

public class ColmenaAgentTest extends BaseColmenaTest {

    private static Stream<Arguments> role_hardware_requirements_work_as_expected_args() {
        return Stream.of(
                Arguments.of(Set.of(),          Set.of(),           true),
                Arguments.of(Set.of(CAMERA),    Set.of(CAMERA),     true),
                Arguments.of(Set.of(CAMERA),    Set.of(),           false),
                Arguments.of(Set.of(),          Set.of(CAMERA),     true)
        );
    }

    @ParameterizedTest
    @MethodSource("role_hardware_requirements_work_as_expected_args")
    public void role_hardware_requirements_work_as_expected(Set<Requirement> roleRequirements,
                                                            Set<Requirement> deviceRequirements,
                                                            boolean expected) throws InterruptedException {
        JavaCodeRoleDefinition roleDefinition = new JavaCodeRoleDefinition("TestRole", BasicRole.class, roleRequirements);

        Device device = TestDeviceFactory.aDevice(Set.of(roleDefinition), deviceRequirements);
        device.start();

        await().untilTrue(device.getHasStarted());
        assert !device.getCurrentRoles().isEmpty() == expected;
        device.disconnect();
    }
}
