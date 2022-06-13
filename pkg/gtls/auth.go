package gtls

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

type Authentication struct {
	User     string
	Password string
}

func (a *Authentication) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{"user": a.User, "password": a.Password}, nil
}

func (a *Authentication) RequireTransportSecurity() bool {
	return true
}

func (a *Authentication) Auth(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return fmt.Errorf("missing credentials")
	}
	var appId, appKey string
	if val, ok := md["user"]; ok {
		appId = val[0]
	}
	if val, ok := md["password"]; ok {
		appKey = val[0]
	}

	if appId != a.User || appKey != a.Password {
		return grpc.Errorf(codes.Unauthenticated, "invalid token")
	}
	return nil
}
