package producer

import (
	"context"
)

type Producer interface {
	Connect(context.Context) error
	Close(context.Context) error
	Publish(context.Context, []byte) error
}
