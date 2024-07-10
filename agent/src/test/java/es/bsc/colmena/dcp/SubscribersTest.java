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

import com.google.protobuf.ByteString;
import es.bsc.colmena.app.routeguide.RouteGuide;
import io.grpc.stub.StreamObserver;
import es.bsc.colmena.dcp.storage.Subscriber;
import es.bsc.colmena.dcp.storage.Subscribers;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.is;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@SuppressWarnings("unchecked")
public class SubscribersTest {

    @Test
    public void exception_thrown_removes_subscriber() {
        Subscriber streamObserver = mock(Subscriber.class);
        Subscribers subscribers = new Subscribers("key", streamObserver);
        Mockito.doThrow(new RuntimeException()).when(streamObserver).onNext(any());
        subscribers.publish(ByteString.EMPTY);
        assertThat(subscribers.numberOfSubscribers(), is(0));
    }

    @Test
    public void two_messages_sent_to_single_subscriber() {
        Subscriber streamObserver = mock(Subscriber.class);
        Subscribers subscribers = new Subscribers("key", streamObserver);
        subscribers.publish(ByteString.EMPTY);
        subscribers.publish(ByteString.EMPTY);
        verify(streamObserver, times(2)).onNext(any());
    }

    @Test
    public void two_messages_sent_to_same_subscriber() {
        Subscriber first = create("first");
        Subscriber secondSubscriber = create("second");
        Subscribers subscribers = new Subscribers("key", first);
        subscribers.publish(ByteString.EMPTY);
        subscribers.add(secondSubscriber);
        subscribers.publish(ByteString.EMPTY);
        verify(first.getStreamObserver(), times(2)).onNext(any());
    }

    @Test
    public void two_messages_sent_to_same_single_subscriber() {
        Subscriber first = createSingle("first");
        Subscriber secondSubscriber = createSingle("second");

        Subscribers subscribers = new Subscribers("key", first);
        subscribers.publish(ByteString.EMPTY);
        subscribers.add(secondSubscriber);
        assertThat(subscribers.numberOfSubscribers(), is(1));

        Subscriber recreatedFirst = createSingle("first");
        subscribers.add(recreatedFirst);
        subscribers.publish(ByteString.EMPTY);
        verify(first.getStreamObserver(), times(1)).onNext(any());
        verify(recreatedFirst.getStreamObserver(), times(1)).onNext(any());
    }

    @Test
    public void single_subscriber_removed_and_completed_after_publish() {
        StreamObserver<RouteGuide.GetStoredResponse> streamObserver = mock(StreamObserver.class);
        Subscriber.SingleSubscriber singleSubscriber = new Subscriber.SingleSubscriber(streamObserver, "127.0.0.1");
        Subscribers subscribers = new Subscribers("key", singleSubscriber);
        subscribers.publish(ByteString.EMPTY);
        subscribers.publish(ByteString.EMPTY);
        verify(streamObserver, times(1)).onNext(any());
        verify(streamObserver, times(1)).onCompleted();
        assertThat(subscribers.numberOfSubscribers(), is(0));
    }

    private Subscriber create(String ipAddress) {
        StreamObserver<RouteGuide.GetStoredResponse> streamObserver = mock(StreamObserver.class);
        return new Subscriber(streamObserver, ipAddress);
    }

    private Subscriber createSingle(String ipAddress) {
        StreamObserver<RouteGuide.GetStoredResponse> streamObserver = mock(StreamObserver.class);
        return new Subscriber.SingleSubscriber(streamObserver, ipAddress);
    }
}
