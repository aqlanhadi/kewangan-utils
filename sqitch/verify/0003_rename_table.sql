-- Verify myduit:0003_rename_table on pg

BEGIN;

-- XXX Add verifications here.
SELECT * FROM myduit.transaction WHERE FALSE;

ROLLBACK;
