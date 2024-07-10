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
package es.bsc.colmena.docker;

import com.github.dockerjava.api.async.ResultCallback;
import com.github.dockerjava.api.model.Event;
import com.github.dockerjava.api.model.EventType;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

import java.util.Objects;
import java.util.function.Consumer;

@Slf4j
@RequiredArgsConstructor
public class DockerMonitor extends ResultCallback.Adapter<Event> {

    private final Consumer<String> roleStopped;

    public void onNext(Event event) {
        EventType type = event.getType();
        if (type == EventType.CONTAINER) {
            containerEvent(event);
        }
    }

    private void containerEvent(Event event) {
        switch (event.getAction()) {
            case "die":
                deadContainer(event);
                break;
            case "destroy":
                destroyedContainer(event);
                break;
            default:
                // Ignore Event
        }
    }


    private void deadContainer(Event event) {
        String id = event.getId();
        roleStopped.accept(id);
    }

    private void destroyedContainer(Event event) {
        String id = event.getId();
        roleStopped.accept(id);
    }

}
