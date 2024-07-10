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
package es.bsc.colmena.dcp.grpc;

import io.grpc.*;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.lang3.StringUtils;

import java.util.Objects;

@Slf4j
public class IpAddressInterceptor implements ServerInterceptor {

    public static final Context.Key<Object> IP_ADDRESS_KEY = Context.key("IP_ADDRESS");

    @Override
    public <ReqT, RespT> ServerCall.Listener<ReqT> interceptCall(
            ServerCall<ReqT, RespT> serverCall,
            Metadata metadata,
            ServerCallHandler<ReqT, RespT> serverCallHandler) {

        String ipAddress = Objects.requireNonNull(serverCall.getAttributes().get(Grpc.TRANSPORT_ATTR_REMOTE_ADDR)).toString();
        String parsed = parse(ipAddress);
        Context context = Context.current().withValue(IP_ADDRESS_KEY, parsed);
        return Contexts.interceptCall(context, serverCall, metadata, serverCallHandler);
    }

    private String parse(String interceptedIpAddress) {
        try {
            String removedStartingSlash = StringUtils.remove(interceptedIpAddress, '/');
            return removedStartingSlash.split(":")[0];
        } catch (Exception e) {
            log.error("Could not parse IP address, returning intercepted. intercepted: {}", interceptedIpAddress);
            return interceptedIpAddress;
        }
    }
}
