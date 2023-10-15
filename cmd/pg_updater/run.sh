#!/bin/bash

set -eux

# Expects env variables:
#
## How to reach postgres
#  - PGHOST
#  - PGPORT
#  - PGUSER
#  - PGPASS
#
## What we're going to set
#  - PGUSER_READWRITE_PASS
#  - PGUSER_READONLY_PASS

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
cd $SCRIPT_DIR

# Load env variables
PGHOST=${PGHOST:-localhost}
PGPORT=${PGPORT:-5432}
PGUSER=${PGUSER:-postgres}

# Template user passwords
RW_PASS=${PGUSER_READWRITE_PASS:-readwrite} \
RO_PASS=${PGUSER_READONLY_PASS:-readonly} \
envsubst <00.tmpl >00.sql

# Create users
PGPASSWORD=${PGPASS:-test} psql -h $PGHOST -p $PGPORT -U $PGUSER -d postgres -f 00.sql

# Create tables
PGPASSWORD=${PGPASS:-test} psql -h $PGHOST -p $PGPORT -U $PGUSER -d postgres -f 01.sql
