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
import es.bsc.colmena.library.Requirement;
import es.bsc.colmena.library.Role;
import lombok.Getter;
import lombok.extern.slf4j.Slf4j;

import java.util.Map;
import java.util.Set;
import java.util.concurrent.atomic.AtomicBoolean;

@Getter
@Slf4j
public class Device {
    private final RoleRunner roleRunner;
    private final RoleFinder roleFinder;
    private final ColmenaAgent colmenaAgent;
    private final HostAndPort dcpHostPort;
    private final Set<Requirement> features;
    private final AtomicBoolean hasStarted;
    private final TryRunRolesTimer tryRunRolesTimer;
    private final Strategy strategy;

    public Device(ColmenaAgent colmenaAgent, HostAndPort dcpHostPort, Set<Requirement> features, Strategy strategy) {
        this.colmenaAgent = colmenaAgent;
        this.dcpHostPort = dcpHostPort;
        this.features = features;
        this.strategy = strategy;
        hasStarted = new AtomicBoolean(false);
        tryRunRolesTimer = new TryRunRolesTimer(this);
        roleRunner = new RoleRunner(dcpHostPort);
        this.roleFinder = new RoleFinder(this, colmenaAgent, roleRunner);
    }

    public void start() {
        try {
            colmenaAgent.register(this, dcpHostPort);
            tryRunRoles();
            hasStarted.set(true);
            log.info("Device started. features: {}", features);
            tryRunRolesTimer.schedule();
        }
        catch (Exception e) {
            log.error("Could not start device", e);
        }
    }

    public synchronized void tryRunRoles() {
        roleFinder.tryRoles();
    }

    public void disconnect() throws InterruptedException {
        log.info("Starting device disconnect. features: {}", features);
        tryRunRolesTimer.disconnect();
        roleRunner.disconnect();
        colmenaAgent.disconnect();
        log.info("Device disconnected. features: {}", features);
    }

    public void blockUntilShutdown() throws InterruptedException{
        roleRunner.blockUntilShutdown();
    }

    public Map<String, Role> getCurrentRoles() {
        return roleRunner.getRunningRoles();
    }

    public enum Strategy {
        EAGER,
        LAZY
    }
}
