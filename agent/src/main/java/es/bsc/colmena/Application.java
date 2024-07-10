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

import com.google.common.net.HostAndPort;
import es.bsc.colmena.agent.ColmenaAgent;
import es.bsc.colmena.agent.Device;
import es.bsc.colmena.library.Requirement;
import es.bsc.colmena.library.requirement.HardwareRequirement;
import lombok.SneakyThrows;
import lombok.extern.slf4j.Slf4j;

import java.time.Instant;
import java.util.Set;

@Slf4j
public class Application {

    @SneakyThrows
    public static void main(String[] args) {
        String colmenaHost = args[0];
        int colmenaPort = Integer.parseInt(args[1]);
        Device.Strategy deviceStrategy = Device.Strategy.valueOf(args[2]);
        Requirement requirement = parseDeviceFeature(args);
        Device device = new Device(
                new ColmenaAgent(),
                HostAndPort.fromParts(colmenaHost, colmenaPort),
                requirement == null ? Set.of() : Set.of(requirement),
                deviceStrategy
        );
        Thread printingHook = new Thread(() -> {
            try {
                log.info("shutting down device");
                device.disconnect();
            } catch (InterruptedException e) {
                System.err.println("could not shut down device");
                throw new RuntimeException(e);
            }
        });
        Runtime.getRuntime().addShutdownHook(printingHook);
        device.start();
        device.blockUntilShutdown();
    }

    private static Requirement parseDeviceFeature(String[] args) {
        if (args.length == 3) {
            return null;
        }
        return HardwareRequirement.valueOf(args[3]);
    }

}
