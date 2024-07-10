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
package es.bsc.colmena.metrics;

import es.bsc.colmena.dcp.metrics.MetricsStorage;

import java.time.Clock;
import java.time.Instant;
import java.time.ZoneId;
import java.time.temporal.ChronoUnit;

public abstract class AbstractTimingTest {

    protected static String KEY = "A.B.C";
    protected MetricsStorage metricsStorage = new MetricsStorage();
    protected Clock testClock = Clock.fixed(Instant.now(), ZoneId.systemDefault());

    Instant minusSecs(int noSecs) {
        return Instant.now(testClock).minus(noSecs, ChronoUnit.SECONDS);
    }

    void storeValue(double value, int minusSecs) {
        metricsStorage.store(new MetricsStorage.Data(KEY, value, minusSecs(minusSecs)));
    }

    void storeValue(int minusSecs) {
        storeValue(1.0, minusSecs);
    }

}
