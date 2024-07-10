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

import es.bsc.colmena.dcp.storage.Subscriber;
import es.bsc.colmena.library.DCPClient;
import io.grpc.Grpc;
import io.grpc.InsecureChannelCredentials;
import io.grpc.ManagedChannel;
import lombok.extern.slf4j.Slf4j;
import org.awaitility.Awaitility;
import org.hamcrest.Matchers;
import org.junit.jupiter.api.Test;

import static org.hamcrest.MatcherAssert.assertThat;

@Slf4j
public class GrpcServerTest {

    @Test
    public void server_can_process_an_incoming_message_and_then_respond() throws InterruptedException {
            DistributedColmenaPlatform dcp = new DistributedColmenaPlatform(5555);
            new Thread(dcp).start();
            ManagedChannel channel = Grpc.newChannelBuilderForAddress("127.0.0.1", 5555, InsecureChannelCredentials.create()).build();
            DCPClient dcpClient = new DCPClient(channel);

            dcpClient.store("Hello", "World");
            String response = (String) dcpClient.getStored("Hello");
            assert response.equals("World");

            dcp.stopServer();
            channel.shutdownNow();
    }

    @Test
    public void ipAddress_interceptor_parses_a_non_empty_ip_address_and_passes_to_subscriber() throws InterruptedException {
        DistributedColmenaPlatform dcp = new DistributedColmenaPlatform(5555);
        new Thread(dcp).start();
        ManagedChannel channel = Grpc.newChannelBuilderForAddress("127.0.0.1", 5555, InsecureChannelCredentials.create()).build();
        DCPClient dcpClient = new DCPClient(channel);

        new Thread(() -> dcpClient.getSubscriptionItem("Hello")).start();

        Awaitility.await().until(() -> dcp.getStorage().getQueueSubscribers().size() == 1);
        Subscriber createdSubscriber = dcp.getStorage().getQueueSubscribers().get("Hello").get().get(0);
        assertThat(createdSubscriber.getIpAddress(), Matchers.is("127.0.0.1"));

        dcp.stopServer();
        channel.shutdownNow();
    }
}
