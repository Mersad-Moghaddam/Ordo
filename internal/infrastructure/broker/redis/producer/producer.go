package producer

import (
	"context"
	"fmt"
	"math"
	"time"

	domaintask "github.com/ordo/backend/internal/domain/task"
	"github.com/ordo/backend/internal/infrastructure/broker/redis"
	repositorytask "github.com/ordo/backend/internal/repository/task"
)

type StreamPublisher interface {
	PublishToStream(requestContext context.Context, streamName string, values map[string]any) (string, error)
}

type OutboxProducer struct {
	outboxRepository repositorytask.OutboxRepository
	streamPublisher  StreamPublisher
	batchSize        int
	baseBackoff      time.Duration
	maxBackoff       time.Duration
	nowFunction      func() time.Time
}

type Option func(producer *OutboxProducer)

func NewOutboxProducer(outboxRepository repositorytask.OutboxRepository, streamPublisher StreamPublisher, options ...Option) *OutboxProducer {
	outboxProducer := &OutboxProducer{outboxRepository: outboxRepository, streamPublisher: streamPublisher, batchSize: 100, baseBackoff: 2 * time.Second, maxBackoff: 5 * time.Minute, nowFunction: time.Now}
	for _, option := range options {
		if option != nil {
			option(outboxProducer)
		}
	}
	return outboxProducer
}

func WithBatchSize(batchSize int) Option {
	return func(producer *OutboxProducer) {
		if batchSize > 0 {
			producer.batchSize = batchSize
		}
	}
}

func (outboxProducer *OutboxProducer) PublishPendingEvents(requestContext context.Context) error {
	pendingEventList, listError := outboxProducer.outboxRepository.ListPendingOutboxEvents(requestContext, outboxProducer.batchSize)
	if listError != nil {
		return fmt.Errorf("pending outbox list failure: %w", listError)
	}
	for _, pendingEvent := range pendingEventList {
		publishError := outboxProducer.publishEvent(requestContext, pendingEvent)
		if publishError != nil {
			return publishError
		}
	}
	return nil
}

func (outboxProducer *OutboxProducer) publishEvent(requestContext context.Context, pendingEvent domaintask.OutboxEvent) error {
	streamValues := map[string]any{"eventId": pendingEvent.EventID, "aggregateType": pendingEvent.AggregateType, "aggregateId": pendingEvent.AggregateID, "eventType": pendingEvent.EventType, "payload": pendingEvent.Payload, "idempotencyKey": pendingEvent.IdempotencyKey}
	_, publishError := outboxProducer.streamPublisher.PublishToStream(requestContext, redis.OrdoEventsStreamName, streamValues)
	if publishError != nil {
		attemptCount := pendingEvent.Attempts + 1
		nextRetryTime := outboxProducer.nowFunction().Add(outboxProducer.calculateBackoff(attemptCount))
		return outboxProducer.outboxRepository.MarkOutboxEventRetry(requestContext, pendingEvent.EventID, attemptCount, nextRetryTime.Unix())
	}
	return outboxProducer.outboxRepository.MarkOutboxEventPublished(requestContext, pendingEvent.EventID)
}

func (outboxProducer *OutboxProducer) calculateBackoff(attemptCount int) time.Duration {
	exponentialSeconds := float64(outboxProducer.baseBackoff) * math.Pow(2, float64(attemptCount-1))
	backoffDuration := time.Duration(exponentialSeconds)
	if backoffDuration > outboxProducer.maxBackoff {
		return outboxProducer.maxBackoff
	}
	return backoffDuration
}
