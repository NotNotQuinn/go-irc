#!/usr/bin/env bash

# exit when any command fails
set -e

# keep track of the last executed command
trap 'last_command=$current_command; current_command=$BASH_COMMAND' DEBUG
# echo an error message before exiting
trap 'echo "\"${last_command}\" command filed with exit code $?."' EXIT


# Wait untill mariadb is running
while ! nc -z mariadb 3306; do
    sleep 0.1
done

# Populate the database with test data.
cd /bot
/bot/bin/populator
# Run all tests.
go test ./...
