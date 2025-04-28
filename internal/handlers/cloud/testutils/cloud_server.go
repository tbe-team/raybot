package testutils

import (
	"net"
	"testing"

	"github.com/jhump/grpctunnel"
	"github.com/jhump/grpctunnel/tunnelpb"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func StartCloudServer(t *testing.T, opt ...grpc.ServerOption) (port int, stop func()) {
	t.Helper()

	server := grpc.NewServer(opt...)
	tunnelSvc := grpctunnel.NewTunnelServiceHandler(grpctunnel.TunnelServiceHandlerOptions{})
	tunnelpb.RegisterTunnelServiceServer(server, tunnelSvc.Service())

	//nolint:gosec
	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	go func() {
		err = server.Serve(listener)
		require.NoError(t, err)
	}()

	addr, ok := listener.Addr().(*net.TCPAddr)
	require.True(t, ok)

	return addr.Port, func() {
		server.Stop()
		listener.Close()
	}
}
