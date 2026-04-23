BEGIN;


DROP MATERIALIZED VIEW cluster_resource_count;
DROP TABLE IF EXISTS deprecated_api;


CREATE TABLE IF NOT EXISTS deprecated_api_group (
    cluster_id UUID PRIMARY KEY REFERENCES cluster(id),
    deprecated_apis JSONB
);


CREATE MATERIALIZED VIEW cluster_resource_count AS
WITH 
count_node AS (
  SELECT 
    cluster_id,
    COUNT(id) AS total_node
  FROM 
    object
  WHERE
    kind = 'Node'
  GROUP BY 
    cluster_id
),
count_namespace AS (
  SELECT 
    cluster_id,
    COUNT(id) AS total_namespace
  FROM 
    object
  WHERE
    kind = 'Namespace'
  GROUP BY 
    cluster_id
),
count_pod AS (
  SELECT
    cluster_id, 
    COUNT(id) AS total_pod
  FROM
    object
  WHERE
    kind = 'Pod'
  GROUP BY cluster_id
),
count_container AS (
  SELECT
    cluster_id,
    SUM(jsonb_array_length(raw -> 'spec' -> 'containers')) AS total_container
  FROM
    object
  WHERE
    kind = 'Pod'
  GROUP BY cluster_id
),
count_deprecated_apis AS (
  SELECT
    cluster_id, 
    jsonb_array_length(deprecated_apis) AS total_deprecated_apis
  FROM
    deprecated_api_group
  GROUP BY cluster_id
)
SELECT
  COALESCE(n.cluster_id, p.cluster_id, c.cluster_id, d.cluster_id)::uuid AS cluster_id,
  COALESCE(n.total_node, 0)::BIGINT AS total_node,
  COALESCE(ns.total_namespace, 0)::BIGINT AS total_namespace,
  COALESCE(p.total_pod, 0)::BIGINT AS total_pod,
  COALESCE(c.total_container, 0)::BIGINT AS total_container,
  COALESCE(d.total_deprecated_apis, 0)::BIGINT total_deprecated_apis
FROM count_node n
FULL OUTER JOIN count_namespace ns ON n.cluster_id = ns.cluster_id
FULL OUTER JOIN count_pod p ON n.cluster_id = p.cluster_id
FULL OUTER JOIN count_container c ON COALESCE(n.cluster_id, p.cluster_id) = c.cluster_id
FULL OUTER JOIN count_deprecated_apis d ON COALESCE(n.cluster_id, p.cluster_id, c.cluster_id) = d.cluster_id
WITH DATA;


CREATE UNIQUE INDEX IF NOT EXISTS uniq_idx_cluster_resource_count_cluster_id 
  ON cluster_resource_count (cluster_id);


END;
