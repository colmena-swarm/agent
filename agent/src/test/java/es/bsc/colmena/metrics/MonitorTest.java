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

import es.bsc.colmena.library.metrics.ThresholdType;
import es.bsc.colmena.dcp.metrics.Monitor;
import org.junit.jupiter.api.Test;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.is;

class MonitorTest extends AbstractTimingTest {

    Monitor monitor = new Monitor(metricsStorage);

    @Test
    public void given_a_metric_and_threshold_then_it_is_monitored_until_broken(){
        //at least one value stored in the last 10 seconds
        storeValue(9);
        assertThat(monitor.metricIsMet(KEY, 1.0, ThresholdType.GREATER_THAN_OR_EQUAL_TO, minusSecs(10)), is(true));
        assertThat(monitor.metricIsMet(KEY, 1.0, ThresholdType.GREATER_THAN_OR_EQUAL_TO, minusSecs(5)), is(false));
        assertThat(monitor.metricIsMet(KEY, 1.0, ThresholdType.LESS_THAN, minusSecs(10)), is(true));
        assertThat(monitor.metricIsMet(KEY, 2.0, ThresholdType.LESS_THAN, minusSecs(10)), is(true));
        assertThat(monitor.metricIsMet(KEY, 1.0, ThresholdType.LESS_THAN, minusSecs(5)), is(true));
        assertThat(monitor.metricIsMet(KEY, 0.0, ThresholdType.LESS_THAN, minusSecs(5)), is(true));
        assertThat(monitor.metricIsMet(KEY, 0.0, ThresholdType.LESS_THAN, minusSecs(10)), is(false));
    }

    @Test
    public void when_no_metric_has_been_recorded_then_monitor_behaves_as_expected(){
        //at least one value stored in the last 10 seconds
        assertThat(monitor.metricIsMet(KEY, 0.0, ThresholdType.GREATER_THAN_OR_EQUAL_TO, minusSecs(5)), is(true));
        assertThat(monitor.metricIsMet(KEY, 1.0, ThresholdType.GREATER_THAN_OR_EQUAL_TO, minusSecs(5)), is(false));
        assertThat(monitor.metricIsMet(KEY, 1.0, ThresholdType.LESS_THAN, minusSecs(5)), is(true));
        assertThat(monitor.metricIsMet(KEY, 0.0, ThresholdType.LESS_THAN, minusSecs(5)), is(true));
    }

}