#! /bin/sh

docker compose -f scripts/sut/mongodb.yaml -f scripts/sut/kafka.yaml up