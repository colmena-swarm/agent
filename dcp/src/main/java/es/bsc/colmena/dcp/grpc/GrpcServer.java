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

import es.bsc.colmena.app.routeguide.ColmenaPlatformGrpc;
import io.grpc.Grpc;
import io.grpc.InsecureServerCredentials;
import io.grpc.Server;
import io.grpc.protobuf.services.ProtoReflectionService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

import java.io.IOException;

import static java.util.concurrent.TimeUnit.*;

@Slf4j
@RequiredArgsConstructor
public class GrpcServer {

    private final int port;
    private final ColmenaPlatformGrpc.ColmenaPlatformImplBase service;
    private Server server;

    public void startServer() throws IOException {
        server = Grpc.newServerBuilderForPort(port, InsecureServerCredentials.create())
                .addService(service)
                .addService(ProtoReflectionService.newInstance())
                .intercept(new IpAddressInterceptor())
                .keepAliveTime(1, DAYS)
                .keepAliveTimeout(1, DAYS)
                .permitKeepAliveTime(1, DAYS)
                .permitKeepAliveWithoutCalls(true)
                .maxConnectionIdle(1, DAYS)
                .maxConnectionAge(1, DAYS)
                .maxConnectionAgeGrace(1, DAYS)
                .build()
                .start();
        log.info("Server started, listening on " + port);
        Runtime.getRuntime().addShutdownHook(new Thread(() -> {
            // Use stderr here since the logger may have been reset by its JVM shutdown hook.
            log.info("shutting down gRPC server since JVM is shutting down");
            try {
                stopServer();
            } catch (InterruptedException e) {
                e.printStackTrace(System.err);
            }
            log.info("server shut down");
        }));
    }

    public void stopServer() throws InterruptedException {
        if (server != null) {
            log.info("Stopping server");
            server.shutdown().awaitTermination(30, SECONDS);
            log.info("Server stopped");
        }
    }

    public void blockUntilShutdown() throws InterruptedException {
        if (server != null) {
            server.awaitTermination();
        }
    }
}
