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
package es.bsc.colmena.library;

import io.grpc.Grpc;
import io.grpc.InsecureChannelCredentials;
import lombok.SneakyThrows;

import java.util.Set;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;

public class ColmenaMain {
    private static final ExecutorService es = Executors.newCachedThreadPool();

    @SneakyThrows
    public static void start(Set<Class<? extends JavaRole>> roles) {
        String dcpHost = System.getenv("DCP_IP_ADDRESS");
        System.out.println(dcpHost);
        int dcpPort = 5555;
        var dcpChannel = Grpc.newChannelBuilderForAddress(dcpHost, dcpPort, InsecureChannelCredentials.create())
                .keepAliveTime(5L, TimeUnit.MINUTES)
                .keepAliveTimeout(1L, TimeUnit.SECONDS)
                .keepAliveWithoutCalls(true)
                .build();
        DCPClient dcpClient = new DCPClient(dcpChannel);
        roles.forEach(role -> instantiateClass(role, dcpClient));
        es.awaitTermination(1, TimeUnit.DAYS);
        dcpChannel.shutdownNow().awaitTermination(30, TimeUnit.SECONDS);
    }

    @SneakyThrows
    private static void instantiateClass(Class<? extends JavaRole> clazz, DCPClient dcpClient) {
        JavaRole javaRole = JavaRoleInstantiation.instantiate(clazz, dcpClient);
        es.submit(javaRole);
    }

}
