package consumer

import (
	"context"
	"testing"
	"time"

	usecaseworker "github.com/ordo/backend/internal/usecase/worker"
)

func TestInMemoryStreamConsumerReadAckAdd(testingSuite *testing.T) {
	streamConsumer := NewInMemoryStreamConsumer()
	streamConsumer.Push("stream-name", usecaseworker.StreamMessage{MessageID: "message-1", Values: map[string]string{"eventType": "task.created"}})
	readMessages, readError := streamConsumer.ReadGroup(context.Background(), "stream-name", "group", "consumer", 10, 1*time.Second)
	if readError != nil {
		testingSuite.Fatalf("read error: %v", readError)
	}
	if len(readMessages) != 1 {
		testingSuite.Fatalf("expected one read message")
	}
	ackError := streamConsumer.Acknowledge(context.Background(), "stream-name", "group", "message-1")
	if ackError != nil {
		testingSuite.Fatalf("ack error: %v", ackError)
	}
	if !streamConsumer.IsAcknowledged("message-1") {
		testingSuite.Fatalf("message should be acknowledged")
	}
	_, addError := streamConsumer.Add(context.Background(), "stream-name", map[string]string{"eventId": "message-2", "eventType": "task.updated"})
	if addError != nil {
		testingSuite.Fatalf("add error: %v", addError)
	}
	readMessagesAfterAdd, secondReadError := streamConsumer.ReadGroup(context.Background(), "stream-name", "group", "consumer", 10, 1*time.Second)
	if secondReadError != nil || len(readMessagesAfterAdd) != 1 {
		testingSuite.Fatalf("expected one message after add")
	}
}
