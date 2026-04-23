# Diagrams (for reviewers)

Module-by-module visuals for the unified Postgres design.

---

## 1. System Overview

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

```mermaid
flowchart TB
  S[services] --> SV[service_versions]
  D[documents] --> SD[service_documents]
  SV --> SD

  U[users] --> SA[service_applications]
  SV --> SA
  C[submisssion_countries] --> SA

  SA --> DTV[dtv_service_applications]
  SA --> DS[document_submissions]
  DS --> DF[document_files]
  SA --> ST[application_steps]
```

---

## 4. AI Document Review Module

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

```mermaid
flowchart TB
  UI[UI application]
  BRIDGE[MCP client service (bridge)]
  GW[API Gateway optional]
  API[Main API]
  PG[(PostgreSQL)]

  MCPIC[Issa Compass MCP server]
  MCP3P[Third-party MCP servers]

  UI -->|staff actions| GW
  UI -->|agent chat + tool calls| BRIDGE
  GW --> API
  BRIDGE -->|domain reads/writes with tenant+RBAC| GW
  BRIDGE -->|or direct path| API
  BRIDGE -->|MCP protocol| MCPIC
  BRIDGE -->|MCP protocol optional| MCP3P
  API --> PG
```

**Agent rollout:**

```mermaid
flowchart LR
  UI2[Staff UI] --> BI[Built-in agents]
  UI2 --> CA[Custom agents]
  BI --> BR[MCP bridge]
  CA --> BR
  BR --> READ[Read tools first]
  READ --> WRITE[Write tools later after validation]
```

---

## How to use this in submission

- Link this file from the main write-up.
- Keep this as the primary visual reference; avoid duplicating many extra diagrams unless explicitly requested.
