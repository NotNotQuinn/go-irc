#/usr/bin/env bash

[[ ! -d config ]] && \
    echo "Please check your working directory: Not a directory: ./config" && \
    exit 1

echo "Creating backup files..."
[[ -d "config/.backup" ]] && \
    echo "Failed to create backup: ./config/.backup/ already exists." && \
    echo "Please delete the existing backup first." && \
    exit 1

mkdir config/.backup
cp config/*.json config/.backup/
[[ -d "config/PRODUCTION" ]] && cp config/PRODUCTION config/.backup/
echo "Config backup is in ./config/.backup/"

for var in "$@"
do
    if [ $var == "--test=true" ]
    then
        echo "Moving tests config to main config..."
        mv ./config/tests_private_conf.json ./config/private_conf.json
    fi
done

echo "All complete - glhf."
exit 0