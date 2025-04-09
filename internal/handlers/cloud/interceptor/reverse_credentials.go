package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	TunnelTokenKey   = "x-tunnel-token"
	OpenTunnelMethod = "/grpctunnel.v1.TunnelService/OpenReverseTunnel"
)

func ReverseCredentialsInterceptor(token string) grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		if method != OpenTunnelMethod {
			return streamer(ctx, desc, cc, method, opts...)
		}

		ctx = metadata.AppendToOutgoingContext(ctx, TunnelTokenKey, token)

		// Invoke RPC, open the tunnel
		stream, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			return nil, err
		}

		_, err = stream.Header()
		if err != nil {
			return nil, err
		}
		// We can handle error code from the metadata here

		// Return the stream
		return stream, nil
	}
}
