-- Revert myduit:0004_modify_columns from pg

BEGIN;

-- XXX Add DDLs here.

ALTER TABLE myduit.transaction DROP COLUMN account_number;
ALTER TABLE myduit.transaction DROP COLUMN account_type;
ALTER TABLE myduit.transaction DROP COLUMN source;

COMMIT;
