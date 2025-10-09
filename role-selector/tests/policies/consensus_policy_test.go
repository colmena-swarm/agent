package policy_test

import (
	"context"
	"log"
	"net"
	"sync"
	"testing"
	"time"

	"colmena.bsc.es/role-selector/policy"
	"colmena.bsc.es/role-selector/policy/grpc/colmena_consensus"
	"colmena.bsc.es/role-selector/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

// RemoteNode simulates a remote node with a SelectionService server
// and a client pointing to the local policy's SchedulingService
type RemoteNode struct {
	colmena_consensus.UnimplementedSelectionServiceServer
	localPort  string
	grpcServer *grpc.Server
}

// NewRemoteNode creates and starts a remote node
func NewRemoteNode(t *testing.T, listenPort, localPort string) *RemoteNode {
	node := &RemoteNode{localPort: localPort}
	lis, err := net.Listen("tcp", listenPort)
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	node.grpcServer = grpc.NewServer()

	colmena_consensus.RegisterSelectionServiceServer(node.grpcServer, node)

	go func() {
		if err := node.grpcServer.Serve(lis); err != nil {
			t.Logf("RemoteNode server stopped: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond) // give server time to start
	return node
}

func (r *RemoteNode) Stop() {
	r.grpcServer.GracefulStop()
}

// Implement SelectionService
func (r *RemoteNode) RequestRoles(ctx context.Context, req *colmena_consensus.RoleRequest) (*emptypb.Empty, error) {
	log.Printf("RemoteNode received RequestRoles: role=%s, startOrStop=%v", req.Role.RoleId, req.StartOrStop)

	// Connect to local policy SchedulingService
	conn, err := grpc.NewClient("localhost:"+r.localPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to local policy: %v", err)
		return &emptypb.Empty{}, nil
	}
	defer conn.Close()

	localClient := colmena_consensus.NewSchedulingServiceClient(conn)

	// Simulate a small delay
	time.Sleep(50 * time.Millisecond)

	// Call TriggerRole on local policy
	_, err = localClient.TriggerRole(ctx, &colmena_consensus.TriggerRoleRequest{
		RoleId:      req.Role.RoleId,
		ServiceId:   req.Role.ServiceId,
		StartOrStop: req.StartOrStop,
	})
	if err != nil {
		log.Printf("Failed to call TriggerRole: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// === Test Example ===
func TestConsensusPolicyRemoteNode(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	var mu sync.Mutex
	var triggered []types.Decision

	trigger := func(dec types.Decision) {
		mu.Lock()
		defer mu.Unlock()
		triggered = append(triggered, dec)
		log.Printf("Triggered decision: %+v", dec)
		wg.Done()
	}

	// Local policy server port
	localPort := "50052"
	// Remote node server port
	remotePort := "50051"

	// Start remote node
	remote := NewRemoteNode(t, ":"+remotePort, localPort)
	defer remote.Stop()

	// Create ConsensusPolicy pointing to remote node
	cp, err := policy.NewConsensusPolicy("localhost:"+remotePort, trigger)
	if err != nil {
		t.Fatalf("Failed to create ConsensusPolicy: %v", err)
	}
	defer cp.Stop()

	// Fake roles & KPIs
	roles := []*types.Role{{Id: "role1", ServiceId: "service1", ImageId: "img1", State: types.Stopped}}
	kpis := []types.KPI{{Query: "cpu_utilization", AssociatedRole: "role1", Level: "Critical"}}

	// Trigger consensus (asynchronous)
	err = cp.Client.RequestRoles(roles, kpis, nil)
	if err != nil {
		t.Fatalf("TriggerConsensus failed: %v", err)
	}

	// wait for the trigger (with timeout)
	done := make(chan struct{})
	go func() {
		wg.Wait()
		t.Logf("Decisions received")
		close(done)
	}()

	select {
	case <-done:
		// Lock the slice to read safely
		mu.Lock()
		defer mu.Unlock()

		found := false
		for _, dec := range triggered {
			if dec.RoleId == "role1" && dec.ServiceId == "service1" {
				found = true
				break
			}
		}

		if !found {
			t.Fatalf("Expected TriggerDecision for role1/service1 but it was not called. Decisions: %+v", triggered)
		} else {
			t.Logf("TriggerDecision correctly called for role1/service1")
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for triggered decision")
	}
}

// === Integration Test ===
// This test assumes a real ConsensusPolicy server is already running and reachable at remotePort.
func TestConsensusPolicyIntegrationRealServer(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	var mu sync.Mutex
	var triggered []types.Decision

	trigger := func(dec types.Decision) {
		mu.Lock()
		defer mu.Unlock()
		triggered = append(triggered, dec)
		log.Printf("Triggered decision: %+v", dec)
		wg.Done()
	}

	// Real server address
	remotePort := "50052" // the server must already be running here

	// Connect to real ConsensusPolicy server
	cp, err := policy.NewConsensusPolicy("localhost:"+remotePort, trigger)
	if err != nil {
		t.Fatalf("Failed to connect to real ConsensusPolicy: %v", err)
	}
	defer cp.Stop()

	// Fake roles & KPIs
	roles := []*types.Role{{Id: "role1", ServiceId: "service1", ImageId: "img1", State: types.Stopped}}
	kpis := []types.KPI{{Query: "cpu_utilization", AssociatedRole: "role1", Level: "Critical"}}

	// Trigger consensus asynchronously
	err = cp.Client.RequestRoles(roles, kpis, nil)
	if err != nil {
		t.Fatalf("Request role call failed: %v", err)
	}

	// Wait for the triggered callback
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		mu.Lock()
		defer mu.Unlock()

		found := false
		for _, dec := range triggered {
			if dec.RoleId == "role1" && dec.ServiceId == "service1" {
				found = true
				break
			}
		}

		if !found {
			t.Fatalf("Expected TriggerDecision for role1/service1 but it was not called. Decisions: %+v", triggered)
		} else {
			t.Logf("TriggerDecision correctly called for role1/service1")
		}

	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for triggered decision from real server")
	}
}
