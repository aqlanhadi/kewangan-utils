-- Revert myduit:expenses from pg

BEGIN;

-- XXX Add DDLs here.
DROP TABLE myduit.expenses;


COMMIT;
