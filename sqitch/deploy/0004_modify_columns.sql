-- Deploy myduit:0004_modify_columns to pg
-- requires: 0003_rename_table

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE myduit.transaction ADD account_number text;
ALTER TABLE myduit.transaction ADD account_type text;
ALTER TABLE myduit.transaction ADD source text;

COMMIT;
