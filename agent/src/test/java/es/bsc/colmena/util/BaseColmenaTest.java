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
package es.bsc.colmena.util;

import io.grpc.Grpc;
import io.grpc.InsecureChannelCredentials;
import io.grpc.ManagedChannel;
import lombok.extern.slf4j.Slf4j;
import es.bsc.colmena.dcp.DistributedColmenaPlatform;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;

@Slf4j
public abstract class BaseColmenaTest {

    protected DistributedColmenaPlatform distributedColmenaPlatform;
    protected ManagedChannel channel;

    @BeforeEach
    public void init() {
        distributedColmenaPlatform = new DistributedColmenaPlatform(5555);
        new Thread(distributedColmenaPlatform).start();
        channel = Grpc.newChannelBuilderForAddress("127.0.0.1", 5555, InsecureChannelCredentials.create()).build();
    }

    @AfterEach
    public void teardown() throws InterruptedException {
        log.info("stopping server");
        distributedColmenaPlatform.stopServer();
    }
}
