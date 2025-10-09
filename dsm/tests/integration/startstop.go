package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"colmena.bsc.es/agent/role"
	"github.com/docker/docker/client"
)

func start(t *testing.T, serviceId string, roleId string, imageId string) {
	// Prepare Docker client
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		t.Fatalf("docker client error: %v", err)
	}
	defer cli.Close()

	// Best-effort cleanup from previous runs
	_ = removeContainerIfExists(ctx, cli, imageId)

	// Trigger start
	postJSON(t, "http://127.0.0.1:50551/start", role.StartRoleCommand{ServiceId: serviceId, RoleId: roleId, ImageId: imageId})

	// Wait for container to be created (running or exited is acceptable given echo command)
	if err := waitForContainerExists(ctx, cli, imageId, 120*time.Second); err != nil {
		t.Fatalf("container did not appear: %v", err)
	}
}

func stop(t *testing.T, serviceId string, roleId string, imageId string, removeContainer bool) {

	// Prepare Docker client
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		t.Fatalf("docker client error: %v", err)
	}
	defer cli.Close()

	stopRoleCommand := role.StopRoleCommand{
		ServiceId: serviceId, 
		RoleId: roleId, 
		ImageId: imageId, 
		RemoveContainer: removeContainer,
	}

	// Trigger stop (idempotent; service returns 200 regardless)
	postJSON(t, "http://127.0.0.1:50551/stop", stopRoleCommand)

	// Verify container is not running (allow exited)
	running, err := isContainerRunning(ctx, cli, imageId)
	if err != nil {
		t.Fatalf("inspect after stop failed: %v", err)
	}
	if running {
		// Wait briefly for stop to take effect
		if err := waitForContainerStopped(ctx, cli, imageId, 30*time.Second); err != nil {
			t.Fatalf("container still running after stop: %v", err)
		}
	}

	if removeContainer {
		//wait for container to be deleted
		if err := waitForContainerDeleted(ctx, cli, imageId, 30*time.Second); err != nil {
			t.Fatalf("container still exists after stop: %v", err)
		}
	}
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
