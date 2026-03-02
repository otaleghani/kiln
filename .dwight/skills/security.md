# Security

## Overview
Security-first coding patterns. Assume all input is hostile and all output is visible.

## Guidelines
1. Never trust user input. Validate, sanitize, and escape at every boundary.
2. Use parameterized queries / prepared statements. Never concatenate strings into SQL/commands.
3. Never hardcode secrets, keys, or passwords. Use environment variables or secret managers.
4. Apply least privilege: functions, users, and services get only the permissions they need.
5. Escape output for its context (HTML-encode for HTML, shell-escape for commands, etc.).
6. Use constant-time comparison for secrets and tokens to prevent timing attacks.
7. Set secure defaults: HTTPS, secure cookies, restrictive CORS, minimal headers.
8. Hash passwords with bcrypt/scrypt/argon2. Never MD5/SHA for passwords. Never roll your own crypto.
9. Validate file paths — prevent directory traversal (../). Resolve to absolute and check prefix.
10. Log security events (auth failures, permission denials) but never log secrets or PII.

## Patterns
```
// GOOD: Parameterized query
const user = await db.query('SELECT * FROM users WHERE id = $1', [userId]);

// GOOD: Path traversal prevention
function safePath(userInput, baseDir) {
  const resolved = path.resolve(baseDir, userInput);
  if (!resolved.startsWith(baseDir)) throw new Error('Path traversal blocked');
  return resolved;
}

// GOOD: Output escaping
function renderComment(text) {
  return escapeHtml(text);  // prevent XSS
}
```

## Anti-Patterns
```
// BAD: SQL injection
db.query(`SELECT * FROM users WHERE name = '${userName}'`);

// BAD: Command injection
exec(`convert ${userFile} output.png`);

// BAD: Hardcoded secret
const API_KEY = "sk-abc123secret";
```
