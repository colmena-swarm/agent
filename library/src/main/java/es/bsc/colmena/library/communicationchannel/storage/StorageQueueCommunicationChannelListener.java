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
package es.bsc.colmena.library.communicationchannel.storage;

import es.bsc.colmena.app.routeguide.RouteGuide;
import es.bsc.colmena.library.communicationchannel.CommunicationChannelListener;
import io.grpc.stub.StreamObserver;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.lang3.SerializationUtils;

import java.io.Serializable;
import java.util.concurrent.BlockingQueue;
import java.util.concurrent.LinkedBlockingQueue;

@Slf4j
@RequiredArgsConstructor
public class StorageQueueCommunicationChannelListener<T extends Serializable>
        implements CommunicationChannelListener<T>, StreamObserver<RouteGuide.GetStoredResponse> {

    private final BlockingQueue<byte[]> messageQueue = new LinkedBlockingQueue<>();

    @Override
    public T take() throws InterruptedException {
        return (T) SerializationUtils.deserialize(messageQueue.take());
    }

    @Override
    public void onNext(RouteGuide.GetStoredResponse value) {
        messageQueue.add(value.getValue().toByteArray());
    }

    @Override
    public void onError(Throwable t) {
        log.error("Error in grpc client streaming", t);
    }

    @Override
    public void onCompleted() {
        log.info("Stream completed. Exiting");
    }
}
