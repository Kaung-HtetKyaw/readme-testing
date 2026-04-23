# Agents + MCP (Future Plan)

This roadmap sits outside the current Postgres core scope and targets higher staff productivity with controlled, auditable agent workflows.

Also yes, this is the classic "agentic AI flow like every app nowadays" chapter, but the rollout is intentionally read-first so production safety stays ahead of hype.

## Decision

- Two agent categories are planned:
  - **Built-in agents** (Assistant, Legal Review, and future product-defined types)
  - **Custom agents** (admin-defined configuration and guardrails)
- Staff interaction is agent-first in the UI: chat for quick answers, context retrieval, and approved actions by agent type.
- Architecture path is:
  - UI application -> MCP client service (bridge)
  - Bridge -> Main API and/or API Gateway for domain reads/writes
  - Bridge -> Issa Compass MCP server (first-party tools)
  - Bridge -> optional third-party MCP servers (allowlisted)
- The bridge remains the control point for tenant scope, RBAC checks, orchestration, and policy enforcement.
- Issa Compass MCP rollout starts with minimal, mostly read-only tools and may include low-risk notification helpers.

## Rationale

- Staff productivity improves by reducing context switching.
- Tool-grounded agents reduce unsupported responses.
- A single MCP integration layer can be reused across multiple UIs and workflows.
- API Gateway compatibility preserves scalability and operational consistency.

## Tradeoffs

- Strict policy enforcement is required to avoid over-permissive tool access.
- Reliability design must handle MCP or model failures safely, especially for writes.
- Custom-agent flexibility increases governance and configuration risk.
- Scope can expand quickly without clear tool boundaries.

## Future Plan (Phased)

1. Launch Bridge plus Issa Compass MCP with a minimal, mostly read-only tool set and optional low-risk notification helpers.
2. Enable built-in agents in the UI using only those low-risk tools.
3. Validate adoption, quality, and safety; then introduce write tools with idempotency and full audit.
4. Add custom agents and optional third-party MCP servers behind tenant-specific allowlists.

## Schema Note

- V1 can ship without new agent tables.
- When implementation starts, dedicated tables can be added for agent definitions, bridge or MCP configuration, sessions, and invocations, or the model can be extended from `audit_logs` with stable action naming.
