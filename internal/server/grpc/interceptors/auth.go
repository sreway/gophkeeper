package interceptors

import (
	"context"

	grpcauth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/sreway/gophkeeper/internal/server/services/auth"
)

const (
	authTokenKey = "authorization"
	userIDKey    = "userID"
	userEmailKey = "userEmail"
)

type AuthInterceptor struct {
	jwtManager auth.JWTManager
}

func (a *AuthInterceptor) ValidateAuthToken(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md[authTokenKey]
	if len(values) == 0 {
		return ctx, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]
	claims, err := a.jwtManager.VerifyToken(accessToken)
	if err != nil {
		return ctx, status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	ctx = metadata.AppendToOutgoingContext(ctx, userIDKey, claims.UserID.String(), userEmailKey, claims.Email)
	return ctx, nil
}

func (a *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return grpcauth.UnaryServerInterceptor(a.ValidateAuthToken)
}

func (a *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return grpcauth.StreamServerInterceptor(a.ValidateAuthToken)
}

func NewAuthInterceptor(jwtManager auth.JWTManager) *AuthInterceptor {
	return &AuthInterceptor{jwtManager: jwtManager}
}
