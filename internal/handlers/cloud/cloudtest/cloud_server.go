package cloudtest

import (
	"net"
	"testing"

	"github.com/jhump/grpctunnel"
	"github.com/jhump/grpctunnel/tunnelpb"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

type Options struct {
	serverOpts              []grpc.ServerOption
	tunnelServiceHandlerOpt grpctunnel.TunnelServiceHandlerOptions
}

var defaultOptions = Options{
	serverOpts:              []grpc.ServerOption{},
	tunnelServiceHandlerOpt: grpctunnel.TunnelServiceHandlerOptions{},
}

type OptionFunc func(opts *Options)

func WithServerOpts(opts ...grpc.ServerOption) OptionFunc {
	return func(o *Options) {
		o.serverOpts = opts
	}
}

func WithTunnelServiceHandlerOpts(opts grpctunnel.TunnelServiceHandlerOptions) OptionFunc {
	return func(o *Options) {
		o.tunnelServiceHandlerOpt = opts
	}
}

// SetupTestCloudServer starts a temporary gRPC server on a random port for testing purposes.
// It registers a TunnelServiceServer (used for reverse gRPC tunnels) with custom handler options,
// which can be customized via the optional OptionFuncs.
func SetupTestCloudServer(t *testing.T, optFuncs ...OptionFunc) (port int, stop func()) {
	t.Helper()

	opts := defaultOptions
	for _, optFunc := range optFuncs {
		optFunc(&opts)
	}

	server := grpc.NewServer(opts.serverOpts...)
	tunnelSvc := grpctunnel.NewTunnelServiceHandler(opts.tunnelServiceHandlerOpt)
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
