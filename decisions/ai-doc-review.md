# AI Document Review Decisions

## Decision

- AI reviews are stored in append-only `ai_doc_reviews`.
- Each review references:
  - `service_application_id`
  - `document_submission_id`
- Reviewed files per run are normalized in `ai_doc_review_run_files`.
- Review output stays focused on business-facing fields (`verdict`, `feedback`, timestamps).

## Rationale

- Re-upload and retry history remains explicit.
- The model avoids ambiguous latest-only state.
- Multi-file submissions are supported with relational integrity.

## Tradeoffs

- `run_status` is omitted in V1 because append-only history is sufficient.
- `model_name` and `model_version` are omitted in V1 to keep the schema lean.
- Join-table modeling adds complexity but improves querying and integrity over arrays.

## Future Plan (V2)

- Add model/version tracking if compliance or analytics requires it.
- Add richer execution telemetry if pipeline debugging demand increases.
- Add archival/partition strategy as review volume grows.
