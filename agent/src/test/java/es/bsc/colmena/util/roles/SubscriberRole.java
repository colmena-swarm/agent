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
import es.bsc.colmena.library.communicationchannel.annotations.ColmenaMessageProcessor;
import lombok.Getter;
import lombok.Setter;
import lombok.extern.slf4j.Slf4j;

import java.util.concurrent.atomic.AtomicBoolean;

@Slf4j
public class SubscriberRole extends JavaRole {
    public static String ROLE_ID = "SUBSCRIBE_TO_STORAGE";
    @Getter private final AtomicBoolean hasProcessedData = new AtomicBoolean(false);
    @Setter private String[] expectedMessages;
    private int i = 0;

    public SubscriberRole(DCPClient dcpClient) {
        super(dcpClient);
    }

    @ColmenaMessageProcessor(key=PublisherRole.TEST_SUBSCRIPTION_KEY)
    public void processMessage(String message) {
        log.info(message);
        assert message.equals(expectedMessages[i++]);
        hasProcessedData.set(true);
    }

    @Override
    public String getRoleId() {
        return ROLE_ID;
    }
}
