package interceptors

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/sreway/gophkeeper/internal/client/services"
)

const (
	authTokenKey = "authorization"
)

type AuthInterceptor struct {
	tokenProvider services.TokenProvider
}

func NewAuthInterceptor(tokenProvider services.TokenProvider) *AuthInterceptor {
	return &AuthInterceptor{
		tokenProvider: tokenProvider,
	}
}

func (auth *AuthInterceptor) Unary() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		if auth.tokenProvider.GetToken() != "" {
			ctx = metadata.AppendToOutgoingContext(ctx, authTokenKey, auth.tokenProvider.GetToken())
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func (auth *AuthInterceptor) Stream() grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		if auth.tokenProvider.GetToken() != "" {
			ctx = metadata.AppendToOutgoingContext(ctx, authTokenKey, auth.tokenProvider.GetToken())
		}
		return streamer(ctx, desc, cc, method, opts...)
	}
}
