-- Revert myduit:0001_appschema from pg

BEGIN;

-- XXX Add DDLs here.
DROP SCHEMA IF EXISTS myduit;

COMMIT;
