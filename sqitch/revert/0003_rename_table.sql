-- Revert myduit:0003_rename_table from pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE myduit.transaction RENAME TO "expenses";

COMMIT;
