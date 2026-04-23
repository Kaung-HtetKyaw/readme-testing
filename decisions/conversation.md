# Conversation Decisions

## Decision

- `conversations` is used as a provider-agnostic current snapshot table.
- Provider identifiers (`provider`, `provider_conversation_id`, `provider_contact_id`) are stored for mapping.
- Lifecycle and assignee state are maintained in the main conversation row.
- Provider-agent to staff mapping uses `staff_provider_agents`.
- Conversation-to-application many-to-many links use `conversation_service_applications`.
- CRM lead data is split into:
  - `conversation_crm_profiles` (current profile)
  - `conversation_crm_profile_history` (append history)

## Rationale

- Dashboard access stays fast with snapshot-first querying.
- Provider-agnostic naming reduces lock-in risk.
- One conversation can support multiple service applications over time.

## Tradeoffs

- V1 is snapshot-first and does not persist full `conversation_events`.
- Raw webhook payload retention is not permanent in V1.
- Some SLA/activity values are derived from events and logs instead of denormalized columns.

## Future Plan (V2)

- Add `conversation_events` when replay/debug/audit depth becomes mandatory.
- Add stronger AI-state automation flows as operations mature.
- Add more composite indexes based on real dashboard filters.
