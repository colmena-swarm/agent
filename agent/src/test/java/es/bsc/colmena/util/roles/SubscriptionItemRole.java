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
package es.bsc.colmena.util.roles;

import es.bsc.colmena.library.DCPClient;
import es.bsc.colmena.library.JavaRole;
import lombok.SneakyThrows;
import lombok.extern.slf4j.Slf4j;

import java.io.Serializable;
import java.util.ArrayList;
import java.util.List;

@Slf4j
public class SubscriptionItemRole extends JavaRole {

    public static final String ROLE_ID = "SUBSCRIPTION_ITEM_ROLE";
    public static final String TEST_SUBSCRIPTION_KEY = "Hello!";
    public final List<Serializable> messages = new ArrayList<>();

    public SubscriptionItemRole(DCPClient dcpClient) {
        super(dcpClient);
    }

    @SneakyThrows
    public Serializable getSubscriptionItem() {
            Serializable received = dcpClient.getSubscriptionItem(TEST_SUBSCRIPTION_KEY);
            messages.add(received);
            return received;
    }

    @Override
    public String getRoleId() {
        return ROLE_ID;
    }

}
