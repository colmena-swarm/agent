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
package es.bsc.colmena.dcp.metrics;

import es.bsc.colmena.library.metrics.ThresholdType;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

import java.time.Instant;

@Slf4j
@RequiredArgsConstructor
public class Monitor {
    private final MetricsStorage metricsStorage;

    public boolean metricIsMet(String key, double threshold, ThresholdType thresholdType, Instant from) {
        switch (thresholdType) {
            case GREATER_THAN_OR_EQUAL_TO -> {
                return metricsStorage.get(key, from, MetricsStorage.SimpleOperation.SUM) >= threshold;
            }
            case LESS_THAN -> {
                return metricsStorage.get(key, from, MetricsStorage.SimpleOperation.SUM) <= threshold;
            }
            default -> {
                log.error("No check configured for thresholdType: {}", thresholdType);
                throw new RuntimeException();
            }
        }
    }

}
