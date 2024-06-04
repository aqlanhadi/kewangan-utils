-- Verify myduit:0004_modify_columns on pg

BEGIN;

-- XXX Add verifications here.
SELECT 
    id,
    account,
    account_number,
    account_type,
    posting_date,
    "date",
    "action",
    beneficiary,
    "description",
    method,
    amount,
    balance,
    source
FROM myduit.transaction
WHERE FALSE;


ROLLBACK;
