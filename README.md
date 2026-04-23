# Issa Compass Design Assessment

This file helps reviewers quickly find the proposed schema, design decisions, and diagrams.
This design is based on current understanding of the application and operations, with assumptions where details were not fully specified.
There may still be some misunderstanding on flow details or business logic interpretation from my side.
Design choices were made iteratively (not in a strict linear order) while refining tradeoffs across modules.

## Where to look for schema

All schema files are in:

- `main/user_account.sql`
- `main/service.sql`
- `main/ai_doc_review.sql`
- `main/conversation.sql`
- `main/audit_log.sql`
- `main/notification.sql`
- `main/payment.sql`

## Where to look for design decisions and tradeoffs

Decision writeups are in the `decisions/` folder:

- `decisions/user-account.md`
- `decisions/service.md`
- `decisions/ai-doc-review.md`
- `decisions/conversation.md`
- `decisions/audit-log.md`
- `decisions/notification.md`
- `decisions/payment.md`
- `decisions/agents-mcp.md` (future direction)

## Where to look for quick visual understanding

For architecture and module-level diagrams:

- `diagrams.md`

`diagrams.md` includes:

- system overview
- user account + RBAC
- service + documents
- AI document review
- conversation
- audit log
- notification
- payment
- agents + MCP (future)
