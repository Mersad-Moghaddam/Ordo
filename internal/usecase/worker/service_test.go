package worker

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ordo/backend/internal/infrastructure/broker/redis"
)

type mockStreamConsumer struct {
	messagesByStream  map[string][]StreamMessage
	acknowledgedIDs   []string
	additionsByStream map[string][]map[string]string
}

func (mockConsumer *mockStreamConsumer) ReadGroup(requestContext context.Context, streamName string, groupName string, consumerName string, count int, blockTimeout time.Duration) ([]StreamMessage, error) {
	messageList := mockConsumer.messagesByStream[streamName]
	if len(messageList) == 0 {
		return []StreamMessage{}, nil
	}
	mockConsumer.messagesByStream[streamName] = []StreamMessage{}
	return messageList, nil
}

func (mockConsumer *mockStreamConsumer) Acknowledge(requestContext context.Context, streamName string, groupName string, messageID string) error {
	mockConsumer.acknowledgedIDs = append(mockConsumer.acknowledgedIDs, messageID)
	return nil
}

func (mockConsumer *mockStreamConsumer) Add(requestContext context.Context, streamName string, values map[string]string) (string, error) {
	mockConsumer.additionsByStream[streamName] = append(mockConsumer.additionsByStream[streamName], values)
	return "added-id", nil
}

type mockIdempotencyStore struct {
	processed map[string]bool
}

func (mockStore *mockIdempotencyStore) HasProcessed(requestContext context.Context, idempotencyKey string) (bool, error) {
	return mockStore.processed[idempotencyKey], nil
}

func (mockStore *mockIdempotencyStore) MarkProcessed(requestContext context.Context, idempotencyKey string, expiration time.Duration) error {
	mockStore.processed[idempotencyKey] = true
	return nil
}

type mockNotificationDispatcher struct {
	shouldFail bool
}

func (mockDispatcher *mockNotificationDispatcher) DispatchNotification(requestContext context.Context, eventType string, payload map[string]any) error {
	if mockDispatcher.shouldFail {
		return errors.New("notification failure")
	}
	return nil
}

func TestPollAndProcessSuccess(testingSuite *testing.T) {
	streamConsumer := &mockStreamConsumer{messagesByStream: map[string][]StreamMessage{redis.OrdoEventsStreamName: {{MessageID: "message-id", Values: map[string]string{"eventType": "task.created", "payload": `{"taskId":"task-id"}`, "idempotencyKey": "key-1"}}}}, additionsByStream: map[string][]map[string]string{}}
	idempotencyStore := &mockIdempotencyStore{processed: map[string]bool{}}
	notificationDispatcher := &mockNotificationDispatcher{}
	workerService := NewService(streamConsumer, idempotencyStore, notificationDispatcher, WithConsumerName("consumer-1"), WithBatchSize(10))

	processError := workerService.PollAndProcess(context.Background())
	if processError != nil {
		testingSuite.Fatalf("process error: %v", processError)
	}
	if len(streamConsumer.acknowledgedIDs) != 1 {
		testingSuite.Fatalf("expected one acknowledgment")
	}
	if !idempotencyStore.processed["key-1"] {
		testingSuite.Fatalf("expected idempotency mark")
	}
}

func TestPollAndProcessRetryAndDLQ(testingSuite *testing.T) {
	streamConsumer := &mockStreamConsumer{messagesByStream: map[string][]StreamMessage{redis.OrdoEventsStreamName: {{MessageID: "message-retry", Values: map[string]string{"eventType": "comment.created", "payload": `{"commentId":"comment-id"}`, "idempotencyKey": "key-retry", "attempts": "1"}}, {MessageID: "message-dlq", Values: map[string]string{"eventType": "comment.deleted", "payload": `{"commentId":"comment-id"}`, "idempotencyKey": "key-dlq", "attempts": "8"}}}}, additionsByStream: map[string][]map[string]string{}}
	idempotencyStore := &mockIdempotencyStore{processed: map[string]bool{}}
	notificationDispatcher := &mockNotificationDispatcher{shouldFail: true}
	workerService := NewService(streamConsumer, idempotencyStore, notificationDispatcher)
	workerService.nowFunction = func() time.Time { return time.Unix(1700000400, 0) }

	processError := workerService.PollAndProcess(context.Background())
	if processError != nil {
		testingSuite.Fatalf("process error: %v", processError)
	}
	if len(streamConsumer.additionsByStream[redis.OrdoEventsStreamName]) != 1 {
		testingSuite.Fatalf("expected one retry stream addition")
	}
	if len(streamConsumer.additionsByStream[redis.OrdoEventsDLQName]) != 1 {
		testingSuite.Fatalf("expected one dlq addition")
	}
	if len(streamConsumer.acknowledgedIDs) != 2 {
		testingSuite.Fatalf("expected two acknowledgments")
	}
}
