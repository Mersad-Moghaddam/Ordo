package platform

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type CacheStore interface {
	Get(requestContext context.Context, cacheKey string) (string, bool, error)
	Set(requestContext context.Context, cacheKey string, cacheValue string, expiration time.Duration) error
}

type RateLimitStore interface {
	IncrementWindow(requestContext context.Context, limiterKey string, windowDuration time.Duration) (int64, error)
}

type CacheService struct {
	cacheStore CacheStore
}

func NewCacheService(cacheStore CacheStore) *CacheService {
	return &CacheService{cacheStore: cacheStore}
}

func ReadThroughCache[CachedDataType any](requestContext context.Context, cacheService *CacheService, cacheKey string, cacheTTL time.Duration, fetchFunction func() (CachedDataType, error)) (CachedDataType, bool, error) {
	cachedValue, hasValue, cacheError := cacheService.cacheStore.Get(requestContext, cacheKey)
	if cacheError != nil {
		return *new(CachedDataType), false, cacheError
	}
	if hasValue {
		var cachedData CachedDataType
		if unmarshalError := json.Unmarshal([]byte(cachedValue), &cachedData); unmarshalError == nil {
			return cachedData, true, nil
		}
	}
	fetchedData, fetchError := fetchFunction()
	if fetchError != nil {
		return *new(CachedDataType), false, fetchError
	}
	encodedBytes, marshalError := json.Marshal(fetchedData)
	if marshalError != nil {
		return *new(CachedDataType), false, marshalError
	}
	setError := cacheService.cacheStore.Set(requestContext, cacheKey, string(encodedBytes), cacheTTL)
	if setError != nil {
		return *new(CachedDataType), false, setError
	}
	return fetchedData, false, nil
}

type RateLimiter struct {
	rateLimitStore RateLimitStore
}

func NewRateLimiter(rateLimitStore RateLimitStore) *RateLimiter {
	return &RateLimiter{rateLimitStore: rateLimitStore}
}

func (rateLimiter *RateLimiter) Enforce(requestContext context.Context, subjectKey string, operationName string, requestLimit int64, windowDuration time.Duration) error {
	limiterKey := fmt.Sprintf("ratelimit:%s:%s", operationName, subjectKey)
	currentCount, incrementError := rateLimiter.rateLimitStore.IncrementWindow(requestContext, limiterKey, windowDuration)
	if incrementError != nil {
		return incrementError
	}
	if currentCount > requestLimit {
		return fmt.Errorf("rate limit exceeded")
	}
	return nil
}
