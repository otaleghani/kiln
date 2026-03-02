# Go API Design

## Overview
Guidelines for building production Go APIs: HTTP handlers, middleware, error handling, and project layout.

## Guidelines
1. Use `http.Handler` interface — keep handlers as methods on a server struct that holds dependencies.
2. Middleware chain: logging → recovery → auth → rate-limit → handler.
3. Always return structured errors with `{ "error": "message", "code": "NOT_FOUND" }`.
4. Use `context.Context` for cancellation, deadlines, and request-scoped values — never store business data in context.
5. Validate input at the handler boundary, not in business logic.
6. Return appropriate HTTP status codes: 400 for bad input, 401 for unauthed, 403 for forbidden, 404 for not found, 500 for internal errors.
7. Use `encoding/json` for simple cases; `json.NewEncoder(w).Encode()` for streaming.
8. Group routes by domain: `/api/v1/users/`, `/api/v1/orders/`.

## Patterns
```go
// Handler method on server struct
func (s *Server) GetUser(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    user, err := s.userStore.Get(r.Context(), id)
    if errors.Is(err, store.ErrNotFound) {
        httpError(w, http.StatusNotFound, "USER_NOT_FOUND", "user not found")
        return
    }
    if err != nil {
        httpError(w, http.StatusInternalServerError, "INTERNAL", "internal error")
        return
    }
    json.NewEncoder(w).Encode(user)
}
```

## Anti-Patterns
- Don't use global variables for database connections — inject via struct fields.
- Don't panic in handlers — use error returns and middleware recovery.
- Don't put business logic in handlers — handlers should only parse input, call services, format output.
