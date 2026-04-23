# Audit Log Policy (V1)

Audit table exists in `main/audit_log.sql`. This note defines operational policy so audit design is complete.

## Immutability

- `audit_logs` is append-only from application perspective.
- Application code must not update/delete existing audit rows.
- Any correction should be written as a new audit row (compensating event), not overwrite.

## Retention

- Keep hot audit data in primary PostgreSQL for 12 months.
- Keep long-term audit history in archive storage for 36 months (or compliance requirement of tenant).
- Retention period can be adjusted per legal/compliance policy later.

## Archive Strategy

- Archive old data by month (time-based batches).
- Move archived rows to cold storage / warehouse table.
- Keep indexable metadata in primary DB (tenant, action, target, created_at) for quick lookups.

## Partitioning Strategy

- V1: single table (current implementation) for simplicity.
- V2: monthly partitioning by `created_at` for `audit_logs` once volume grows.
- Drop/archive old partitions based on retention rules.

## Access Control

- Only privileged internal roles can query full audit logs.
- Client-facing roles cannot access cross-account audit records.
- Sensitive values in `metadata` should be masked/redacted in read APIs where needed.

## Correlation and Traceability

- Use `correlation_id` when available (request id, webhook event id, job run id).
- Allow null for legacy/manual paths, but new flows should provide it.

## Expected Write Sources

- API writes by staff/client actions
- Webhook ingestion handlers
- Async workers (AI/document processing)
- Scheduled jobs (SLA checks, lifecycle automation)

## Technology Fit (for this assessment stack)

- Application layer: writes audit rows at write points.
- Database: PostgreSQL is the primary audit store in V1.
- Queue/async processing: use existing worker queue/event bus used by document AI and webhook processing.
- Archive target (V2): object storage (for cold JSON/CSV export) or analytics warehouse table for long retention and reporting.
- Observability: request/job/webhook ids should be propagated as `correlation_id`.

## Implementation Approach

1. Add shared audit writer utility in backend:
   - Input: `tenant_id`, `actor_type`, `actor_user_id`, `action`, `target_type`, `target_id`, `metadata`, `correlation_id`
   - Output: append one row to `audit_logs`

2. Call audit writer from all mutation paths:
   - API handlers after successful write
   - Webhook processors after normalized updates
   - Workers/jobs after state transitions

3. Keep write consistency:
   - Prefer writing audit row in the same DB transaction as business mutation when possible
   - If same transaction is not possible, use outbox/event-forwarding pattern in V2

4. Enforce append-only behavior in app:
   - No update/delete code paths for `audit_logs`
   - Corrections must be new audit events

5. Add periodic retention/archive jobs:
   - V1: simple scheduled SQL job for age-based export/delete policy
   - V2: partition-aware archive and partition drop strategy
