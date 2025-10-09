package integration

import (
	"context"
	"errors"
	"time"

	"colmena.bsc.es/agent/docker"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func dockerFindContainerByName(ctx context.Context, cli *client.Client, imageId string) (string, error) {
	containerName := docker.ContainerName(imageId)
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return "", err
	}
	want := "/" + containerName
	for _, c := range containers {
		for _, n := range c.Names {
			if n == want {
				return c.ID, nil
			}
		}
	}
	return "", errors.New("not found")
}

func waitForContainerExists(ctx context.Context, cli *client.Client, name string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if _, err := dockerFindContainerByName(ctx, cli, name); err == nil {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return errors.New("timeout waiting for container")
}

func waitForContainerDeleted(ctx context.Context, cli *client.Client, name string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if _, err := dockerFindContainerByName(ctx, cli, name); err != nil {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return errors.New("timeout waiting for container")
}

func isContainerRunning(ctx context.Context, cli *client.Client, name string) (bool, error) {
	id, err := dockerFindContainerByName(ctx, cli, name)
	if err != nil {
		return false, err
	}
	ins, err := cli.ContainerInspect(ctx, id)
	if err != nil {
		return false, err
	}
	if ins.ContainerJSONBase == nil || ins.State == nil {
		return false, nil
	}
	return ins.State.Running, nil
}

func waitForContainerStopped(ctx context.Context, cli *client.Client, name string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		running, err := isContainerRunning(ctx, cli, name)
		if err != nil {
			return err
		}
		if !running {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return errors.New("timeout waiting for stop")
}

func waitForContainerStarted(ctx context.Context, cli *client.Client, name string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		running, err := isContainerRunning(ctx, cli, name)
		if err != nil {
			return err
		}
		if running {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return errors.New("timeout waiting for start")
}

func removeContainerIfExists(ctx context.Context, cli *client.Client, name string) error {
	id, err := dockerFindContainerByName(ctx, cli, name)
	if err == nil {
		return cli.ContainerRemove(ctx, id, types.ContainerRemoveOptions{Force: true, RemoveVolumes: true})
	}
	return nil
}


