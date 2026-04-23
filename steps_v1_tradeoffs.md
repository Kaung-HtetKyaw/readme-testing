# Steps Design: V1 (Code List + Append API)

## V1 Model

For each service workflow, we keep a known list of step codes in app logic.

Example:

```text
["begin", "upload", "staff_review", "request_changes", "approved"]
```

Each time state changes for an application, we append a step event.

Pseudocode:

```text
appendStep(service_application_id, code, by)
```

Where:

- `service_application_id` identifies the active application
- `code` is one of the known step codes
- `by` identifies actor (`client`, `staff`, or `ai`)

## Runtime Table (Append-First)

`application_steps` stores event-style rows.

- `id`
- `tenant_id`
- `service_application_id` (FK)
- `code`
- `by_actor_type` (`client`, `staff`, `ai`, `system`)
- `by_user_id` (nullable FK for machine/system events)
- `created_at`
- `metadata` (optional JSONB for details)

This allows repeated steps, such as multiple `request_changes` cycles.

## Example Events

For one application:

1. `appendStep(app_1, "begin", "client")`
2. `appendStep(app_1, "upload", "client")`
3. `appendStep(app_1, "staff_review", "staff")`
4. `appendStep(app_1, "request_changes", "staff")`
5. `appendStep(app_1, "upload", "client")`
6. `appendStep(app_1, "staff_review", "staff")`
7. `appendStep(app_1, "approved", "staff")`

## Tradeoffs

- Fast and simple for V1 delivery.
- Easy to reason about event timeline.
- Step codes are app-managed constants, so code discipline is required.
- Harder to let non-engineering teams change workflows without deploy.

## V2 Refinement (Later)

Move step definitions to DB config (for example `service_version_steps`) while keeping append events in `application_steps`.

This enables per-service-version configurable workflows without changing the append API pattern.
