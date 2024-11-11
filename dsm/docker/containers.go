/*
 *  Copyright 2002-2025 Barcelona Supercomputing Center (www.bsc.es)
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
package docker

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type ContainerEngine interface {
	RunContainer(containerId string, imageId string, agentId string, interfc string) (string, error)
	StopContainer(containerId string) error
	Subscribe(stopped chan string)
}

type DockerContainerEngine struct {}

func (DockerContainerEngine) RunContainer(containerId string, imageId string, agentId string, interfc string) (string, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	rc, err := cli.ImagePull(ctx, imageId, types.ImagePullOptions{})
	if err != nil {
		return "", err
	}
	defer rc.Close()
	err = waitForImagePull(imageId, cli)
	if err != nil {
		return "", err
	}

	log.Printf("Starting container: imageName: %v", imageId)
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageId,
		Cmd:   []string{"echo", "hello world"},
		Env:   []string{"PEER_DISCOVERY_INTERFACE="+interfc, "HOSTNAME="+agentId, "AGENT_ID="+agentId},
	}, &container.HostConfig{NetworkMode: "host"}, nil, nil, containerId)
	if err != nil {
		return "", err
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", err
	}

	return resp.ID, nil
}

func (DockerContainerEngine) StopContainer(containerId string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	err = cli.ContainerStop(ctx, containerId, container.StopOptions{})
	log.Printf("Stopped container: %v", containerId)
	return err
}

func waitForImagePull(imageId string, cli *client.Client) error {
	wg := sync.WaitGroup{}
    wg.Add(1)
	timeout := 180 // 3 minutes
    go func(wg *sync.WaitGroup, timeout int) {
		for {
			list, err := cli.ImageList(context.Background(), types.ImageListOptions{})
			if err != nil {
				panic(err)
			}
			for _, ims := range list {
				for _, tag := range ims.RepoTags {
					if strings.Contains(tag, imageId) {
						wg.Done()
						return
					}
				}
			}
			if timeout == 0 {
				log.Panicf("timed out while waiting for image pull. imageId: %v", imageId)
			}
			timeout--
			time.Sleep(1 * time.Second)
		}
    }(&wg, timeout)
    wg.Wait()
	log.Printf("Pulled imageId: %v", imageId)
	return nil
}

