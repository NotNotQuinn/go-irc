#!/usr/bin/env bash

# exit when any command fails
set -e

# keep track of the last executed command
trap 'last_command=$current_command; current_command=$BASH_COMMAND' DEBUG
# echo an error message before exiting
trap 'echo "\"${last_command}\" command failed with exit code $?."' EXIT

# Wait untill mariadb is running
echo "Waiting for database to open..."
while ! nc -z mariadb 3306; do
    sleep 0.1
done

# Populate the database with test data.
# TODO(possibly): Make populator accept an option to omit data.
# Poissbly could be used here.
echo "Populating database..."
/bot/bin/populator

cd /bot
if [ "$WB_TEST" == "true" ]; then
    # Run all tests.
    echo "Running tests..."
    go test ./...
else
    echo "Starting bot..."
    /bot/bin/wanductbot
fi