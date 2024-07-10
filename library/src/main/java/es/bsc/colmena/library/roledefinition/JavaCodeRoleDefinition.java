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
package es.bsc.colmena.library.roledefinition;

import es.bsc.colmena.library.JavaRole;
import es.bsc.colmena.library.Requirement;
import es.bsc.colmena.library.Role;
import es.bsc.colmena.library.RoleDefinition;
import lombok.Getter;
import lombok.RequiredArgsConstructor;
import lombok.Value;

import java.time.temporal.ChronoUnit;
import java.util.Set;

@Getter
public class JavaCodeRoleDefinition extends RoleDefinition {
    private final Class<? extends JavaRole> role;

    public JavaCodeRoleDefinition(String roleId, Class<? extends JavaRole> role) {
        super(roleId, Set.of(), null);
        this.role = role;
    }

    public JavaCodeRoleDefinition(String roleId, Class<? extends JavaRole> role, Set<Requirement> requirements) {
        super(roleId, requirements, null);
        this.role = role;
    }

    public JavaCodeRoleDefinition(String roleId, Class<? extends JavaRole> role, Set<Requirement> requirements, Set<Metric> metrics) {
        super(roleId, requirements, metrics);
        this.role = role;
    }
}
