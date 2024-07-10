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

import es.bsc.colmena.library.Requirement;
import es.bsc.colmena.library.RoleDefinition;
import es.bsc.colmena.library.proto.metric.KpiConversion;
import es.bsc.colmena.library.roledefinition.DockerImageRoleDefinition;
import es.bsc.colmena.library.roledefinition.JavaCodeRoleDefinition;
import es.bsc.colmena.library.ServiceDescription;
import es.bsc.colmena.library.requirement.HardwareRequirement;
import es.bsc.colmena.app.routeguide.RouteGuide;
import lombok.extern.slf4j.Slf4j;

import java.util.Set;
import java.util.stream.Collectors;

@Slf4j
class ToGrpc {
    public static RouteGuide.ServiceDescription convert(ServiceDescription serviceDescription) {
        return RouteGuide.ServiceDescription.newBuilder()
                .setId(RouteGuide.ServiceDescriptionId.newBuilder().setValue(serviceDescription.getId()).build())
                .addAllRoleDefinitions(convertJavaCodeRoleDefinitions(serviceDescription))
                .addAllDockerRoleDefinitions(convertDockerRoleDefinitions(serviceDescription))
                .build();
    }

    private static Set<RouteGuide.DockerRoleDefinition> convertDockerRoleDefinitions(ServiceDescription serviceDescription) {
        return serviceDescription.getRoles().stream()
                .filter(roleDefinition -> roleDefinition instanceof DockerImageRoleDefinition)
                .map(roleDefinition -> (DockerImageRoleDefinition) roleDefinition)
                .map(ToGrpc::convertDockerRoleDefinition)
                .collect(Collectors.toSet());
    }

    private static Set<RouteGuide.RoleDefinition> convertJavaCodeRoleDefinitions(ServiceDescription serviceDescription) {
        return serviceDescription.getRoles().stream()
                .filter(roleDefinition -> roleDefinition instanceof JavaCodeRoleDefinition)
                .map(roleDefinition -> (JavaCodeRoleDefinition) roleDefinition)
                .map(ToGrpc::convertRoleDefinition)
                .collect(Collectors.toSet());
    }

    private static RouteGuide.DockerRoleDefinition convertDockerRoleDefinition(DockerImageRoleDefinition roleDefinition) {
        return RouteGuide.DockerRoleDefinition.newBuilder()
                .setId(roleDefinition.getRoleId())
                .setImageId(roleDefinition.getImageId())
                .addAllKpis(roleDefinition.getMetrics().stream().map(ToGrpc::convertMetric).collect(Collectors.toList()))
                .addAllHardwareRequirements(roleDefinition.getRequirements().stream().map(ToGrpc::convertHardwareRequirement).collect(Collectors.toList()))
                .build();
    }

    private static RouteGuide.RoleDefinition convertRoleDefinition(JavaCodeRoleDefinition roleDefinition) {
        return RouteGuide.RoleDefinition.newBuilder()
                .setId(roleDefinition.getRoleId())
                .setClassName(roleDefinition.getRole().getName())
                .addAllKpis(roleDefinition.getMetrics().stream().map(ToGrpc::convertMetric).collect(Collectors.toList()))
                .addAllHardwareRequirements(roleDefinition.getRequirements().stream().map(ToGrpc::convertHardwareRequirement).collect(Collectors.toList()))
                .build();
    }

    private static RouteGuide.HardwareRequirement convertHardwareRequirement(Requirement requirement) {
        try {
            HardwareRequirement hardwareRequirement = (HardwareRequirement) requirement;
            return RouteGuide.HardwareRequirement.valueOf(hardwareRequirement.name());
        }
        catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

    private static RouteGuide.Kpi convertMetric(RoleDefinition.Metric metric) {
        return RouteGuide.Kpi.newBuilder()
                .setValue(KpiConversion.convert(metric))
                .build();
    }
}
