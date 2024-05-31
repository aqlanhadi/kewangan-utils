-- Deploy myduit:0002_expenses_table to pg

BEGIN;

-- XXX Add DDLs here.
SET client_min_messages = 'warning';

CREATE TABLE myduit.expenses (
    id SERIAL PRIMARY KEY,

    account text,
    posting_date date,
    "date" date,
    "action" text,
    beneficiary text,
    "description" text,
    method text,

    amount money,
    balance money
);

COMMIT;
