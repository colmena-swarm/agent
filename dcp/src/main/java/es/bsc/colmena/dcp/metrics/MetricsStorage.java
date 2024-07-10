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

import lombok.Value;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.collections4.Trie;
import org.apache.commons.collections4.trie.PatriciaTrie;

import java.time.Instant;
import java.util.Queue;
import java.util.concurrent.ConcurrentLinkedDeque;
import java.util.function.Function;
import java.util.stream.Collectors;
import java.util.stream.Stream;

@Slf4j
public class MetricsStorage {

    private final Trie<String, Queue<Data>> trie = new PatriciaTrie<>();

    public synchronized void store(Data data) {
        Queue<Data> stored = trie.get(data.key);
        if (stored == null) {
            stored = new ConcurrentLinkedDeque<>();
        }
        stored.add(data);
        trie.put(data.key, stored);
        log.info("MetricStored: {}, {}", data.key, data.value);
    }

    public synchronized double get(String key, Instant from, Operation op) {
        Stream<Data> filtered = trie.prefixMap(key).values()
                .stream()
                .flatMap(Queue::stream)
                .filter(each -> each.getTimestamp().isAfter(from));
        double value = op.apply(filtered);
        log.info("Metric queried: {}, {}", key, value);
        return value;
    }

    @Value
    public static class Data {
        String key;
        double value;
        Instant timestamp;
    }

    public interface Operation {
        double apply(Stream<Data> stream);
    }

    public enum SimpleOperation implements Operation {
        SUM(stream -> stream.collect(Collectors.summarizingDouble(Data::getValue)).getSum()),
        AVG(stream -> stream.collect(Collectors.summarizingDouble(Data::getValue)).getAverage());

        private final Function<Stream<Data>, Double> collector;

        SimpleOperation(Function<Stream<Data>, Double> collector) {
            this.collector = collector;
        }

        @Override
        public double apply(Stream<Data> stream) {
            return collector.apply(stream);
        }
    }
}
