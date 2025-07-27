#!/bin/sh

docker compose -f scripts/sut/mongodb.yaml -f scripts/sut/kafka.yaml down --remove-orphans --volumes

# Remove all subfolders in the storage directory
rm -rf producer/storage

# Remove the resume token
rm -f miner/app/data/resume_token.bin
