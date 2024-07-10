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

import es.bsc.colmena.infrastructure.Host;
import es.bsc.colmena.library.Role;
import es.bsc.colmena.library.RoleDefinition;
import es.bsc.colmena.library.metrics.ThresholdType;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

import java.time.temporal.ChronoUnit;
import java.util.Map;
import java.util.Set;
import java.util.function.Function;
import java.util.stream.Collectors;

import static es.bsc.colmena.agent.Device.Strategy.EAGER;

@Slf4j
@RequiredArgsConstructor
public class RoleFinder {
    private final Device device;
    private final ColmenaAgent colmenaAgent;
    private final RoleRunner roleRunner;

    private void startRoles() {
        Set<RoleDefinition> all = colmenaAgent.allRoleDefinitions();
        Map<String, Role> currentRoles = roleRunner.getRunningRoles();
        all.stream()
                .filter(rd -> device.getFeatures().containsAll(rd.getRequirements()))
                .filter(rd -> deviceIsEager() || noMetrics(rd) || metricsBroken(rd))
                .filter(rd -> !currentRoles.containsKey(rd.getRoleId()))
                .map(roleDefinition -> RoleInstantiation.instantiate(roleDefinition, colmenaAgent))
                .forEach(roleRunner::startRole);
    }

    private void stopRoles() {
        if (deviceIsEager()) {
            return; //do not stop roles
        }

        Map<String, Role> currentRoles = roleRunner.getRunningRoles();
        Map<String, RoleDefinition> allRoleDefinitions = colmenaAgent.allRoleDefinitions().stream()
                .collect(Collectors.toMap(RoleDefinition::getRoleId, Function.identity()));
        currentRoles.values().stream()
                .map(role -> allRoleDefinitions.get(role.getRoleId()))
                .filter(rd -> hasMetrics(rd) && !metricsBroken(rd) && !hasWorkload(rd))
                .map(rd -> currentRoles.get(rd.getRoleId()))
                .forEach(roleRunner::stopRole);
    }

    private boolean hasWorkload(RoleDefinition roleDefinition) {
        String roleId = roleDefinition.getRoleId();
        String hostname = Host.hostname();
        String metricName = "num_executions_"+hostname+"_"+roleId;
        RoleDefinition.Metric query = new RoleDefinition.Metric(metricName, 1, ThresholdType.GREATER_THAN_OR_EQUAL_TO, 5, ChronoUnit.SECONDS);
        boolean hasWorkload = colmenaAgent.getDcpClient().metricMet(query);
        log.info("hasWorkload: {}", hasWorkload);
        return hasWorkload;
    }

    public synchronized void tryRoles() {
        startRoles();
        stopRoles();
    }

    private boolean deviceIsEager() {
        return EAGER.equals(device.getStrategy());
    }

    private boolean noMetrics(RoleDefinition roleDefinition) {
        return roleDefinition.getMetrics() == null || roleDefinition.getMetrics().isEmpty();
    }

    private boolean hasMetrics(RoleDefinition roleDefinition) {
        return !noMetrics(roleDefinition);
    }

    private boolean metricsBroken(RoleDefinition roleDefinition) {
        boolean metricsBroken = roleDefinition.getMetrics().stream()
                .anyMatch(each -> !colmenaAgent.getDcpClient().metricMet(each));
        log.info("metrics broken: {}", metricsBroken);
        return metricsBroken;
    }

}
