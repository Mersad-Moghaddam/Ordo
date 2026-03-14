package redis

import (
	"context"
	"sync"
	"time"
)

type cachedEntry struct {
	cacheValue string
	expiresAt  time.Time
}

type inMemoryRateEntry struct {
	count     int64
	expiresAt time.Time
}

type InMemoryCacheStore struct {
	mutex        sync.Mutex
	cachedValues map[string]cachedEntry
	rateValues   map[string]inMemoryRateEntry
	nowFunction  func() time.Time
}

func NewInMemoryCacheStore() *InMemoryCacheStore {
	return &InMemoryCacheStore{cachedValues: map[string]cachedEntry{}, rateValues: map[string]inMemoryRateEntry{}, nowFunction: time.Now}
}

func (inMemoryCacheStore *InMemoryCacheStore) Get(requestContext context.Context, cacheKey string) (string, bool, error) {
	inMemoryCacheStore.mutex.Lock()
	defer inMemoryCacheStore.mutex.Unlock()
	entry, hasEntry := inMemoryCacheStore.cachedValues[cacheKey]
	if !hasEntry {
		return "", false, nil
	}
	if inMemoryCacheStore.nowFunction().After(entry.expiresAt) {
		delete(inMemoryCacheStore.cachedValues, cacheKey)
		return "", false, nil
	}
	return entry.cacheValue, true, nil
}

func (inMemoryCacheStore *InMemoryCacheStore) Set(requestContext context.Context, cacheKey string, cacheValue string, expiration time.Duration) error {
	inMemoryCacheStore.mutex.Lock()
	defer inMemoryCacheStore.mutex.Unlock()
	inMemoryCacheStore.cachedValues[cacheKey] = cachedEntry{cacheValue: cacheValue, expiresAt: inMemoryCacheStore.nowFunction().Add(expiration)}
	return nil
}

func (inMemoryCacheStore *InMemoryCacheStore) IncrementWindow(requestContext context.Context, limiterKey string, windowDuration time.Duration) (int64, error) {
	inMemoryCacheStore.mutex.Lock()
	defer inMemoryCacheStore.mutex.Unlock()
	entry, hasEntry := inMemoryCacheStore.rateValues[limiterKey]
	nowValue := inMemoryCacheStore.nowFunction()
	if !hasEntry || nowValue.After(entry.expiresAt) {
		inMemoryCacheStore.rateValues[limiterKey] = inMemoryRateEntry{count: 1, expiresAt: nowValue.Add(windowDuration)}
		return 1, nil
	}
	entry.count++
	inMemoryCacheStore.rateValues[limiterKey] = entry
	return entry.count, nil
}
