# Error Handling

## Overview
Patterns for robust, debuggable error handling. Errors should be explicit, contextual, and never silenced.

## Guidelines
1. Never swallow errors silently. Every catch must log, rethrow, or return an error value.
2. Include context: what was attempted, with what inputs, and why it failed.
3. Validate inputs at boundaries (function entry, API endpoints) — not deep in business logic.
4. Fail fast: detect problems early and surface them immediately.
5. Distinguish recoverable (retry, fallback) from fatal (crash with clear message).
6. Use typed/custom errors for different failure categories when the language supports it.
7. Always clean up resources in error paths (finally/defer/RAII/with).
8. Never expose internal details (stack traces, queries) to end users.
9. Prefer Result/Either/Option types over exceptions for expected failures.
10. Error messages are for humans. Error codes are for machines. Provide both when applicable.

## Patterns
```
// GOOD: Custom error with context
class ValidationError extends Error {
  constructor(field, value, constraint) {
    super(`Validation failed: ${field} value '${value}' violates ${constraint}`);
    this.field = field;
  }
}

// GOOD: Guard clause validation at boundary
function createUser(input) {
  if (!input.email?.includes('@')) {
    throw new ValidationError('email', input.email, 'must be valid email');
  }
  // ... business logic with validated data
}
```

## Anti-Patterns
```
// BAD: Silent catch
try { doSomething(); } catch (e) { /* ignore */ }

// BAD: Useless message
throw new Error("Something went wrong");

// BAD: Validation deep inside a call chain
function saveToDb(data) {
  // Too late to validate here — should have been caught at the API boundary
  if (!data.email) throw new Error("missing email");
}
```
