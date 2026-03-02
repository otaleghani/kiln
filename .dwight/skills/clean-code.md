# Clean Code

## Overview
Guidelines for producing clean, readable, maintainable code regardless of language.

## Guidelines
1. Functions should do ONE thing. If you can describe it using "and", split it.
2. Names should reveal intent. Avoid abbreviations. `getUserAccountBalance` > `getUsrAcctBal`.
3. Functions should have 3 or fewer parameters. Use an options/config object for more.
4. Avoid boolean flags as parameters — they signal the function does two things.
5. Early returns over nested conditionals. Guard clauses first.
6. No magic numbers or strings. Use named constants.
7. Comments explain WHY, never WHAT. If code needs a WHAT comment, refactor it.
8. Keep functions under 20 lines. If longer, extract helpers.
9. Prefer immutability. Use `const`/`final`/`let` over mutable bindings where possible.
10. One level of abstraction per function. Don't mix high-level orchestration with low-level details.

## Patterns
```
// GOOD: Early return, clear naming, single responsibility
function calculateMonthlyInterest(principal, annualRate) {
  if (principal <= 0) return 0;

  const MONTHS_PER_YEAR = 12;
  return principal * (annualRate / MONTHS_PER_YEAR);
}
```

## Anti-Patterns
```
// BAD: Nested conditionals, magic numbers, vague names
function proc(u) {
  if (u) {
    if (u.a) {
      return u.b * 0.05 / 12;
    }
  }
  return null;
}
```
