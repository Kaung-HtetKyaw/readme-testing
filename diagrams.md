# Diagrams (for reviewers)

Short visuals for the unified Postgres design. GitHub (and many editors) render Mermaid in preview.

---

## 1. High-level: who talks to what

Shows the main runtime boundaries: clients and staff use your API; Chatwoot and the payment provider push webhooks; async workers read/write Postgres.

```mermaid
flowchart TB
  subgraph clients["Clients"]
    MOBILE[Mobile app]
    WEB[Web app]
  end

  subgraph staff["Internal"]
    LEGAL[Legal dashboard]
    COMMS[Comms / Chatwoot UI]
  end

  subgraph platform["Your backend"]
    API[HTTP API]
    WH[Webhook handlers]
    JOBS[Workers / jobs]
  end

  subgraph external["External systems"]
    CW[Chatwoot]
    PAY[Payment provider]
    AI[AI doc review pipeline]
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
  API --> AI
  JOBS --> AI
```

---

## 2. Tenancy + identity + RBAC

`tenants` is the isolation root. Every `users` row belongs to one tenant. Role **assignments** are tenant-scoped: `user_roles` must match that user’s `(id, tenant_id)` (composite FK in `main/user_account.sql`). `roles` / `permissions` stay a global catalog in V1.

```mermaid
erDiagram
  tenants ||--o{ users : "has"
  tenants ||--o{ user_roles : "scopes assignment"
  users ||--o{ user_roles : "has roles"
  users ||--o{ password_auth : "optional"
  users ||--o{ oauth_auth : "optional"
  roles ||--o{ user_roles : "assigned"
  roles ||--o{ role_permissions : "grants"
  permissions ||--o{ role_permissions : ""

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
    timestamptz revoked_at
  }

  roles {
    uuid id PK
    text code
    text name
  }

  permissions {
    uuid id PK
    text code
    text resource
    text action
  }

  role_permissions {
    uuid role_id FK
    uuid permission_id FK
  }
```

---

## 3. Payments: tables and money flow

**Intent** lives on `payment_orders`. **One active charge path** goes through `payment_attempts` (idempotency: do not open a second “in flight” attempt for the same order). **Provider noise** lands in `payment_provider_events`. **Immutable ledger** for money movement is `payment_transactions`. **Refund workflow** is `payment_refunds` (status + provider ids), while ledger rows record the financial outcome.

```mermaid
flowchart LR
  PO[payment_orders]
  PA[payment_attempts]
  PPE[payment_provider_events]
  PT[payment_transactions]
  PR[payment_refunds]

  PO -->|1..n attempts over time| PA
  PA -->|ingest / dedupe| PPE
  PPE -->|derive facts| PT
  PO --> PR
  PA --> PR
  PR -->|ledger entries| PT
```

**Typical sequence** (happy path, simplified):

```mermaid
sequenceDiagram
  participant App as API
  participant PG as PostgreSQL
  participant Prov as Payment provider

  App->>PG: insert payment_order
  App->>PG: insert payment_attempt (open)
  App->>Prov: create payment intent
  Prov-->>App: provider_payment_intent_id
  App->>PG: update payment_attempt ids / status
  Prov->>App: webhook event
  App->>PG: insert payment_provider_event (idempotent)
  App->>PG: append payment_transactions
  App->>PG: update payment_order status / paid_at
```

---

## 4. Conversations: mirror + link to internal work

Chatwoot (or another provider) stays the **message UI source of truth**. Postgres stores a **provider-agnostic snapshot** (`conversations`) plus **CRM profile** tables and **links** to `service_applications` when the same person starts product work. Staff in the provider map to internal users via `staff_provider_agents`.

```mermaid
flowchart TB
  CW[Chatwoot conversation]

  subgraph pg["Postgres"]
    CONV[conversations]
    SPA[staff_provider_agents]
    CSA[conversation_service_applications]
    CRM[conversation_crm_profiles]
    CRMH[conversation_crm_profile_history]
    APP[service_applications]
    U[users]
  end

  CW -->|webhooks update snapshot| CONV
  CONV -->|optional link| U
  CONV --> CSA
  APP --> CSA
  CONV --> CRM
  CRM --> CRMH
  U --> SPA
  SPA -->|maps provider agent id| CONV
```

---

## 5. Core “service work” chain (documents + AI review)

How an application ties to submissions, files, and append-only AI history (names from `main/service.sql` / `main/ai_doc_review.sql`).

```mermaid
flowchart TB
  SV[service_versions]
  SA[service_applications]
  DS[document_submissions]
  DF[document_files]
  AIR[ai_doc_reviews]
  AIF[ai_doc_review_run_files]

  SV --> SA
  SA --> DS
  DS --> DF
  DS --> AIR
  AIR --> AIF
  DF --> AIF
```

---

## 6. Future: staff agents (bridge + Main API / gateway + MCP)

**Goal:** staff chat with **agents** (built-in like Assistant / Legal Review, or **custom** agents admins define) to get answers fast and run **allowed actions** per agent. Nothing bypasses tenancy, RBAC, or audit.

**Flow (conceptual):** UI → **MCP client service (bridge)** → **Main API** and/or **API Gateway** (scalable, normal platform path) → **Postgres**. Same bridge → **MCP protocol** → **Issa Compass MCP server** (first-party tools) and optionally **other third-party MCP servers** (separate allowlists per agent / tenant).

```mermaid
flowchart TB
  subgraph staff["Staff"]
    UI[UI application]
  end

  subgraph bridge_layer["MCP client service - bridge"]
    BRIDGE[Bridge orchestration]
  end

  subgraph api_layer["Domain access - scalable path"]
    GW[API Gateway optional]
    API[Main API]
  end

  subgraph mcp_layer["MCP servers"]
    MCP_IC[Issa Compass MCP server]
    MCP_3P[Third-party MCP servers]
  end

  PG[(PostgreSQL)]

  UI -->|staff product actions| GW
  UI -->|agent chat and tool runs| BRIDGE
  GW --> API
  BRIDGE -->|same RBAC and tenant rules| GW
  BRIDGE -->|or direct if small deploy| API
  BRIDGE -->|MCP| MCP_IC
  BRIDGE -->|MCP optional| MCP_3P
  API --> PG
```

**Agent types (who configures what):**

```mermaid
flowchart LR
  UI2[Staff UI]
  BI[Built-in agents - Assistant Legal Review etc]
  CA[Custom agents - admin-defined]
  BR2[MCP bridge]

  UI2 --> BI
  UI2 --> CA
  BI --> BR2
  CA --> BR2
```

**V1 schema:** no extra tables required for this roadmap; first cut can lean on `audit_logs` and add `agent_*` tables when product scope is fixed. Detail: `decisions/agents-mcp.md`.

---

## How to use this in the submission

- Point reviewers here from the main write-up (one line in email or README).
- You do not need more than these; extra ER diagrams for every table usually add noise unless the brief asks for full ERD.
