package interfaces

import "context"

type MessageConsumer interface {
	StartConsume(ctx context.Context) error
	StopConsume(ctx context.Context) error
}
