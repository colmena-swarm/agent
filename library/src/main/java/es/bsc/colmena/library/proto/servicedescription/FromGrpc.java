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
package es.bsc.colmena.library.proto.servicedescription;

import com.google.common.collect.Sets;
import es.bsc.colmena.library.*;
import es.bsc.colmena.library.metrics.ThresholdType;
import es.bsc.colmena.library.proto.metric.KpiConversion;
import es.bsc.colmena.library.roledefinition.DockerImageRoleDefinition;
import es.bsc.colmena.library.roledefinition.JavaCodeRoleDefinition;
import es.bsc.colmena.library.requirement.HardwareRequirement;
import es.bsc.colmena.app.routeguide.RouteGuide;
import lombok.extern.slf4j.Slf4j;

import java.time.temporal.ChronoUnit;
import java.util.List;
import java.util.Set;
import java.util.stream.Collectors;

@Slf4j
class FromGrpc {
    static ServiceDescription convert(RouteGuide.ServiceDescription serviceDescription) {
        return ServiceDescription.builder()
                .id(serviceDescription.getId().getValue())
                .roles(parseRoles(serviceDescription))
                .build();
    }

    private static Set<RoleDefinition> parseRoles(RouteGuide.ServiceDescription serviceDescription) {
        Set<RoleDefinition> javaRoleDefinitions = serviceDescription
                .getRoleDefinitionsList().stream()
                .map(FromGrpc::parseRoleDefinition)
                .collect(Collectors.toSet());

        Set<RoleDefinition> dockerImageRoleDefinitions = serviceDescription
                .getDockerRoleDefinitionsList().stream()
                .map(FromGrpc::parseDockerRoleDefinition)
                .collect(Collectors.toSet());

        return Sets.union(javaRoleDefinitions, dockerImageRoleDefinitions);
    }

    private static DockerImageRoleDefinition parseDockerRoleDefinition(RouteGuide.DockerRoleDefinition dockerRoleDefinition) {
        return new DockerImageRoleDefinition(
                dockerRoleDefinition.getId(),
                dockerRoleDefinition.getImageId(),
                parseHardwareRequirements(dockerRoleDefinition.getHardwareRequirementsList()),
                parseMetrics(dockerRoleDefinition.getKpisList())
        );
    }

    private static JavaCodeRoleDefinition parseRoleDefinition(RouteGuide.RoleDefinition roleDefinition) {
        return new JavaCodeRoleDefinition(
                roleDefinition.getId(),
                parseClass(roleDefinition.getClassName()),
                parseHardwareRequirements(roleDefinition.getHardwareRequirementsList()),
                parseMetrics(roleDefinition.getKpisList()));
    }

    private static Set<RoleDefinition.Metric> parseMetrics(List<RouteGuide.Kpi> kpis) {
        return kpis.stream()
                .map(RouteGuide.Kpi::getValue)
                .map(KpiConversion::parse)
                .collect(Collectors.toSet());
    }

    //must be fully qualified
    private static Class<? extends JavaRole> parseClass(String className) {
        try {
            return (Class<? extends JavaRole>) Class.forName(className);
        } catch (ClassNotFoundException e) {
            log.error("Could not create class of {}, is the parameter is a fully qualified class name?", className, e);
            throw new RuntimeException(e);
        }
    }

    private static Set<Requirement> parseHardwareRequirements(List<RouteGuide.HardwareRequirement> hardwareRequirements) {
        return hardwareRequirements.stream()
                .map(hr -> HardwareRequirement.valueOf(hr.name()))
                .collect(Collectors.toSet());
    }
}
