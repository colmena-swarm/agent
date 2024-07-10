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
package es.bsc.colmena.library;

import com.google.protobuf.ByteString;
import com.google.protobuf.Timestamp;
import es.bsc.colmena.app.routeguide.ColmenaPlatformGrpc;
import es.bsc.colmena.app.routeguide.RouteGuide;
import es.bsc.colmena.library.communicationchannel.CommunicationChannelListener;
import es.bsc.colmena.library.communicationchannel.storage.StorageQueueCommunicationChannelListener;
import es.bsc.colmena.library.proto.servicedescription.ServiceDescriptionConverter;
import io.grpc.ManagedChannel;
import io.grpc.StatusRuntimeException;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.lang3.SerializationUtils;

import java.io.Serializable;
import java.time.Instant;

@Slf4j
public class DCPClient {

    private final ColmenaPlatformGrpc.ColmenaPlatformBlockingStub blockingStub;
    private final ColmenaPlatformGrpc.ColmenaPlatformStub stub;

    public DCPClient(ManagedChannel channel) {
        blockingStub = ColmenaPlatformGrpc.newBlockingStub(channel);
        stub = ColmenaPlatformGrpc.newStub(channel);
    }

    public void store(String key, Serializable object) {
        RouteGuide.StorageRequest request = RouteGuide.StorageRequest.newBuilder()
                .setKey(key)
                .setValue(ByteString.copyFrom(SerializationUtils.serialize(object)))
                .build();
        requestWithoutResponse(() -> blockingStub.store(request));
    }

    public void storeMetrics(String key, double value, Instant timestamp) {
        Timestamp converted = Timestamp.newBuilder().setSeconds(timestamp.getEpochSecond())
             .setNanos(timestamp.getNano()).build();
        RouteGuide.MetricsStorageRequest request = RouteGuide.MetricsStorageRequest.newBuilder()
                .setKey(key)
                .setValue(value)
                .setTimestamp(converted)
                .build();
        requestWithoutResponse(() -> blockingStub.storeMetrics(request));
    }

    public boolean metricMet(RoleDefinition.Metric metric) {
        Instant cutoff = Instant.now().minus(metric.getAmountOfTime(), metric.getUnit());
        Timestamp converted = Timestamp.newBuilder().setSeconds(cutoff.getEpochSecond())
                .setNanos(cutoff.getNano()).build();
        RouteGuide.MetricsQueryRequest request = RouteGuide.MetricsQueryRequest.newBuilder()
                .setKey(metric.getKey())
                .setThreshold(metric.getThreshold())
                .setThresholdType(RouteGuide.ThresholdType.valueOf(metric.getThresholdType().name()))
                .setFrom(converted)
                .build();
        RouteGuide.MetricsQueryResponse metricsQueryResponse = requestWithResponse(() -> blockingStub.queryMetrics(request));
        return metricsQueryResponse.getMet();
    }

    public void publish(String key, Serializable object) {
        RouteGuide.StorageRequest request = RouteGuide.StorageRequest.newBuilder()
                .setKey(key)
                .setValue(ByteString.copyFrom(SerializationUtils.serialize(object)))
                .build();
        requestWithoutResponse(() -> blockingStub.publish(request));
    }

    public Serializable getStored(String key) {
        RouteGuide.GetStoredRequest request = RouteGuide.GetStoredRequest.newBuilder()
                .setKey(key)
                .build();
        RouteGuide.GetStoredResponse response = requestWithResponse(() -> blockingStub.getStored(request));
        return SerializationUtils.deserialize(response.getValue().toByteArray());
    }

    public void addService(ServiceDescription serviceDescription) {
        RouteGuide.ServiceDescription request = ServiceDescriptionConverter.convert(serviceDescription);
        requestWithoutResponse(() -> blockingStub.addService(request));
    }

    public <T extends Serializable> CommunicationChannelListener<T> subscribe(String key) {
        RouteGuide.GetStoredRequest request = RouteGuide.GetStoredRequest.newBuilder().setKey(key).build();
        StorageQueueCommunicationChannelListener<T> streamObserver = new StorageQueueCommunicationChannelListener<>();
        stub.subscribe(request, streamObserver);
        return streamObserver;
    }

    public Serializable getSubscriptionItem(String key) {
        RouteGuide.GetStoredRequest request = RouteGuide.GetStoredRequest.newBuilder()
                .setKey(key)
                .build();
        RouteGuide.GetStoredResponse response = requestWithResponse(() -> blockingStub.getSubscriptionItem(request));
        return SerializationUtils.deserialize(response.getValue().toByteArray());
    }

    private void requestWithoutResponse(NoResponse request) {
        try {
            request.run();
        } catch (StatusRuntimeException e) {
            log.error("RPC failed: status: {0}, exception: {}", e.getStatus(), e);
            throw new RuntimeException();
        }
    }

    private <T> T requestWithResponse(WithResponse<T> request) {
        try {
            return request.run();
        } catch (StatusRuntimeException e) {
            log.error("RPC failed: status: {0}, exception: {}", e.getStatus(), e);
            throw new RuntimeException();
        }
    }

    @FunctionalInterface
    interface NoResponse {
        void run() throws StatusRuntimeException;
    }

    @FunctionalInterface
    interface WithResponse<T> {
        T run() throws StatusRuntimeException;
    }
}
