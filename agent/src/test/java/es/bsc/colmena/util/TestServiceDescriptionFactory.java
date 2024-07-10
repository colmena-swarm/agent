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

import es.bsc.colmena.app.routeguide.RouteGuide;
import es.bsc.colmena.util.roles.BasicRole;

public class TestServiceDescriptionFactory {

    public static RouteGuide.ServiceDescription create() {
        return  RouteGuide.ServiceDescription.newBuilder()
                .setId(RouteGuide.ServiceDescriptionId.newBuilder().setValue("Basic Role").build())
                .addRoleDefinitions(basicRoleDefinition())
                .build();
    }

    private static RouteGuide.RoleDefinition basicRoleDefinition() {
        return RouteGuide.RoleDefinition.newBuilder()
                .setId("TestRole")
                .setClassName(BasicRole.class.getName())
                .addHardwareRequirements(RouteGuide.HardwareRequirement.CAMERA)
                .addKpis(basicMetric())
                .build();
    }

    private static RouteGuide.Kpi basicMetric() {
        return RouteGuide.Kpi.newBuilder()
                .setValue("test[5h] >= 12.0 ")
                .build();
    }}
