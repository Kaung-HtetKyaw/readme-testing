# Notification Decisions

## Decision

- Notification responsibilities are split across:
  - `notifications` (intent and lifecycle state)
  - `notification_targets` (channel destination)
  - `notification_preferences` (opt-in/opt-out)
  - `notification_deliveries` (attempts and retries)
- `resource_type` plus `resource_id` is used for deep-link context.
- Delivery state transitions are stored explicitly.

## Rationale

- The design remains reusable across modules.
- Multi-channel delivery can scale without schema redesign.
- Retry/failure behavior remains auditable.

## Tradeoffs

- V1 skips heavy template/version modeling in the database.
- Generic resource references require rendering logic in the application.
- Clear table separation increases join complexity.

## Future Plan (V2)

- Add template version tracking if governance needs increase.
- Add fallback policies across channels (email to push to SMS).
- Add analytics tables for engagement and delivery quality.
