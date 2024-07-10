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
package es.bsc.colmena.docker;

import es.bsc.colmena.library.Role;
import lombok.Getter;
import lombok.RequiredArgsConstructor;
import lombok.Setter;

import java.util.Objects;

@RequiredArgsConstructor
public class DockerRole implements Role {

    private final String roleId;
    @Getter private final String imageId;
    @Setter @Getter private String containerId;

    @Override
    public String getRoleId() {
        return roleId;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.getRoleId());
    }

    @Override
    public boolean equals(Object other) {
        if (other == this) return true;
        if (!(other instanceof DockerRole otherCast)) return false;
        return this.getRoleId().equals(otherCast.getRoleId());
    }

    @Override
    public String toString() {
        return getRoleId();
    }
}
