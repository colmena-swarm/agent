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

import es.bsc.colmena.app.routeguide.RouteGuide;
import io.grpc.stub.StreamObserver;
import lombok.Getter;
import lombok.extern.slf4j.Slf4j;

@Slf4j
public class Subscriber implements StreamObserver<RouteGuide.GetStoredResponse> {

    @Getter
    protected final StreamObserver<RouteGuide.GetStoredResponse> streamObserver;
    @Getter
    protected final String ipAddress;

    public Subscriber(StreamObserver<RouteGuide.GetStoredResponse> streamObserver, String ipAddress) {
        this.streamObserver = streamObserver;
        this.ipAddress = ipAddress;
    }

    @Override
    public void onNext(RouteGuide.GetStoredResponse o) {
        streamObserver.onNext(o);
    }

    @Override
    public void onError(Throwable throwable) {
        streamObserver.onError(throwable);
    }

    @Override
    public void onCompleted() {
        try {
            streamObserver.onCompleted();
        } catch (Exception e) {
            log.warn("Could not complete subscriber", e);
        }
    }


    public static class SingleSubscriber extends Subscriber {

        public SingleSubscriber(StreamObserver<RouteGuide.GetStoredResponse> subscriber, String ipAddress) {
            super(subscriber, ipAddress);
        }

        @Override
        public void onNext(RouteGuide.GetStoredResponse o) {
            streamObserver.onNext(o);
            streamObserver.onCompleted();
        }
    }
}
