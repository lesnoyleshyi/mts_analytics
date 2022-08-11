package interfaces

import (
	"context"
	ports "gitlab.com/g6834/team17/analytics-service/internal/ports/output"
)

type Storage interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	ports.EventStorage
}
