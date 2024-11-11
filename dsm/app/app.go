package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"colmena.bsc.es/agent/colmenacontext"
	"colmena.bsc.es/agent/docker"
	"colmena.bsc.es/agent/role"
)

// Run starts the HTTP server and blocks until the provided context is canceled.
// It returns an error only if the server fails unexpectedly before shutdown.
func Run(ctx context.Context) error {
	agentId := agentId()
	interfc := os.Getenv("PEER_DISCOVERY_INTERFACE")

	colmenacontext.PublishColmenaServiceDefinition(agentId)

	roleRunner := role.CommandListener{
		ContainerEngine: docker.DockerContainerEngine{},
		AgentId:         agentId,
		Interfc:         interfc,
	}

	srv := &http.Server{
		Addr:    ":50551",
		Handler: roleRunner.Endpoints(),
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("Listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
		return nil
	case err := <-errCh:
		return err
	}
}

func agentId() string {
	agentId := os.Getenv("AGENT_ID")
	if len(agentId) == 0 {
		log.Printf("Agent ID not set, using hostname")
		hostname, err := os.Hostname()
		if err != nil {
			log.Fatalf("Error getting hostname: %v", err)
		}
		agentId = hostname
	}
	log.Printf("Agent ID: %v", agentId)
	return agentId
}
