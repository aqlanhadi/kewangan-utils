-- Verify myduit:0002_expenses_table on pg

BEGIN;

-- XXX Add verifications here.
SELECT 
    id,
    account,
    posting_date,
    "date",
    "action",
    beneficiary,
    "description",
    method,
    amount,
    balance 
FROM myduit.expenses
WHERE FALSE;

ROLLBACK;
