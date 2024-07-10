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

import java.util.List;
import java.util.stream.Stream;

public class PercentileOperation implements MetricsStorage.Operation {

    private final double percentile;

    private PercentileOperation(double percentile) {
        this.percentile = percentile;
    }

    public static PercentileOperation of(double percentile) {
        return new PercentileOperation(percentile);
    }

    @Override
    public double apply(Stream<MetricsStorage.Data> stream) {
        List<Double> sorted = stream.map(MetricsStorage.Data::getValue).sorted().toList();
        int index = (int) Math.ceil(percentile / 100.0 * sorted.size());
        return sorted.get(index-1);
    }
}
