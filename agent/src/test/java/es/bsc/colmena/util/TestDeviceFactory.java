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

import com.google.common.net.HostAndPort;
import es.bsc.colmena.agent.ColmenaAgent;
import es.bsc.colmena.agent.Device;
import es.bsc.colmena.library.Requirement;
import es.bsc.colmena.library.RoleDefinition;
import es.bsc.colmena.library.ServiceDescription;

import java.net.DatagramSocket;
import java.net.InetAddress;
import java.net.SocketException;
import java.net.UnknownHostException;
import java.util.Set;

public class TestDeviceFactory {

    // may not work on mac
    // https://stackoverflow.com/questions/9481865/getting-the-ip-address-of-the-current-machine-using-java
    private static String machineIpAddress() {
        try(final DatagramSocket socket = new DatagramSocket()){
            socket.connect(InetAddress.getByName("8.8.8.8"), 10002);
            return socket.getLocalAddress().getHostAddress();
        } catch (SocketException | UnknownHostException e) {
            throw new RuntimeException(e);
        }
    }

    public static final String SERVICE_DESCRIPTION_ID = "TEST_SERVICE_DESCRIPTION";

    public static Device aDevice(Set<RoleDefinition> roleDefinitions) {
        return new Device(
                new ColmenaAgent(ServiceDescription.builder().id(SERVICE_DESCRIPTION_ID).roles(roleDefinitions).build()),
                HostAndPort.fromParts(machineIpAddress(), 5555),
                Set.of(),
                Device.Strategy.EAGER
        );
    }

    public static Device aLazyDevice(Set<RoleDefinition> roleDefinitions) {
        return new Device(
                new ColmenaAgent(ServiceDescription.builder().id(SERVICE_DESCRIPTION_ID).roles(roleDefinitions).build()),
                HostAndPort.fromParts(machineIpAddress(), 5555),
                Set.of(),
                Device.Strategy.LAZY
        );
    }

    public static Device aDevice(Set<RoleDefinition> roleDefinitions, Set<Requirement> features) {
        return
                new Device(
                new ColmenaAgent(ServiceDescription.builder().id(SERVICE_DESCRIPTION_ID).roles(roleDefinitions).build()),
                HostAndPort.fromParts(machineIpAddress(), 5555),
                features, Device.Strategy.EAGER
        );
    }

    public static Device aLazyDevice(Set<RoleDefinition> roleDefinitions, Set<Requirement> features) {
        return
                new Device(
                        new ColmenaAgent(ServiceDescription.builder().id(SERVICE_DESCRIPTION_ID).roles(roleDefinitions).build()),
                        HostAndPort.fromParts(machineIpAddress(), 5555),
                        features, Device.Strategy.LAZY
                );
    }

    public static Device aDeviceWithRequirements(Set<Requirement> requirements) {
        return new Device(
                new ColmenaAgent(),
                HostAndPort.fromParts(machineIpAddress(), 5555),
                requirements, Device.Strategy.EAGER
        );
    }
}
