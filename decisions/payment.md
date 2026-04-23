# Payment Decisions

## Decision

- Payments are split into:
  - `payment_orders` (business payment intent)
  - `payment_attempts` (provider attempt lifecycle)
  - `payment_refunds` (refund request lifecycle)
  - `payment_transactions` (immutable money ledger)
  - `payment_provider_events` (webhook ingestion and idempotency)
- V1 idempotency uses one open attempt per order.
- Only provider references (intent, charge, refund, event IDs) are stored; card data is excluded.

## Rationale

- Retry, failure, and out-of-order webhook scenarios remain manageable.
- Mutable operational state is separated from immutable financial records.
- The schema stays aligned with PCI-aware boundaries.

## Tradeoffs

- More tables increase schema complexity.
- `payment_refunds` and `payment_transactions` look similar but serve different concerns (workflow vs ledger).
- V1 idempotency does not include a full generic idempotency-key subsystem.

## Future Plan (V2)

- Add stronger reconciliation jobs across provider events and ledger totals.
- Add partitioning and archival for high-volume transaction and event tables.
- Add richer chargeback and dispute workflows as payment operations expand.
