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
package es.bsc.colmena.agent;

import com.google.common.net.HostAndPort;
import es.bsc.colmena.app.routeguide.ColmenaPlatformGrpc;
import es.bsc.colmena.infrastructure.GrpcChannel;
import es.bsc.colmena.library.*;
import es.bsc.colmena.library.DCPClient;
import io.grpc.ConnectivityState;
import io.grpc.Grpc;
import io.grpc.InsecureChannelCredentials;
import io.grpc.ManagedChannel;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.extern.slf4j.Slf4j;

import java.net.DatagramSocket;
import java.util.*;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.TimeUnit;
import java.util.function.Function;
import java.util.stream.Collectors;

@Getter
@Slf4j
@NoArgsConstructor
public class ColmenaAgent {

    private static final Map<HostAndPort, DatagramSocket> testMap = new ConcurrentHashMap<>();
    private final Map<String, ServiceDescription> services = new ConcurrentHashMap<>();
    private DCPClient dcpClient;
    private ManagedChannel dcpChannel;
    private Device device;

    public ColmenaAgent(ServiceDescription serviceDescription) {
        this.services.put(serviceDescription.getId(), serviceDescription);
    }

    public void register(Device device, HostAndPort dcpHostPort) {
        this.device = device;
        dcpChannel = GrpcChannel.newChannel(dcpHostPort);
        dcpClient = new DCPClient(dcpChannel);
        ColmenaPlatformGrpc.ColmenaPlatformStub stub = ColmenaPlatformGrpc.newStub(dcpChannel);
        new ServiceDescriptionSubscriber(stub, this::addServiceDescription).subscribeToChanges();
        new ServiceDescriptionSubscriber(stub, this::addServiceDescription).getAllServices();
    }

    public Set<RoleDefinition> allRoleDefinitions() {
        return services.values().stream()
                .map(ServiceDescription::getRoles)
                .flatMap(Collection::stream)
                .collect(Collectors.toSet());
    }

    public void addServiceDescription(ServiceDescription serviceDescription) {
        services.put(serviceDescription.getId(), serviceDescription);
    }

    public void disconnect() throws InterruptedException {
        dcpChannel.shutdownNow().awaitTermination(30, TimeUnit.SECONDS);
    }

}
