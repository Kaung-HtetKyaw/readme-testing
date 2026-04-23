For conversation module, V1 will store only normalized current state in `conversations` (snapshot style). I will not add `conversation_events` yet to keep storage and implementation smaller.

I treat future `audit_log` as generic cross-domain log. If we need conversation-specific lifecycle timeline (assignment/handoff/reopen/close sequence), then I will add `conversation_events` in V2.

Main tradeoff: snapshot-only is simpler now, but it cannot show full routing sequence history. If SLA analytics or replay/debug become important later, `conversation_events` will be added as append-only timeline table.

Question:
What does this `is_waiting_for_legal_review` mean ?
Isn't legal review handled by the application flow instead of by conversation (messaging provider) flow ?

Things to keep in mind:
If user keep coming back for more services via the same messaging platform with the same account, the conversation id will never change
And one conversation will be tied to a lot of conversation_service_applications
And

Conversation Table
id (UUID PK)

tenant_id (FK)

provider (TEXT)

provider_conversation_id (TEXT)

provider_contact_id (TEXT)

channel_id (TEXT)

channel_name (TEXT)

channel_source (TEXT)

channel_meta (JSONB, optional provider-specific extras)

assignee_id (UUID FK -> users.id, nullable)

provider_assignee_id (TEXT, nullable)

lifecycle (TEXT: open/pending/resolved/closed)

contact_status (TEXT, nullable)

opened_at (TIMESTAMPTZ, nullable)

closed_at (TIMESTAMPTZ, nullable)

user_id (UUID FK -> users.id, nullable)

// INSTEAD WE WILL HAVE A JOINT TABLE

<!-- service_application_id (UUID FK -> service_applications.id, nullable) -->

created_at (TIMESTAMPTZ)

updated_at (TIMESTAMPTZ)

conversation_service_applications
id
conversation_id
service_application_id
created_at (TIMESTAMPTZ)
updated_at (TIMESTAMPTZ)

conversation_crm_profiles
id
tenant_id
conversation_id FK
location
nationality
current_visa
purpose
package
submission_country
interested_in
urgency
client_buying_segment // for example: budget_sensitive, premium, family, student, business_owner.

client_buying_intent // low, medium, high
current_step // new_inquiry, waiting_for_client_response, etc... ??
source (agent, ai, system)
updated_by_user_id nullable
updated_at

// This Table is to keep history of `what changed, from what to what, by whom, and when?`
Why useful:
no overwrite loss when fields change repeatedly because conversation is long lived
explainability for sales/comms decisions
supports analytics (how often intent changes, how fast urgency escalates)
clean trail even if snapshot table only stores latest values

conversation_crm_profile_history
id UUID PRIMARY KEY
tenant_id UUID NOT NULL (FK -> tenants.id)
conversation_id UUID NOT NULL (FK -> conversations.id)
field_name TEXT NOT NULL
(example: client_interested_in, client_urgency, submission_country)
value TEXT
(latest value recorded at that change point; nullable if clearing value)
source TEXT NOT NULL
(agent, ai, system, webhook)
created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
created_by
updated_by
