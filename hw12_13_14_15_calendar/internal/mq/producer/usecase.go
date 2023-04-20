package producer

import (
	"context"
)

type ProducerUseCase struct { //nolint:revive
	producer Producer
}

func (p *ProducerUseCase) Publish(ctx context.Context, data []byte) error {
	return p.producer.Publish(ctx, data)
}

func New(producer Producer) *ProducerUseCase {
	return &ProducerUseCase{
		producer: producer,
	}
}
