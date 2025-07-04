# Development Guidelines for bookid

## Architecture
Follow Ben Johnson's Standard Package Layout (see ai_docs/ben-johnson-standard-package-layout.md when making architectural decisions).
Packages are layers, not groups - dependencies flow inward toward the domain.

## Test-Driven Development (Required)
- Write failing tests FIRST - no exceptions
- Test through public APIs only (use `package_test` convention)
- External API contracts: use golden files pattern (see ai_docs/golden-files-testing-pattern.md when testing external APIs)
- Testing difficulties = design feedback opportunity
- ALWAYS use t.Parallel() in all tests and subtests to detect data races with -race flag

## Validation Workflow
Before ANY commit or when stuck:
```bash
make validate
```
This runs all checks: formatting, vetting, tests, go.mod tidy, and linting.

## Core Principles
1. Domain purity: Root package contains only types and interfaces
2. Dependencies isolated: Each external dependency gets its own package
3. Manual mocks: Simple, explicit mocks in mock/ package
4. No circular dependencies: Achieved through proper layering

When uncertain about correctness criteria, ASK before implementing.