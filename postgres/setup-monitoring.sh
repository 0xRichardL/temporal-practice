#!/bin/sh

set -eu

: "${POSTGRES_HOST:?POSTGRES_HOST is required}"
: "${POSTGRES_PORT:?POSTGRES_PORT is required}"
: "${POSTGRES_USER:?POSTGRES_USER is required}"
: "${PGPASSWORD:?PGPASSWORD is required}"
: "${POSTGRES_MONITOR_PASSWORD:?POSTGRES_MONITOR_PASSWORD is required}"

monitor_user=alloy_monitor

psql \
  --host "$POSTGRES_HOST" \
  --port "$POSTGRES_PORT" \
  --username "$POSTGRES_USER" \
  --dbname postgres \
  --set ON_ERROR_STOP=1 \
  --set monitor_user="$monitor_user" \
  --set monitor_password="$POSTGRES_MONITOR_PASSWORD" <<'SQL'
SELECT format('CREATE ROLE %I LOGIN PASSWORD %L', :'monitor_user', :'monitor_password')
WHERE NOT EXISTS (
  SELECT 1
  FROM pg_catalog.pg_roles
  WHERE rolname = :'monitor_user'
) \gexec

SELECT format('ALTER ROLE %I WITH LOGIN NOSUPERUSER NOCREATEDB NOCREATEROLE NOREPLICATION PASSWORD %L', :'monitor_user', :'monitor_password') \gexec
SELECT format('GRANT pg_monitor TO %I', :'monitor_user') \gexec
SQL

for database in postgres temporal temporal_visibility; do
  psql \
    --host "$POSTGRES_HOST" \
    --port "$POSTGRES_PORT" \
    --username "$POSTGRES_USER" \
    --dbname "$database" \
    --set ON_ERROR_STOP=1 \
    --set monitor_user="$monitor_user" <<'SQL'
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;
SELECT format('GRANT CONNECT ON DATABASE %I TO %I', current_database(), :'monitor_user') \gexec
SQL
done

echo "PostgreSQL monitoring role and extensions are ready"
