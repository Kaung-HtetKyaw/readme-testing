# Diagrams (for reviewers)

Module-by-module visuals for the unified Postgres design.

---

## 1. System Overview

Flow summary: client and staff-facing apps call the Main API for core product actions. External systems (Chatwoot and payment provider) send webhooks into handlers, and both synchronous API paths and async workers persist state in Postgres. AI review runs are triggered by API or jobs and their outcomes are written back through the backend.

```mermaid
flowchart TB
  subgraph clients[Clients]
    MOBILE[Mobile app]
    WEB[Web app]
  end

  subgraph staff[Internal Staff]
    LEGAL[Legal dashboard]
    COMMS[Comms / Chatwoot UI]
  end

  subgraph backend[Platform Backend]
    API[Main API]
    WH[Webhook handlers]
    JOBS[Workers / jobs]
  end

  subgraph external[External Systems]
    CW[Chatwoot]
    PAY[Payment provider]
    AIPIPE[AI review pipeline]
  end

  PG[(PostgreSQL)]

  MOBILE --> API
  WEB --> API
  LEGAL --> API
  COMMS --> CW

  CW -->|webhooks| WH
  PAY -->|webhooks| WH
  API --> PG
  WH --> PG
  JOBS --> PG
  API --> AIPIPE
  JOBS --> AIPIPE
```

---

## 2. User Account + RBAC Module

Flow summary: tenants are the isolation root, users are created under a tenant, auth identities are attached per user, and role assignments are written in `user_roles`. Permission checks resolve from `user_roles` -> `roles` -> `role_permissions` -> `permissions` for authorization decisions in the app.

```mermaid
erDiagram
  tenants ||--o{ users : has
  tenants ||--o{ user_roles : scopes
  users ||--o{ password_auth : optional
  users ||--o{ oauth_auth : optional
  users ||--o{ user_roles : assigned
  roles ||--o{ user_roles : grants
  roles ||--o{ role_permissions : maps
  permissions ||--o{ role_permissions : maps

  tenants {
    uuid id PK
    text code
    text status
  }

  users {
    uuid id PK
    uuid tenant_id FK
    text email
    text user_type
    text auth_type
  }

  user_roles {
    uuid id PK
    uuid tenant_id FK
    uuid user_id FK
    uuid role_id FK
  }
```

---

## 3. Service + Documents Module

`dtv_service_applications` is separated from `service_applications` on purpose:

- `service_applications` stays generic for all service types.
- DTV-only fields (for example travel/arrival-specific fields) do not pollute the base table for every other service.
- New service-specific fields can be added without risky schema bloat on the shared table.
- Querying remains clean: common queries hit `service_applications`; DTV flows join `dtv_service_applications` only when needed.

Flow summary: service definitions are configured through `services`, `service_versions`, `documents`, and `service_documents`. A user starts work through `service_applications`, then uploads files via `document_submissions` and `document_files`, while process checkpoints are appended in `application_steps`.

```mermaid
flowchart TB
  S[services] --> SV[service_versions]
  D[documents] --> SD[service_documents]
  SV --> SD

  U[users] --> SA[service_applications]
  SV --> SA
  C[submission_countries] --> SA

  SA --> DTV[dtv_service_applications]
  SA --> DS[document_submissions]
  DS --> DF[document_files]
  SA --> ST[application_steps]
```

---

## 4. AI Document Review Module

Flow summary: when files are submitted for an application, a new AI review run is appended in `ai_doc_reviews`. Reviewed files for that run are linked in `ai_doc_review_run_files`, allowing many files per run and full retry history without overwriting prior outcomes.

```mermaid
flowchart LR
  SA[service_applications] --> DS[document_submissions]
  DS --> DF[document_files]

  SA --> AIR[ai_doc_reviews]
  DS --> AIR
  AIR --> AIRF[ai_doc_review_run_files]
  DF --> AIRF
```

---

## 5. Conversation Module

Flow summary: provider webhooks update the `conversations` snapshot, staff mapping is resolved through `staff_provider_agents`, and conversations can be linked to one or more service applications through `conversation_service_applications`. CRM qualifiers stay current in `conversation_crm_profiles` and every change is appended to `conversation_crm_profile_history`.

```mermaid
flowchart TB
  CW[Chatwoot / provider conversation]

  subgraph pg[Postgres]
    CONV[conversations]
    SPA[staff_provider_agents]
    CSA[conversation_service_applications]
    CRM[conversation_crm_profiles]
    CRMH[conversation_crm_profile_history]
    SA[service_applications]
    U[users]
  end

  CW -->|webhooks| CONV
  CONV --> U
  U --> SPA
  SPA --> CONV

  CONV --> CSA
  SA --> CSA

  CONV --> CRM
  CRM --> CRMH
```

---

## 6. Audit Log Module

Flow summary: each domain module emits audit entries into `audit_logs` for important actions. The audit stream is append-only and cross-module, so reviewers and support teams can trace what happened, who performed it, and when it happened using one canonical log path.

```mermaid
flowchart LR
  AUTH[User Account]
  SERVICE[Service]
  AI[AI Review]
  CONV[Conversation]
  NOTI[Notification]
  PAY[Payment]

  AUTH --> AUDIT[audit_logs]
  SERVICE --> AUDIT
  AI --> AUDIT
  CONV --> AUDIT
  NOTI --> AUDIT
  PAY --> AUDIT
```

---

## 7. Notification Module

Flow summary: domain events create notification intents in `notifications`, channel targets and user preferences shape delivery behavior, and each send/retry attempt is recorded in `notification_deliveries` for operational visibility and retry safety.

```mermaid
flowchart LR
  N[notifications] --> NT[notification_targets]
  N --> ND[notification_deliveries]
  U[users] --> NP[notification_preferences]
  U --> N

  EV[Domain events]
  EV --> N
  ND --> CH[Email / Push / SMS providers]
```

---

## 8. Payment Module

Flow summary: payment intent starts at `payment_orders`, execution attempts are tracked in `payment_attempts`, provider webhooks are ingested in `payment_provider_events`, and immutable money movement is appended to `payment_transactions`. Refund lifecycle is tracked in `payment_refunds` and also reflected in ledger rows.

```mermaid
flowchart LR
  PO[payment_orders]
  PA[payment_attempts]
  PR[payment_refunds]
  PPE[payment_provider_events]
  PT[payment_transactions]

  PO --> PA
  PO --> PR
  PA --> PR

  PA --> PPE
  PPE --> PT
  PR --> PT
```

**Operational sequence (happy path):**

```mermaid
sequenceDiagram
  participant API as Main API
  participant PG as PostgreSQL
  participant PSP as Payment provider

  API->>PG: create payment_order
  API->>PG: create payment_attempt (open)
  API->>PSP: create payment intent
  PSP-->>API: payment_intent_id
  API->>PG: update attempt refs/status
  PSP->>API: webhook event
  API->>PG: insert payment_provider_event (idempotent)
  API->>PG: append payment_transactions
  API->>PG: update payment_order status
```

---

## 9. Agents + MCP Future Module

Flow summary: staff interacts with agents from the UI, the bridge enforces tenant and RBAC context, then routes domain reads/writes through Main API or API Gateway and routes tool calls to Issa Compass MCP (and optional third-party MCP servers). Rollout starts read-first, then moves to guarded writes after validation.

```mermaid
flowchart TB
  UI["UI application"]
  BRIDGE["MCP client service - bridge"]
  GW["API Gateway optional"]
  API["Main API"]
  PG[(PostgreSQL)]

  MCPIC["Issa Compass MCP server"]
  MCP3P["Third-party MCP servers"]

  UI -->|staff actions| GW
  UI -->|agent chat + tool calls| BRIDGE
  GW --> API
  BRIDGE -->|domain reads and writes with tenant RBAC| GW
  BRIDGE -->|or direct path| API
  BRIDGE -->|MCP protocol| MCPIC
  BRIDGE -->|MCP protocol optional| MCP3P
  API --> PG
```

**Agent rollout:**

```mermaid
flowchart LR
  UI2["Staff UI"] --> BI["Built-in agents"]
  UI2 --> CA["Custom agents"]
  BI --> BR["MCP bridge"]
  CA --> BR
  BR --> READ["Read tools first"]
  READ --> WRITE["Write tools later after validation"]
```
