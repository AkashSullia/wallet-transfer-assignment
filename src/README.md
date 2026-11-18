## Wallet Transfer Assignment

### Approach Taken
- The system requirements were read and understood based on [ASSIGNMENT.md](../ASSIGNMENT.md).
- General idea of using clean architecture and technologies used were decided on. The technologies/libraries used along with the reasoning behind them are as follows:
  - `mise-en-place (mise)` - for environment variable management and task runner.
    reason: simplifies management of environment and automating tasks for building, testing etc.
  - `PostgreSQL` - a SQL database.
    reason: suggested by assignment but also being an industry standard would help greatly.
  - `fiber framework` - backend framework (with semantics similar to express from node.js).
    reason: simplify HTTP request/response handling.
  - `sqlc` -  to help compile SQL to type-safe code, have advantages of both raw SQL (speed, less complexity)
    and ORM (type-safety)
  - `goose` - database migration
- Used help of AI (`deepseek`) to:
  - refine the tasks assigned
  - research on various things: such as get a better understanding on transaction approach for double ledger entries.
  - understand how to handle sqlc in the context of clean architecture.
