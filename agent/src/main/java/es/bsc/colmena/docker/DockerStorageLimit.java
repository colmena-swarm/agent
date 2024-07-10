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


import com.github.dockerjava.api.model.Container;
import com.github.dockerjava.api.model.Image;
import com.github.dockerjava.api.model.PruneType;
import lombok.extern.slf4j.Slf4j;

import java.util.Objects;

@Slf4j
public class DockerStorageLimit {

    public static boolean SHOULD_MANAGE_STORAGE = false;
    private static final long threshold = 2_000_000_000L; //2GB
    private static final String UNTIL_FILTER = "1h";
    private final com.github.dockerjava.api.DockerClient dockerClient;

    public DockerStorageLimit(com.github.dockerjava.api.DockerClient dockerClient) {
        this.dockerClient = dockerClient;
    }

    public void manage() {
        if (!SHOULD_MANAGE_STORAGE) {
            log.info("Docker objects will not be pruned");
            return;
        }
        if (!overLimit()) {
            return;
        }
        Long containerSpaceReclaimed = dockerClient.pruneCmd(PruneType.CONTAINERS).withUntilFilter(UNTIL_FILTER).exec().getSpaceReclaimed();
        Long imageSpaceReclaimed = dockerClient.pruneCmd(PruneType.IMAGES).withDangling(false).withUntilFilter(UNTIL_FILTER).exec().getSpaceReclaimed();
        log.info("Pruned docker images. untilFilter: {}, containerSpaceReclaimed: {}, imageSpaceReclaimed: {}",
                UNTIL_FILTER, containerSpaceReclaimed, imageSpaceReclaimed);
    }

    private boolean overLimit() {
        long combinedImageSize = dockerClient.listImagesCmd().exec().stream()
                .map(Image::getSize).reduce(0L, Long::sum);
        long combinedContainerSize = dockerClient.listContainersCmd().exec().stream()
                .map(Container::getSizeRootFs).filter(Objects::nonNull).reduce(0L, Long::sum);
        boolean overStorageLimit = combinedImageSize + combinedContainerSize > threshold;
        log.info("Over docker storage limit {}. combinedImageSize: {}, combinedContainerSize: {}", overStorageLimit, combinedImageSize, combinedContainerSize);
        return overStorageLimit;
    }
}
