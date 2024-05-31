-- Deploy myduit:0001_appschema to pg

BEGIN;

-- XXX Add DDLs here.
CREATE SCHEMA IF NOT EXISTS myduit;

COMMIT;
