package integration

import (
	"testing"
	"time"
)

func startDsm(t *testing.T) func() {
	// Prevent publishing colmena service definition
	t.Setenv("ZENOH_URL", "")

	// Start service in-process
	cancelSvc := startService(t)
	
	waitForHealthcheck(t, "http://127.0.0.1:50551", 30*time.Second)

	return cancelSvc
}