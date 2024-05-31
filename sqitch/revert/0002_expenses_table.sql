-- Revert myduit:0002_expenses_table from pg

BEGIN;

-- XXX Add DDLs here.
DROP TABLE myduit.expenses;

COMMIT;
