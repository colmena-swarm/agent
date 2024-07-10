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

import es.bsc.colmena.library.communicationchannel.CommunicationChannelListener;
import lombok.Getter;
import lombok.Setter;
import lombok.extern.slf4j.Slf4j;

import java.io.IOException;
import java.io.Serializable;
import java.util.Map;
import java.util.Objects;
import java.util.Set;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.function.Consumer;

@Getter
@Slf4j
public abstract class JavaRole extends Thread implements Role {
    protected final DCPClient dcpClient;
    protected AtomicBoolean running = new AtomicBoolean(true);
    protected final ExecutorService es = Executors.newCachedThreadPool();
    @Setter private Set<Map.Entry<String, Consumer<? extends Serializable>>> messageChannelSubscriptions;

    public JavaRole(DCPClient dcpClient) {
        this.dcpClient = dcpClient;
    }

    @Override
    public void run() {
        try {
            messageChannelSubscriptions.forEach(entry -> es.submit(() -> messageChannelSubscription(entry.getKey(), entry.getValue())));
            es.awaitTermination(1, TimeUnit.DAYS);
        } catch (InterruptedException e) {
            running.set(false);
        }
    }

    private <T extends Serializable> void messageChannelSubscription(String subscriptionKey, Consumer<T> messageProcessor) {
        CommunicationChannelListener<T> channelListener = dcpClient.subscribe(subscriptionKey);
        while(running.get()) {
            try {
                T message = channelListener.take();
                messageProcessor.accept(message);
            } catch (InterruptedException e) {
                log.info("Interrupted thread. subscriptionKey={}, roleId={}", subscriptionKey, getRoleId(), e);
                return;
            }
            catch (IOException e) {
                log.error("IO failed. subscriptionKey={}, roleId={}", subscriptionKey, getRoleId(), e);
                return;
            }
        }
    }

    public void stopExecuting() {
        running.set(false);
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.getRoleId());
    }

    @Override
    public boolean equals(Object other) {
        if (other == this) return true;
        if (!(other instanceof JavaRole otherCast)) return false;
        return this.getRoleId().equals(otherCast.getRoleId());
    }

    @Override
    public String toString() {
        return getRoleId();
    }
}
