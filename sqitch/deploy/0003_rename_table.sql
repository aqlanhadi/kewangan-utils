-- Deploy myduit:0003_rename_table to pg
-- requires: 0002_expenses_table

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE myduit.expenses RENAME TO "transaction";

COMMIT;
