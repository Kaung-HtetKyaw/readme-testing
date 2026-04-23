# Audit Log Decisions

## Decision

- A single system-wide `audit_logs` table is used.
- Audit rows are append-only and treated as immutable from application workflows.
- Core fields include actor, action, target, correlation id, and `metadata` JSONB.
- Indexes are added for common query paths (`tenant`, actor, action, time).

## Rationale

- Cross-module action traceability stays centralized.
- Compliance, incident review, and debugging are supported with one canonical log.
- Correlation id enables end-to-end flow tracing.

## Tradeoffs

- V1 keeps a single table without partitioning.
- Some domain-level history remains in module tables in addition to audit logs.
- Query cost will rise over time without archival or partitioning.

## Future Plan (V2)

- Partition by time once volume thresholds are reached.
- Archive older partitions by retention policy.
- Add stricter role-based access policies for audit-log reads.
