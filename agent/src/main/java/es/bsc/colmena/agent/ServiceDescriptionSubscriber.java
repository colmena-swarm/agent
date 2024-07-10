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
package es.bsc.colmena.agent;

import com.google.protobuf.Empty;
import es.bsc.colmena.library.ServiceDescription;
import es.bsc.colmena.library.proto.servicedescription.ServiceDescriptionConverter;
import es.bsc.colmena.app.routeguide.ColmenaPlatformGrpc;
import es.bsc.colmena.app.routeguide.RouteGuide;
import io.grpc.Status;
import io.grpc.StatusRuntimeException;
import io.grpc.stub.StreamObserver;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

import java.util.function.Consumer;

@Slf4j
@RequiredArgsConstructor
public class ServiceDescriptionSubscriber implements StreamObserver<RouteGuide.ServiceDescription> {

    private final ColmenaPlatformGrpc.ColmenaPlatformStub stub;
    private final Consumer<ServiceDescription> listener;

    public void subscribeToChanges() {
        log.info("Subscribing to service description changes");
        stub.subscribeToServiceChanges(Empty.newBuilder().build(), this);
    }

    public void getAllServices() {
        log.info("Requesting all services from DCP");
        stub.getAllServices(Empty.newBuilder().build(), this);
    }

    @Override
    public void onNext(RouteGuide.ServiceDescription value) {
        log.info("New service description: {}", value.getId().getValue());
        listener.accept(ServiceDescriptionConverter.convert(value));
    }

    @Override
    public void onError(Throwable t) {
        if (t instanceof StatusRuntimeException statusRuntimeException) {
            if (statusRuntimeException.getStatus() == Status.UNAVAILABLE) {
                log.info("Stopping ServiceDescriptionSubscriber. gRPC connection closed.");
            }
        }
        else {
            log.error("Caught exception", t);
        }
    }

    @Override
    public void onCompleted() {
        log.info("Closed Service Description subscriber channel.");
    }
}
