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
package es.bsc.colmena.kpi;

import es.bsc.colmena.library.RoleDefinition;
import es.bsc.colmena.library.proto.metric.KpiConversion;
import lombok.val;
import org.junit.jupiter.api.Test;

import static es.bsc.colmena.library.metrics.ThresholdType.GREATER_THAN_OR_EQUAL_TO;
import static java.time.temporal.ChronoUnit.HOURS;
import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.MatcherAssert.assertThat;

class KpiConversionTest {

    @Test
    public void test_conversion() {
        val metric = new RoleDefinition.Metric("processing_time", 60, GREATER_THAN_OR_EQUAL_TO, 1, HOURS);
        String converted = KpiConversion.convert(metric);
        assertThat(converted, is("processing_time[1h] >= 60.0"));
    }

    @Test
    public void test_regex() {
        String toParse = "processing_time[1h] >= 60.0 ";
        RoleDefinition.Metric parsed = KpiConversion.parse(toParse);
        assertThat(parsed.getKey(), is("processing_time"));
        assertThat(parsed.getUnit(), is(HOURS));
        assertThat(parsed.getThreshold(), is(60d));
        assertThat(parsed.getThresholdType(), is(GREATER_THAN_OR_EQUAL_TO));
        assertThat(parsed.getAmountOfTime(), is(1));
    }

}