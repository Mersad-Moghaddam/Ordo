package consumer

import (
	"context"
	"fmt"
	"sync"
	"time"

	usecaseworker "github.com/ordo/backend/internal/usecase/worker"
)

type InMemoryStreamConsumer struct {
	mutex           sync.Mutex
	streamMessages  map[string][]usecaseworker.StreamMessage
	acknowledgedIDs map[string]bool
}

func NewInMemoryStreamConsumer() *InMemoryStreamConsumer {
	return &InMemoryStreamConsumer{streamMessages: map[string][]usecaseworker.StreamMessage{}, acknowledgedIDs: map[string]bool{}}
}

func (inMemoryStreamConsumer *InMemoryStreamConsumer) ReadGroup(requestContext context.Context, streamName string, groupName string, consumerName string, count int, blockTimeout time.Duration) ([]usecaseworker.StreamMessage, error) {
	inMemoryStreamConsumer.mutex.Lock()
	defer inMemoryStreamConsumer.mutex.Unlock()
	messageList := inMemoryStreamConsumer.streamMessages[streamName]
	if len(messageList) == 0 {
		return []usecaseworker.StreamMessage{}, nil
	}
	if count <= 0 || count > len(messageList) {
		count = len(messageList)
	}
	resultList := make([]usecaseworker.StreamMessage, count)
	copy(resultList, messageList[:count])
	inMemoryStreamConsumer.streamMessages[streamName] = messageList[count:]
	return resultList, nil
}

func (inMemoryStreamConsumer *InMemoryStreamConsumer) Acknowledge(requestContext context.Context, streamName string, groupName string, messageID string) error {
	inMemoryStreamConsumer.mutex.Lock()
	defer inMemoryStreamConsumer.mutex.Unlock()
	inMemoryStreamConsumer.acknowledgedIDs[messageID] = true
	return nil
}

func (inMemoryStreamConsumer *InMemoryStreamConsumer) Add(requestContext context.Context, streamName string, values map[string]string) (string, error) {
	inMemoryStreamConsumer.mutex.Lock()
	defer inMemoryStreamConsumer.mutex.Unlock()
	messageID := values["eventId"]
	if messageID == "" {
		messageID = fmt.Sprintf("generated-%d", time.Now().UnixNano())
	}
	inMemoryStreamConsumer.streamMessages[streamName] = append(inMemoryStreamConsumer.streamMessages[streamName], usecaseworker.StreamMessage{MessageID: messageID, Values: values})
	return messageID, nil
}

func (inMemoryStreamConsumer *InMemoryStreamConsumer) Push(streamName string, message usecaseworker.StreamMessage) {
	inMemoryStreamConsumer.mutex.Lock()
	defer inMemoryStreamConsumer.mutex.Unlock()
	inMemoryStreamConsumer.streamMessages[streamName] = append(inMemoryStreamConsumer.streamMessages[streamName], message)
}

func (inMemoryStreamConsumer *InMemoryStreamConsumer) IsAcknowledged(messageID string) bool {
	inMemoryStreamConsumer.mutex.Lock()
	defer inMemoryStreamConsumer.mutex.Unlock()
	return inMemoryStreamConsumer.acknowledgedIDs[messageID]
}
