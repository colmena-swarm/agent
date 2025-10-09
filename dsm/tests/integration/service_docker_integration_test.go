package integration

import (
	"testing"
)

const SERVICE_ID = "itest"
const ROLE_NAME = "itest-role"
const IMAGE = "gcr.io/google-containers/pause:3.2"

func TestService_StartAndStopContainer(t *testing.T) {
	stopDsm := startDsm(t)
	defer stopDsm()

	start(t, SERVICE_ID, ROLE_NAME, IMAGE)

	stop(t, SERVICE_ID, ROLE_NAME, IMAGE, true)
}

func TestService_RestartContainer(t *testing.T) {
	stopDsm := startDsm(t)
	defer stopDsm()

	start(t, SERVICE_ID, ROLE_NAME, IMAGE)

	stop(t, SERVICE_ID, ROLE_NAME, IMAGE, false)

	start(t, SERVICE_ID, ROLE_NAME, IMAGE)

	stop(t, SERVICE_ID, ROLE_NAME, IMAGE, true)
}




