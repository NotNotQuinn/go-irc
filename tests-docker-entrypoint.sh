#!/usr/bin/env bash

# exit when any command fails
set -e

# keep track of the last executed command
trap 'last_command=$current_command; current_command=$BASH_COMMAND' DEBUG
# echo an error message before exiting
trap 'echo "\"${last_command}\" command filed with exit code $?."' EXIT


# I dont know if mariadb will always be working before this starts.
sleep 3
# I dont know if WORKDIR is preserved when executing the entrypoint.
cd /bot
# Populate the database with test data.
/bot/bin/populator
# Run all tests.
go test ./...
