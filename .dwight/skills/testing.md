# Testing Patterns

## Overview
Guidelines for writing tests that are readable, reliable, and actually catch bugs.

## Guidelines
1. Test names describe behavior, not implementation: `rejects_expired_tokens` not `test_validate_3`.
2. One assertion per test (or one logical assertion). Tests that check 5 things are 5 tests.
3. Follow Arrange-Act-Assert (AAA) structure with blank lines between sections.
4. Tests must be independent. No test should depend on another test running first.
5. Don't test implementation details — test behavior and outputs. Refactoring shouldn't break tests.
6. Use factories/builders for test data, not copy-pasted object literals.
7. Test edge cases explicitly: empty inputs, nulls, boundaries, off-by-one, concurrent access.
8. Mocks should be minimal. If you're mocking 5 things, the code has too many dependencies.
9. Negative tests matter: verify that invalid inputs are rejected, not just that valid ones work.
10. Flaky tests are worse than no tests. Fix or delete them immediately.

## Patterns
```
// GOOD: Clear name, AAA structure, focused assertion
test('rejects login with expired token', () => {
  const token = createToken({ expiresAt: yesterday() });

  const result = authenticate(token);

  expect(result.success).toBe(false);
  expect(result.reason).toBe('token_expired');
});

// GOOD: Factory for test data
function createUser(overrides = {}) {
  return { id: 'test-1', name: 'Test User', email: 'test@example.com', ...overrides };
}
```

## Anti-Patterns
```
// BAD: Vague name, tests implementation, multiple concerns
test('test user', () => {
  const u = { name: 'Bob', email: 'bob@test.com', role: 'admin' };
  const result = processUser(u);
  expect(result.internalId).toBeTruthy();     // testing implementation
  expect(result.name).toBe('Bob');             // different concern
  expect(db.users.length).toBe(1);             // side effect dependency
});
```
