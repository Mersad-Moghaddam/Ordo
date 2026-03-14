package producer

import (
	"context"
	"errors"
	"testing"
	"time"

	domaintask "github.com/ordo/backend/internal/domain/task"
)

type mockOutboxRepository struct {
	pendingEvents     []domaintask.OutboxEvent
	publishedEventIDs []string
	retryEventIDs     []string
}

func (mockRepository *mockOutboxRepository) CreateOutboxEvent(requestContext context.Context, event domaintask.OutboxEvent) error {
	return nil
}

func (mockRepository *mockOutboxRepository) ListPendingOutboxEvents(requestContext context.Context, batchSize int) ([]domaintask.OutboxEvent, error) {
	return mockRepository.pendingEvents, nil
}

func (mockRepository *mockOutboxRepository) MarkOutboxEventPublished(requestContext context.Context, eventID string) error {
	mockRepository.publishedEventIDs = append(mockRepository.publishedEventIDs, eventID)
	return nil
}

func (mockRepository *mockOutboxRepository) MarkOutboxEventRetry(requestContext context.Context, eventID string, attempts int, nextRetryUnixTimestamp int64) error {
	mockRepository.retryEventIDs = append(mockRepository.retryEventIDs, eventID)
	return nil
}

type mockStreamPublisher struct {
	shouldFail bool
}

func (mockPublisher *mockStreamPublisher) PublishToStream(requestContext context.Context, streamName string, values map[string]any) (string, error) {
	if mockPublisher.shouldFail {
		return "", errors.New("publish failure")
	}
	return "stream-message-id", nil
}

func TestPublishPendingEvents(testingSuite *testing.T) {
	outboxRepository := &mockOutboxRepository{pendingEvents: []domaintask.OutboxEvent{{EventID: "event-success", AggregateType: "task", AggregateID: "task-id", EventType: "task.created", Payload: "{}", IdempotencyKey: "key-success", Attempts: 0}, {EventID: "event-retry", AggregateType: "task", AggregateID: "task-id", EventType: "task.updated", Payload: "{}", IdempotencyKey: "key-retry", Attempts: 1}}}
	streamPublisher := &mockStreamPublisher{}
	outboxProducer := NewOutboxProducer(outboxRepository, streamPublisher)

	publishError := outboxProducer.PublishPendingEvents(context.Background())
	if publishError != nil {
		testingSuite.Fatalf("publish pending events failure: %v", publishError)
	}
	if len(outboxRepository.publishedEventIDs) != 2 {
		testingSuite.Fatalf("expected published events count 2")
	}
}

func TestPublishRetryFlow(testingSuite *testing.T) {
	outboxRepository := &mockOutboxRepository{pendingEvents: []domaintask.OutboxEvent{{EventID: "event-retry", AggregateType: "task", AggregateID: "task-id", EventType: "task.updated", Payload: "{}", IdempotencyKey: "key-retry", Attempts: 3}}}
	streamPublisher := &mockStreamPublisher{shouldFail: true}
	outboxProducer := NewOutboxProducer(outboxRepository, streamPublisher)
	outboxProducer.nowFunction = func() time.Time { return time.Unix(1700000300, 0) }

	publishError := outboxProducer.PublishPendingEvents(context.Background())
	if publishError != nil {
		testingSuite.Fatalf("publish retry flow should persist retry without failure: %v", publishError)
	}
	if len(outboxRepository.retryEventIDs) != 1 {
		testingSuite.Fatalf("expected one retry mark")
	}
}
