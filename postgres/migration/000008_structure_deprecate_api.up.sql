BEGIN;


DROP MATERIALIZED VIEW IF EXISTS cluster_resource_count;
DROP TABLE IF EXISTS deprecated_api_group;


CREATE TABLE IF NOT EXISTS deprecated_api (
    cluster_id UUID REFERENCES cluster(id),
    organization_id UUID NOT NULL REFERENCES organization(id),
    current_group_version TEXT,
    kind TEXT,
    name TEXT NOT NULL,
    cluster_k8s_version TEXT NOT NULL,
    deprecated BOOLEAN NOT NULL,
    deprecated_in TEXT NOT NULL,
    removed_in TEXT NOT NULL,
    replacement_version TEXT NOT NULL,
    PRIMARY KEY (cluster_id, current_group_version, kind)
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
count_deprecated_api AS (
  SELECT
    cluster_id,
    count(*) AS total_deprecated_api 
  FROM
    deprecated_api
  GROUP BY cluster_id
)
SELECT
  COALESCE(n.cluster_id, p.cluster_id, c.cluster_id, d.cluster_id)::uuid AS cluster_id,
  COALESCE(n.total_node, 0)::BIGINT AS total_node,
  COALESCE(ns.total_namespace, 0)::BIGINT AS total_namespace,
  COALESCE(p.total_pod, 0)::BIGINT AS total_pod,
  COALESCE(c.total_container, 0)::BIGINT AS total_container,
  COALESCE(d.total_deprecated_api, 0)::BIGINT total_deprecated_api
FROM count_node n
FULL OUTER JOIN count_namespace ns ON n.cluster_id = ns.cluster_id
FULL OUTER JOIN count_pod p ON n.cluster_id = p.cluster_id
FULL OUTER JOIN count_container c ON COALESCE(n.cluster_id, p.cluster_id) = c.cluster_id
FULL OUTER JOIN count_deprecated_api d ON COALESCE(n.cluster_id, p.cluster_id, c.cluster_id) = d.cluster_id
WITH DATA;


CREATE UNIQUE INDEX IF NOT EXISTS uniq_idx_cluster_resource_count_cluster_id 
  ON cluster_resource_count (cluster_id);


END;
