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
package es.bsc.colmena.util.roles;

import es.bsc.colmena.library.DCPClient;
import es.bsc.colmena.library.JavaRole;
import lombok.extern.slf4j.Slf4j;

@Slf4j
public class ExitsRole extends JavaRole {
    public static String ROLE_ID = "Exits";

    public ExitsRole(DCPClient dcpClient) {
        super(dcpClient);
    }

    @Override
    public void run() {
        dcpClient.store("ExitingJavaRole", "hasStarted");
        log.info("Exiting");
    }

    @Override
    public String getRoleId() {
        return ROLE_ID;
    }
}
