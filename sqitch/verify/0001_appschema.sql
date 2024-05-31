-- Verify myduit:0001_appschema on pg

BEGIN;

-- XXX Add verifications here.
SELECT pg_catalog.has_schema_privilege('myduit', 'usage');

ROLLBACK;
