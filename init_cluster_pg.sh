#!/bin/bash
psql -d postgres <<EOF

\c gecko_local

CREATE TABLE IF NOT EXISTS documents (
    name VARCHAR(255) PRIMARY KEY,
    content JSONB
);

EOF
