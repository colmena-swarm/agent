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
package es.bsc.colmena.rolerunner;

import com.google.common.util.concurrent.*;
import es.bsc.colmena.library.JavaRole;
import es.bsc.colmena.library.Role;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

import java.util.HashSet;
import java.util.Set;
import java.util.concurrent.*;
import java.util.function.Consumer;

import static java.util.concurrent.TimeUnit.SECONDS;

@Slf4j
public class JavaRoleRunner {
    private final ThreadPoolExecutor executor = (ThreadPoolExecutor) Executors.newFixedThreadPool(10);
    private final ExecutorService executorService = MoreExecutors.getExitingExecutorService(executor, 1, SECONDS);
    private final ListeningExecutorService listeningExecutorService = MoreExecutors.listeningDecorator(executorService);
    private final ConcurrentMap<String, JavaRole> runningRoles = new ConcurrentHashMap<>();

    public synchronized void startRole(JavaRole javaRole) {
        ListenableFuture<JavaRole> submitted = listeningExecutorService.submit(javaRole, javaRole);
        Futures.addCallback(submitted, new JavaRoleFinished(javaRole, this::roleStopped), executor);
        runningRoles.put(javaRole.getRoleId(), javaRole);
        log.info("Started roleId: {}", javaRole.getRoleId());
    }

    public synchronized void stopRole(JavaRole javaRole) {
        javaRole.stopExecuting();
        roleStopped(javaRole);
        log.info("Stopped roleId: {}", javaRole.getRoleId());
    }

    private synchronized void roleStopped(JavaRole javaRole) {
        log.info("Role stopped. Removing from running roles. roleId: {}", javaRole.getRoleId());
        runningRoles.remove(javaRole.getRoleId());
    }

    public synchronized Set<Role> getRunningRoles() {
        return new HashSet<>(runningRoles.values());
    }

    public synchronized void disconnect() throws InterruptedException {
        log.info("Disconnecting. Stopping {} roles", runningRoles.size());
        executor.shutdownNow();
        executor.awaitTermination(5, TimeUnit.SECONDS);
    }

    public void blockUntilShutdown() throws InterruptedException {
        executor.awaitTermination(1, TimeUnit.DAYS);
    }

    @RequiredArgsConstructor
    private static class JavaRoleFinished implements FutureCallback<JavaRole> {
        public final JavaRole startedRole;
        public final Consumer<JavaRole> roleStopped;

        @Override
        public void onSuccess(JavaRole result) {
            roleStopped.accept(startedRole);
        }

        @Override
        public void onFailure(Throwable t) {
            roleStopped.accept(startedRole);
        }
    }
}
