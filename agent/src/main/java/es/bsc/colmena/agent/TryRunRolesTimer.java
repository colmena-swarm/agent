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
package es.bsc.colmena.agent;

import lombok.extern.slf4j.Slf4j;

import java.util.Timer;
import java.util.TimerTask;
import java.util.concurrent.atomic.AtomicBoolean;

@Slf4j
public class TryRunRolesTimer {
    private final Timer timer = new Timer();
    private final TryRunRolesTimerTask timerTask;

    public TryRunRolesTimer(Device device) {
        this.timerTask = new TryRunRolesTimerTask(device);
    }

    public void schedule() {
        timer.scheduleAtFixedRate(timerTask, 5000, 5000);
    }

    public void disconnect() {
        timerTask.stop();
        timer.cancel();
        timer.purge();
    }

    @Slf4j
    private static class TryRunRolesTimerTask extends TimerTask {

        private final Device device;
        private final AtomicBoolean runTask = new AtomicBoolean(false);

        public TryRunRolesTimerTask(Device device) {
            this.device = device;
        }

        @Override
        public void run() {
            runTask.set(true);
            if (runTask.get()) {
                log.info("Periodic role check");
                device.tryRunRoles();
            } else {
                log.info("Not running periodic role check");
            }
        }

        public void stop() {
            runTask.set(false);
            log.info("Stopped periodic role check");
        }
    }
}
