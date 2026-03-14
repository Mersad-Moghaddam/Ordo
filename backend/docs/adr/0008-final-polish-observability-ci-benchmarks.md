# ADR 0008: Final Polishing for Observability, CI, and Benchmarks

## Status
Accepted

## Context
Phase 7 requires release-readiness checks beyond feature delivery, including measurable performance signals and enforced quality gates in automation.

## Decision
Extend CI to run tests, revive linting, and benchmark target in one pipeline.

Add benchmark coverage for cache read-through and rate-limiter operations to track performance drift over time.

Keep worker and API observability endpoints represented in OpenAPI snippets for operational integration consistency.

## Consequences
Positive outcomes:
- Quality gates become reproducible in pull request automation.
- Benchmark baseline can be compared across commits to catch regressions.
- Operational teams receive consistent endpoint contracts for health/metrics.

Tradeoffs:
- CI duration increases due to benchmark execution.
- Benchmarks require controlled environment interpretation for strict comparisons.
