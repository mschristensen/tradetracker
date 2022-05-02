#!/bin/bash
source ./scripts/env.sh $1
shift

psql -U postgres postgres -c "CREATE USER $POSTGRES_USER WITH SUPERUSER LOGIN PASSWORD '$POSTGRES_PASSWORD'"
psql -U postgres postgres -c "ALTER ROLE $POSTGRES_USER SUPERUSER"
psql -U postgres postgres -c "CREATE DATABASE $POSTGRES_DATABASE OWNER $POSTGRES_USER"
