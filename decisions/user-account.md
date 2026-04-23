# User Account + Access Control Decisions

## Decision

- The schema uses a `tenants` table and carries `tenant_id` across core identity and access tables.
- Identity remains centralized in a single `users` table.
- Authentication supports both password and OAuth login.
- Authentication data is split into:
  - `password_auth` for password users
  - `oauth_auth` for OAuth users
- RBAC uses:
  - `roles`
  - `permissions`
  - `role_permissions`
  - `user_roles` (tenant-scoped assignments)

## Rationale

- This structure meets multi-tenant requirements.
- Access control remains explicit and query-friendly.
- Authentication and authorization stay separated while still simple to operate.

## Tradeoffs

- Password/OAuth identity consistency is enforced in application logic in V1.
- `users.email` remains globally unique in V1.
- `roles` and `permissions` remain global catalogs; only assignment is tenant-scoped in V1.

## Future Plan (V2)

- Add finer-grained scope above tenant level (team, region, application) if needed.
- Add stricter DB-level auth consistency checks if complexity grows.
- Revisit global email uniqueness if cross-tenant staffing becomes common.
