# Performance

## Overview
Write code that's fast by default. Prefer algorithmic wins over micro-optimization.

## Guidelines
1. Choose the right data structure first. A hashmap lookup is O(1) vs O(n) linear scan.
2. Avoid N+1 queries: batch DB/API calls instead of looping individual requests.
3. Don't allocate in hot loops. Preallocate buffers, reuse objects where possible.
4. Lazy-load and defer: don't compute what isn't needed yet.
5. Cache expensive results, but always have an invalidation strategy.
6. Prefer streaming/iteration over loading everything into memory for large datasets.
7. String concatenation in loops is O(n²) in many languages. Use builders/join.
8. Profile before optimizing. Measure the actual bottleneck, don't guess.
9. Avoid unnecessary copies of large data. Use references, slices, or views.
10. Set reasonable timeouts and limits on all I/O operations.

## Patterns
```
// GOOD: Batch query instead of N+1
const userIds = orders.map(o => o.userId);
const users = await db.query('SELECT * FROM users WHERE id = ANY($1)', [userIds]);
const userMap = new Map(users.map(u => [u.id, u]));

// GOOD: Preallocate
const results = new Array(items.length);
for (let i = 0; i < items.length; i++) {
  results[i] = transform(items[i]);
}

// GOOD: Stream large files
const stream = fs.createReadStream(path);
for await (const chunk of stream) { process(chunk); }
```

## Anti-Patterns
```
// BAD: N+1 — one query per iteration
for (const order of orders) {
  const user = await db.query('SELECT * FROM users WHERE id = $1', [order.userId]);
}

// BAD: String concatenation in loop
let result = '';
for (const item of items) { result += item.toString(); }  // O(n²)
```
