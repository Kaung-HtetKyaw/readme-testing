Right now, I assume that there will be no staff transfer from one tenant (Thailand) to another (Japan). So I make the email globally unique.

For submission country, I use `submission_countries` table instead of direct text code in `service_applications`. It adds one more join, but gives cleaner data, tenant-level country control, and place to keep pros/cons + metadata.

I keep DTV-only fields in `dtv_service_applications` (extension table) instead of putting them in `service_applications`. It makes queries slightly more complex for DTV, but keeps core application table generic for other services.

For conversation webhook ingestion, V1 will not store full raw webhook payload retention to reduce storage and keep ingestion simple. We only store normalized records needed by the application flow. If replay/debug/compliance needs grow, we can add raw-event retention later with TTL/partitioned archive.
