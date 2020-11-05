#!/usr/bin/env bash

USER="dazwilkin" # Or your GitHub username
REPO="akri-http" # Or your preferred GHCR repo
TAGS="$(git rev-parse HEAD)"

echo "Ignoring cmd/client and cmd/server"

for IMAGE in "devices" "device" "discovery"
do
  docker build \
  --tag=ghcr.io/${USER}/${REPO}-${IMAGE}:${TAGS} \
  --file=./cmd/${IMAGE}/Dockerfile \
  .
done
