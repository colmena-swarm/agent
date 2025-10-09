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
	RunContainer(serviceId string, roleId string, imageId string, agentId string, interfc string) (string, error)
	StopContainer(roleId string, imageId string, removeContainer bool) error
	Subscribe(ctx context.Context)
}

type DockerContainerEngine struct {
	Context context.Context
}

func (dce DockerContainerEngine) RunContainer(serviceId string, roleId string, imageId string, agentId string, interfc string) (string, error) {
	go dce.Subscribe(dce.Context)
	ctx := dce.Context
	containerName := ContainerName(imageId)
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Printf("Error creating client: %v", err)
		panic(err)
	}
	defer cli.Close()

	exists, err := checkImageExists(imageId, cli)
	if err != nil {
		log.Printf("Error checking if image exists: %v", err)
		return "", err
	}
	if !exists {	
		rc, err := cli.ImagePull(ctx, imageId, types.ImagePullOptions{})
		if err != nil {
			log.Printf("Error pulling image: %v", err)
			return "", err
		}
		defer rc.Close()
		err = waitForImagePull(imageId, cli)
		if err != nil {
			log.Printf("Error waiting for image pull: %v", err)
			return "", err
		}
	}

	containerExists, err := containerExists(cli, containerName)
	if err != nil {
		log.Printf("Error checking if container exists: %v", err)
		return "", err
	}

	if !containerExists {
		log.Printf("Starting container: imageName: %v", imageId)
		_, err := cli.ContainerCreate(ctx, &container.Config{
			Image: imageId,
			Env:   []string{"PEER_DISCOVERY_INTERFACE="+interfc, "HOSTNAME="+agentId, "AGENT_ID="+agentId},
			Labels: map[string]string{
				"es.bsc.colmena.roleId": roleId,
				"es.bsc.colmena.serviceId": serviceId,
				"es.bsc.colmena.imageId": imageId},
		}, &container.HostConfig{
			NetworkMode: "host",
			Binds: []string{
				"/tmp:/tmp",
				"/var/run/docker.sock:/var/run/docker.sock",
			},
		}, nil, nil, containerName)
		if err != nil {
			log.Printf("Error creating container: %v", err)
			return "", err
		}
	}

	err = cli.ContainerStart(ctx, containerName, container.StartOptions{}) 
	if err != nil {
		log.Printf("Error starting container: %v", err)
		return "", err
	}

	log.Printf("Started container: %v", containerName)
	return containerName, nil
}

func containerExists(cli *client.Client, id string) (bool, error) {
	_, err := cli.ContainerInspect(context.Background(), id)
	if err != nil {
		// Not found error means container doesn't exist
		if client.IsErrNotFound(err) {
			return false, nil
		}
		// Any other error is unexpected
		return false, err
	}
	return true, nil
}

func (DockerContainerEngine) StopContainer(roleId string, imageId string, removeContainer bool) error {
	containerName := ContainerName(imageId)
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	err = cli.ContainerStop(ctx, containerName, container.StopOptions{})
	log.Printf("Stopped container: %v", containerName)
	if removeContainer {
		err = cli.ContainerRemove(ctx, containerName, types.ContainerRemoveOptions{Force: true, RemoveVolumes: true})
		if err != nil {
			log.Printf("Error removing container: %v", err)
			return err
		}
		log.Printf("Removed container: %v", containerName)
	}
	return err
}

func waitForImagePull(imageId string, cli *client.Client) error {
	wg := sync.WaitGroup{}
    wg.Add(1)
	timeout := 600 // 10 minutes
    go func(wg *sync.WaitGroup, timeout int) {
		for {
			exists, err := checkImageExists(imageId, cli)
			if exists {
				wg.Done()
				return
			}
			if err != nil {
				panic(err)
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


func checkImageExists(imageId string, cli *client.Client) (bool, error) {
	ctx := context.Background()
	list, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return false, err
	}
	for _, ims := range list {
		for _, tag := range ims.RepoTags {
			if strings.Contains(tag, imageId) {
				return true, nil
			}
		}
	}
	return false, nil
}

func ContainerName(imageId string) string {
	slash_removed := strings.ReplaceAll(imageId, "/", "-")
	return strings.ReplaceAll(slash_removed, ":", "-")
}
