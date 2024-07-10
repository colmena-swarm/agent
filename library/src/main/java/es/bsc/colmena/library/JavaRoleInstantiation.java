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

import es.bsc.colmena.library.communicationchannel.annotations.ColmenaMessageProcessor;
import es.bsc.colmena.library.communicationchannel.annotations.ColmenaMessagePublisher;
import es.bsc.colmena.library.communicationchannel.storage.StorageQueueCommunicationChannelPublisher;
import lombok.SneakyThrows;
import org.apache.commons.lang3.tuple.Pair;

import java.io.Serializable;
import java.lang.reflect.Field;
import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.util.Arrays;
import java.util.Map;
import java.util.function.Consumer;
import java.util.stream.Collectors;

public class JavaRoleInstantiation {

    @SneakyThrows
    public static <T extends JavaRole> JavaRole instantiate(Class<T> clazz, DCPClient dcpClient) {
        T javaRole = clazz.getDeclaredConstructor(DCPClient.class).newInstance(dcpClient);
        setMessageChannelSubscribers(javaRole);
        setMessageChannelPublishers(javaRole, dcpClient);
        return javaRole;
    }

    private static <T extends JavaRole> void setMessageChannelPublishers(T javaRole, DCPClient dcpClient) {
        Class<? extends JavaRole> clazz = javaRole.getClass();
        Arrays.stream(clazz.getDeclaredFields())
                .filter(field -> field.isAnnotationPresent(ColmenaMessagePublisher.class))
                .forEach(field -> set(field, javaRole, dcpClient));
    }

    @SneakyThrows
    private static <T extends JavaRole> void set(Field field, T javaRole, DCPClient dcpClient) {
        String key = field.getAnnotation(ColmenaMessagePublisher.class).key();
        var communicationChannel = new StorageQueueCommunicationChannelPublisher<>(key, dcpClient);
        field.set(javaRole, communicationChannel);
    }

    private static <T extends JavaRole> void setMessageChannelSubscribers(T javaRole) {
        Class<? extends JavaRole> clazz = javaRole.getClass();
        javaRole.setMessageChannelSubscriptions(
                Arrays.stream(clazz.getDeclaredMethods())
                        .filter(method -> method.isAnnotationPresent(ColmenaMessageProcessor.class))
                        .map(method -> parseFromMethod(javaRole, method))
                        .collect(Collectors.toSet()));
    }


    private static Map.Entry<String, Consumer<? extends Serializable>> parseFromMethod(Role role, Method method) {
        String subscriptionKey = method.getAnnotation(ColmenaMessageProcessor.class).key();
        var consumer = methodAsConsumer(role, method);
        return Pair.of(subscriptionKey, consumer);
    }

    private static Consumer<? extends Serializable> methodAsConsumer(Object annotated, Method m) {
        return param -> {
            try {
                m.invoke(annotated, param);
            } catch (IllegalAccessException | InvocationTargetException e) {
                throw new RuntimeException(e);
            }
        };
    }
}
