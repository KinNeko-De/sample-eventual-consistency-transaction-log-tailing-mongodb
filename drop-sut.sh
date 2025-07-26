#!/bin/sh

docker compose -f scripts/sut/mongodb.yaml down --remove-orphans --volumes

# Remove all subfolders in the storage directory
rm -rf producer/storage

