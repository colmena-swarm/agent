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

import com.google.common.annotations.VisibleForTesting;
import com.google.protobuf.ByteString;
import es.bsc.colmena.app.routeguide.RouteGuide.GetStoredResponse;
import lombok.extern.slf4j.Slf4j;

import java.util.*;

@Slf4j
public class Subscribers {

    private final String key;
    private final Map<String, Subscriber> subscribers;
    private String lastPublishedIpAddress;

    public Subscribers(String key, Subscriber initial) {
        this.key = key;
        subscribers = new HashMap<>();
        subscribers.put(initial.getIpAddress(), initial);
    }

    public synchronized void add(Subscriber streamObserver) {
        subscribers.put(streamObserver.getIpAddress(), streamObserver);
    }

    public synchronized void remove(Subscriber subscriber) {
        subscribers.remove(subscriber.getIpAddress());
    }

    public synchronized int numberOfSubscribers() {
        return subscribers.size();
    }

    public synchronized void completeStreamObservers() {
        subscribers.values().forEach(Subscriber::onCompleted);
    }

    public synchronized void publish(ByteString byteString) {
        if (subscribers.isEmpty() || byteString == null) {
            return;
        }

        Subscriber subscriber = Optional.ofNullable(subscribers.get(lastPublishedIpAddress))
                .or(this::random).orElse(null);
        if (subscriber == null) {
            return;
        }
        try {
            GetStoredResponse message = GetStoredResponse.newBuilder().setValue(byteString).build();
            subscriber.onNext(message);
            lastPublishedIpAddress = subscriber.getIpAddress();
            if (subscriber instanceof Subscriber.SingleSubscriber) {
                remove(subscriber);
            }
            log.info("Message published. " +
                    "key: {}, subscriberCount: {}", key, numberOfSubscribers());
        }
        catch (Exception e) {
            log.warn("Caught exception while publishing", e);
            subscriber.onError(e);
            subscribers.remove(subscriber.getIpAddress());
        }
    }

    private Optional<Subscriber> random() {
        return subscribers.values().stream().findFirst();
    }

    @VisibleForTesting
    public List<Subscriber> get() {
        return subscribers.values().stream().toList();
    }
}
