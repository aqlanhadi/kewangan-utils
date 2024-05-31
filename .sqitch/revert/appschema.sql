-- Revert myduit:appschema from pg

BEGIN;

-- XXX Add DDLs here.
DROP SCHEMA myduit;

COMMIT;
