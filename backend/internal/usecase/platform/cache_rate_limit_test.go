package platform

import (
	"context"
	"errors"
	"testing"
	"time"

	cacheredis "github.com/ordo/backend/internal/infrastructure/cache/redis"
)

type cachePayload struct {
	Name string `json:"name"`
}

func TestReadThroughCache(testingSuite *testing.T) {
	cacheStore := cacheredis.NewInMemoryCacheStore()
	cacheService := NewCacheService(cacheStore)
	fetchCallCount := 0

	firstValue, firstHit, firstError := ReadThroughCache(context.Background(), cacheService, "workspace:1", 5*time.Minute, func() (cachePayload, error) {
		fetchCallCount++
		return cachePayload{Name: "platform"}, nil
	})
	if firstError != nil || firstHit {
		testingSuite.Fatalf("expected first fetch miss without error")
	}
	if firstValue.Name != "platform" {
		testingSuite.Fatalf("unexpected first value")
	}

	secondValue, secondHit, secondError := ReadThroughCache(context.Background(), cacheService, "workspace:1", 5*time.Minute, func() (cachePayload, error) {
		fetchCallCount++
		return cachePayload{Name: "other"}, nil
	})
	if secondError != nil || !secondHit {
		testingSuite.Fatalf("expected second fetch hit without error")
	}
	if secondValue.Name != "platform" {
		testingSuite.Fatalf("unexpected cached second value")
	}
	if fetchCallCount != 1 {
		testingSuite.Fatalf("fetch should be called once")
	}
}

func TestReadThroughCacheFetchError(testingSuite *testing.T) {
	cacheStore := cacheredis.NewInMemoryCacheStore()
	cacheService := NewCacheService(cacheStore)
	_, _, readError := ReadThroughCache(context.Background(), cacheService, "workspace:2", 5*time.Minute, func() (cachePayload, error) {
		return cachePayload{}, errors.New("fetch failure")
	})
	if readError == nil {
		testingSuite.Fatalf("expected fetch failure")
	}
}

func TestRateLimiter(testingSuite *testing.T) {
	cacheStore := cacheredis.NewInMemoryCacheStore()
	rateLimiter := NewRateLimiter(cacheStore)
	requestContext := context.Background()

	firstError := rateLimiter.Enforce(requestContext, "user-1", "create-task", 2, time.Minute)
	secondError := rateLimiter.Enforce(requestContext, "user-1", "create-task", 2, time.Minute)
	thirdError := rateLimiter.Enforce(requestContext, "user-1", "create-task", 2, time.Minute)
	if firstError != nil || secondError != nil {
		testingSuite.Fatalf("first two requests should pass")
	}
	if thirdError == nil {
		testingSuite.Fatalf("third request should be rate limited")
	}
}
