package interceptor_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jhump/grpctunnel"
	"github.com/jhump/grpctunnel/tunnelpb"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/tbe-team/raybot/internal/handlers/cloud/cloudtest"
	"github.com/tbe-team/raybot/internal/handlers/cloud/interceptor"
)

func TestIntegrationReverseCredentialsInterceptor_AuthenticateSuccessfully(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	headerKey := "x-tunnel-token"
	requireToken := "test-token"

	port, stop := cloudtest.SetupTestCloudServer(
		t,
		cloudtest.WithServerOpts(
			grpc.ChainUnaryInterceptor(
				cloudtest.UnaryTunneledAuthnInterceptor(headerKey, requireToken),
			),
			grpc.ChainStreamInterceptor(
				cloudtest.StreamTunneledAuthnInterceptor(headerKey, requireToken),
			),
		),
	)
	defer stop()

	conn, err := grpc.NewClient(
		fmt.Sprintf("localhost:%d", port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStreamInterceptor(interceptor.ReverseCredentialsInterceptor(requireToken)),
	)
	require.NoError(t, err)

	tunnel := tunnelpb.NewTunnelServiceClient(conn)
	reverseTunnelServer := grpctunnel.NewReverseTunnelServer(tunnel)

	go func() {
		time.Sleep(100 * time.Millisecond)
		reverseTunnelServer.Stop()
	}()

	started, err := reverseTunnelServer.Serve(context.Background())
	require.NoError(t, err)
	require.True(t, started)
}

func TestIntegrationReverseCredentialsInterceptor_AuthenticateFailed(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	headerKey := "x-tunnel-token"
	requireToken := "test-token"

	port, stop := cloudtest.SetupTestCloudServer(
		t,
		cloudtest.WithServerOpts(
			grpc.ChainUnaryInterceptor(
				cloudtest.UnaryTunneledAuthnInterceptor(headerKey, requireToken),
			),
			grpc.ChainStreamInterceptor(
				cloudtest.StreamTunneledAuthnInterceptor(headerKey, requireToken),
			),
		),
	)
	defer stop()

	conn, err := grpc.NewClient(
		fmt.Sprintf("localhost:%d", port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStreamInterceptor(interceptor.ReverseCredentialsInterceptor("wrong-token")),
	)
	require.NoError(t, err)

	tunnel := tunnelpb.NewTunnelServiceClient(conn)
	reverseTunnelServer := grpctunnel.NewReverseTunnelServer(tunnel)

	go func() {
		time.Sleep(100 * time.Millisecond)
		stop()
	}()

	started, err := reverseTunnelServer.Serve(context.Background())
	require.ErrorIs(t, err, status.Error(codes.Unauthenticated, "invalid token"))
	require.True(t, started)
}
