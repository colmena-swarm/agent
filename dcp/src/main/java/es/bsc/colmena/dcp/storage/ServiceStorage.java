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

import com.google.common.collect.Sets;
import es.bsc.colmena.library.ServiceDescription;
import es.bsc.colmena.library.proto.servicedescription.ServiceDescriptionConverter;
import es.bsc.colmena.app.routeguide.RouteGuide;
import io.grpc.stub.StreamObserver;
import lombok.Getter;
import lombok.extern.slf4j.Slf4j;

import java.util.Map;
import java.util.Set;
import java.util.concurrent.ConcurrentHashMap;

@Slf4j
@Getter
public class ServiceStorage {

    private final Map<String, ServiceDescription> services = new ConcurrentHashMap<>();
    private final Set<StreamObserver<RouteGuide.ServiceDescription>> subscribers = Sets.newConcurrentHashSet();

    public void addSubscriber(StreamObserver<RouteGuide.ServiceDescription> subscriber) {
        subscribers.add(subscriber);
        log.info("New colmena agent. totalSubscribers: {}", subscribers.size());
    }

    public void addService(ServiceDescription serviceDescription) {
        services.put(serviceDescription.getId(), serviceDescription);
        subscribers.forEach(subscriber -> publishServiceDescription(subscriber, serviceDescription));
    }

    public void publishAllServices(StreamObserver<RouteGuide.ServiceDescription> subscriber) {
        services.values().forEach(service -> publishServiceDescription(subscriber, service));
    }

    private void publishServiceDescription(StreamObserver<RouteGuide.ServiceDescription> subscriber,
                                           ServiceDescription serviceDescription) {
        try {
            subscriber.onNext(ServiceDescriptionConverter.convert(serviceDescription));
        }
        catch (Exception e) {
            log.info("Caught exception while writing to {}. Dropping subscriber", subscriber, e);
            subscribers.remove(subscriber);
            subscriber.onError(e);
        }
    }
}
