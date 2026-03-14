package platform

import (
	"context"
	"fmt"
	"testing"
	"time"

	cacheredis "github.com/ordo/backend/internal/infrastructure/cache/redis"
)

type benchmarkPayload struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

func BenchmarkReadThroughCacheMissThenHit(benchmarkContext *testing.B) {
	cacheStore := cacheredis.NewInMemoryCacheStore()
	cacheService := NewCacheService(cacheStore)
	requestContext := context.Background()
	fetchFunction := func() (benchmarkPayload, error) {
		return benchmarkPayload{Identifier: "workspace-1", Name: "Platform"}, nil
	}
	benchmarkContext.ResetTimer()
	for iterationIndex := 0; iterationIndex < benchmarkContext.N; iterationIndex++ {
		cacheKey := fmt.Sprintf("benchmark-key-%d", iterationIndex%26)
		_, _, _ = ReadThroughCache(requestContext, cacheService, cacheKey, 10*time.Minute, fetchFunction)
	}
}

func BenchmarkRateLimiterEnforce(benchmarkContext *testing.B) {
	cacheStore := cacheredis.NewInMemoryCacheStore()
	rateLimiter := NewRateLimiter(cacheStore)
	requestContext := context.Background()
	benchmarkContext.ResetTimer()
	for iterationIndex := 0; iterationIndex < benchmarkContext.N; iterationIndex++ {
		_ = rateLimiter.Enforce(requestContext, "subject", "operation", 1000000, time.Minute)
	}
}
