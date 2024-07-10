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
package es.bsc.colmena.infrastructure;

import com.google.common.net.HostAndPort;
import io.grpc.InsecureChannelCredentials;
import io.grpc.ManagedChannel;

import java.util.concurrent.TimeUnit;

import static io.grpc.Grpc.newChannelBuilderForAddress;

public class GrpcChannel {

    public static ManagedChannel newChannel(HostAndPort dcpHostPort) {
        return newChannelBuilderForAddress(dcpHostPort.getHost(), dcpHostPort.getPort(), InsecureChannelCredentials.create())
                .keepAliveTime(1, TimeUnit.DAYS)
                .keepAliveTimeout(1, TimeUnit.DAYS)
                .keepAliveWithoutCalls(true)
                .build();
    }

}