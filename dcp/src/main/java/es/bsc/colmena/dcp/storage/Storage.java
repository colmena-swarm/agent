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
package es.bsc.colmena.dcp.storage;

import com.google.protobuf.ByteString;
import lombok.Getter;
import lombok.extern.slf4j.Slf4j;
import es.bsc.colmena.dcp.metrics.MetricsStorage;

import java.time.Instant;
import java.util.Map;
import java.util.Queue;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.LinkedBlockingQueue;

@Slf4j
public class Storage {
    private final Map<String, ByteString> storage = new ConcurrentHashMap<>();
    @Getter
    private final Map<String, Queue<ByteString>> queues = new ConcurrentHashMap<>();
    @Getter
    private final Map<String, Subscribers> queueSubscribers = new ConcurrentHashMap<>();
    private final MetricsStorage metricsStorage;

    public Storage(MetricsStorage metricsStorage) {
        this.metricsStorage = metricsStorage;
    }

    public void store(String key, ByteString value) {
        storage.put(key, value);
    }

    public ByteString get(String key) {
        return storage.get(key);
    }

    public synchronized void addSubscriber(String key, Subscriber subscriber) {
        Subscribers subscribers = queueSubscribers.get(key);
        if (subscribers == null) {
            subscribers = new Subscribers(key, subscriber);
            queueSubscribers.put(key, subscribers);
        } else {
            subscribers.add(subscriber);
        }
        log.info("Added subscriber. key: {}", key);

        if (queues.get(key) != null && !queues.get(key).isEmpty()) {
            ByteString polled = queues.get(key).poll();
            decrementMetrics(key);
            subscribers.publish(polled);
        }
    }

    public synchronized void publish(String key, ByteString value) {
        Subscribers subscribers = queueSubscribers.get(key);
        if (subscribers == null || subscribers.numberOfSubscribers() == 0) {
            enqueue(key, value);
        }
        else {
            subscribers.publish(value);
        }
    }

    private void enqueue(String key, ByteString value) {
        Queue<ByteString> byteStrings = queues.get(key);
        if (byteStrings == null) {
            byteStrings = new LinkedBlockingQueue<>();
            queues.put(key, byteStrings);
        }
        byteStrings.add(value);
        incrementMetrics(key);
        log.info("Message added to queue, no subscribers. key: {}", key);
    }

    private void incrementMetrics(String queue) {
        metricsStorage.store(new MetricsStorage.Data(queue + "_queue_size", 1, Instant.now()));
    }

    private void decrementMetrics(String queue) {
        metricsStorage.store(new MetricsStorage.Data(queue + "_queue_size", -1, Instant.now()));
    }
}
