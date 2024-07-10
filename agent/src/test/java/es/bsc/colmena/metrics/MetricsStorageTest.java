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
import es.bsc.colmena.dcp.metrics.PercentileOperation;
import org.junit.jupiter.api.Test;

import static org.hamcrest.CoreMatchers.equalTo;
import static org.hamcrest.MatcherAssert.assertThat;

class MetricsStorageTest extends AbstractTimingTest {

    @Test
    public void data_aggregated_and_filtered_by_time_then_summed() {
        storeValue(1);
        storeValue(2);
        storeValue(0);
        //should not be counted
        storeValue(10);
        storeValue(12);
        double sum = metricsStorage.get(KEY, minusSecs(5), MetricsStorage.SimpleOperation.SUM);
        assertThat(sum, equalTo(3.0));
    }

    @Test
    public void given_stored_metrics_then_average_can_be_calculated() {
        storeValue(5.0, 3);
        storeValue(3.0, 3);
        double avg = metricsStorage.get(KEY, minusSecs(4), MetricsStorage.SimpleOperation.AVG);
        assertThat(avg, equalTo(avg));
    }

    @Test
    public void given_stored_metrics_then_percentiles_can_be_calculated() {
        storeValue(3.0, 3);
        storeValue(6.0, 3);
        storeValue(7.0, 3);
        storeValue(8.0, 3);
        storeValue(8.0, 3);
        storeValue(9.0, 3);
        storeValue(10.0, 3);
        storeValue(13.0, 3);
        storeValue(15.0, 3);
        storeValue(16.0, 3);
        storeValue(20.0, 3);
        assertThat(percentile(25), equalTo(7.0));
        assertThat(percentile(50), equalTo(9.0));
        assertThat(percentile(75), equalTo(15.0));
        assertThat(percentile(100), equalTo(20.0));
    }

    private double percentile(double i) {
        return metricsStorage.get(KEY, minusSecs(4), PercentileOperation.of(i));
    }
}