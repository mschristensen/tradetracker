#!/bin/bash

source ./scripts/env.sh $1
shift

SCHEMA_DUMP=./internal/pkg/db/schema.sql

PGPASSWORD=$POSTGRES_PASSWORD pg_dump -s -d $POSTGRES_DATABASE -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER > $SCHEMA_DUMP

echo "Wrote schema dump to $SCHEMA_DUMP"
