-- Deploy myduit:expenses to pg
-- requires: appschema

BEGIN;

-- XXX Add DDLs here.
SET client_min_messages = 'warning';

CREATE TABLE myduit.expenses (
    id SERIAL PRIMARY KEY,
    origin TEXT
);

COMMIT;
