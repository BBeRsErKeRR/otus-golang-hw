package consumer

import (
	"context"
	"fmt"
)

type ConsumerUseCase struct { //nolint:revive
	consumer Consumer
}

func (c *ConsumerUseCase) Consume(ctx context.Context, f func(ctx context.Context, msg []byte)) error {
	return c.consumer.Consume(ctx, f)
}

func (c *ConsumerUseCase) Connect(ctx context.Context) error {
	err := c.consumer.Connect(ctx)
	if err != nil {
		return fmt.Errorf("ConsumerUseCase - Connect - u.consumer.Connect: %w", err)
	}
	return nil
}

func (c *ConsumerUseCase) Close(ctx context.Context) error {
	err := c.consumer.Close(ctx)
	if err != nil {
		return fmt.Errorf("ConsumerUseCase - Connect - u.consumer.Close: %w", err)
	}
	return nil
}

func New(c Consumer) *ConsumerUseCase {
	return &ConsumerUseCase{
		consumer: c,
	}
}
