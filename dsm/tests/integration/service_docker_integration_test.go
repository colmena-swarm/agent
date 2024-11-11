package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"colmena.bsc.es/agent/app"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type roleCommand struct {
	ServiceId string `json:"serviceId"`
	RoleId    string `json:"roleId"`
	ImageId   string `json:"imageId"`
}

func TestService_StartAndStopContainer(t *testing.T) {
	// Ensure we don't inherit any conflicting env that affects service publishing
	t.Setenv("ZENOH_URL", "")

	// Start service in-process
	cancelSvc := startService(t)
	defer cancelSvc()

	waitForHealthcheck(t, "http://127.0.0.1:50551", 30*time.Second)

	// Prepare Docker client
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		t.Fatalf("docker client error: %v", err)
	}
	defer cli.Close()

	serviceId := "itest"
	roleName := "itest-role"
	image := "busybox:latest"

	// Best-effort cleanup from previous runs
	_ = removeContainerIfExists(ctx, cli, roleName)

	// Trigger start
	postJSON(t, "http://127.0.0.1:50551/start", roleCommand{ServiceId: serviceId, RoleId: roleName, ImageId: image})

	// Wait for container to be created (running or exited is acceptable given echo command)
	if err := waitForContainerExists(ctx, cli, roleName, 120*time.Second); err != nil {
		t.Fatalf("container did not appear: %v", err)
	}

	// Optionally observe running state (non-fatal if not running due to short-lived command)
	_, _ = isContainerRunning(ctx, cli, roleName)

	// Trigger stop (idempotent; service returns 200 regardless)
	postJSON(t, "http://127.0.0.1:50551/stop", roleCommand{ServiceId: serviceId, RoleId: roleName, ImageId: image})

	// Verify container is not running (allow exited)
	running, err := isContainerRunning(ctx, cli, roleName)
	if err != nil {
		t.Fatalf("inspect after stop failed: %v", err)
	}
	if running {
		// Wait briefly for stop to take effect
		if err := waitForContainerStopped(ctx, cli, roleName, 30*time.Second); err != nil {
			t.Fatalf("container still running after stop: %v", err)
		}
	}

	// Cleanup container to allow re-runs
	if err := removeContainerIfExists(ctx, cli, roleName); err != nil {
		t.Fatalf("cleanup failed: %v", err)
	}
}

func startService(t *testing.T) func() {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	go func() { _ = app.Run(ctx) }()
	return cancel
}

func waitForHealthcheck(t *testing.T, baseURL string, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	client := http.Client{Timeout: 2 * time.Second}
	for time.Now().Before(deadline) {
		resp, err := client.Get(baseURL + "/healthz")
		if err == nil {
			if resp.StatusCode == http.StatusOK {
				_ = resp.Body.Close()
				return
			}
			_ = resp.Body.Close()
		}
		time.Sleep(200 * time.Millisecond)
	}
	t.Fatalf("service did not become ready within %s", timeout)
}

func postJSON(t *testing.T, url string, body any) {
	t.Helper()
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	client := &http.Client{Timeout: 10 * time.Second}
	deadline := time.Now().Add(10 * time.Second)
	for {
		resp, err := client.Post(url, "application/json", bytes.NewReader(b))
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				t.Fatalf("unexpected status %d for %s", resp.StatusCode, url)
			}
			return
		}
		// retry on transient errors while server stabilizes
		if time.Now().After(deadline) {
			t.Fatalf("post %s: %v", url, err)
		}
		t.Fatalf("post %s: %v", url, err)
	}
}

func dockerFindContainerByName(ctx context.Context, cli *client.Client, name string) (string, error) {
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return "", err
	}
	want := "/" + name
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

func removeContainerIfExists(ctx context.Context, cli *client.Client, name string) error {
	id, err := dockerFindContainerByName(ctx, cli, name)
	if err == nil {
		return cli.ContainerRemove(ctx, id, types.ContainerRemoveOptions{Force: true, RemoveVolumes: true})
	}
	return nil
}
