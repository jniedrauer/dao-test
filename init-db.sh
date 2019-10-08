#!/usr/bin/env bash

set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE TABLE test_table (
        id serial PRIMARY KEY,
        data VARCHAR (5) NOT NULL
    );
EOSQL
