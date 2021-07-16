#!/usr/bin/env bash
# I dont know if mariadb will always be working before this starts.
sleep 3
# I dont know if WORKDIR is preserved when executing the entrypoint.
cd /bot
# Populate the database with test data.
bin/populator
# Run all tests.
go test ./...
