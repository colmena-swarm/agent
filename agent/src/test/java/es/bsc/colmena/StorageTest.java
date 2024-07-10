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
import es.bsc.colmena.util.roles.UseStorageRole;
import es.bsc.colmena.util.BaseColmenaTest;
import es.bsc.colmena.util.TestDeviceFactory;
import lombok.extern.slf4j.Slf4j;
import org.junit.jupiter.api.Test;

import java.util.Set;

import static java.util.concurrent.TimeUnit.SECONDS;
import static org.awaitility.Awaitility.await;

@Slf4j
public class StorageTest extends BaseColmenaTest {

    @Test
    public void two_roles_can_share_data_stored_by_a_third_role() throws InterruptedException {
        Set<RoleDefinition> roleDefs = Set.of(new JavaCodeRoleDefinition(UseStorageRole.ROLE_ID, UseStorageRole.class));
        Device firstDevice = TestDeviceFactory.aDevice(roleDefs);
        Device secondDevice = TestDeviceFactory.aDevice(roleDefs);

        firstDevice.start();
        secondDevice.start();
        await().atMost(1, SECONDS).untilTrue(secondDevice.getHasStarted());
        // both devices should be able to use the storage
        assert firstDevice.getCurrentRoles().containsKey(UseStorageRole.ROLE_ID);
        assert secondDevice.getCurrentRoles().containsKey(UseStorageRole.ROLE_ID);

        ((UseStorageRole) firstDevice.getCurrentRoles().get(UseStorageRole.ROLE_ID)).store("Hello", "World");
        assert ((UseStorageRole) firstDevice.getCurrentRoles().get(UseStorageRole.ROLE_ID)).get("Hello").equals("World");
        assert ((UseStorageRole) secondDevice.getCurrentRoles().get(UseStorageRole.ROLE_ID)).get("Hello").equals("World");

        firstDevice.disconnect();
        secondDevice.disconnect();
    }



}
