For service_account, should we adopt like doc_review_token ? so that normal user account with enough privileges cannot do ai doc review

Question: How does AI review docs ? Is this an agent that can query ? Is there Issa Compass MCP Server and has prompt to respond in certian response structure ???
Assumption: Whatever ai services the applicatio used, we should have a adapter layer.
It could be third party AI service
It could be Issa Compass MCP Server. There is a service that has MCP clients lives. When user upload(re-uploads) docs, an endpoint to that service will be hit and get a response back which will be stored as run records

Question:
How will chatwoot agent id and application staff id be linked ?
Assumption:
staff account is created in the system
admin (or onboarding flow) links that staff to Chatwoot agent ID

Question:
Where does AI policy layer live?
A third party service ? An MCP server ? or Chatwoot handles it ?
Assumption: N/A

Choice (SLA / lifecycle automation): For V1, do not add `last_client_activity_at` / `last_staff_activity_at` columns yet. Scheduled jobs will derive activity from existing operational tables + audit logs.
Tradeoff: This keeps schema simpler and avoids duplicated data now, but job queries will be heavier (more joins/aggregation) and business rules for "activity" can be harder to maintain.
V2 plan: If complexity or volume grows, add materialized activity fields and optional controls (like `lifecycle_automation_disabled`, reminder cooldown timestamps) to simplify and speed up SLA jobs.
