package integration

import (
	"testing"
)

func TestDockerEventsIntegration(t *testing.T) {

	// Setup mock role selector
	ms := NewMockRoleSelector(t, Expectation{
		Method: "PUT",
		Path:   "/stopped",
		BodyJSON: map[string]any{
			"roleId":    ROLE_NAME,
			"serviceId": SERVICE_ID,
			"imageId":   IMAGE,
		},
	})
	defer ms.Close()

	// Publish to mock role selector
	t.Setenv("ROLE_SELECTOR_URL", ms.URL())

	stopDsm := startDsm(t)
	defer stopDsm()

	start(t, SERVICE_ID, ROLE_NAME, IMAGE)

	stop(t, SERVICE_ID, ROLE_NAME, IMAGE, true)

	ms.Verify()
}
