#!/bin/bash

# Start PostgreSQL service
brew services start postgresql

# Wait for PostgreSQL to start (adjust sleep time if needed)
sleep 5

psql postgres <<EOF

ALTER USER postgres WITH PASSWORD 'your_strong_password';

CREATE DATABASE testdb;

GRANT ALL PRIVILEGES ON DATABASE testdb TO postgres;

\c testdb;

CREATE TABLE IF NOT EXISTS documents (
    name VARCHAR(255) PRIMARY KEY,
    content JSONB
);

\q
EOF
