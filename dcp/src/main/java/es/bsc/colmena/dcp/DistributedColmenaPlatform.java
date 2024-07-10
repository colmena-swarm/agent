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
package es.bsc.colmena.dcp;

import com.google.common.annotations.VisibleForTesting;
import com.google.protobuf.ByteString;
import com.google.protobuf.Empty;
import com.google.protobuf.Timestamp;
import es.bsc.colmena.dcp.grpc.GrpcServer;
import es.bsc.colmena.dcp.grpc.IpAddressInterceptor;
import es.bsc.colmena.dcp.metrics.MetricsStorage;
import es.bsc.colmena.dcp.metrics.Monitor;
import es.bsc.colmena.dcp.storage.ServiceStorage;
import es.bsc.colmena.dcp.storage.Storage;
import es.bsc.colmena.dcp.storage.Subscriber;
import es.bsc.colmena.dcp.storage.Subscribers;
import es.bsc.colmena.library.metrics.ThresholdType;
import es.bsc.colmena.library.proto.servicedescription.ServiceDescriptionConverter;
import es.bsc.colmena.app.routeguide.ColmenaPlatformGrpc;
import es.bsc.colmena.app.routeguide.RouteGuide;
import io.grpc.stub.StreamObserver;
import lombok.Getter;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.lang3.NotImplementedException;

import java.io.IOException;
import java.time.Instant;

@Slf4j
// each colmena agent shares an instance of this class so this can be thought of as shared data store
public class DistributedColmenaPlatform extends ColmenaPlatformGrpc.ColmenaPlatformImplBase implements Runnable {
    
    private final ServiceStorage serviceStorage = new ServiceStorage();
    private final int port;
    private final Storage storage;
    @Getter private final MetricsStorage metricsStorage;
    @Getter private final Monitor monitor;
    private GrpcServer server;

    public DistributedColmenaPlatform(int port) {
        this.port = port;
        metricsStorage = new MetricsStorage();
        monitor = new Monitor(metricsStorage);
        storage = new Storage(metricsStorage);
    }

    private void startServer() throws IOException {
        server = new GrpcServer(port, this);
        server.startServer();
    }

    public void stopServer() throws InterruptedException {
        storage.getQueueSubscribers().values().forEach(Subscribers::completeStreamObservers);
        server.stopServer();
    }

    public void blockUntilShutdown() throws InterruptedException {
        server.blockUntilShutdown();
    }

    public void run() {
        try {
            startServer();
            blockUntilShutdown();
        } catch (IOException | InterruptedException e) {
            throw new RuntimeException(e);
        }
        log.info("Server stopped, terminating thread");
    }

    @Override
    public void store(RouteGuide.StorageRequest storageRequest, StreamObserver<Empty> responseObserver) {
        storage.store(storageRequest.getKey(), storageRequest.getValue());
        responseObserver.onNext(Empty.newBuilder().build());
        responseObserver.onCompleted();
    }

    @Override
    public void storeMetrics(RouteGuide.MetricsStorageRequest metricsStorageRequest, StreamObserver<Empty> responseObserver) {
        Instant timestamp = metricsStorageRequest.hasTimestamp() ? convert(metricsStorageRequest.getTimestamp()) : Instant.now();
        MetricsStorage.Data toStore = new MetricsStorage.Data(
                metricsStorageRequest.getKey(),
                metricsStorageRequest.getValue(),
                timestamp);
        metricsStorage.store(toStore);
        responseObserver.onNext(Empty.newBuilder().build());
        responseObserver.onCompleted();
    }

    @Override
    public void queryMetrics(RouteGuide.MetricsQueryRequest metricsQueryRequest, StreamObserver<RouteGuide.MetricsQueryResponse> responseObserver) {
        Instant from = metricsQueryRequest.hasFrom() ? convert(metricsQueryRequest.getFrom()) : Instant.now();
        boolean met = monitor.metricIsMet(metricsQueryRequest.getKey(),
                metricsQueryRequest.getThreshold(),
                ThresholdType.valueOf(metricsQueryRequest.getThresholdType().name()),
                from);
        log.info("Metric: {}, met: {}", metricsQueryRequest.getKey(), met);
        responseObserver.onNext(RouteGuide.MetricsQueryResponse.newBuilder().setMet(met).build());
        responseObserver.onCompleted();
    }

    @Override
    public void getStored(RouteGuide.GetStoredRequest getStoredRequest, StreamObserver<RouteGuide.GetStoredResponse> responseObserver) {
        ByteString serializable = storage.get(getStoredRequest.getKey());
        responseObserver.onNext(RouteGuide.GetStoredResponse.newBuilder().setValue(serializable).build());
        responseObserver.onCompleted();
    }

    @Override
    public void publish(RouteGuide.StorageRequest storageRequest, StreamObserver<Empty> responseObserver) {
        storage.publish(storageRequest.getKey(), storageRequest.getValue());
        responseObserver.onNext(Empty.newBuilder().build());
        responseObserver.onCompleted();
    }

    @Override
    public void subscribe(RouteGuide.GetStoredRequest getStoredRequest, StreamObserver<RouteGuide.GetStoredResponse> responseObserver) {
        String ipAddress = (String) IpAddressInterceptor.IP_ADDRESS_KEY.get();
        storage.addSubscriber(getStoredRequest.getKey(), new Subscriber(responseObserver, ipAddress));
    }

    @Override
    public void getSubscriptionItem(RouteGuide.GetStoredRequest getStoredRequest, StreamObserver<RouteGuide.GetStoredResponse> responseObserver) {
        String ipAddress = (String) IpAddressInterceptor.IP_ADDRESS_KEY.get();
        storage.addSubscriber(getStoredRequest.getKey(), new Subscriber.SingleSubscriber(responseObserver, ipAddress));
    }

    @Override
    public void addService(RouteGuide.ServiceDescription serviceDescription, StreamObserver<Empty> responseObserver) {
        log.info("Added service: " + serviceDescription.getId().getValue());
        serviceStorage.addService(ServiceDescriptionConverter.convert(serviceDescription));
        responseObserver.onNext(Empty.newBuilder().build());
        responseObserver.onCompleted();
    }

    @Override
    public void getAllServices(Empty empty, StreamObserver<RouteGuide.ServiceDescription> responseObserver) {
        serviceStorage.publishAllServices(responseObserver);
        responseObserver.onCompleted();
    }

    @Override
    public void subscribeToServiceChanges(Empty empty, StreamObserver<RouteGuide.ServiceDescription> responseObserver) {
        log.info("New agent subscription");
        serviceStorage.addSubscriber(responseObserver);
    }

    @Override
    public void removeService(RouteGuide.ServiceDescriptionId serviceDescriptionId, StreamObserver<Empty> responseObserver) {
        throw new NotImplementedException();
    }

    @VisibleForTesting
    public Storage getStorage() {
        return storage;
    }

    @VisibleForTesting
    public ServiceStorage getServiceStorage() {
        return serviceStorage;
    }

    private Instant convert(Timestamp input) {
        return Instant.ofEpochSecond(input.getSeconds(), input.getNanos());
    }
}
