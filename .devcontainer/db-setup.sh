#!/bin/sh

# Runs during MariaDB first-time initialization (docker-entrypoint-initdb.d).
# Creates the schema, applies migrations, and inserts development test data.

set -e

SCRIPTS_DIR="/docker-entrypoint-initdb.d/scripts"

echo "Initializing database ..."
mariadb -u root -p"$MYSQL_ROOT_PASSWORD" < $SCRIPTS_DIR/db-init.sql

echo "Creating schema ..."
mariadb -u root -p"$MYSQL_ROOT_PASSWORD" work_log < $SCRIPTS_DIR/db/db_create.sql

for f in $SCRIPTS_DIR/db/db_update_v*.sql; do
  i=$(echo "$f" | sed 's/.*db_update_v\([0-9]*\)\.sql/\1/')
  echo "Applying migration v${i} ..."
  mariadb -u root -p"$MYSQL_ROOT_PASSWORD" work_log < "$f"
  mariadb -u root -p"$MYSQL_ROOT_PASSWORD" work_log \
    -e "UPDATE setting SET setting_value = '${i}' WHERE setting_key = 'db_version';"
done

echo "Inserting test data ..."
mariadb -u root -p"$MYSQL_ROOT_PASSWORD" work_log < $SCRIPTS_DIR/db-dev-setup.sql

echo "Database setup complete."
