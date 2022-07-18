package client

import (
	"context"
	"errors"
	"fmt"
	auth "gitlab.com/g6834/team17/analytics-service/api/auth_service"
	"gitlab.com/g6834/team17/analytics-service/internal/adapters/http/dto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type AuthClient struct {
	client *auth.AuthServiceClient
	conn   *grpc.ClientConn
}

const gRPCAddr = `localhost:8082`

var (
	ErrInvalidToken  = errors.New("invalid token")
	ErrExpiredToken  = errors.New("token expired")
	ErrUnknownStatus = errors.New("unknown status")
)

func NewGrpcAuth() AuthClient {
	return AuthClient{
		client: nil,
		conn:   nil,
	}
}

func (a *AuthClient) Connect(ctx context.Context) error {
	const ctxTimeOut = 10

	timeOutCtx, cancel := context.WithTimeout(ctx, time.Second*ctxTimeOut)
	defer cancel()

	conn, err := grpc.DialContext(timeOutCtx, gRPCAddr, grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("unable connect auth gRPC server: %w", err)
	}

	a.conn = conn
	client := auth.NewAuthServiceClient(conn)
	a.client = &client

	return nil
}

func (a *AuthClient) Disconnect(ctx context.Context) error {
	if a.conn != nil {
		return a.conn.Close()
	}

	return nil
}

func (a AuthClient) Validate(ctx context.Context, t dto.TokenPair) (dto.TokenPair, error) {
	var newTokens dto.TokenPair

	resp, err := (*a.client).Validate(ctx, &auth.ValidateTokenRequest{
		AccessToken:  t.Access,
		RefreshToken: t.Refresh,
	})
	if err != nil {
		return newTokens, err
	}

	switch resp.Status {
	case auth.Statuses_valid:
		newTokens.Access = resp.AccessToken
		newTokens.Refresh = resp.RefreshToken
		return newTokens, nil
	case auth.Statuses_invalid:
		return newTokens, ErrInvalidToken
	case auth.Statuses_expired:
		return newTokens, ErrExpiredToken
	default:
		return newTokens, ErrUnknownStatus
	}
}
