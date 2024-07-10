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
import es.bsc.colmena.library.metrics.Metrics;
import lombok.extern.slf4j.Slf4j;

import java.util.concurrent.atomic.AtomicBoolean;

@Slf4j
public class ThroughputRole extends JavaRole {

    public final static String ROLE_ID = "THROUGHPUT_ROLE";
    public final static String METRICS_KEY = "A.B.C";
    private final Metrics metrics;
    public AtomicBoolean incrementedMetrics = new AtomicBoolean(false);

    public ThroughputRole(DCPClient dcpClient) {
        super(dcpClient);
        metrics = new Metrics(dcpClient);
    }

    @Override
    public void run() {
        log.info("Started inc metrics role");
        running.set(true);
        try {
            while (running.get()) {
                metrics.increment(METRICS_KEY);
                incrementedMetrics.set(true);
                Thread.sleep(100);
            }
        } catch (InterruptedException e) {
            running.set(false);
            log.info("Interrupted. Stopping");
        } catch (Exception e) {
            running.set(false);
            log.error("Caught", e);
        }
    }

    @Override
    public String getRoleId() {
        return ROLE_ID;
    }
}
