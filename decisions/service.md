# Service + Documents Decisions

## Decision

- The service domain is split into:
  - `services`
  - `service_versions`
  - `documents`
  - `service_documents`
  - `service_applications`
- `service_applications` stays generic; service-specific fields are placed in extension tables.
- DTV-specific fields are stored in `dtv_service_applications`.
- Document workflow is normalized through:
  - `document_submissions`
  - `document_files`
- Submission country is modeled via `submission_countries` (not hardcoded text).
- Workflow steps are append-only in `application_steps` for V1.

## Rationale

- Multiple concurrent services per user are supported cleanly.
- The model avoids one-table-per-service coupling.
- Document and workflow history remain scalable and auditable.

## Tradeoffs

- Step codes are application-managed in V1.
- `service_applications` does not directly store every service-specific field.
- `submission_countries` adds joins, but enables admin control and metadata.

## Future Plan (V2)

- Add `service_version_steps` if dynamic workflow definition is required.
- Add stronger transition validation in DB plus application logic.
- Add new extension tables only when service-specific complexity appears.
