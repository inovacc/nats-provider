package testutils

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// NATSServer wraps the running NATS container information
type NATSServer struct {
	Container testcontainers.Container
	URL       string
}

// StartNATSServer starts a disposable NATS server with JetStream enabled for integration testing.
func StartNATSServer(t *testing.T) *NATSServer {
	ctx := context.Background()

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "nats:latest",
			ExposedPorts: []string{"4222/tcp"},
			Cmd:          []string{"-js"}, // ✅ Enable JetStream explicitly
			WaitingFor:   wait.ForListeningPort("4222/tcp"),
		},
		Started: true,
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, container.Terminate(ctx))
	})

	port, err := container.MappedPort(ctx, "4222")
	require.NoError(t, err)

	host, err := container.Host(ctx)
	require.NoError(t, err)

	url := fmt.Sprintf("nats://%s:%s", host, port.Port())

	return &NATSServer{
		Container: container,
		URL:       url,
	}
}
