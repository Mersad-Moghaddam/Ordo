package redis

import (
	"context"
	"testing"
	"time"
)

func TestInMemoryCacheStoreSetGetAndIncrement(testingSuite *testing.T) {
	cacheStore := NewInMemoryCacheStore()
	setError := cacheStore.Set(context.Background(), "key", "value", time.Minute)
	if setError != nil {
		testingSuite.Fatalf("set error: %v", setError)
	}
	cacheValue, hasValue, getError := cacheStore.Get(context.Background(), "key")
	if getError != nil || !hasValue || cacheValue != "value" {
		testingSuite.Fatalf("unexpected get result")
	}
	firstCount, firstError := cacheStore.IncrementWindow(context.Background(), "limiter", time.Minute)
	secondCount, secondError := cacheStore.IncrementWindow(context.Background(), "limiter", time.Minute)
	if firstError != nil || secondError != nil || firstCount != 1 || secondCount != 2 {
		testingSuite.Fatalf("unexpected increment counts")
	}
}
