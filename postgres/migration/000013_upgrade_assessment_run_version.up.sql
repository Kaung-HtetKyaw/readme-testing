BEGIN;


ALTER TABLE IF EXISTS upgradeable_component 
ADD COLUMN run_version UUID NOT NULL DEFAULT gen_random_uuid();


ALTER TABLE IF EXISTS unmatched_component 
ADD COLUMN run_version UUID NOT NULL DEFAULT gen_random_uuid();


UPDATE upgradeable_component AS uc
SET 
  run_version = o.run_version
FROM object AS o
WHERE uc.object_id = o.id;


UPDATE unmatched_component AS uc
SET 
  run_version = o.run_version
FROM object AS o
WHERE uc.object_id = o.id;
  

END;
