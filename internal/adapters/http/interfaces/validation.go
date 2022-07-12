package interfaces

import (
	"context"
	"gitlab.com/g6834/team17/analytics-service/internal/adapters/http/dto"
	"net/http"
)

type MiddlewareValidator interface {
	Validate(next http.Handler) http.Handler
}

type JWTValidator interface {
	Validate(ctx context.Context, tokens dto.TokenPair) (dto.TokenPair, error)
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
}
