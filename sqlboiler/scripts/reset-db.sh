#!/bin/bash

set -e
BASEDIR=$(dirname "$0")/..

echo "dropping db..."
dropdb -f --if-exists $PSQL_DBNAME

echo "creating db..."
createdb $PSQL_DBNAME

echo "generating models..."
psql -d $PSQL_DBNAME -q -f $BASEDIR/schema.sql
sqlboiler \
    --no-tests \
    --templates $BASEDIR/templates \
    --templates $GOPATH/pkg/mod/github.com/volatiletech/sqlboiler/v4@v4.6.0/templates \
    --templates $GOPATH/pkg/mod/github.com/volatiletech/sqlboiler/v4@v4.6.0/templates_test \
    psql
