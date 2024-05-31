-- Verify myduit:expenses on pg

BEGIN;

-- XXX Add verifications here.
SELECT id FROM myduit.expenses;

ROLLBACK;
