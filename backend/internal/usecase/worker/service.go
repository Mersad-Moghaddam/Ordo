package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/ordo/backend/internal/infrastructure/broker/redis"
)

type StreamMessage struct {
	MessageID string
	Values    map[string]string
}

type StreamConsumer interface {
	ReadGroup(requestContext context.Context, streamName string, groupName string, consumerName string, count int, blockTimeout time.Duration) ([]StreamMessage, error)
	Acknowledge(requestContext context.Context, streamName string, groupName string, messageID string) error
	Add(requestContext context.Context, streamName string, values map[string]string) (string, error)
}

type IdempotencyStore interface {
	HasProcessed(requestContext context.Context, idempotencyKey string) (bool, error)
	MarkProcessed(requestContext context.Context, idempotencyKey string, expiration time.Duration) error
}

type NotificationDispatcher interface {
	DispatchNotification(requestContext context.Context, eventType string, payload map[string]any) error
}

type Service struct {
	streamConsumer          StreamConsumer
	idempotencyStore        IdempotencyStore
	notificationDispatcher  NotificationDispatcher
	consumerName            string
	batchSize               int
	maxAttempts             int
	baseBackoff             time.Duration
	maxBackoff              time.Duration
	processedEventRetention time.Duration
	nowFunction             func() time.Time
}

type Option func(service *Service)

func NewService(streamConsumer StreamConsumer, idempotencyStore IdempotencyStore, notificationDispatcher NotificationDispatcher, options ...Option) *Service {
	service := &Service{streamConsumer: streamConsumer, idempotencyStore: idempotencyStore, notificationDispatcher: notificationDispatcher, consumerName: "worker-default", batchSize: 50, maxAttempts: 8, baseBackoff: 2 * time.Second, maxBackoff: 10 * time.Minute, processedEventRetention: 24 * time.Hour, nowFunction: time.Now}
	for _, option := range options {
		if option != nil {
			option(service)
		}
	}
	return service
}

func WithConsumerName(consumerName string) Option {
	return func(service *Service) {
		if consumerName != "" {
			service.consumerName = consumerName
		}
	}
}

func WithBatchSize(batchSize int) Option {
	return func(service *Service) {
		if batchSize > 0 {
			service.batchSize = batchSize
		}
	}
}

func (service *Service) PollAndProcess(requestContext context.Context) error {
	messageList, readError := service.streamConsumer.ReadGroup(requestContext, redis.OrdoEventsStreamName, redis.OrdoWorkersGroupName, service.consumerName, service.batchSize, 2*time.Second)
	if readError != nil {
		return fmt.Errorf("stream read failure: %w", readError)
	}
	for _, streamMessage := range messageList {
		if processError := service.processMessage(requestContext, streamMessage); processError != nil {
			return processError
		}
	}
	return nil
}

func (service *Service) processMessage(requestContext context.Context, streamMessage StreamMessage) error {
	idempotencyKey := streamMessage.Values["idempotencyKey"]
	if idempotencyKey == "" {
		idempotencyKey = streamMessage.MessageID
	}
	processed, processedError := service.idempotencyStore.HasProcessed(requestContext, idempotencyKey)
	if processedError != nil {
		return fmt.Errorf("idempotency lookup failure: %w", processedError)
	}
	if processed {
		return service.streamConsumer.Acknowledge(requestContext, redis.OrdoEventsStreamName, redis.OrdoWorkersGroupName, streamMessage.MessageID)
	}
	attemptCount := parseAttemptCount(streamMessage.Values["attempts"])
	payloadMap, payloadError := parsePayload(streamMessage.Values["payload"])
	if payloadError != nil {
		return service.failMessage(requestContext, streamMessage, attemptCount, payloadError.Error())
	}
	if dispatchError := service.notificationDispatcher.DispatchNotification(requestContext, streamMessage.Values["eventType"], payloadMap); dispatchError != nil {
		return service.failMessage(requestContext, streamMessage, attemptCount, dispatchError.Error())
	}
	if markError := service.idempotencyStore.MarkProcessed(requestContext, idempotencyKey, service.processedEventRetention); markError != nil {
		return fmt.Errorf("idempotency mark failure: %w", markError)
	}
	return service.streamConsumer.Acknowledge(requestContext, redis.OrdoEventsStreamName, redis.OrdoWorkersGroupName, streamMessage.MessageID)
}

func (service *Service) failMessage(requestContext context.Context, streamMessage StreamMessage, attemptCount int, failureReason string) error {
	nextAttemptCount := attemptCount + 1
	if nextAttemptCount >= service.maxAttempts {
		dlqValues := map[string]string{"originalMessageId": streamMessage.MessageID, "payload": streamMessage.Values["payload"], "eventType": streamMessage.Values["eventType"], "idempotencyKey": streamMessage.Values["idempotencyKey"], "failureReason": failureReason, "failedAt": service.nowFunction().Format(time.RFC3339Nano)}
		if _, addError := service.streamConsumer.Add(requestContext, redis.OrdoEventsDLQName, dlqValues); addError != nil {
			return fmt.Errorf("dlq publish failure: %w", addError)
		}
		return service.streamConsumer.Acknowledge(requestContext, redis.OrdoEventsStreamName, redis.OrdoWorkersGroupName, streamMessage.MessageID)
	}
	retryValues := map[string]string{"eventId": streamMessage.Values["eventId"], "aggregateType": streamMessage.Values["aggregateType"], "aggregateId": streamMessage.Values["aggregateId"], "eventType": streamMessage.Values["eventType"], "payload": streamMessage.Values["payload"], "idempotencyKey": streamMessage.Values["idempotencyKey"], "attempts": fmt.Sprintf("%d", nextAttemptCount), "nextRetryAt": service.nowFunction().Add(service.calculateBackoff(nextAttemptCount)).Format(time.RFC3339Nano)}
	if _, addError := service.streamConsumer.Add(requestContext, redis.OrdoEventsStreamName, retryValues); addError != nil {
		return fmt.Errorf("retry publish failure: %w", addError)
	}
	return service.streamConsumer.Acknowledge(requestContext, redis.OrdoEventsStreamName, redis.OrdoWorkersGroupName, streamMessage.MessageID)
}

func (service *Service) calculateBackoff(attemptCount int) time.Duration {
	exponentialValue := float64(service.baseBackoff) * math.Pow(2, float64(attemptCount-1))
	backoffDuration := time.Duration(exponentialValue)
	if backoffDuration > service.maxBackoff {
		return service.maxBackoff
	}
	return backoffDuration
}

func parseAttemptCount(attemptValue string) int {
	if attemptValue == "" {
		return 0
	}
	var parsedAttempt int
	_, _ = fmt.Sscanf(attemptValue, "%d", &parsedAttempt)
	if parsedAttempt < 0 {
		return 0
	}
	return parsedAttempt
}

func parsePayload(payloadValue string) (map[string]any, error) {
	if payloadValue == "" {
		return map[string]any{}, nil
	}
	payloadMap := make(map[string]any)
	if unmarshalError := json.Unmarshal([]byte(payloadValue), &payloadMap); unmarshalError != nil {
		return nil, unmarshalError
	}
	return payloadMap, nil
}
